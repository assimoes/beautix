# BeautyBiz Development Tasks - Customer Journey Focused

This document outlines the development tasks organized by customer journeys, prioritizing the most essential user flows first.

## Phase 1: Core Business Setup & Basic Booking (MVP)

### Common Architecture Components

#### Base Models (All models extend this)
```go
type BaseModel struct {
    ID        string    `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
    CreatedAt time.Time `gorm:"not null;default:now()"`
    CreatedBy *string   `gorm:"type:uuid"`
    UpdatedAt time.Time `gorm:"not null;default:now()"`
    UpdatedBy *string   `gorm:"type:uuid"`
    DeletedAt *time.Time `gorm:"index"`
    DeletedBy *string   `gorm:"type:uuid"`
}
```

#### Common Repository Interface
```go
type BaseRepository[T any] interface {
    Create(entity *T) error
    GetByID(id string) (*T, error)
    Update(entity *T) error
    Delete(id string) error
    List(page, pageSize int) ([]*T, int64, error)
    FindBy(criteria map[string]interface{}) ([]*T, error)
    ExistsByID(id string) (bool, error)
}
```

#### Common Service Interface
```go
type BaseService[CreateDTO, UpdateDTO, ResponseDTO any] interface {
    Create(dto CreateDTO) (*ResponseDTO, error)
    GetByID(id string) (*ResponseDTO, error)
    Update(id string, dto UpdateDTO) (*ResponseDTO, error)
    Delete(id string) error
    List(page, pageSize int) ([]*ResponseDTO, int64, error)
    FindBy(criteria map[string]interface{}) ([]*ResponseDTO, error)
}
```

### 1. Business Owner Onboarding Journey
**Goal**: Business owner can set up their beauty business on the platform

#### **Task 1.1**: Authentication & Business Registration
**Models Needed:**
- `User` model (maps to users table)
- `UserConnectedAccount` model (maps to user_connected_accounts table)
- `Business` model (maps to businesses table)

**Repositories:**
- `UserRepository` extending `BaseRepository[User]`
  - `FindByEmail(email string) (*User, error)`
  - `FindByClerkID(clerkID string) (*User, error)`
  - `UpdateClerkID(userID, clerkID string) error`

- `BusinessRepository` extending `BaseRepository[Business]`
  - `FindByUserID(userID string) ([]*Business, error)`
  - `FindByName(name string) ([]*Business, error)`
  - `ExistsByName(name string) (bool, error)`

**Services:**
- `AuthService`
  - `RegisterWithClerk(clerkUserData ClerkUserDTO) (*UserResponseDTO, error)`
  - `SyncClerkUser(clerkID string, userData ClerkUserDTO) (*UserResponseDTO, error)`
  - `GetCurrentUser(clerkID string) (*UserResponseDTO, error)`

- `BusinessService` extending `BaseService[CreateBusinessDTO, UpdateBusinessDTO, BusinessResponseDTO]`
  - `CreateForOwner(userID string, dto CreateBusinessDTO) (*BusinessResponseDTO, error)`
  - `GetByOwner(userID string) ([]*BusinessResponseDTO, error)`
  - `VerifyBusiness(businessID string) error`

**DTOs:**
```go
type CreateBusinessDTO struct {
    Name           string            `json:"name" validate:"required,min=2,max=100"`
    DisplayName    *string           `json:"display_name,omitempty"`
    BusinessType   *string           `json:"business_type,omitempty"`
    Email          string            `json:"email" validate:"required,email"`
    Website        *string           `json:"website,omitempty"`
    Currency       string            `json:"currency" validate:"required,len=3"`
    TimeZone       string            `json:"time_zone" validate:"required"`
}

