package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Asset represents an asset following ISO 55001 standards
type Asset struct {
	ID uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`

	// Basic Information (ISO 55001 - Asset Identification)
	Name         string      `json:"name" gorm:"type:varchar(255);not null"`
	Description  string      `json:"description" gorm:"type:text"`
	CategoryID   *uuid.UUID  `json:"category_id" gorm:"type:uuid"`
	Category     *Category   `json:"category,omitempty" gorm:"foreignKey:CategoryID"`
	DepartmentID *uuid.UUID  `json:"department_id" gorm:"type:uuid"`
	Department   *Department `json:"department,omitempty" gorm:"foreignKey:DepartmentID"`

	// Technical Specifications
	Type         string `json:"type" gorm:"type:varchar(100)"`
	Model        string `json:"model" gorm:"type:varchar(100)"`
	SerialNumber string `json:"serial_number" gorm:"type:varchar(100);unique"`
	Manufacturer string `json:"manufacturer" gorm:"type:varchar(100)"`

	// Financial Information
	AcquisitionCost  float64 `json:"acquisition_cost" gorm:"type:decimal(15,2)"`
	CurrentValue     float64 `json:"current_value" gorm:"type:decimal(15,2)"`
	DepreciationRate float64 `json:"depreciation_rate" gorm:"type:decimal(5,2)"`

	// Operational Status
	Status      string `json:"status" gorm:"type:varchar(50);default:'active';check:status IN ('active', 'inactive', 'maintenance', 'disposed')"`
	Condition   string `json:"condition" gorm:"type:varchar(50);default:'good';check:condition IN ('excellent', 'good', 'fair', 'poor', 'critical')"`
	Criticality string `json:"criticality" gorm:"type:varchar(50);default:'low';check:criticality IN ('low', 'medium', 'high', 'critical')"`

	// Location Information
	Latitude     *float64 `json:"latitude" gorm:"type:decimal(10,8)"`
	Longitude    *float64 `json:"longitude" gorm:"type:decimal(11,8)"`
	Address      string   `json:"address" gorm:"type:text"`
	BuildingRoom string   `json:"building_room" gorm:"type:varchar(100)"`

	// Lifecycle Information
	AcquisitionDate     *time.Time `json:"acquisition_date" gorm:"type:date"`
	ExpectedLifeYears   *int       `json:"expected_life_years" gorm:"type:integer"`
	MaintenanceSchedule string     `json:"maintenance_schedule" gorm:"type:text"`

	// Compliance and Standards
	Certifications string `json:"certifications" gorm:"type:text"`
	Standards      string `json:"standards" gorm:"type:text"`
	AuditInfo      string `json:"audit_info" gorm:"type:text"`

	// Metadata
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`

	// Relationships
	// Category Category `json:"category,omitempty" gorm:"foreignKey:CategoryID"`
}

// BeforeCreate will set a UUID rather than numeric ID
func (a *Asset) BeforeCreate(tx *gorm.DB) error {
	if a.ID == uuid.Nil {
		a.ID = uuid.New()
	}
	return nil
}

// TableName specifies the table name for Asset
func (Asset) TableName() string {
	return "assets"
}

// AssetSummary represents a simplified view of an asset for AI queries
type AssetSummary struct {
	ID                uuid.UUID  `json:"id"`
	Name              string     `json:"name"`
	Description       string     `json:"description"`
	CategoryName      string     `json:"category_name"`
	Type              string     `json:"type"`
	Model             string     `json:"model"`
	SerialNumber      string     `json:"serial_number"`
	Status            string     `json:"status"`
	Condition         string     `json:"condition"`
	Criticality       string     `json:"criticality"`
	Latitude          *float64   `json:"latitude"`
	Longitude         *float64   `json:"longitude"`
	Address           string     `json:"address"`
	BuildingRoom      string     `json:"building_room"`
	CurrentValue      float64    `json:"current_value"`
	AcquisitionDate   *time.Time `json:"acquisition_date"`
	ExpectedLifeYears *int       `json:"expected_life_years"`
}
