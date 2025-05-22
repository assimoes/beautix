package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// DBAdapter interface abstracts database operations to support both regular DB and transactions
type DBAdapter interface {
	WithContext(ctx context.Context) *gorm.DB
}

// BaseRepository provides common functionality for all repositories
type BaseRepository struct {
	db DBAdapter
}

// NewBaseRepository creates a new base repository
func NewBaseRepository(db DBAdapter) *BaseRepository {
	return &BaseRepository{db: db}
}

// DB returns the database connection
func (r *BaseRepository) DB() DBAdapter {
	return r.db
}

// WithContext returns the database connection with context
func (r *BaseRepository) WithContext(ctx context.Context) *gorm.DB {
	return r.db.WithContext(ctx)
}

// HandleNotFound converts GORM's record not found error to a more descriptive error
func (r *BaseRepository) HandleNotFound(err error, entityType string, id uuid.UUID) error {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return fmt.Errorf("%s with id %s not found", entityType, id.String())
	}
	return err
}

// ApplyAuditFields applies standard audit fields for create operations
func (r *BaseRepository) ApplyAuditFields(model interface{}, createdBy *uuid.UUID) {
	if auditable, ok := model.(AuditableCreate); ok {
		auditable.SetCreatedBy(createdBy)
		auditable.SetCreatedAt(time.Now())
	}
}

// ApplyUpdateAuditFields applies standard audit fields for update operations
func (r *BaseRepository) ApplyUpdateAuditFields(model interface{}, updatedBy uuid.UUID) {
	if auditable, ok := model.(AuditableUpdate); ok {
		now := time.Now()
		auditable.SetUpdatedBy(&updatedBy)
		auditable.SetUpdatedAt(&now)
	}
}

// ApplyDeleteAuditFields applies standard audit fields for delete operations
func (r *BaseRepository) ApplyDeleteAuditFields(model interface{}, deletedBy uuid.UUID) {
	if auditable, ok := model.(AuditableDelete); ok {
		now := time.Now()
		auditable.SetDeletedBy(&deletedBy)
		auditable.SetDeletedAt(&now)
	}
}

// CalculateOffset calculates the offset for pagination
func (r *BaseRepository) CalculateOffset(page, pageSize int) int {
	if page <= 0 {
		page = 1
	}
	return (page - 1) * pageSize
}

// CreateWithAudit creates a new record with audit information
func (r *BaseRepository) CreateWithAudit(ctx context.Context, model interface{}, createdBy *uuid.UUID) error {
	r.ApplyAuditFields(model, createdBy)
	return r.WithContext(ctx).Create(model).Error
}

// UpdateWithAudit updates a record with audit information
func (r *BaseRepository) UpdateWithAudit(ctx context.Context, model interface{}, updates interface{}, updatedBy uuid.UUID) error {
	// First apply audit fields to the updates
	if updateMap, ok := updates.(map[string]interface{}); ok {
		updateMap["updated_by"] = &updatedBy
		updateMap["updated_at"] = time.Now()
	} else {
		// If updates is a struct, apply audit fields
		r.ApplyUpdateAuditFields(updates, updatedBy)
	}
	
	return r.WithContext(ctx).Model(model).Updates(updates).Error
}

// SoftDeleteWithAudit performs soft delete with audit information
func (r *BaseRepository) SoftDeleteWithAudit(ctx context.Context, model interface{}, deletedBy uuid.UUID) error {
	r.ApplyDeleteAuditFields(model, deletedBy)
	return r.WithContext(ctx).Save(model).Error
}

// AuditableCreate interface for entities that support creation audit fields
type AuditableCreate interface {
	SetCreatedBy(*uuid.UUID)
	SetCreatedAt(time.Time)
}

// AuditableUpdate interface for entities that support update audit fields
type AuditableUpdate interface {
	SetUpdatedBy(*uuid.UUID)
	SetUpdatedAt(*time.Time)
}

// AuditableDelete interface for entities that support delete audit fields
type AuditableDelete interface {
	SetDeletedBy(*uuid.UUID)
	SetDeletedAt(*time.Time)
}

// Auditable interface combines all audit interfaces
type Auditable interface {
	AuditableCreate
	AuditableUpdate
	AuditableDelete
}

// StandardCRUD defines common CRUD operations
type StandardCRUD[T any] interface {
	Create(ctx context.Context, entity *T) error
	GetByID(ctx context.Context, id uuid.UUID) (*T, error)
	Update(ctx context.Context, id uuid.UUID, updates interface{}, updatedBy uuid.UUID) error
	Delete(ctx context.Context, id uuid.UUID, deletedBy uuid.UUID) error
	Count(ctx context.Context) (int64, error)
}

// StandardList defines common list operations
type StandardList[T any] interface {
	List(ctx context.Context, page, pageSize int) ([]*T, error)
}

// Paginated represents a paginated result
type Paginated[T any] struct {
	Data       []*T `json:"data"`
	Total      int64 `json:"total"`
	Page       int   `json:"page"`
	PageSize   int   `json:"page_size"`
	TotalPages int   `json:"total_pages"`
}

// NewPaginated creates a new paginated result
func NewPaginated[T any](data []*T, total int64, page, pageSize int) *Paginated[T] {
	totalPages := int((total + int64(pageSize) - 1) / int64(pageSize))
	return &Paginated[T]{
		Data:       data,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}
}