package models

import (
	"database/sql/driver"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

// StaffEmploymentType represents the type of employment for a staff member
type StaffEmploymentType string

const (
	// StaffEmploymentTypeFull represents full-time employment
	StaffEmploymentTypeFull StaffEmploymentType = "full-time"
	// StaffEmploymentTypePart represents part-time employment
	StaffEmploymentTypePart StaffEmploymentType = "part-time"
	// StaffEmploymentTypeContract represents contract/freelance employment
	StaffEmploymentTypeContract StaffEmploymentType = "contract"
	// StaffEmploymentTypeIntern represents internship
	StaffEmploymentTypeIntern StaffEmploymentType = "intern"
)

// ExceptionType represents the type of availability exception
type ExceptionType string

const (
	// ExceptionTypeTimeOff represents time off or vacation
	ExceptionTypeTimeOff ExceptionType = "time_off"
	// ExceptionTypeHoliday represents a holiday
	ExceptionTypeHoliday ExceptionType = "holiday"
	// ExceptionTypeCustomHours represents custom working hours
	ExceptionTypeCustomHours ExceptionType = "custom_hours"
)

// PerformancePeriod represents the period for staff performance metrics
type PerformancePeriod string

const (
	// PerformancePeriodDaily represents daily metrics
	PerformancePeriodDaily PerformancePeriod = "daily"
	// PerformancePeriodWeekly represents weekly metrics
	PerformancePeriodWeekly PerformancePeriod = "weekly"
	// PerformancePeriodMonthly represents monthly metrics
	PerformancePeriodMonthly PerformancePeriod = "monthly"
	// PerformancePeriodYearly represents yearly metrics
	PerformancePeriodYearly PerformancePeriod = "yearly"
)

// Staff represents a staff member in a business
type Staff struct {
	BaseModel
	BusinessID      uuid.UUID           `gorm:"type:uuid;not null;index" json:"business_id"`
	Business        Business            `gorm:"foreignKey:BusinessID" json:"business"`
	UserID          uuid.UUID           `gorm:"type:uuid;not null;index" json:"user_id"`
	User            User                `gorm:"foreignKey:UserID" json:"user"`
	Position        string              `gorm:"not null" json:"position"`
	Bio             string              `json:"bio"`
	SpecialtyAreas  SpecialtyAreas      `gorm:"type:text[]" json:"specialty_areas"`
	ProfileImageURL string              `json:"profile_image_url"`
	WorkingHours    WorkingHours        `gorm:"type:jsonb" json:"working_hours"`
	IsActive        bool                `gorm:"not null;default:true" json:"is_active"`
	EmploymentType  StaffEmploymentType `gorm:"type:text;not null" json:"employment_type"`
	JoinDate        time.Time           `gorm:"not null" json:"join_date"`
	EndDate         *time.Time          `json:"end_date"`
	CommissionRate  float64             `gorm:"type:decimal(5,2)" json:"commission_rate"`
}

// TableName overrides the table name
func (Staff) TableName() string {
	return "staff"
}

// SpecialtyAreas is a string array that can be stored as PostgreSQL text array
type SpecialtyAreas []string

// Scan implements the sql.Scanner interface for SpecialtyAreas
func (sa *SpecialtyAreas) Scan(value interface{}) error {
	if value == nil {
		*sa = SpecialtyAreas{}
		return nil
	}

	switch v := value.(type) {
	case []byte:
		// For PostgreSQL array format like {value1,value2}
		str := string(v)
		return sa.parsePostgreSQLArray(str)
	case string:
		// For PostgreSQL array format like {value1,value2}
		return sa.parsePostgreSQLArray(v)
	case []string:
		*sa = SpecialtyAreas(v)
		return nil
	default:
		return fmt.Errorf("cannot scan %T into SpecialtyAreas", value)
	}
}

// parsePostgreSQLArray parses a PostgreSQL array string format
func (sa *SpecialtyAreas) parsePostgreSQLArray(str string) error {
	if str == "{}" || str == "" {
		*sa = SpecialtyAreas{}
		return nil
	}
	// Remove the braces and split by comma
	str = strings.Trim(str, "{}")
	if str == "" {
		*sa = SpecialtyAreas{}
		return nil
	}
	parts := strings.Split(str, ",")
	result := make(SpecialtyAreas, len(parts))
	for i, part := range parts {
		result[i] = strings.Trim(part, `"`)
	}
	*sa = result
	return nil
}

// Value implements the driver.Valuer interface for SpecialtyAreas
func (sa SpecialtyAreas) Value() (driver.Value, error) {
	if len(sa) == 0 {
		return "{}", nil
	}
	// Format as PostgreSQL array: {"value1","value2"}
	result := "{"
	for i, item := range sa {
		if i > 0 {
			result += ","
		}
		result += `"` + item + `"`
	}
	result += "}"
	return result, nil
}

// ServiceAssignment represents the assignment of services to staff members
type ServiceAssignment struct {
	BaseModel
	BusinessID uuid.UUID `gorm:"type:uuid;not null;index" json:"business_id"`
	Business   Business  `gorm:"foreignKey:BusinessID" json:"business"`
	StaffID    uuid.UUID `gorm:"type:uuid;not null;index" json:"staff_id"`
	Staff      Staff     `gorm:"foreignKey:StaffID" json:"staff"`
	ServiceID  uuid.UUID `gorm:"type:uuid;not null;index" json:"service_id"`
	Service    Service   `gorm:"foreignKey:ServiceID" json:"service"`
	IsActive   bool      `gorm:"not null;default:true" json:"is_active"`
}

// TableName overrides the table name
func (ServiceAssignment) TableName() string {
	return "service_assignment"
}

// AvailabilityException represents exceptions to a staff member's regular working hours
type AvailabilityException struct {
	BaseModel
	BusinessID     uuid.UUID     `gorm:"type:uuid;not null;index" json:"business_id"`
	Business       Business      `gorm:"foreignKey:BusinessID" json:"business"`
	StaffID        uuid.UUID     `gorm:"type:uuid;not null;index" json:"staff_id"`
	Staff          Staff         `gorm:"foreignKey:StaffID" json:"staff"`
	ExceptionType  ExceptionType `gorm:"type:text;not null" json:"exception_type"`
	StartTime      time.Time     `gorm:"not null" json:"start_time"`
	EndTime        time.Time     `gorm:"not null" json:"end_time"`
	IsFullDay      bool          `gorm:"not null;default:false" json:"is_full_day"`
	IsRecurring    bool          `gorm:"not null;default:false" json:"is_recurring"`
	RecurrenceRule string        `json:"recurrence_rule"`
	Notes          string        `json:"notes"`
}

// TableName overrides the table name
func (AvailabilityException) TableName() string {
	return "availability_exception"
}

// StaffPerformance represents performance metrics for a staff member
type StaffPerformance struct {
	ID                    uuid.UUID         `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	BusinessID            uuid.UUID         `gorm:"type:uuid;not null;index" json:"business_id"`
	Business              Business          `gorm:"foreignKey:BusinessID" json:"business"`
	StaffID               uuid.UUID         `gorm:"type:uuid;not null;index" json:"staff_id"`
	Staff                 Staff             `gorm:"foreignKey:StaffID" json:"staff"`
	Period                PerformancePeriod `gorm:"type:text;not null" json:"period"`
	StartDate             time.Time         `gorm:"not null" json:"start_date"`
	EndDate               time.Time         `gorm:"not null" json:"end_date"`
	TotalAppointments     int               `gorm:"not null;default:0" json:"total_appointments"`
	CompletedAppointments int               `gorm:"not null;default:0" json:"completed_appointments"`
	CanceledAppointments  int               `gorm:"not null;default:0" json:"canceled_appointments"`
	NoShowAppointments    int               `gorm:"not null;default:0" json:"no_show_appointments"`
	TotalRevenue          float64           `gorm:"type:decimal(10,2);not null;default:0" json:"total_revenue"`
	AverageRating         float64           `gorm:"type:decimal(3,2)" json:"average_rating"`
	ClientRetentionRate   float64           `gorm:"type:decimal(5,2)" json:"client_retention_rate"`
	NewClients            int               `gorm:"not null;default:0" json:"new_clients"`
	ReturnClients         int               `gorm:"not null;default:0" json:"return_clients"`
	CreatedAt             time.Time         `gorm:"not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt             time.Time         `gorm:"not null;default:CURRENT_TIMESTAMP" json:"updated_at"`
}

// TableName overrides the table name
func (StaffPerformance) TableName() string {
	return "staff_performance"
}
