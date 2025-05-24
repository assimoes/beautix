package dto

import (
	"time"

	"github.com/assimoes/beautix/internal/domain"
)

// CreateBusinessDTO represents the data for creating a business
type CreateBusinessDTO struct {
	Name          string  `json:"name" validate:"required,min=2,max=100"`
	DisplayName   *string `json:"display_name,omitempty" validate:"omitempty,min=2,max=100"`
	BusinessType  *string `json:"business_type,omitempty" validate:"omitempty,max=50"`
	TaxID         *string `json:"tax_id,omitempty" validate:"omitempty,max=50"`
	Email         string  `json:"email" validate:"required,email"`
	Website       *string `json:"website,omitempty" validate:"omitempty,url"`
	Currency      string  `json:"currency" validate:"required,currency"`
	TimeZone      string  `json:"time_zone" validate:"required"`
}

// UpdateBusinessDTO represents the data for updating a business
type UpdateBusinessDTO struct {
	Name          *string `json:"name,omitempty" validate:"omitempty,min=2,max=100"`
	DisplayName   *string `json:"display_name,omitempty" validate:"omitempty,min=2,max=100"`
	BusinessType  *string `json:"business_type,omitempty" validate:"omitempty,max=50"`
	TaxID         *string `json:"tax_id,omitempty" validate:"omitempty,max=50"`
	Email         *string `json:"email,omitempty" validate:"omitempty,email"`
	Website       *string `json:"website,omitempty" validate:"omitempty,url"`
	LogoURL       *string `json:"logo_url,omitempty" validate:"omitempty,url"`
	CoverPhotoURL *string `json:"cover_photo_url,omitempty" validate:"omitempty,url"`
	Currency      *string `json:"currency,omitempty" validate:"omitempty,currency"`
	TimeZone      *string `json:"time_zone,omitempty"`
	IsActive      *bool   `json:"is_active,omitempty"`
}

// BusinessResponseDTO represents the response data for a business
type BusinessResponseDTO struct {
	BaseResponse
	UserID           string    `json:"user_id"`
	Name             string    `json:"name"`
	DisplayName      *string   `json:"display_name,omitempty"`
	BusinessType     *string   `json:"business_type,omitempty"`
	TaxID            *string   `json:"tax_id,omitempty"`
	Email            string    `json:"email"`
	Website          *string   `json:"website,omitempty"`
	LogoURL          *string   `json:"logo_url,omitempty"`
	CoverPhotoURL    *string   `json:"cover_photo_url,omitempty"`
	IsVerified       bool      `json:"is_verified"`
	Currency         string    `json:"currency"`
	TimeZone         string    `json:"time_zone"`
	IsActive         bool      `json:"is_active"`
	SubscriptionTier string    `json:"subscription_tier"`
	TrialEndsAt      *time.Time `json:"trial_ends_at,omitempty"`
	DisplayNameValue string    `json:"display_name_value"`
}

// BusinessWithLocationsDTO represents a business with its locations
type BusinessWithLocationsDTO struct {
	BusinessResponseDTO
	Locations []*LocationResponseDTO `json:"locations,omitempty"`
}

// CreateLocationDTO represents the data for creating a business location
type CreateLocationDTO struct {
	Name       string  `json:"name" validate:"required,min=2,max=100"`
	Address    *string `json:"address,omitempty" validate:"omitempty,max=255"`
	City       *string `json:"city,omitempty" validate:"omitempty,max=100"`
	State      *string `json:"state,omitempty" validate:"omitempty,max=100"`
	PostalCode *string `json:"postal_code,omitempty" validate:"omitempty,max=20"`
	Country    string  `json:"country" validate:"required,max=50"`
	Phone      *string `json:"phone,omitempty" validate:"omitempty,min=10,max=20"`
	Email      *string `json:"email,omitempty" validate:"omitempty,email"`
	Timezone   string  `json:"timezone" validate:"required"`
	IsMain     bool    `json:"is_main"`
}

// UpdateLocationDTO represents the data for updating a business location
type UpdateLocationDTO struct {
	Name       *string `json:"name,omitempty" validate:"omitempty,min=2,max=100"`
	Address    *string `json:"address,omitempty" validate:"omitempty,max=255"`
	City       *string `json:"city,omitempty" validate:"omitempty,max=100"`
	State      *string `json:"state,omitempty" validate:"omitempty,max=100"`
	PostalCode *string `json:"postal_code,omitempty" validate:"omitempty,max=20"`
	Country    *string `json:"country,omitempty" validate:"omitempty,max=50"`
	Phone      *string `json:"phone,omitempty" validate:"omitempty,min=10,max=20"`
	Email      *string `json:"email,omitempty" validate:"omitempty,email"`
	Timezone   *string `json:"timezone,omitempty"`
	IsActive   *bool   `json:"is_active,omitempty"`
	IsMain     *bool   `json:"is_main,omitempty"`
}

