package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	Username  string         `gorm:"uniqueIndex;not null" json:"username"`
	Email     string         `gorm:"uniqueIndex;not null" json:"email"`
	Password  string         `gorm:"not null" json:"-"`
	IsAdmin   bool           `gorm:"default:false" json:"is_admin"`
}

type Share struct {
	ID          uint           `gorm:"primarykey" json:"id"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
	ShareCode   string         `gorm:"uniqueIndex;not null" json:"share_code"`
	UserID      uint           `json:"user_id"`
	User        User           `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Type        string         `gorm:"not null" json:"type"` // "file" or "text"
	FileName    string         `json:"file_name,omitempty"`
	FileSize    int64          `json:"file_size,omitempty"`
	FilePath    string         `json:"-"`
	TextContent string         `gorm:"type:text" json:"text_content,omitempty"`
	Downloads   int            `gorm:"default:0" json:"downloads"`
	ExpiresAt   *time.Time     `json:"expires_at,omitempty"`
}
