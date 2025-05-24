package testdb

import (
	"fmt"

	"github.com/assimoes/beautix/internal/domain"
	"github.com/assimoes/beautix/internal/infrastructure/database"
	"github.com/google/uuid"
)

// FixtureBuilder helps create test data
type FixtureBuilder struct {
	db *database.DB
}

// NewFixtureBuilder creates a new fixture builder
func NewFixtureBuilder(db *database.DB) *FixtureBuilder {
	return &FixtureBuilder{db: db}
}

// CreateUser creates a test user with default values
func (fb *FixtureBuilder) CreateUser(overrides ...func(*domain.User)) (*domain.User, error) {
	phone := "+1234567890"
	user := &domain.User{
		Email:     fmt.Sprintf("test-%s@example.com", uuid.New().String()[:8]),
		FirstName: "Test",
		LastName:  "User",
		Phone:     &phone,
		IsActive:  true,
	}

	// Apply overrides
	for _, override := range overrides {
		override(user)
	}

	if err := fb.db.Create(user).Error; err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return user, nil
}

// CreateBusiness creates a test business with default values
func (fb *FixtureBuilder) CreateBusiness(ownerID string, overrides ...func(*domain.Business)) (*domain.Business, error) {
	settings := "{}"
	business := &domain.Business{
		UserID:   ownerID,
		Name:     fmt.Sprintf("Test Business %s", uuid.New().String()[:8]),
		Email:    fmt.Sprintf("business-%s@example.com", uuid.New().String()[:8]),
		IsActive: true,
		Settings: &settings,
		Currency: "EUR",
		TimeZone: "Europe/Lisbon",
	}

	// Apply overrides
	for _, override := range overrides {
		override(business)
	}

	if err := fb.db.Create(business).Error; err != nil {
		return nil, fmt.Errorf("failed to create business: %w", err)
	}

	return business, nil
}

// CreateStaff creates a test staff member
func (fb *FixtureBuilder) CreateStaff(businessID, userID string, overrides ...func(*domain.Staff)) (*domain.Staff, error) {
	permissions := `{"canManageAppointments": true, "canViewClients": true}`
	staff := &domain.Staff{
		BusinessID:  businessID,
		UserID:      userID,
		Role:        domain.BusinessRoleEmployee,
		IsActive:    true,
		Permissions: &permissions,
	}

	// Apply overrides
	for _, override := range overrides {
		override(staff)
	}

	if err := fb.db.Create(staff).Error; err != nil {
		return nil, fmt.Errorf("failed to create staff: %w", err)
	}

	return staff, nil
}

// UserBuilder provides a fluent interface for creating users
type UserBuilder struct {
	user *domain.User
	fb   *FixtureBuilder
}

// NewUserBuilder creates a new user builder
func (fb *FixtureBuilder) NewUser() *UserBuilder {
	return &UserBuilder{
		user: &domain.User{
			Email:     fmt.Sprintf("test-%s@example.com", uuid.New().String()[:8]),
			FirstName: "Test",
			LastName:  "User",
			IsActive:  true,
		},
		fb: fb,
	}
}

// WithEmail sets the email
func (ub *UserBuilder) WithEmail(email string) *UserBuilder {
	ub.user.Email = email
	return ub
}

// WithName sets the first and last name
func (ub *UserBuilder) WithName(firstName, lastName string) *UserBuilder {
	ub.user.FirstName = firstName
	ub.user.LastName = lastName
	return ub
}

// WithPhone sets the phone number
func (ub *UserBuilder) WithPhone(phone string) *UserBuilder {
	ub.user.Phone = &phone
	return ub
}

// WithClerkID sets the Clerk ID
func (ub *UserBuilder) WithClerkID(clerkID string) *UserBuilder {
	ub.user.ClerkID = &clerkID
	return ub
}

// WithActive sets the active status
func (ub *UserBuilder) WithActive(active bool) *UserBuilder {
	ub.user.IsActive = active
	return ub
}

// Build creates the user in the database
func (ub *UserBuilder) Build() (*domain.User, error) {
	if err := ub.fb.db.Create(ub.user).Error; err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}
	return ub.user, nil
}

// BusinessBuilder provides a fluent interface for creating businesses
type BusinessBuilder struct {
	business *domain.Business
	fb       *FixtureBuilder
}

