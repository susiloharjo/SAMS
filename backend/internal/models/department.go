package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Department represents an organizational department
type Department struct {
	ID          uuid.UUID      `gorm:"type:uuid;primary_key;" json:"id"`
	Name        string         `gorm:"type:varchar(100);not null;unique" json:"name"`
	Description string         `gorm:"type:text" json:"description"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

// BeforeCreate will set a UUID rather than numeric ID.
func (d *Department) BeforeCreate(tx *gorm.DB) (err error) {
	d.ID = uuid.New()
	return
}

// TableName returns the name of the table for the Department model
func (d *Department) TableName() string {
	return "departments"
}
