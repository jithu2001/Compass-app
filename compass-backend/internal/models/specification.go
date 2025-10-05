package models

import (
	"time"

	"gorm.io/gorm"
)

type ProjectSpecification struct {
	SpecificationID     uint64    `gorm:"primaryKey;autoIncrement" json:"specification_id"`
	ProjectID           uint64    `gorm:"not null;uniqueIndex:idx_project_version" json:"project_id"`
	Project             *Project  `gorm:"foreignKey:ProjectID;references:ProjectID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"project,omitempty"`
	VersionNo           int       `gorm:"not null;uniqueIndex:idx_project_version" json:"version_no"`
	Colour              string    `gorm:"size:100" json:"colour"`
	ColourAttachment    *string   `gorm:"type:text" json:"colour_attachment,omitempty"`
	Ironmongery         string    `gorm:"size:150" json:"ironmongery"`
	IronmongeryAttachment *string `gorm:"type:text" json:"ironmongery_attachment,omitempty"`
	UValue              *string   `gorm:"size:100" json:"u_value"`
	UValueAttachment    *string   `gorm:"type:text" json:"u_value_attachment,omitempty"`
	GValue              *string   `gorm:"size:100" json:"g_value"`
	GValueAttachment    *string   `gorm:"type:text" json:"g_value_attachment,omitempty"`
	Vents               string    `gorm:"size:100" json:"vents"`
	VentsAttachment     *string   `gorm:"type:text" json:"vents_attachment,omitempty"`
	Acoustics           string    `gorm:"size:100" json:"acoustics"`
	AcousticsAttachment *string   `gorm:"type:text" json:"acoustics_attachment,omitempty"`
	SBD                 string    `gorm:"size:100" json:"sbd"`
	SBDAttachment       *string   `gorm:"type:text" json:"sbd_attachment,omitempty"`
	PAS24               string    `gorm:"size:100" json:"pas24"`
	PAS24Attachment     *string   `gorm:"type:text" json:"pas24_attachment,omitempty"`
	Restrictors         string    `gorm:"size:100" json:"restrictors"`
	RestrictorsAttachment *string `gorm:"type:text" json:"restrictors_attachment,omitempty"`
	SpecialComments     string    `gorm:"type:text" json:"special_comments"`
	AttachmentURL       string    `gorm:"type:text" json:"attachment_url"`
	CreatedBy           uint64    `gorm:"not null" json:"created_by"`
	Creator             *User     `gorm:"foreignKey:CreatedBy" json:"creator,omitempty"`
	CreatedAt           time.Time `json:"created_at"`
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