// LocationResponseDTO represents the response data for a business location
type LocationResponseDTO struct {
	BaseResponse
	BusinessID string  `json:"business_id"`
	Name       string  `json:"name"`
	Address    *string `json:"address,omitempty"`
	City       *string `json:"city,omitempty"`
	State      *string `json:"state,omitempty"`
	PostalCode *string `json:"postal_code,omitempty"`
	Country    string  `json:"country"`
	Phone      *string `json:"phone,omitempty"`
	Email      *string `json:"email,omitempty"`
	IsActive   bool    `json:"is_active"`
	IsMain     bool    `json:"is_main"`
	Timezone   string  `json:"timezone"`
}

// CreateSettingsDTO represents the data for creating business settings
type CreateSettingsDTO struct {
	CalendarStartHour            int    `json:"calendar_start_hour" validate:"min=0,max=23"`
	CalendarEndHour              int    `json:"calendar_end_hour" validate:"min=1,max=24"`
	AppointmentBufferMinutes     int    `json:"appointment_buffer_minutes" validate:"min=0"`
	AllowOnlineBooking           bool   `json:"allow_online_booking"`
	DefaultAppointmentDuration   int    `json:"default_appointment_duration" validate:"min=15"`
	Currency                     string `json:"currency" validate:"required,currency"`
	DateFormat                   string `json:"date_format" validate:"required"`
	TimeFormat                   string `json:"time_format" validate:"required,oneof='12h' '24h'"`
}

// UpdateSettingsDTO represents the data for updating business settings
type UpdateSettingsDTO struct {
	CalendarStartHour            *int    `json:"calendar_start_hour,omitempty" validate:"omitempty,min=0,max=23"`
	CalendarEndHour              *int    `json:"calendar_end_hour,omitempty" validate:"omitempty,min=1,max=24"`
	AppointmentBufferMinutes     *int    `json:"appointment_buffer_minutes,omitempty" validate:"omitempty,min=0"`
	AllowOnlineBooking           *bool   `json:"allow_online_booking,omitempty"`
	DefaultAppointmentDuration   *int    `json:"default_appointment_duration,omitempty" validate:"omitempty,min=15"`
	Currency                     *string `json:"currency,omitempty" validate:"omitempty,currency"`
	DateFormat                   *string `json:"date_format,omitempty"`
	TimeFormat                   *string `json:"time_format,omitempty" validate:"omitempty,oneof='12h' '24h'"`
}

// SettingsResponseDTO represents the response data for business settings
type SettingsResponseDTO struct {
	BaseResponse
	BusinessID                   string `json:"business_id"`
	CalendarStartHour            int    `json:"calendar_start_hour"`
	CalendarEndHour              int    `json:"calendar_end_hour"`
	AppointmentBufferMinutes     int    `json:"appointment_buffer_minutes"`
	AllowOnlineBooking           bool   `json:"allow_online_booking"`
	DefaultAppointmentDuration   int    `json:"default_appointment_duration"`
	Currency                     string `json:"currency"`
	DateFormat                   string `json:"date_format"`
	TimeFormat                   string `json:"time_format"`
}

// BusinessHoursDTO represents business hours configuration
type BusinessHoursDTO struct {
	Monday    *DayHoursDTO `json:"monday,omitempty"`
	Tuesday   *DayHoursDTO `json:"tuesday,omitempty"`
	Wednesday *DayHoursDTO `json:"wednesday,omitempty"`
	Thursday  *DayHoursDTO `json:"thursday,omitempty"`
	Friday    *DayHoursDTO `json:"friday,omitempty"`
	Saturday  *DayHoursDTO `json:"saturday,omitempty"`
	Sunday    *DayHoursDTO `json:"sunday,omitempty"`
}

// DayHoursDTO represents hours for a specific day
type DayHoursDTO struct {
	IsOpen    bool   `json:"is_open"`
	OpenTime  string `json:"open_time,omitempty" validate:"omitempty,datetime=15:04"`
	CloseTime string `json:"close_time,omitempty" validate:"omitempty,datetime=15:04"`
}

