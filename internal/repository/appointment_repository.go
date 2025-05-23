package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/assimoes/beautix/internal/domain"
	"github.com/assimoes/beautix/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// AppointmentRepository implements the domain.AppointmentRepository interface using GORM
type AppointmentRepository struct {
	*BaseRepository
}

// NewAppointmentRepository creates a new instance of AppointmentRepository
func NewAppointmentRepository(db DBAdapter) domain.AppointmentRepository {
	return &AppointmentRepository{
		BaseRepository: NewBaseRepository(db),
	}
}

// Create creates a new appointment
func (r *AppointmentRepository) Create(ctx context.Context, appointment *domain.Appointment) error {
	appointmentModel := mapDomainToAppointment(appointment)

	err := r.WithContext(ctx).Create(&appointmentModel).Error
	if err != nil {
		return fmt.Errorf("failed to create appointment: %w", err)
	}

	// Update the domain entity with any generated fields
	appointment.ID = appointmentModel.ID
	appointment.CreatedAt = appointmentModel.CreatedAt

	return nil
}

// GetByID retrieves an appointment by ID
func (r *AppointmentRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Appointment, error) {
	var appointmentModel models.Appointment

	err := r.WithContext(ctx).
		Preload("Business").
		Preload("Client").
		Preload("Staff").
		Preload("Service").
		Where("deleted_at IS NULL").
		First(&appointmentModel, "id = ?", id).Error

	if err != nil {
		return nil, r.HandleNotFound(err, "appointment", id)
	}

	return mapAppointmentToDomain(&appointmentModel), nil
}

// Update updates an appointment
func (r *AppointmentRepository) Update(ctx context.Context, id uuid.UUID, input *domain.UpdateAppointmentInput, updatedBy uuid.UUID) error {
	// First find the appointment to ensure it exists
	var appointmentModel models.Appointment
	err := r.WithContext(ctx).
		Where("deleted_at IS NULL").
		First(&appointmentModel, "id = ?", id).Error
	if err != nil {
		return r.HandleNotFound(err, "appointment", id)
	}

	// Apply updates from the input
	updates := make(map[string]interface{})
	now := time.Now()

	if input.StartTime != nil {
		updates["start_time"] = *input.StartTime
	}

	if input.EndTime != nil {
		updates["end_time"] = *input.EndTime
	}

	if input.Status != nil {
		updates["status"] = *input.Status
	}

	if input.Notes != nil {
		updates["notes"] = *input.Notes
	}

	// Always set updated fields
	updates["updated_at"] = &now
	updates["updated_by"] = &updatedBy

	// Perform the update
	err = r.WithContext(ctx).Model(&appointmentModel).Updates(updates).Error
	if err != nil {
		return fmt.Errorf("failed to update appointment: %w", err)
	}

	return nil
}

// Delete soft deletes an appointment
func (r *AppointmentRepository) Delete(ctx context.Context, id uuid.UUID, deletedBy uuid.UUID) error {
	// First find the appointment to ensure it exists
	var appointmentModel models.Appointment
	err := r.WithContext(ctx).
		Where("deleted_at IS NULL").
		First(&appointmentModel, "id = ?", id).Error
	if err != nil {
		return r.HandleNotFound(err, "appointment", id)
	}

	// Perform soft delete
	now := time.Now()
	updates := map[string]interface{}{
		"deleted_at": &now,
		"deleted_by": &deletedBy,
		"updated_at": &now,
		"updated_by": &deletedBy,
	}

	err = r.WithContext(ctx).Model(&appointmentModel).Updates(updates).Error
	if err != nil {
		return fmt.Errorf("failed to delete appointment: %w", err)
	}

	return nil
}

// ListByBusiness retrieves a paginated list of appointments by business ID within a date range
func (r *AppointmentRepository) ListByBusiness(ctx context.Context, businessID uuid.UUID, startDate, endDate time.Time, page, pageSize int) ([]*domain.Appointment, error) {
	var appointmentModels []models.Appointment

	offset := r.CalculateOffset(page, pageSize)

	err := r.WithContext(ctx).
		Preload("Business").
		Preload("Client").
		Preload("Staff").
		Preload("Service").
		Where("business_id = ?", businessID).
		Where("start_time >= ? AND start_time <= ?", startDate, endDate).
		Where("deleted_at IS NULL").
		Offset(offset).
		Limit(pageSize).
		Order("start_time ASC").
		Find(&appointmentModels).Error

	if err != nil {
		return nil, fmt.Errorf("failed to list appointments by business: %w", err)
	}

	return mapAppointmentSliceToDomain(appointmentModels), nil
}