// NewBusinessBuilder creates a new business builder
func (fb *FixtureBuilder) NewBusiness(ownerID string) *BusinessBuilder {
	settings := "{}"
	return &BusinessBuilder{
		business: &domain.Business{
			UserID:   ownerID,
			Name:     fmt.Sprintf("Test Business %s", uuid.New().String()[:8]),
			Email:    fmt.Sprintf("business-%s@example.com", uuid.New().String()[:8]),
			IsActive: true,
			Settings: &settings,
			Currency: "EUR",
			TimeZone: "Europe/Lisbon",
		},
		fb: fb,
	}
}

// WithName sets the business name
func (bb *BusinessBuilder) WithName(name string) *BusinessBuilder {
	bb.business.Name = name
	return bb
}

// WithEmail sets the business email
func (bb *BusinessBuilder) WithEmail(email string) *BusinessBuilder {
	bb.business.Email = email
	return bb
}

// WithTimeZone sets the time zone
func (bb *BusinessBuilder) WithTimeZone(tz string) *BusinessBuilder {
	bb.business.TimeZone = tz
	return bb
}

// WithCurrency sets the currency
func (bb *BusinessBuilder) WithCurrency(currency string) *BusinessBuilder {
	bb.business.Currency = currency
	return bb
}

// WithSettings sets business settings as JSON string
func (bb *BusinessBuilder) WithSettings(settings string) *BusinessBuilder {
	bb.business.Settings = &settings
	return bb
}

// Build creates the business in the database
func (bb *BusinessBuilder) Build() (*domain.Business, error) {
	if err := bb.fb.db.Create(bb.business).Error; err != nil {
		return nil, fmt.Errorf("failed to create business: %w", err)
	}
	return bb.business, nil
}

// TestData holds commonly used test data
type TestData struct {
	Users      []*domain.User
	Businesses []*domain.Business
	Staff      []*domain.Staff
}

// CreateBasicTestData creates a basic set of test data
func (fb *FixtureBuilder) CreateBasicTestData() (*TestData, error) {
	td := &TestData{
		Users:      make([]*domain.User, 0),
		Businesses: make([]*domain.Business, 0),
		Staff:      make([]*domain.Staff, 0),
	}

	// Create users
	owner, err := fb.NewUser().
		WithEmail("owner@example.com").
		WithName("Business", "Owner").
		Build()
	if err != nil {
		return nil, err
	}
	td.Users = append(td.Users, owner)

	provider, err := fb.NewUser().
		WithEmail("provider@example.com").
		WithName("Service", "Provider").
		Build()
	if err != nil {
		return nil, err
	}
	td.Users = append(td.Users, provider)

	// Create business
	business, err := fb.NewBusiness(owner.ID).
		WithName("Test Beauty Salon").
		WithEmail("salon@example.com").
		Build()
	if err != nil {
		return nil, err
	}
	td.Businesses = append(td.Businesses, business)

	// Create staff
	ownerStaff, err := fb.CreateStaff(business.ID, owner.ID, func(s *domain.Staff) {
		s.Role = domain.BusinessRoleOwner
	})
	if err != nil {
		return nil, err
	}
	td.Staff = append(td.Staff, ownerStaff)

	providerStaff, err := fb.CreateStaff(business.ID, provider.ID, func(s *domain.Staff) {
		s.Role = domain.BusinessRoleEmployee
	})
	if err != nil {
		return nil, err
	}
	td.Staff = append(td.Staff, providerStaff)

	return td, nil
}

// CleanupTestData removes all test data created by a specific test
func (fb *FixtureBuilder) CleanupTestData(data *TestData) error {
	// Delete in reverse order of dependencies

	// Delete staff
	for _, staff := range data.Staff {
		if err := fb.db.Delete(staff).Error; err != nil {
			return fmt.Errorf("failed to delete staff: %w", err)
		}
	}

	// Delete businesses
	for _, business := range data.Businesses {
		if err := fb.db.Delete(business).Error; err != nil {
			return fmt.Errorf("failed to delete business: %w", err)
		}
	}

	// Delete users
	for _, user := range data.Users {
		if err := fb.db.Delete(user).Error; err != nil {
			return fmt.Errorf("failed to delete user: %w", err)
		}
	}

	return nil
}
