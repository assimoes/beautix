package repository

import (
	"github.com/assimoes/beautix/internal/domain"
	"github.com/assimoes/beautix/internal/models"
	"github.com/google/uuid"
)

// Common mapping functions to avoid duplication across repositories

// mapUserModelToDomain converts a model User to a domain User entity
func mapUserModelToDomain(u *models.User) *domain.User {
	if u == nil || u.ID == uuid.Nil {
		return nil
	}
	
	return &domain.User{
		UserID:    u.ID,
		Email:     u.Email,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Phone:     u.Phone,
		IsActive:  u.IsActive,
		Role:      string(u.Role),
	}
}

// mapBusinessModelToDomain converts a model Business to a domain Business entity
func mapBusinessModelToDomain(b *models.Business) *domain.Business {
	if b == nil || b.ID == uuid.Nil {
		return nil
	}
	
	return &domain.Business{
		BusinessID:   b.ID,
		OwnerID:      b.UserID,
		BusinessName: b.Name,
		City:         b.City,
		Country:      b.Country,
		Phone:        b.Phone,
		Email:        b.Email,
		IsActive:     b.IsActive,
	}
}

// mapProviderModelToDomain converts a model Business to a domain Provider entity
// This is an alias for business since a provider is essentially a business in our domain
func mapProviderModelToDomain(b *models.Business) *domain.Provider {
	if b == nil || b.ID == uuid.Nil {
		return nil
	}
	
	return &domain.Provider{
		ID:               b.ID,
		UserID:           b.UserID,
		BusinessName:     b.Name,
		Description:      b.Description,
		Address:          b.Address,
		City:             b.City,
		PostalCode:       b.PostalCode,
		Country:          b.Country,
		Website:          b.Website,
		LogoURL:          b.LogoURL,
		SubscriptionTier: string(b.SubscriptionTier),
	}
}

// mapServiceCategoryModelToDomain converts a model ServiceCategory to a domain ServiceCategory entity
func mapServiceCategoryModelToDomain(sc *models.ServiceCategory) *domain.ServiceCategory {
	if sc == nil {
		return nil
	}
	
	category := &domain.ServiceCategory{
		ID:          sc.ID,
		Name:        sc.Name,
		Description: sc.Description,
		CreatedAt:   sc.CreatedAt,
	}
	
	// Handle optional/pointer fields
	if sc.CreatedBy != nil {
		category.CreatedBy = sc.CreatedBy
	}
	
	if !sc.UpdatedAt.IsZero() {
		category.UpdatedAt = &sc.UpdatedAt
	}
	
	if sc.UpdatedBy != nil {
		category.UpdatedBy = sc.UpdatedBy
	}
	
	if sc.DeletedAt.Valid {
		deletedAt := sc.DeletedAt.Time
		category.DeletedAt = &deletedAt
	}
	
	if sc.DeletedBy != nil {
		category.DeletedBy = sc.DeletedBy
	}
	
	return category
}