// ListByStaff retrieves a paginated list of appointments by staff ID within a date range
func (r *AppointmentRepository) ListByStaff(ctx context.Context, staffID uuid.UUID, startDate, endDate time.Time, page, pageSize int) ([]*domain.Appointment, error) {
	var appointmentModels []models.Appointment

	offset := r.CalculateOffset(page, pageSize)

	err := r.WithContext(ctx).
		Preload("Business").
		Preload("Client").
		Preload("Staff").
		Preload("Service").
		Where("staff_id = ?", staffID).
		Where("start_time >= ? AND start_time <= ?", startDate, endDate).
		Where("deleted_at IS NULL").
		Offset(offset).
		Limit(pageSize).
		Order("start_time ASC").
		Find(&appointmentModels).Error

	if err != nil {
		return nil, fmt.Errorf("failed to list appointments by staff: %w", err)
	}

	return mapAppointmentSliceToDomain(appointmentModels), nil
}

// ListByClient retrieves a paginated list of appointments by client ID within a date range
func (r *AppointmentRepository) ListByClient(ctx context.Context, clientID uuid.UUID, startDate, endDate time.Time, page, pageSize int) ([]*domain.Appointment, error) {
	var appointmentModels []models.Appointment

	offset := r.CalculateOffset(page, pageSize)

	err := r.WithContext(ctx).
		Preload("Business").
		Preload("Client").
		Preload("Staff").
		Preload("Service").
		Where("client_id = ?", clientID).
		Where("start_time >= ? AND start_time <= ?", startDate, endDate).
		Where("deleted_at IS NULL").
		Offset(offset).
		Limit(pageSize).
		Order("start_time ASC").
		Find(&appointmentModels).Error

	if err != nil {
		return nil, fmt.Errorf("failed to list appointments by client: %w", err)
	}

	return mapAppointmentSliceToDomain(appointmentModels), nil
}

// Count counts all appointments
func (r *AppointmentRepository) Count(ctx context.Context) (int64, error) {
	var count int64

	err := r.WithContext(ctx).
		Model(&models.Appointment{}).
		Where("deleted_at IS NULL").
		Count(&count).Error
	if err != nil {
		return 0, fmt.Errorf("failed to count appointments: %w", err)
	}

	return count, nil
}

// CountByBusiness counts appointments by business ID
func (r *AppointmentRepository) CountByBusiness(ctx context.Context, businessID uuid.UUID) (int64, error) {
	var count int64

	err := r.WithContext(ctx).
		Model(&models.Appointment{}).
		Where("business_id = ?", businessID).
		Where("deleted_at IS NULL").
		Count(&count).Error
	if err != nil {
		return 0, fmt.Errorf("failed to count appointments by business: %w", err)
	}

	return count, nil
}

// CountByStaff counts appointments by staff ID
func (r *AppointmentRepository) CountByStaff(ctx context.Context, staffID uuid.UUID) (int64, error) {
	var count int64

	err := r.WithContext(ctx).
		Model(&models.Appointment{}).
		Where("staff_id = ?", staffID).
		Where("deleted_at IS NULL").
		Count(&count).Error
	if err != nil {
		return 0, fmt.Errorf("failed to count appointments by staff: %w", err)
	}

	return count, nil
}

// CountByBusinessAndDateRange counts appointments by business ID within a date range
func (r *AppointmentRepository) CountByBusinessAndDateRange(ctx context.Context, businessID uuid.UUID, startDate, endDate time.Time) (int64, error) {
	var count int64

	err := r.WithContext(ctx).
		Model(&models.Appointment{}).
		Where("business_id = ?", businessID).
		Where("start_time >= ? AND start_time <= ?", startDate, endDate).
		Where("deleted_at IS NULL").
		Count(&count).Error
	if err != nil {
		return 0, fmt.Errorf("failed to count appointments by business and date range: %w", err)
	}

	return count, nil
}

