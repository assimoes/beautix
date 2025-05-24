# Migration Design Decisions

## Overview
This document explains the design decisions made in migrations up to 000002, particularly regarding the separation of data integrity constraints from business logic and the architectural change from global user roles to business-context roles.

## Migration 000015: Add Missing Audit Fields
This migration standardizes audit fields across all tables by adding:
- `created_by`, `updated_by`, `deleted_by` - Foreign keys to users table
- `deleted_at` - For soft delete support
- Proper indexes on all audit fields

**Rationale**: Consistent audit trails are essential for compliance, debugging, and data integrity.

## Migration 000016: Fix Redundant Fields (Revised)

### What We Added
1. **Basic Data Integrity Constraints**:
   - Time order checks (end_time > start_time)
   - Value range checks (ratings 1-5, non-negative quantities)
   - Simple data validation that prevents invalid data

2. **Performance Indexes**:
   - Indexes on all foreign key columns
   - Partial indexes on deleted_at for soft delete performance

3. **Configuration Tables**:
   - `appointment_booking_rules` table for business-specific settings

### What We Intentionally Did NOT Add

1. **Loyalty Points Calculation Constraint**
   - Original: `CHECK (points_after = points_before + points_change)`
   - **Why removed**: This prevents manual adjustments, corrections, or bulk imports
   - **Better approach**: Validate in application with override capability

2. **Appointment Overlap Triggers**
   - Original: Database triggers to prevent overlapping appointments
   - **Why removed**: 
     - Business rules vary (buffer times, multi-location, travel time)
     - Hard to maintain and debug database triggers
     - Can't easily override for special cases
   - **Better approach**: Service layer validation with configurable rules

3. **Resource Booking Conflict Triggers**
   - Similar reasoning as appointment overlaps
   - **Better approach**: Application-level validation with business rules

## Design Principles

### Keep in Database:
- **Data type constraints** (positive numbers, valid dates)
- **Referential integrity** (foreign keys)
- **Simple validations** that are universally true
- **Indexes** for performance

### Move to Application:
- **Business rules** that might change or have exceptions
- **Complex validations** requiring external data
- **Workflow logic** (approval processes, notifications)
- **Calculations** that might need adjustment capabilities

## Benefits of This Approach

1. **Flexibility**: Business rules can be changed without database migrations
2. **Testability**: Application logic is easier to unit test than triggers
3. **Debugging**: Application logs are more accessible than database trigger errors
4. **Performance**: Avoid complex trigger logic on every insert/update
5. **Override Capability**: Special cases can bypass rules when needed

## Configuration-Driven Rules

The `appointment_booking_rules` table allows businesses to configure:
- Buffer time between appointments
- Advance booking limits
- Double booking permissions
- Staff overlap rules

This makes the system adaptable to different business models without code changes.

## Implementation Guidelines

When implementing the application layer:

1. **Create Service Layer Validators**:
   ```go
   type AppointmentValidator struct {
       rules AppointmentBookingRules
   }
   
   func (v *AppointmentValidator) ValidateOverlap(appointment *Appointment) error {
       // Check for conflicts with buffer time
       // Consider multi-location scenarios
       // Apply business-specific rules
   }
   ```

2. **Allow Override Mechanisms**:
   ```go
   type CreateAppointmentOptions struct {
       SkipOverlapCheck bool // For admin overrides
       ForceBooking     bool // For special cases
   }
   ```

3. **Provide Clear Error Messages**:
   - "Cannot book appointment: conflicts with existing booking at 2:00 PM"
   - "Cannot book: requires 24 hours advance notice"

## Migration Rollback

Both migrations include complete rollback scripts that:
- Remove added columns and constraints
- Restore original schema
- Preserve data where possible

## Migration 000002: Remove User Role and Add Business-Context Roles

### Major Architectural Change
This migration implements a fundamental shift from global user roles to business-context roles:

**Before**: Users had a single global role (admin, owner, staff, user)
**After**: Users have roles only within the context of specific businesses

### Changes Made

1. **Removed Global User Role**:
   - Dropped `role` column from `users` table
   - Removed role-based indexes and constraints
   - Updated application code to remove role dependencies

2. **Enhanced Staff Table for Business Roles**:
   - Added `role` column with business-context values: owner, manager, employee, assistant
   - Added `permissions` JSONB column for granular permissions
   - Added `start_date` and `end_date` for role temporal tracking
   - Created unique constraint ensuring one active role per user per business

### Design Rationale

1. **Multi-Business Support**: Users can now have different roles in different businesses
2. **Granular Permissions**: JSON permissions allow fine-grained access control
3. **Temporal Tracking**: Start/end dates enable role history and transitions
4. **Scalability**: Supports complex organizational structures

### Migration Strategy

- **Data Preservation**: Existing staff position data mapped to appropriate business roles
- **Default Role Assignment**: Users without staff positions default to basic access
- **Rollback Safety**: Complete rollback script restores original structure

### Business Logic Updates

The application now determines user permissions by:
1. Checking if user owns the business (automatic owner role)
2. Looking up active staff record for the business
3. Evaluating role and specific permissions JSON

### User Creation Flow

When a new user is created:
1. User record created without global role
2. Default business created with user's name
3. Staff record created with owner role for default business

## Future Considerations

1. **Audit Log Table**: Consider a separate audit log table for detailed change tracking
2. **Event Sourcing**: For complex business logic, consider event sourcing patterns
3. **Rule Engine**: For very complex rules, consider a rule engine approach
4. **Role Templates**: Pre-defined permission sets for common business roles
5. **Role Hierarchy**: Consider implementing role inheritance (manager includes employee permissions)