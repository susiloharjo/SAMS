package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// User represents a system user
type User struct {
	ID           uuid.UUID      `gorm:"type:uuid;primary_key;" json:"id"`
	Username     string         `gorm:"type:varchar(50);not null;unique" json:"username"`
	Email        string         `gorm:"type:varchar(100);not null;unique" json:"email"`
	FirstName    string         `gorm:"type:varchar(50);not null" json:"first_name"`
	LastName     string         `gorm:"type:varchar(50);not null" json:"last_name"`
	Password     string         `gorm:"type:varchar(255);not null" json:"-"`                  // Hidden from JSON
	Role         string         `gorm:"type:varchar(20);not null;default:'user'" json:"role"` // admin, manager, user
	DepartmentID *uuid.UUID     `gorm:"type:uuid" json:"department_id"`
	Department   *Department    `gorm:"foreignKey:DepartmentID" json:"department,omitempty"`
	IsActive     bool           `gorm:"default:true" json:"is_active"`
	LastLogin    *time.Time     `json:"last_login"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

// BeforeCreate will set a UUID rather than numeric ID.
func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	u.ID = uuid.New()
	return
}

// TableName returns the name of the table for the User model
func (u *User) TableName() string {
	return "users"
}

// UserCreateRequest represents the data needed to create a new user
type UserCreateRequest struct {
	Username     string     `json:"username" validate:"required,min=3,max=50"`
	Email        string     `json:"email" validate:"required,email"`
	FirstName    string     `json:"first_name" validate:"required,min=2,max=50"`
	LastName     string     `json:"last_name" validate:"required,min=2,max=50"`
	Password     string     `json:"password" validate:"required,min=6"`
	Role         string     `json:"role" validate:"required,oneof=admin manager user"`
	DepartmentID *uuid.UUID `json:"department_id"`
}

// UserUpdateRequest represents the data needed to update a user
type UserUpdateRequest struct {
	Username     *string    `json:"username" validate:"omitempty,min=3,max=50"`
	Email        *string    `json:"email" validate:"omitempty,email"`
	FirstName    *string    `json:"first_name" validate:"omitempty,min=2,max=50"`
	LastName     *string    `json:"last_name" validate:"omitempty,min=2,max=50"`
	Role         *string    `json:"role" validate:"omitempty,oneof=admin manager user"`
	DepartmentID *uuid.UUID `json:"department_id"`
	IsActive     *bool      `json:"is_active"`
}

// UserResponse represents the user data sent to the frontend
type UserResponse struct {
	ID           uuid.UUID   `json:"id"`
	Username     string      `json:"username"`
	Email        string      `json:"email"`
	FirstName    string      `json:"first_name"`
	LastName     string      `json:"last_name"`
	Role         string      `json:"role"`
	DepartmentID *uuid.UUID  `json:"department_id"`
	Department   *Department `json:"department,omitempty"`
	IsActive     bool        `json:"is_active"`
	LastLogin    *time.Time  `json:"last_login"`
	CreatedAt    time.Time   `json:"created_at"`
	UpdatedAt    time.Time   `json:"updated_at"`
}

// UserPasswordUpdateRequest represents the request to update a user's password by an admin
type UserPasswordUpdateRequest struct {
	NewPassword string `json:"new_password" validate:"required,min=8"`
}