type BusinessResponseDTO struct {
    ID             string            `json:"id"`
    Name           string            `json:"name"`
    DisplayName    *string           `json:"display_name"`
    Email          string            `json:"email"`
    IsVerified     bool              `json:"is_verified"`
    SubscriptionTier string          `json:"subscription_tier"`
    CreatedAt      time.Time         `json:"created_at"`
}
```

#### **Task 1.2**: Business Profile Setup
**Models Needed:**
- `BusinessLocation` model (maps to business_locations table)
- `BusinessSettings` model (maps to business_settings table)

**Repositories:**
- `BusinessLocationRepository` extending `BaseRepository[BusinessLocation]`
  - `FindByBusinessID(businessID string) ([]*BusinessLocation, error)`
  - `GetMainLocation(businessID string) (*BusinessLocation, error)`
  - `SetMainLocation(businessID, locationID string) error`

- `BusinessSettingsRepository` extending `BaseRepository[BusinessSettings]`
  - `GetByBusinessID(businessID string) (*BusinessSettings, error)`
  - `UpdateByBusinessID(businessID string, settings *BusinessSettings) error`

**Services:**
- `BusinessLocationService` extending `BaseService[CreateLocationDTO, UpdateLocationDTO, LocationResponseDTO]`
  - `CreateForBusiness(businessID string, dto CreateLocationDTO) (*LocationResponseDTO, error)`
  - `GetByBusiness(businessID string) ([]*LocationResponseDTO, error)`
  - `SetAsMain(businessID, locationID string) error`

- `BusinessSettingsService` extending `BaseService[CreateSettingsDTO, UpdateSettingsDTO, SettingsResponseDTO]`
  - `GetOrCreateForBusiness(businessID string) (*SettingsResponseDTO, error)`
  - `UpdateBusinessHours(businessID string, hours BusinessHoursDTO) error`

#### **Task 1.3**: Service Catalog Management
**Models Needed:**
- `ServiceCategory` model (maps to service_categories table)
- `Service` model (maps to services table)

**Repositories:**
- `ServiceCategoryRepository` extending `BaseRepository[ServiceCategory]`
  - `FindByBusinessID(businessID string) ([]*ServiceCategory, error)`
  - `ExistsByNameAndBusiness(name, businessID string) (bool, error)`
  - `GetByDisplayOrder(businessID string) ([]*ServiceCategory, error)`

- `ServiceRepository` extending `BaseRepository[Service]`
  - `FindByBusinessID(businessID string) ([]*Service, error)`
  - `FindByCategory(businessID, category string) ([]*Service, error)`
  - `FindActiveByBusiness(businessID string) ([]*Service, error)`
  - `UpdatePricing(serviceID string, price decimal.Decimal) error`

**Services:**
- `ServiceCategoryService` extending `BaseService[CreateCategoryDTO, UpdateCategoryDTO, CategoryResponseDTO]`
  - `CreateForBusiness(businessID string, dto CreateCategoryDTO) (*CategoryResponseDTO, error)`
  - `GetByBusiness(businessID string) ([]*CategoryResponseDTO, error)`
  - `ReorderCategories(businessID string, categoryOrders []CategoryOrderDTO) error`

- `ServiceManagementService` extending `BaseService[CreateServiceDTO, UpdateServiceDTO, ServiceResponseDTO]`
  - `CreateForBusiness(businessID string, dto CreateServiceDTO) (*ServiceResponseDTO, error)`
  - `GetByBusiness(businessID string) ([]*ServiceResponseDTO, error)`
  - `GetByCategory(businessID, categoryID string) ([]*ServiceResponseDTO, error)`
  - `UpdateAvailability(serviceID string, isActive bool) error`
  - `BulkUpdatePricing(businessID string, updates []ServicePricingUpdateDTO) error`

#### **Task 1.4**: Staff Management
**Models Needed:**
- `Staff` model (maps to staff table)
- `ServiceAssignment` model (maps to service_assignment table)

**Repositories:**
- `StaffRepository` extending `BaseRepository[Staff]`
  - `FindByBusinessID(businessID string) ([]*Staff, error)`
  - `FindByUserID(userID string) ([]*Staff, error)`
  - `FindActiveByBusiness(businessID string) ([]*Staff, error)`
  - `ExistsByUserAndBusiness(userID, businessID string) (bool, error)`

- `ServiceAssignmentRepository` extending `BaseRepository[ServiceAssignment]`
  - `FindByStaffID(staffID string) ([]*ServiceAssignment, error)`
  - `FindByServiceID(serviceID string) ([]*ServiceAssignment, error)`
  - `FindByBusinessID(businessID string) ([]*ServiceAssignment, error)`
  - `ExistsByStaffAndService(staffID, serviceID string) (bool, error)`
  - `DeleteByStaffAndService(staffID, serviceID string) error`

**Services:**
- `StaffService` extending `BaseService[CreateStaffDTO, UpdateStaffDTO, StaffResponseDTO]`
  - `InviteStaff(businessID string, dto InviteStaffDTO) (*StaffResponseDTO, error)`
  - `GetByBusiness(businessID string) ([]*StaffResponseDTO, error)`
  - `UpdateProfile(staffID string, dto UpdateStaffProfileDTO) (*StaffResponseDTO, error)`
  - `SetAvailability(staffID string, isActive bool) error`
  - `GetStaffWithServices(staffID string) (*StaffWithServicesDTO, error)`

- `ServiceAssignmentService`
  - `AssignServiceToStaff(staffID, serviceID string) error`
  - `UnassignServiceFromStaff(staffID, serviceID string) error`
  - `GetStaffServices(staffID string) ([]*ServiceResponseDTO, error)`
  - `GetServiceStaff(serviceID string) ([]*StaffResponseDTO, error)`
  - `BulkAssignServices(staffID string, serviceIDs []string) error`
  - `BulkUnassignServices(staffID string, serviceIDs []string) error`

### 2. Client Discovery & Booking Journey
**Goal**: Clients can find and book appointments with beauty service providers

#### **Task 2.1**: Service Discovery
**Models Needed:**
- `Client` model (maps to clients table)
- Reuse: `Business`, `Service`, `Staff`, `BusinessLocation`

**Repositories:**
- `ClientRepository` extending `BaseRepository[Client]`
  - `FindByBusinessID(businessID string) ([]*Client, error)`
  - `FindByEmail(email string) (*Client, error)`
  - `FindByUserID(userID string) ([]*Client, error)`
  - `ExistsByEmailAndBusiness(email, businessID string) (bool, error)`

- Enhanced `BusinessRepository` methods:
  - `FindActiveBusinesses(page, pageSize int) ([]*Business, int64, error)`
  - `SearchByLocation(city, country string) ([]*Business, error)`
  - `SearchByService(serviceName string) ([]*Business, error)`
  - `GetBusinessWithDetails(businessID string) (*Business, error)`

**Services:**
- `PublicBusinessService` (read-only public access)
  - `SearchBusinesses(criteria SearchCriteriaDTO) ([]*PublicBusinessDTO, error)`
  - `GetBusinessDetails(businessID string) (*BusinessDetailsDTO, error)`
  - `GetBusinessServices(businessID string) ([]*PublicServiceDTO, error)`
  - `GetBusinessStaff(businessID string) ([]*PublicStaffDTO, error)`
  - `FindServicesByCategory(businessID, category string) ([]*PublicServiceDTO, error)`

- `ClientService` extending `BaseService[CreateClientDTO, UpdateClientDTO, ClientResponseDTO]`
  - `CreateOrGetByEmail(businessID string, dto CreateClientDTO) (*ClientResponseDTO, error)`
  - `GetByBusiness(businessID string) ([]*ClientResponseDTO, error)`
  - `LinkToUser(clientID, userID string) error`

#### **Task 2.2**: Appointment Booking Flow
**Models Needed:**
- `Appointment` model (maps to appointments table)
- `AppointmentServices` model (maps to appointment_services table)
- `AppointmentBookingRules` model (maps to appointment_booking_rules table)

**Repositories:**
- `AppointmentRepository` extending `BaseRepository[Appointment]`
  - `FindByBusinessID(businessID string, filters AppointmentFilters) ([]*Appointment, error)`
  - `FindByClientID(clientID string) ([]*Appointment, error)`
  - `FindByStaffID(staffID string, dateRange DateRange) ([]*Appointment, error)`
  - `FindByDateRange(staffID string, start, end time.Time) ([]*Appointment, error)`
  - `CheckOverlap(staffID string, start, end time.Time, excludeID *string) (bool, error)`
  - `GetUpcomingByStaff(staffID string, limit int) ([]*Appointment, error)`

- `AppointmentServicesRepository` extending `BaseRepository[AppointmentServices]`
  - `FindByAppointmentID(appointmentID string) ([]*AppointmentServices, error)`
  - `DeleteByAppointmentID(appointmentID string) error`

- `AppointmentBookingRulesRepository` extending `BaseRepository[AppointmentBookingRules]`
  - `GetByBusinessID(businessID string) (*AppointmentBookingRules, error)`

**Services:**
- `AvailabilityService`
  - `GetAvailableSlots(staffID, serviceID string, date time.Time) ([]*TimeSlotDTO, error)`
  - `CheckSlotAvailability(staffID string, start, end time.Time) (bool, error)`
  - `GetStaffAvailability(staffID string, dateRange DateRange) (*AvailabilityResponseDTO, error)`
  - `CalculateAppointmentEnd(serviceID string, start time.Time) (time.Time, error)`

- `AppointmentBookingService`
  - `CreateAppointment(dto CreateAppointmentDTO) (*AppointmentResponseDTO, error)`
  - `ValidateBookingRequest(dto CreateAppointmentDTO) error`
  - `GetBookingOptions(businessID, serviceID string, date time.Time) (*BookingOptionsDTO, error)`
  - `CheckBusinessRules(businessID string, dto CreateAppointmentDTO) error`

- `AppointmentService` extending `BaseService[CreateAppointmentDTO, UpdateAppointmentDTO, AppointmentResponseDTO]`
  - `GetByBusiness(businessID string, filters AppointmentFilters) ([]*AppointmentResponseDTO, error)`
  - `GetByClient(clientID string) ([]*AppointmentResponseDTO, error)`
  - `GetByStaff(staffID string, dateRange DateRange) ([]*AppointmentResponseDTO, error)`
  - `UpdateStatus(appointmentID string, status AppointmentStatus) error`
  - `AddServices(appointmentID string, services []AppointmentServiceDTO) error`

#### **Task 2.3**: Client Account Management
**Models Needed:**
- Reuse: `Client`, `User`, `Appointment`

**Services Enhanced:**
- Enhanced `ClientService`:
  - `GetClientAppointments(clientID string, status *string) ([]*AppointmentResponseDTO, error)`
  - `UpdatePreferences(clientID string, preferences ClientPreferencesDTO) error`
  - `GetClientHistory(clientID string, businessID string) (*ClientHistoryDTO, error)`

- `ClientProfileService`
  - `CreateProfile(userID string, dto CreateClientProfileDTO) (*ClientProfileDTO, error)`
  - `UpdateProfile(clientID string, dto UpdateClientProfileDTO) (*ClientProfileDTO, error)`
  - `GetProfile(userID string) (*ClientProfileDTO, error)`
  - `AddBusinessToProfile(userID, businessID string) error`

### 3. Appointment Management Journey
**Goal**: Both business owners/staff and clients can manage appointments efficiently

#### **Task 3.1**: Business Appointment Dashboard
**Models Needed:**
- `AppointmentNotes` model (maps to appointment_notes table)
- Reuse: `Appointment`, `Client`, `Staff`, `Service`

**Repositories:**
- `AppointmentNotesRepository` extending `BaseRepository[AppointmentNotes]`
  - `FindByAppointmentID(appointmentID string) ([]*AppointmentNotes, error)`
  - `FindByStaffID(staffID string, dateRange DateRange) ([]*AppointmentNotes, error)`

- Enhanced `AppointmentRepository`:
  - `GetDashboardData(businessID string, date time.Time) (*DashboardData, error)`
  - `GetCalendarView(businessID string, start, end time.Time) ([]*CalendarAppointment, error)`
  - `GetByStatus(businessID string, status AppointmentStatus) ([]*Appointment, error)`
  - `GetRecentByBusiness(businessID string, limit int) ([]*Appointment, error)`

**Services:**
- `AppointmentDashboardService`
  - `GetDashboardOverview(businessID string, date time.Time) (*DashboardOverviewDTO, error)`
  - `GetCalendarView(businessID string, view CalendarView, date time.Time) (*CalendarViewDTO, error)`
  - `GetAppointmentDetails(appointmentID string) (*AppointmentDetailsDTO, error)`
  - `GetDailySchedule(businessID string, date time.Time) ([]*ScheduleItemDTO, error)`

- `AppointmentNotesService` extending `BaseService[CreateNoteDTO, UpdateNoteDTO, NoteResponseDTO]`
  - `AddNote(appointmentID string, dto CreateNoteDTO) (*NoteResponseDTO, error)`
  - `GetByAppointment(appointmentID string) ([]*NoteResponseDTO, error)`
  - `UpdateNote(noteID string, content string) error`

- Enhanced `AppointmentService`:
  - `RescheduleAppointment(appointmentID string, newStart time.Time) (*AppointmentResponseDTO, error)`
  - `CancelAppointment(appointmentID string, reason string) error`
  - `CompleteAppointment(appointmentID string, dto CompleteAppointmentDTO) error`
  - `MarkNoShow(appointmentID string) error`

#### **Task 3.2**: Client Appointment Management
**Models Needed:**
- `AppointmentReminders` model (maps to appointment_reminders table)

**Repositories:**
- `AppointmentRemindersRepository` extending `BaseRepository[AppointmentReminders]`
  - `FindByAppointmentID(appointmentID string) ([]*AppointmentReminders, error)`
  - `FindPendingReminders(beforeTime time.Time) ([]*AppointmentReminders, error)`
  - `MarkAsSent(reminderID string) error`

**Services:**
- `ClientAppointmentService`
  - `GetUpcomingAppointments(clientID string) ([]*ClientAppointmentDTO, error)`
  - `GetAppointmentHistory(clientID string, page, pageSize int) ([]*ClientAppointmentDTO, int64, error)`
  - `RequestReschedule(appointmentID string, dto RescheduleRequestDTO) error`
  - `CancelAppointment(appointmentID string, reason string) error`
  - `ConfirmAppointment(appointmentID string) error`

- `ReminderService` extending `BaseService[CreateReminderDTO, UpdateReminderDTO, ReminderResponseDTO]`
  - `ScheduleReminders(appointmentID string) error`
  - `ProcessPendingReminders() error`
  - `SendReminder(reminderID string) error`
  - `CancelAppointmentReminders(appointmentID string) error`

#### **Task 3.3**: Staff Schedule Management
**Models Needed:**
- `AvailabilityException` model (maps to availability_exception table)

**Repositories:**
- `AvailabilityExceptionRepository` extending `BaseRepository[AvailabilityException]`
  - `FindByStaffID(staffID string, dateRange DateRange) ([]*AvailabilityException, error)`
  - `FindByBusinessID(businessID string, dateRange DateRange) ([]*AvailabilityException, error)`
  - `CheckConflicts(staffID string, start, end time.Time) ([]*AvailabilityException, error)`

**Services:**
- `StaffScheduleService`
  - `GetSchedule(staffID string, dateRange DateRange) (*StaffScheduleDTO, error)`
  - `GetAvailability(staffID string, date time.Time) (*DailyAvailabilityDTO, error)`
  - `RequestTimeOff(staffID string, dto TimeOffRequestDTO) error`
  - `UpdateWorkingHours(staffID string, hours WorkingHoursDTO) error`

- `AvailabilityExceptionService` extending `BaseService[CreateExceptionDTO, UpdateExceptionDTO, ExceptionResponseDTO]`
  - `CreateForStaff(staffID string, dto CreateExceptionDTO) (*ExceptionResponseDTO, error)`
  - `GetByStaff(staffID string, dateRange DateRange) ([]*ExceptionResponseDTO, error)`
  - `ValidateException(dto CreateExceptionDTO) error`
  - `DeleteRecurring(exceptionID string) error`

---

## Common DTOs and Types

```go
type AppointmentStatus string
const (
    StatusScheduled AppointmentStatus = "scheduled"
    StatusConfirmed AppointmentStatus = "confirmed"
    StatusCompleted AppointmentStatus = "completed"
    StatusCancelled AppointmentStatus = "cancelled"
    StatusNoShow    AppointmentStatus = "no-show"
)

