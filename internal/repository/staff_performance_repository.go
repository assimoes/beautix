package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/assimoes/beautix/internal/domain"
	"github.com/assimoes/beautix/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// StaffPerformanceRepository implements the domain.StaffPerformanceRepository interface using GORM
type StaffPerformanceRepository struct {
	*BaseRepository
}

// NewStaffPerformanceRepository creates a new instance of StaffPerformanceRepository
func NewStaffPerformanceRepository(db DBAdapter) domain.StaffPerformanceRepository {
	return &StaffPerformanceRepository{
		BaseRepository: NewBaseRepository(db),
	}
}

// Create creates a new staff performance record
func (r *StaffPerformanceRepository) Create(ctx context.Context, performance *domain.StaffPerformance) error {
	performanceModel := mapPerformanceDomainToModel(performance)

	if err := r.WithContext(ctx).Create(&performanceModel).Error; err != nil {
		return fmt.Errorf("failed to create staff performance: %w", err)
	}

	// Update the domain entity with any generated fields
	performance.PerformanceID = performanceModel.ID
	performance.CreatedAt = performanceModel.CreatedAt
	performance.UpdatedAt = performanceModel.UpdatedAt

	return nil
}

// GetByID retrieves a staff performance record by ID
func (r *StaffPerformanceRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.StaffPerformance, error) {
	var performanceModel models.StaffPerformance

	err := r.WithContext(ctx).
		Preload("Staff").
		Preload("Staff.User").
		First(&performanceModel, "id = ?", id).Error

	if err != nil {
		return nil, r.HandleNotFound(err, "staff performance", id)
	}

	return mapPerformanceModelToDomain(&performanceModel), nil
}

// GetByStaffAndPeriod retrieves a staff performance record by staff ID, period type, and start date
func (r *StaffPerformanceRepository) GetByStaffAndPeriod(ctx context.Context, staffID uuid.UUID, period string, startDate time.Time) (*domain.StaffPerformance, error) {
	var performanceModel models.StaffPerformance

	// Format the date to truncate time component for more reliable matching
	formattedDate := startDate.Format("2006-01-02")

	err := r.WithContext(ctx).
		Preload("Staff").
		Preload("Staff.User").
		Where("staff_id = ? AND period = ? AND DATE(start_date) = ?",
			staffID, models.PerformancePeriod(period), formattedDate).
		First(&performanceModel).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("staff performance not found for staff %s and period %s starting on %s",
				staffID, period, formattedDate)
		}
		return nil, fmt.Errorf("failed to get staff performance by period: %w", err)
	}

	return mapPerformanceModelToDomain(&performanceModel), nil
}

// GetByStaffAndDateRange retrieves staff performance records by staff ID and date range
func (r *StaffPerformanceRepository) GetByStaffAndDateRange(ctx context.Context, staffID uuid.UUID, startDate, endDate time.Time) ([]*domain.StaffPerformance, error) {
	var performanceModels []models.StaffPerformance

	err := r.WithContext(ctx).
		Preload("Staff").
		Preload("Staff.User").
		Where("staff_id = ?", staffID).
		Where("start_date BETWEEN ? AND ? OR end_date BETWEEN ? AND ? OR (start_date <= ? AND end_date >= ?)",
			startDate, endDate, startDate, endDate, startDate, endDate).
		Order("start_date ASC").
		Find(&performanceModels).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get staff performance by date range: %w", err)
	}

	return mapPerformanceModelsToDomainSlice(performanceModels), nil
}

// Update updates a staff performance record
func (r *StaffPerformanceRepository) Update(ctx context.Context, id uuid.UUID, performance *domain.StaffPerformance) error {
	// First find the performance record to ensure it exists
	var performanceModel models.StaffPerformance
	err := r.WithContext(ctx).First(&performanceModel, "id = ?", id).Error
	if err != nil {
		return r.HandleNotFound(err, "staff performance", id)
	}

	// Map domain entity to model
	updatedModel := mapPerformanceDomainToModel(performance)
	updatedModel.ID = id
	updatedModel.UpdatedAt = time.Now()

	// Perform the update
	err = r.WithContext(ctx).Model(&performanceModel).Updates(map[string]interface{}{
		"period":                 updatedModel.Period,
		"start_date":             updatedModel.StartDate,
		"end_date":               updatedModel.EndDate,
		"total_appointments":     updatedModel.TotalAppointments,
		"completed_appointments": updatedModel.CompletedAppointments,
		"canceled_appointments":  updatedModel.CanceledAppointments,
		"no_show_appointments":   updatedModel.NoShowAppointments,
		"total_revenue":          updatedModel.TotalRevenue,
		"average_rating":         updatedModel.AverageRating,
		"client_retention_rate":  updatedModel.ClientRetentionRate,
		"new_clients":            updatedModel.NewClients,
		"return_clients":         updatedModel.ReturnClients,
		"updated_at":             updatedModel.UpdatedAt,
	}).Error

	if err != nil {
		return fmt.Errorf("failed to update staff performance: %w", err)
	}

	return nil
}

