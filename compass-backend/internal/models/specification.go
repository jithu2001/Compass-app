package models

import (
	"time"

	"gorm.io/gorm"
)

type ProjectSpecification struct {
	SpecificationID uint64    `gorm:"primaryKey;autoIncrement" json:"specification_id"`
	ProjectID       uint64    `gorm:"not null;uniqueIndex:idx_project_version" json:"project_id"`
	Project         *Project  `gorm:"foreignKey:ProjectID" json:"project,omitempty"`
	VersionNo       int       `gorm:"not null;uniqueIndex:idx_project_version" json:"version_no"`
	Colour          string    `gorm:"size:100" json:"colour"`
	Ironmongery     string    `gorm:"size:150" json:"ironmongery"`
	UValue          *float64  `gorm:"type:decimal(5,2)" json:"u_value"`
	GValue          *float64  `gorm:"type:decimal(5,2)" json:"g_value"`
	Vents           string    `gorm:"size:100" json:"vents"`
	Acoustics       string    `gorm:"size:100" json:"acoustics"`
	SBD             string    `gorm:"size:100" json:"sbd"`
	PAS24           string    `gorm:"size:100" json:"pas24"`
	Restrictors     string    `gorm:"size:100" json:"restrictors"`
	SpecialComments string    `gorm:"type:text" json:"special_comments"`
	AttachmentURL   string    `gorm:"type:text" json:"attachment_url"`
	CreatedBy       uint64    `gorm:"not null" json:"created_by"`
	Creator         *User     `gorm:"foreignKey:CreatedBy" json:"creator,omitempty"`
	CreatedAt       time.Time `json:"created_at"`
}

func (ProjectSpecification) TableName() string {
	return "project_specifications"
}

func (ps *ProjectSpecification) BeforeCreate(tx *gorm.DB) error {
	ps.CreatedAt = time.Now()
	
	// Auto-increment version number
	var lastVersion int
	tx.Model(&ProjectSpecification{}).
		Where("project_id = ?", ps.ProjectID).
		Select("COALESCE(MAX(version_no), 0)").
		Scan(&lastVersion)
	
	ps.VersionNo = lastVersion + 1
	return nil
}