type CalendarView string
const (
    ViewDay   CalendarView = "day"
    ViewWeek  CalendarView = "week"
    ViewMonth CalendarView = "month"
)

type DateRange struct {
    Start time.Time `json:"start"`
    End   time.Time `json:"end"`
}

type AppointmentFilters struct {
    Status    *AppointmentStatus `json:"status,omitempty"`
    StaffID   *string           `json:"staff_id,omitempty"`
    ServiceID *string           `json:"service_id,omitempty"`
    DateRange *DateRange        `json:"date_range,omitempty"`
}
```

---

## Phase 1 Implementation Order

### **Sprint 1** (Foundation - 2 weeks)
1. Set up project structure with Clean Architecture
2. Implement BaseModel, BaseRepository, BaseService patterns
3. Create and test database connection
4. Implement User and Business models with basic CRUD
5. Set up Clerk authentication integration
6. Create basic GraphQL schema and resolvers

**Deliverable**: Basic authentication and business registration working

### **Sprint 2** (Core Business Setup - 2 weeks)
1. Complete Business Profile management (Task 1.2)
2. Implement BusinessLocation and BusinessSettings
3. Create Service and ServiceCategory management (Task 1.3)
4. Add file upload capabilities for business photos
5. Create basic business dashboard

**Deliverable**: Business owners can fully set up their businesses

### **Sprint 3** (Staff & Service Management - 2 weeks)
1. Implement Staff management system (Task 1.4)
2. Create ServiceAssignment functionality
3. Add staff working hours and availability
4. Implement staff invitation system
5. Create staff dashboard view

**Deliverable**: Business owners can manage staff and assign services

### **Sprint 4** (Client & Booking Foundation - 2 weeks)
1. Implement Client model and management
2. Create public business discovery (Task 2.1)
3. Build appointment availability checking
4. Implement basic appointment booking (Task 2.2)
5. Add appointment overlap validation in application layer

**Deliverable**: Clients can discover businesses and book appointments

### **Sprint 5** (Appointment Management - 2 weeks)
1. Complete appointment management dashboard (Task 3.1)
2. Implement appointment status updates
3. Add appointment notes functionality
4. Create client appointment management (Task 3.2)
5. Build basic reminder system

**Deliverable**: Full appointment lifecycle management

### **Sprint 6** (Staff Scheduling & Polish - 2 weeks)
1. Implement staff schedule management (Task 3.3)
2. Add availability exceptions
3. Create comprehensive appointment views
4. Add data validation and error handling
5. Performance optimization and testing

**Deliverable**: Complete Phase 1 MVP ready for production

---

## Testing Strategy for Phase 1

### Unit Tests
- Every service method must have unit tests
- Use testify for assertions and mockery for mocks
- Test all business logic edge cases
- Aim for >90% code coverage

### Integration Tests
- Use transaction-based testing pattern
- Test repository methods with real database
- Test service integration with multiple repositories
- Test GraphQL resolvers end-to-end

### Business Logic Tests
- Test appointment overlap prevention
- Test availability calculation algorithms
- Test business rule validation
- Test multi-tenant data isolation

---

## Key Architectural Decisions

1. **Multi-tenancy**: All business logic must respect business_id boundaries
2. **Audit Trail**: All entities track who created/modified/deleted them
3. **Soft Deletes**: Use deleted_at for all user-facing deletions
4. **UUID Keys**: All primary keys are UUIDs for security and scalability
5. **Application-level Validation**: Business rules enforced in services, not database
6. **Clean Architecture**: Strict separation between domain, service, repository, and transport layers
7. **Interface-driven**: All dependencies use interfaces for testability
8. **Transaction Management**: Repository methods are transaction-aware
9. **Error Handling**: Structured error responses with proper HTTP status codes
10. **Performance**: Eager loading for N+1 prevention, pagination for large datasets

## Phase 2: Enhanced Business Operations

### 4. Advanced Scheduling & Availability Journey
**Goal**: Sophisticated scheduling with conflict prevention and optimization

- [ ] **Task 4.1**: Advanced Availability Management
  - Recurring availability patterns
  - Exception handling (holidays, time off)
  - Buffer time between appointments
  - Preparation and cleanup time consideration

- [ ] **Task 4.2**: Smart Booking System
  - Appointment overlap prevention (application-level validation)
  - Automatic booking confirmations
  - Waiting list functionality
  - Multi-service appointment booking

- [ ] **Task 4.3**: Resource Management
  - Equipment/room booking integration
  - Resource availability tracking
  - Resource assignment to appointments

### 5. Client Experience Enhancement Journey
**Goal**: Improved client satisfaction and engagement

- [ ] **Task 5.1**: Communication System
  - Automated appointment reminders (email, SMS)
  - Booking confirmation notifications
  - Appointment status updates
  - Two-way messaging between clients and businesses

- [ ] **Task 5.2**: Feedback & Rating System
  - Post-appointment rating requests
  - Review collection and display
  - Feedback management for businesses
  - Public rating display on business profiles

- [ ] **Task 5.3**: Client Relationship Management
  - Client notes and preferences tracking
  - Appointment history with service details
  - Client health information and allergies tracking
  - Referral source tracking

## Phase 3: Business Growth & Analytics

### 6. Business Performance Journey
**Goal**: Business owners can monitor and optimize their operations

- [ ] **Task 6.1**: Analytics Dashboard
  - Appointment statistics (bookings, cancellations, no-shows)
  - Revenue tracking and reporting
  - Staff performance metrics
  - Client retention analysis

- [ ] **Task 6.2**: Financial Management
  - Service completion tracking
  - Payment method recording
  - Revenue reports by service/staff/period
  - Commission calculation for staff

- [ ] **Task 6.3**: Business Insights
  - Popular services analysis
  - Peak hours identification
  - Client behavior patterns
  - Staff efficiency metrics

### 7. Marketing & Customer Retention Journey
**Goal**: Businesses can attract and retain clients through targeted marketing

- [ ] **Task 7.1**: Loyalty Program Management
  - Program setup (points, visits, spending-based)
  - Client enrollment and tracking
  - Reward redemption system
  - Loyalty program analytics

- [ ] **Task 7.2**: Marketing Campaigns
  - Campaign creation and management
  - Client segmentation and targeting
  - Automated campaign triggers (birthdays, re-engagement)
  - Campaign performance tracking

- [ ] **Task 7.3**: Promotional Tools
  - Discount code creation and management
  - Special offer campaigns
  - Seasonal promotions
  - Referral programs

## Phase 4: Advanced Features & Platform Growth

### 8. Multi-Location Business Journey
**Goal**: Businesses can manage multiple locations efficiently

- [ ] **Task 8.1**: Multi-Location Management
  - Location-specific settings and staff
  - Cross-location appointment booking
  - Location-specific analytics
  - Centralized business management

### 9. Inventory & Product Sales Journey
**Goal**: Businesses can manage inventory and sell products to clients

- [ ] **Task 9.1**: Inventory Management
  - Product catalog management
  - Stock tracking and alerts
  - Supplier management
  - Usage tracking for services

- [ ] **Task 9.2**: Retail Sales
  - Point-of-sale integration
  - Product recommendations
  - Client purchase history
  - Inventory depletion on sales

### 10. Platform Administration Journey
**Goal**: Platform administrators can monitor and manage the entire platform

- [ ] **Task 10.1**: Platform Monitoring
  - System health monitoring
  - Usage analytics across all businesses
  - Performance metrics
  - Error tracking and alerting

- [ ] **Task 10.2**: Support & Moderation
  - Business verification workflow
  - Dispute resolution system
  - Content moderation
  - Support ticket management

## Implementation Notes

### Technology Stack Considerations
- **Backend**: Go with Clean Architecture
- **Database**: PostgreSQL with the consolidated schema
- **Authentication**: Clerk for user management
- **API**: GraphQL for flexible client queries
- **Frontend**: TBD based on client requirements
- **Notifications**: Email and SMS service integration
- **File Storage**: Cloud storage for images/documents

### Development Principles
- Test-Driven Development (TDD) approach
- Domain-driven design for complex business logic
- API-first development for frontend flexibility
- Comprehensive error handling and validation
- Performance optimization for high-traffic scenarios
- Security best practices throughout

### Success Metrics per Journey
- **Business Onboarding**: Time to first appointment booking
- **Client Booking**: Booking conversion rate and completion time
- **Appointment Management**: Appointment adherence rate and satisfaction
- **Business Performance**: Monthly recurring revenue and client retention
- **Marketing**: Campaign engagement and client lifetime value

---

**Current Priority**: Phase 1 - Core Business Setup & Basic Booking (MVP)
**Next Phase**: To be determined based on user feedback and business priorities