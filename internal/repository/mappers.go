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

	// Parse ClerkID as UUID
	clerkUUID := uuid.Nil
	if u.ClerkID != "" {
		if parsed, err := uuid.Parse(u.ClerkID); err == nil {
			clerkUUID = parsed
		}
	}

	user := &domain.User{
		UserID:             u.ID,
		ClerkID:            clerkUUID,
		Email:              u.Email,
		FirstName:          u.FirstName,
		LastName:           u.LastName,
		Phone:              u.Phone,
		IsActive:           u.IsActive,
		Role:               string(u.Role),
		CreatedAt:          u.CreatedAt,
		UpdatedAt:          u.UpdatedAt,
		EmailVerified:      false, // Not in models.User, would need separate table
		LanguagePreference: "en",  // Default value, not in models.User
	}

	// Handle optional/pointer fields
	if u.CreatedBy != nil {
		// Note: domain.User doesn't have CreatedBy as uuid.UUID, not *uuid.UUID
		// We'll need to handle this appropriately in the domain model
	}

	if !u.UpdatedAt.IsZero() {
		// UpdatedAt is already set above
	}

	if u.DeletedAt.Valid {
		// Handle soft delete information if needed in domain
	}

	return user
}

// mapUserDomainToModel converts a domain User to a model User entity
func mapUserDomainToModel(u *domain.User) *models.User {
	if u == nil {
		return nil
	}

	userModel := &models.User{
		BaseModel: models.BaseModel{
			ID:        u.UserID,
			CreatedAt: u.CreatedAt,
			UpdatedAt: u.UpdatedAt,
		},
		Email:     u.Email,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Phone:     u.Phone,
		Role:      models.UserRole(u.Role),
		IsActive:  u.IsActive,
	}

	// Handle ClerkID
	if u.ClerkID != uuid.Nil {
		userModel.ClerkID = u.ClerkID.String()
	} else {
		// Generate a unique ClerkID if not provided
		userModel.ClerkID = "clerk_" + uuid.New().String()
	}

	return userModel
}

// mapBusinessModelToDomain converts a model Business to a domain Business entity
func mapBusinessModelToDomain(b *models.Business) *domain.Business {
	if b == nil || b.ID == uuid.Nil {
		return nil
	}

	business := &domain.Business{
		BusinessID:       b.ID,
		OwnerID:          b.UserID,
		BusinessName:     b.Name,
		BusinessType:     "beauty", // Default type since it's not in models.Business
		Phone:            b.Phone,
		Email:            b.Email,
		AddressLine1:     b.Address,
		City:             b.City,
		Region:           b.State,
		PostalCode:       b.PostalCode,
		Country:          b.Country,
		TimeZone:         b.TimeZone,
		CreatedAt:        b.CreatedAt,
		UpdatedAt:        b.UpdatedAt,
		IsActive:         b.IsActive,
		SubscriptionPlan: string(b.SubscriptionTier),
	}

	// Handle business hours - convert from models.WorkingHours to JSONB
	if b.Settings.WorkingHours != (models.WorkingHours{}) {
		// Convert WorkingHours to JSON bytes for domain
		businessHoursJSON, _ := b.Settings.WorkingHours.Value()
		if jsonBytes, ok := businessHoursJSON.([]byte); ok {
			business.BusinessHours = jsonBytes
		}
	}

	return business
}

// mapBusinessDomainToModel converts a domain Business to a model Business entity
func mapBusinessDomainToModel(b *domain.Business) *models.Business {
	if b == nil {
		return nil
	}

	businessModel := &models.Business{
		BaseModel: models.BaseModel{
			ID:        b.BusinessID,
			CreatedAt: b.CreatedAt,
			UpdatedAt: b.UpdatedAt,
		},
		UserID:           b.OwnerID,
		Name:             b.BusinessName,
		DisplayName:      b.BusinessName, // Use business name as display name
		Address:          b.AddressLine1,
		City:             b.City,
		State:            b.Region,
		PostalCode:       b.PostalCode,
		Country:          b.Country,
		Phone:            b.Phone,
		Email:            b.Email,
		IsActive:         b.IsActive,
		TimeZone:         b.TimeZone,
		SubscriptionTier: models.SubscriptionTier(b.SubscriptionPlan),
	}

	// Set default values for required fields
	if businessModel.Country == "" {
		businessModel.Country = "Portugal"
	}
	if businessModel.TimeZone == "" {
		businessModel.TimeZone = "Europe/Lisbon"
	}
	if businessModel.SubscriptionTier == "" {
		businessModel.SubscriptionTier = models.SubscriptionTierFree
	}

	// Handle business hours - convert from JSONB to models.WorkingHours
	if len(b.BusinessHours) > 0 {
		var workingHours models.WorkingHours
		err := workingHours.Scan(b.BusinessHours)
		if err == nil {
			businessModel.Settings.WorkingHours = workingHours
		}
	}

	return businessModel
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
		BusinessID:  sc.BusinessID,
		Name:        sc.Name,
		Description: sc.Description,
		CreatedAt:   sc.CreatedAt,
	}

	// Handle optional/pointer fields
	if !sc.UpdatedAt.IsZero() {
		category.UpdatedAt = &sc.UpdatedAt
	}

	return category
}
