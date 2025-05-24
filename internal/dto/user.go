package dto

import (
	"github.com/assimoes/beautix/internal/domain"
)

// CreateUserDTO represents the data for creating a user
type CreateUserDTO struct {
	Email     string  `json:"email" validate:"required,email"`
	ClerkID   *string `json:"clerk_id,omitempty"`
	FirstName string  `json:"first_name" validate:"required,min=2,max=100"`
	LastName  string  `json:"last_name" validate:"required,min=2,max=100"`
	Phone     *string `json:"phone,omitempty" validate:"omitempty,min=10,max=20"`
}

// UpdateUserDTO represents the data for updating a user
type UpdateUserDTO struct {
	FirstName *string `json:"first_name,omitempty" validate:"omitempty,min=2,max=100"`
	LastName  *string `json:"last_name,omitempty" validate:"omitempty,min=2,max=100"`
	Phone     *string `json:"phone,omitempty" validate:"omitempty,min=10,max=20"`
	IsActive  *bool   `json:"is_active,omitempty"`
}

// UserResponseDTO represents the response data for a user
type UserResponseDTO struct {
	BaseResponse
	Email     string  `json:"email"`
	ClerkID   *string `json:"clerk_id,omitempty"`
	FirstName string  `json:"first_name"`
	LastName  string  `json:"last_name"`
	Phone     *string `json:"phone,omitempty"`
	IsActive  bool    `json:"is_active"`
	FullName  string  `json:"full_name"`
}

// UserWithBusinessesDTO represents a user with their businesses
type UserWithBusinessesDTO struct {
	UserResponseDTO
	Businesses []*BusinessResponseDTO `json:"businesses,omitempty"`
}

// ClerkUserDTO represents user data from Clerk webhook
type ClerkUserDTO struct {
	ClerkID          string `json:"clerk_id"`
	Email            string `json:"email"`
	FirstName        string `json:"first_name"`
	LastName         string `json:"last_name"`
	Phone            string `json:"phone,omitempty"`
	EmailVerified    bool   `json:"email_verified"`
	ProfileImageURL  string `json:"profile_image_url,omitempty"`
}

// ToCreateUserDTO converts ClerkUserDTO to CreateUserDTO
func (c ClerkUserDTO) ToCreateUserDTO() CreateUserDTO {
	var phone *string
	if c.Phone != "" {
		phone = &c.Phone
	}

	return CreateUserDTO{
		Email:     c.Email,
		ClerkID:   &c.ClerkID,
		FirstName: c.FirstName,
		LastName:  c.LastName,
		Phone:     phone,
	}
}

// ToUserResponseDTO converts a User domain model to UserResponseDTO
func ToUserResponseDTO(user *domain.User) *UserResponseDTO {
	if user == nil {
		return nil
	}

	return &UserResponseDTO{
		BaseResponse: BaseResponse{
			ID:        user.ID,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		},
		Email:     user.Email,
		ClerkID:   user.ClerkID,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Phone:     user.Phone,
		IsActive:  user.IsActive,
		FullName:  user.GetFullName(),
	}
}

// ToUserResponseDTOs converts a slice of User domain models to UserResponseDTOs
func ToUserResponseDTOs(users []*domain.User) []*UserResponseDTO {
	result := make([]*UserResponseDTO, len(users))
	for i, user := range users {
		result[i] = ToUserResponseDTO(user)
	}
	return result
}

// ToUserWithBusinessesDTO converts a User with businesses to UserWithBusinessesDTO
func ToUserWithBusinessesDTO(user *domain.User, businesses []*BusinessResponseDTO) *UserWithBusinessesDTO {
	return &UserWithBusinessesDTO{
		UserResponseDTO: *ToUserResponseDTO(user),
		Businesses:      businesses,
	}
}

// UserSearchResultDTO represents search results for users
type UserSearchResultDTO struct {
	Users []*UserResponseDTO `json:"users"`
	Total int               `json:"total"`
}