// CountByStaffAndDateRange counts appointments by staff ID within a date range
func (r *AppointmentRepository) CountByStaffAndDateRange(ctx context.Context, staffID uuid.UUID, startDate, endDate time.Time) (int64, error) {
	var count int64

	err := r.WithContext(ctx).
		Model(&models.Appointment{}).
		Where("staff_id = ?", staffID).
		Where("start_time >= ? AND start_time <= ?", startDate, endDate).
		Where("deleted_at IS NULL").
		Count(&count).Error
	if err != nil {
		return 0, fmt.Errorf("failed to count appointments by staff and date range: %w", err)
	}

	return count, nil
}

// Helper functions to map between domain entities and database models

// mapDomainToAppointment converts a domain Appointment entity to a database Appointment model
func mapDomainToAppointment(a *domain.Appointment) *models.Appointment {
	if a == nil {
		return nil
	}

	appointmentModel := &models.Appointment{
		BusinessID: a.BusinessID,
		ClientID:   a.ClientID,
		StaffID:    a.StaffID,
		ServiceID:  a.ServiceID,
		StartTime:  a.StartTime,
		EndTime:    a.EndTime,
		Status:     models.AppointmentStatus(a.Status),
		Notes:      a.Notes,
	}

	// Set the ID if it exists (for updates)
	if a.ID != uuid.Nil {
		appointmentModel.ID = a.ID
	}

	// Handle created by
	if a.CreatedBy != nil {
		appointmentModel.CreatedBy = a.CreatedBy
	}

	// Handle other audit fields
	if !a.CreatedAt.IsZero() {
		appointmentModel.CreatedAt = a.CreatedAt
	}

	if a.UpdatedAt != nil {
		appointmentModel.UpdatedAt = *a.UpdatedAt
	}

	if a.UpdatedBy != nil {
		appointmentModel.UpdatedBy = a.UpdatedBy
	}

	if a.DeletedAt != nil {
		appointmentModel.DeletedAt = gorm.DeletedAt{Time: *a.DeletedAt, Valid: true}
	}

	if a.DeletedBy != nil {
		appointmentModel.DeletedBy = a.DeletedBy
	}

	return appointmentModel
}

// mapAppointmentToDomain converts a database Appointment model to a domain Appointment entity
func mapAppointmentToDomain(a *models.Appointment) *domain.Appointment {
	if a == nil {
		return nil
	}

	appointment := &domain.Appointment{
		ID:         a.ID,
		BusinessID: a.BusinessID,
		ClientID:   a.ClientID,
		StaffID:    a.StaffID,
		ServiceID:  a.ServiceID,
		StartTime:  a.StartTime,
		EndTime:    a.EndTime,
		Status:     string(a.Status),
		Notes:      a.Notes,
		CreatedAt:  a.CreatedAt,
		CreatedBy:  a.CreatedBy,
	}

	// Handle optional fields
	if !a.UpdatedAt.IsZero() {
		appointment.UpdatedAt = &a.UpdatedAt
	}

	if a.UpdatedBy != nil {
		appointment.UpdatedBy = a.UpdatedBy
	}

	if !a.DeletedAt.Time.IsZero() {
		deletedAt := a.DeletedAt.Time
		appointment.DeletedAt = &deletedAt
	}

	if a.DeletedBy != nil {
		appointment.DeletedBy = a.DeletedBy
	}

	// Map related entities if loaded - use functions from individual repositories
	if a.Client != nil && a.Client.ID != uuid.Nil {
		appointment.Client = mapClientModelToDomain(a.Client)
	}

	if a.Business != nil && a.Business.ID != uuid.Nil {
		appointment.Business = mapBusinessModelToDomain(a.Business)
	}

	if a.Staff != nil && a.Staff.ID != uuid.Nil {
		appointment.Staff = mapStaffModelToDomain(a.Staff)
	}

	if a.Service != nil && a.Service.ID != uuid.Nil {
		appointment.Service = mapServiceModelToDomain(a.Service)
	}

	return appointment
}

// mapAppointmentSliceToDomain converts a slice of database Appointment models to a slice of domain Appointment entities
func mapAppointmentSliceToDomain(appointmentModels []models.Appointment) []*domain.Appointment {
	result := make([]*domain.Appointment, len(appointmentModels))
	for i, model := range appointmentModels {
		modelCopy := model // create a copy to avoid pointer issues
		result[i] = mapAppointmentToDomain(&modelCopy)
	}
	return result
}
