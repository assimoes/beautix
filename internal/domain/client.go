package domain

import (
	"context"
	"time"
	"github.com/shopspring/decimal"
)

// Client represents a client of a business
type Client struct {
	BaseModel
	BusinessID   string     `gorm:"not null;type:uuid;index" json:"business_id"`
	UserID       *string    `gorm:"type:uuid;index" json:"user_id,omitempty"` // Optional link to user account
	FirstName    string     `gorm:"not null;size:100" json:"first_name"`
	LastName     string     `gorm:"not null;size:100" json:"last_name"`
	Email        string     `gorm:"not null;size:255;index" json:"email"`
	Phone        *string    `gorm:"size:20" json:"phone,omitempty"`
	DateOfBirth  *time.Time `gorm:"" json:"date_of_birth,omitempty"`
	Gender       *string    `gorm:"size:20" json:"gender,omitempty"`
	Notes        *string    `gorm:"type:text" json:"notes,omitempty"`
	Preferences  *string    `gorm:"type:jsonb;default:'{}'" json:"preferences,omitempty"` // JSON for client preferences
	Allergies    *string    `gorm:"type:text" json:"allergies,omitempty"`
	IsActive     bool       `gorm:"not null;default:true" json:"is_active"`
	ReferralSource *string  `gorm:"size:100" json:"referral_source,omitempty"`
	LastVisit    *time.Time `gorm:"" json:"last_visit,omitempty"`
	TotalVisits  int        `gorm:"not null;default:0" json:"total_visits"`
	TotalSpent   decimal.Decimal `gorm:"type:decimal(10,2);not null;default:0" json:"total_spent"`

	// Relationships
	Business     Business     `gorm:"foreignKey:BusinessID;constraint:OnDelete:CASCADE" json:"business"`
	User         *User        `gorm:"foreignKey:UserID;constraint:OnDelete:SET NULL" json:"user,omitempty"`
	Appointments []Appointment `gorm:"foreignKey:ClientID" json:"appointments,omitempty"`
}

// TableName returns the table name for Client
func (Client) TableName() string { return "clients" }

// Validate validates the client model
func (c *Client) Validate() error {
	if c.BusinessID == "" {
		return ErrValidation
	}
	if c.FirstName == "" {
		return ErrValidation
	}
	if c.LastName == "" {
		return ErrValidation
	}
	if c.Email == "" {
		return ErrValidation
	}
	return nil
}

// GetFullName returns the client's full name
func (c *Client) GetFullName() string {
	return c.FirstName + " " + c.LastName
}

// UpdateVisitStats updates the last visit and total visits
func (c *Client) UpdateVisitStats(visitTime time.Time, amount decimal.Decimal) {
	c.LastVisit = &visitTime
	c.TotalVisits++
	c.TotalSpent = c.TotalSpent.Add(amount)
}

// ClientRepository defines the repository interface for Client
type ClientRepository interface {
	BaseRepository[Client]
	FindByBusinessID(ctx context.Context, businessID string) ([]*Client, error)
	FindByEmail(ctx context.Context, email string) (*Client, error)
	FindByUserID(ctx context.Context, userID string) ([]*Client, error)
	FindByBusinessAndEmail(ctx context.Context, businessID, email string) (*Client, error)
	ExistsByEmailAndBusiness(ctx context.Context, email, businessID string) (bool, error)
	UpdateVisitStats(ctx context.Context, clientID string, visitTime time.Time, amount decimal.Decimal) error
}