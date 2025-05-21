package models

import (
	"fmt"
	"time"

	"github.com/assimoes/beautix/internal/infrastructure/database"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// BaseModel defines common fields for all models
type BaseModel struct {
	ID        uuid.UUID      `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	CreatedAt time.Time      `gorm:"not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time      `gorm:"not null;default:CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	CreatedBy *uuid.UUID     `json:"created_by,omitempty"`
	UpdatedBy *uuid.UUID     `json:"updated_by,omitempty"`
	DeletedBy *uuid.UUID     `json:"deleted_by,omitempty"`
}

// BeforeCreate hook that automatically sets ID if not provided
func (b *BaseModel) BeforeCreate(tx *gorm.DB) error {
	if b.ID == uuid.Nil {
		b.ID = uuid.New()
	}
	return nil
}

// TenantModel extends BaseModel with business_id for multi-tenancy
type TenantModel struct {
	BaseModel
	BusinessID uuid.UUID `gorm:"type:uuid;not null;index" json:"business_id"`
}

// AuditModel is used for models that require detailed audit information
type AuditModel struct {
	BaseModel
	LastAccessedAt *time.Time  `json:"last_accessed_at,omitempty"`
	LastAccessedBy *uuid.UUID  `json:"last_accessed_by,omitempty"`
}

// RegisterModels registers all models with GORM
func RegisterModels(db interface{}) error {
	var gormDB *gorm.DB
	
	// Handle different types of database connections
	switch dbConn := db.(type) {
	case *gorm.DB:
		gormDB = dbConn
	case *database.DB:
		gormDB = dbConn.DB
	default:
		return fmt.Errorf("invalid database connection type")
	}
	
	// Auto-migrate models - only for development and testing
	// In production, changes should be made via proper migrations
	err := gormDB.AutoMigrate(
		&User{},
		&UserConnectedAccount{},
		&Business{},
		&BusinessLocation{},
		&Staff{},
		&ServiceAssignment{},
		&AvailabilityException{},
		&StaffPerformance{},
		&ServiceCategory{},
		&Service{},
		&ServiceVariant{},
		&ServiceOption{},
		&ServiceOptionChoice{},
		&ServiceBundle{},
		&ServiceBundleItem{},
		&Client{},
		&ClientNote{},
		&ClientDocument{},
		&Appointment{},
		&AppointmentService{},
		&AppointmentPayment{},
		&AppointmentForm{},
		&AppointmentNote{},
		&AppointmentReminder{},
		&AppointmentFeedback{},
	)
	if err != nil {
		return err
	}

	// Note: Custom indexes should be defined in migrations rather than in code
	// For example, see migrations/000007_user_indexes.up.sql

	return nil
}