// SearchBusinessCriteriaDTO represents search criteria for businesses
type SearchBusinessCriteriaDTO struct {
	Query       string  `json:"query,omitempty"`
	City        *string `json:"city,omitempty"`
	Country     *string `json:"country,omitempty"`
	ServiceName *string `json:"service_name,omitempty"`
	IsVerified  *bool   `json:"is_verified,omitempty"`
	Page        int     `json:"page" validate:"min=1"`
	PageSize    int     `json:"page_size" validate:"min=1,max=100"`
}

// PublicBusinessDTO represents public business information for discovery
type PublicBusinessDTO struct {
	ID               string    `json:"id"`
	Name             string    `json:"name"`
	DisplayName      *string   `json:"display_name,omitempty"`
	BusinessType     *string   `json:"business_type,omitempty"`
	Email            string    `json:"email"`
	Website          *string   `json:"website,omitempty"`
	LogoURL          *string   `json:"logo_url,omitempty"`
	CoverPhotoURL    *string   `json:"cover_photo_url,omitempty"`
	IsVerified       bool      `json:"is_verified"`
	SubscriptionTier string    `json:"subscription_tier"`
	CreatedAt        time.Time `json:"created_at"`
}

// ToBusinessResponseDTO converts a Business domain model to BusinessResponseDTO
func ToBusinessResponseDTO(business *domain.Business) *BusinessResponseDTO {
	if business == nil {
		return nil
	}

	return &BusinessResponseDTO{
		BaseResponse: BaseResponse{
			ID:        business.ID,
			CreatedAt: business.CreatedAt,
			UpdatedAt: business.UpdatedAt,
		},
		UserID:           business.UserID,
		Name:             business.Name,
		DisplayName:      business.DisplayName,
		BusinessType:     business.BusinessType,
		TaxID:            business.TaxID,
		Email:            business.Email,
		Website:          business.Website,
		LogoURL:          business.LogoURL,
		CoverPhotoURL:    business.CoverPhotoURL,
		IsVerified:       business.IsVerified,
		Currency:         business.Currency,
		TimeZone:         business.TimeZone,
		IsActive:         business.IsActive,
		SubscriptionTier: business.SubscriptionTier,
		TrialEndsAt:      business.TrialEndsAt,
		DisplayNameValue: business.GetDisplayName(),
	}
}

// ToBusinessResponseDTOs converts a slice of Business domain models to BusinessResponseDTOs
func ToBusinessResponseDTOs(businesses []*domain.Business) []*BusinessResponseDTO {
	result := make([]*BusinessResponseDTO, len(businesses))
	for i, business := range businesses {
		result[i] = ToBusinessResponseDTO(business)
	}
	return result
}

// ToLocationResponseDTO converts a BusinessLocation domain model to LocationResponseDTO
func ToLocationResponseDTO(location *domain.BusinessLocation) *LocationResponseDTO {
	if location == nil {
		return nil
	}

	return &LocationResponseDTO{
		BaseResponse: BaseResponse{
			ID:        location.ID,
			CreatedAt: location.CreatedAt,
			UpdatedAt: location.UpdatedAt,
		},
		BusinessID: location.BusinessID,
		Name:       location.Name,
		Address:    location.Address,
		City:       location.City,
		State:      location.State,
		PostalCode: location.PostalCode,
		Country:    location.Country,
		Phone:      location.Phone,
		Email:      location.Email,
		IsActive:   location.IsActive,
		IsMain:     location.IsMain,
		Timezone:   location.Timezone,
	}
}

// ToLocationResponseDTOs converts a slice of BusinessLocation domain models to LocationResponseDTOs
func ToLocationResponseDTOs(locations []*domain.BusinessLocation) []*LocationResponseDTO {
	result := make([]*LocationResponseDTO, len(locations))
	for i, location := range locations {
		result[i] = ToLocationResponseDTO(location)
	}
	return result
}

// ToSettingsResponseDTO converts a BusinessSettings domain model to SettingsResponseDTO
func ToSettingsResponseDTO(settings *domain.BusinessSettings) *SettingsResponseDTO {
	if settings == nil {
		return nil
	}

	return &SettingsResponseDTO{
		BaseResponse: BaseResponse{
			ID:        settings.ID,
			CreatedAt: settings.CreatedAt,
			UpdatedAt: settings.UpdatedAt,
		},
		BusinessID:                   settings.BusinessID,
		CalendarStartHour:            settings.CalendarStartHour,
		CalendarEndHour:              settings.CalendarEndHour,
		AppointmentBufferMinutes:     settings.AppointmentBufferMinutes,
		AllowOnlineBooking:           settings.AllowOnlineBooking,
		DefaultAppointmentDuration:   settings.DefaultAppointmentDuration,
		Currency:                     settings.Currency,
		DateFormat:                   settings.DateFormat,
		TimeFormat:                   settings.TimeFormat,
	}
}

