package domain

import (
	"context"
	"time"

	"gorm.io/gorm"
)

// BaseModel contains common fields for all entities
type BaseModel struct {
	ID        string         `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	CreatedAt time.Time      `gorm:"not null;default:now()" json:"created_at"`
	CreatedBy *string        `gorm:"type:uuid" json:"created_by,omitempty"`
	UpdatedAt time.Time      `gorm:"not null;default:now()" json:"updated_at"`
	UpdatedBy *string        `gorm:"type:uuid" json:"updated_by,omitempty"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	DeletedBy *string        `gorm:"type:uuid" json:"deleted_by,omitempty"`
}

// Entity interface that all domain models should implement
type Entity interface {
	GetID() string
	GetCreatedAt() time.Time
	GetUpdatedAt() time.Time
	IsDeleted() bool
}

// GetID returns the entity ID
func (b BaseModel) GetID() string {
	return b.ID
}

// GetCreatedAt returns the creation timestamp
func (b BaseModel) GetCreatedAt() time.Time {
	return b.CreatedAt
}

// GetUpdatedAt returns the last update timestamp
func (b BaseModel) GetUpdatedAt() time.Time {
	return b.UpdatedAt
}

// IsDeleted returns true if the entity is soft deleted
func (b BaseModel) IsDeleted() bool {
	return b.DeletedAt.Valid
}

// SetAuditFields sets the audit fields for creation
func (b *BaseModel) SetAuditFields(userID *string) {
	now := time.Now()
	if b.ID == "" {
		// New entity
		b.CreatedAt = now
		b.CreatedBy = userID
	}
	b.UpdatedAt = now
	b.UpdatedBy = userID
}

// SetDeletedFields sets the deleted audit fields
func (b *BaseModel) SetDeletedFields(userID *string) {
	b.DeletedAt = gorm.DeletedAt{Time: time.Now(), Valid: true}
	b.DeletedBy = userID
}

// BaseRepository defines common repository operations
type BaseRepository[T any] interface {
	// Basic CRUD operations
	Create(ctx context.Context, entity *T) error
	GetByID(ctx context.Context, id string) (*T, error)
	Update(ctx context.Context, entity *T) error
	Delete(ctx context.Context, id string) error
	
	// List and search operations
	List(ctx context.Context, page, pageSize int) ([]*T, int64, error)
	FindBy(ctx context.Context, criteria map[string]any) ([]*T, error)
	ExistsByID(ctx context.Context, id string) (bool, error)
	
	// Transaction operations
	WithTx(tx *gorm.DB) BaseRepository[T]
	GetDB() *gorm.DB
}

// BaseService defines common service operations
type BaseService[CreateDTO, UpdateDTO, ResponseDTO any] interface {
	Create(ctx context.Context, dto CreateDTO) (*ResponseDTO, error)
	GetByID(ctx context.Context, id string) (*ResponseDTO, error)
	Update(ctx context.Context, id string, dto UpdateDTO) (*ResponseDTO, error)
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, page, pageSize int) ([]*ResponseDTO, int64, error)
}