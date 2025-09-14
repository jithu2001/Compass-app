package models

import (
	"time"

	"gorm.io/gorm"
)

type AnswerValue string

const (
	AnswerYes AnswerValue = "yes"
	AnswerNo  AnswerValue = "no"
)

type ProjectRFI struct {
	RFIID        uint64       `gorm:"primaryKey;autoIncrement" json:"rfi_id"`
	ProjectID    uint64       `gorm:"not null" json:"project_id"`
	Project      *Project     `gorm:"foreignKey:ProjectID" json:"project,omitempty"`
	QuestionText string       `gorm:"type:text;not null" json:"question_text"`
	AnswerValue  *AnswerValue `gorm:"type:varchar(10);check:answer_value IN ('yes','no')" json:"answer_value,omitempty"`
	AnsweredBy   *uint64      `json:"answered_by,omitempty"`
	Answerer     *User        `gorm:"foreignKey:AnsweredBy" json:"answerer,omitempty"`
	CreatedAt    time.Time    `json:"created_at"`
}

func (ProjectRFI) TableName() string {
	return "project_rfis"
}

func (rfi *ProjectRFI) BeforeCreate(tx *gorm.DB) error {
	rfi.CreatedAt = time.Now()
	return nil
}