// ToPublicBusinessDTO converts a Business domain model to PublicBusinessDTO
func ToPublicBusinessDTO(business *domain.Business) *PublicBusinessDTO {
	if business == nil {
		return nil
	}

	return &PublicBusinessDTO{
		ID:               business.ID,
		Name:             business.Name,
		DisplayName:      business.DisplayName,
		BusinessType:     business.BusinessType,
		Email:            business.Email,
		Website:          business.Website,
		LogoURL:          business.LogoURL,
		CoverPhotoURL:    business.CoverPhotoURL,
		IsVerified:       business.IsVerified,
		SubscriptionTier: business.SubscriptionTier,
		CreatedAt:        business.CreatedAt,
	}
}

// Staff DTOs

// CreateStaffDTO represents the data for creating a staff member
type CreateStaffDTO struct {
	UserID      string              `json:"user_id" validate:"required,uuid"`
	Role        domain.BusinessRole `json:"role" validate:"required,business_role"`
	Permissions *string             `json:"permissions,omitempty"`
	StartDate   *time.Time          `json:"start_date,omitempty"`
}

// UpdateStaffDTO represents the data for updating a staff member
type UpdateStaffDTO struct {
	Role        *domain.BusinessRole `json:"role,omitempty" validate:"omitempty,business_role"`
	IsActive    *bool                `json:"is_active,omitempty"`
	Permissions *string              `json:"permissions,omitempty"`
	EndDate     *time.Time           `json:"end_date,omitempty"`
}

// StaffResponseDTO represents the response data for a staff member
type StaffResponseDTO struct {
	BaseResponse
	BusinessID  string              `json:"business_id"`
	UserID      string              `json:"user_id"`
	Role        domain.BusinessRole `json:"role"`
	IsActive    bool                `json:"is_active"`
	Permissions *string             `json:"permissions,omitempty"`
	StartDate   *time.Time          `json:"start_date,omitempty"`
	EndDate     *time.Time          `json:"end_date,omitempty"`
}

// StaffWithUserDTO represents a staff member with user details
type StaffWithUserDTO struct {
	StaffResponseDTO
	User *UserResponseDTO `json:"user,omitempty"`
}

// ToStaffResponseDTO converts a Staff domain model to StaffResponseDTO
func ToStaffResponseDTO(staff *domain.Staff) *StaffResponseDTO {
	if staff == nil {
		return nil
	}

	return &StaffResponseDTO{
		BaseResponse: BaseResponse{
			ID:        staff.ID,
			CreatedAt: staff.CreatedAt,
			UpdatedAt: staff.UpdatedAt,
		},
		BusinessID:  staff.BusinessID,
		UserID:      staff.UserID,
		Role:        staff.Role,
		IsActive:    staff.IsActive,
		Permissions: staff.Permissions,
		StartDate:   staff.StartDate,
		EndDate:     staff.EndDate,
	}
}

// ToStaffResponseDTOs converts a slice of Staff domain models to StaffResponseDTOs
func ToStaffResponseDTOs(staff []*domain.Staff) []*StaffResponseDTO {
	result := make([]*StaffResponseDTO, len(staff))
	for i, s := range staff {
		result[i] = ToStaffResponseDTO(s)
	}
	return result
}

// ToStaffWithUserDTO converts a Staff domain model with User to StaffWithUserDTO
func ToStaffWithUserDTO(staff *domain.Staff, user *UserResponseDTO) *StaffWithUserDTO {
	return &StaffWithUserDTO{
		StaffResponseDTO: *ToStaffResponseDTO(staff),
		User:             user,
	}
}

// ParseBusinessRole converts a string to BusinessRole
func ParseBusinessRole(role string) domain.BusinessRole {
	switch role {
	case "owner":
		return domain.BusinessRoleOwner
	case "manager":
		return domain.BusinessRoleManager
	case "employee":
		return domain.BusinessRoleEmployee
	case "assistant":
		return domain.BusinessRoleAssistant
	default:
		return domain.BusinessRoleEmployee
	}
}