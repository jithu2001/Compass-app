package models

import (
	"time"

	"gorm.io/gorm"
)

type UserRole string
type AccountStatus string

const (
	RoleAdmin UserRole = "admin"
	RoleUser  UserRole = "user"

	StatusPending  AccountStatus = "pending"
	StatusActive   AccountStatus = "active"
	StatusDisabled AccountStatus = "disabled"
)

type User struct {
	UserID        uint64        `gorm:"primaryKey;autoIncrement" json:"user_id"`
	FullName      string        `gorm:"size:150;not null" json:"full_name"`
	Email         string        `gorm:"size:150;uniqueIndex;not null" json:"email"`
	PasswordHash  *string       `gorm:"size:255" json:"-"`
	Role          UserRole      `gorm:"type:varchar(20);default:'user';check:role IN ('admin','user')" json:"role"`
	AccountStatus AccountStatus `gorm:"type:varchar(20);default:'pending';check:account_status IN ('pending','active','disabled')" json:"account_status"`
	InvitedBy     *uint64       `json:"invited_by,omitempty"`
	InvitedByUser *User         `gorm:"foreignKey:InvitedBy" json:"invited_by_user,omitempty"`
	CreatedAt     time.Time     `json:"created_at"`
	UpdatedAt     time.Time     `json:"updated_at"`
}

func (User) TableName() string {
	return "users"
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	u.CreatedAt = time.Now()
	u.UpdatedAt = time.Now()
	return nil
}

func (u *User) BeforeUpdate(tx *gorm.DB) error {
	u.UpdatedAt = time.Now()
	return nil
}