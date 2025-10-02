package models

import (
	"time"

	"gorm.io/gorm"
)

type ProjectStatus string
type ProjectType string

const (
	StatusNotYetStarted ProjectStatus = "not_yet_started"
	StatusProgress      ProjectStatus = "progress"
	StatusCompleted     ProjectStatus = "completed"

	TypeWindows ProjectType = "windows"
	TypeDoors   ProjectType = "doors"
)

type Project struct {
	ProjectID      uint64                  `gorm:"primaryKey;autoIncrement" json:"project_id"`
	ProjectName    string                  `gorm:"size:200;not null" json:"project_name"`
	CompanyName    string                  `gorm:"size:200" json:"company_name"`
	CompanyAddress string                  `gorm:"type:text" json:"company_address"`
	ProjectStatus  ProjectStatus           `gorm:"type:varchar(20);default:'not_yet_started';check:project_status IN ('not_yet_started','progress','completed')" json:"project_status"`
	ProjectType    ProjectType             `gorm:"type:varchar(20);check:project_type IN ('windows','doors')" json:"project_type"`
	CreatedBy      uint64                  `gorm:"not null" json:"created_by"`
	Creator        *User                   `gorm:"foreignKey:CreatedBy" json:"creator,omitempty"`
	LastUpdatedBy  *uint64                 `json:"last_updated_by,omitempty"`
	LastUpdater    *User                   `gorm:"foreignKey:LastUpdatedBy" json:"last_updater,omitempty"`
	Specifications []ProjectSpecification  `gorm:"foreignKey:ProjectID;references:ProjectID" json:"specifications,omitempty"`
	RFIs           []ProjectRFI            `gorm:"foreignKey:ProjectID;references:ProjectID" json:"rfis,omitempty"`
	CreatedAt      time.Time               `json:"created_at"`
	UpdatedAt      time.Time               `json:"updated_at"`
}

func (Project) TableName() string {
	return "projects"
}

func (p *Project) BeforeCreate(tx *gorm.DB) error {
	p.CreatedAt = time.Now()
	p.UpdatedAt = time.Now()
	return nil
}

func (p *Project) BeforeUpdate(tx *gorm.DB) error {
	p.UpdatedAt = time.Now()
	return nil
}