// Delete deletes a staff performance record
func (r *StaffPerformanceRepository) Delete(ctx context.Context, id uuid.UUID) error {
	// First find the performance record to ensure it exists
	var performanceModel models.StaffPerformance
	err := r.WithContext(ctx).First(&performanceModel, "id = ?", id).Error
	if err != nil {
		return r.HandleNotFound(err, "staff performance", id)
	}

	// Perform the delete (hard delete since these are metrics that can be recalculated)
	err = r.WithContext(ctx).Unscoped().Delete(&performanceModel).Error
	if err != nil {
		return fmt.Errorf("failed to delete staff performance: %w", err)
	}

	return nil
}

// ListByBusiness retrieves a paginated list of staff performance records by business ID and period
func (r *StaffPerformanceRepository) ListByBusiness(ctx context.Context, businessID uuid.UUID, period string, page, pageSize int) ([]*domain.StaffPerformance, error) {
	var performanceModels []models.StaffPerformance

	// Apply pagination
	offset := r.CalculateOffset(page, pageSize)

	query := r.WithContext(ctx).
		Preload("Staff").
		Preload("Staff.User").
		Where("business_id = ?", businessID)

	// Filter by period if provided
	if period != "" {
		query = query.Where("period = ?", models.PerformancePeriod(period))
	}

	err := query.
		Offset(offset).
		Limit(pageSize).
		Order("start_date DESC").
		Find(&performanceModels).Error

	if err != nil {
		return nil, fmt.Errorf("failed to list staff performance by business: %w", err)
	}

	return mapPerformanceModelsToDomainSlice(performanceModels), nil
}

// Helper functions to map between domain entities and models

// mapPerformanceDomainToModel converts a domain StaffPerformance entity to a model StaffPerformance
func mapPerformanceDomainToModel(p *domain.StaffPerformance) *models.StaffPerformance {
	if p == nil {
		return nil
	}

	performanceModel := &models.StaffPerformance{
		ID:                    p.PerformanceID,
		CreatedAt:             p.CreatedAt,
		UpdatedAt:             p.UpdatedAt,
		BusinessID:            p.BusinessID,
		StaffID:               p.StaffID,
		Period:                models.PerformancePeriod(p.Period),
		StartDate:             p.StartDate,
		EndDate:               p.EndDate,
		TotalAppointments:     p.TotalAppointments,
		CompletedAppointments: p.CompletedAppointments,
		CanceledAppointments:  p.CanceledAppointments,
		NoShowAppointments:    p.NoShowAppointments,
		TotalRevenue:          p.TotalRevenue,
		AverageRating:         p.AverageRating,
		ClientRetentionRate:   p.ClientRetentionRate,
		NewClients:            p.NewClients,
		ReturnClients:         p.ReturnClients,
	}

	return performanceModel
}

// mapPerformanceModelToDomain converts a model StaffPerformance to a domain StaffPerformance entity
func mapPerformanceModelToDomain(p *models.StaffPerformance) *domain.StaffPerformance {
	if p == nil {
		return nil
	}

	performance := &domain.StaffPerformance{
		PerformanceID:         p.ID,
		BusinessID:            p.BusinessID,
		StaffID:               p.StaffID,
		Period:                string(p.Period),
		StartDate:             p.StartDate,
		EndDate:               p.EndDate,
		TotalAppointments:     p.TotalAppointments,
		CompletedAppointments: p.CompletedAppointments,
		CanceledAppointments:  p.CanceledAppointments,
		NoShowAppointments:    p.NoShowAppointments,
		TotalRevenue:          p.TotalRevenue,
		AverageRating:         p.AverageRating,
		ClientRetentionRate:   p.ClientRetentionRate,
		NewClients:            p.NewClients,
		ReturnClients:         p.ReturnClients,
		CreatedAt:             p.CreatedAt,
		UpdatedAt:             p.UpdatedAt,
	}

	// Map related entities if loaded
	if p.Staff.ID != uuid.Nil {
		performance.Staff = mapStaffModelToDomain(&p.Staff)
	}

	return performance
}

// mapPerformanceModelsToDomainSlice converts a slice of model StaffPerformance to a slice of domain StaffPerformance entities
func mapPerformanceModelsToDomainSlice(performanceModels []models.StaffPerformance) []*domain.StaffPerformance {
	result := make([]*domain.StaffPerformance, len(performanceModels))
	for i, model := range performanceModels {
		modelCopy := model // create a copy to avoid pointer issues
		result[i] = mapPerformanceModelToDomain(&modelCopy)
	}
	return result
}
