# BeautyBiz: Phase 1 (MVP) Data Model with Multi-Tenancy Focus

## Overview

This document presents the essential data model for the BeautyBiz MVP phase, with multi-tenancy as a core architectural principle. The model prioritizes the critical entities required for the initial launch while ensuring proper data isolation between tenants (beauty businesses).

## Multi-Tenancy Approach

BeautyBiz will implement a **row-level multi-tenancy** approach, where:

1. Each table includes a `business_id` column that serves as a tenant identifier
2. Database policies enforce tenant isolation at the query level
3. Application logic maintains tenant context throughout all operations
4. Indexes are optimized for tenant-scoped queries

This approach offers several advantages for our MVP:
- Simplified schema compared to schema-per-tenant
- Easier maintenance and upgrades
- Efficient resource utilization
- Balanced security and performance

## Core Domain Concepts for MVP

For Phase 1, we focus only on the essential entities:

1. **Users & Businesses** - Account management and multi-tenancy foundation
2. **Clients** - Basic client information management
3. **Services** - Beauty services offered by the business
4. **Appointments** - Core scheduling functionality

## MVP Entity Relationship Diagram

```
┌─────────────────┐       ┌─────────────────┐       ┌─────────────────┐
│                 │       │                 │       │                 │
│    Business     │─────┬─│     Users       │───────│  Subscription   │
│                 │     │ │                 │       │                 │
└─────────────────┘     │ └─────────────────┘       └─────────────────┘
        │               │         │
        │               │         │
        │               │         ▼
        │               │ ┌─────────────────┐
        │               │ │                 │
        │               └─│     Staff       │
        │                 │                 │
        │                 └─────────────────┘
        │                         │
        ▼                         │
┌─────────────────┐               │
│                 │               │
│    Clients      │               │
│                 │               │
└─────────────────┘               │
        │                         │
        │                         │
        ▼                         ▼
┌─────────────────┐       ┌─────────────────┐       ┌─────────────────┐
│                 │       │                 │       │                 │
│  Appointments   │───────│    Services     │───────│ Service Category│
│                 │       │                 │       │                 │
└─────────────────┘       └─────────────────┘       └─────────────────┘
```

## Database Schema (DBML) for MVP

```dbml
// DBML for BeautyBiz Platform - MVP with Multi-tenancy Focus

// Enable Row-Level Security for Multi-tenancy
// Note: This would be implemented in the actual database as policies

// Users & Authentication
Table users {
  user_id UUID [pk]
  email VARCHAR(255) [unique, not null]
  password_hash VARCHAR(255) [not null]
  first_name VARCHAR(100) [not null]
  last_name VARCHAR(100) [not null]
  phone VARCHAR(20)
  role VARCHAR(20) [not null, note: 'admin, owner, staff']
  created_at TIMESTAMP [not null, default: `now()`]
  updated_at TIMESTAMP [not null, default: `now()`]
  last_login TIMESTAMP
  is_active BOOLEAN [not null, default: true]
  email_verified BOOLEAN [not null, default: false]
  language_preference VARCHAR(10) [not null, default: 'pt-PT']
  indexes {
    email [unique]
  }
}

// Business Details - Primary tenant entity
Table businesses {
  business_id UUID [pk]
  owner_id UUID [ref: > users.user_id, not null]
  business_name VARCHAR(255) [not null]
  business_type VARCHAR(50) [not null, note: 'salon, spa, individual, mobile, etc.']
  tax_id VARCHAR(50)
  phone VARCHAR(20) [not null]
  email VARCHAR(255) [not null]
  address_line1 VARCHAR(255)
  city VARCHAR(100)
  region VARCHAR(100)
  postal_code VARCHAR(20)
  country VARCHAR(50) [not null, default: 'Portugal']
  time_zone VARCHAR(50) [not null, default: 'Europe/Lisbon']
  business_hours JSON [note: 'Hours of operation for each day']
  created_at TIMESTAMP [not null, default: `now()`]
  updated_at TIMESTAMP [not null, default: `now()`]
  is_active BOOLEAN [not null, default: true]
  subscription_plan VARCHAR(20) [not null, default: 'trial', note: 'trial, solo, salon, premium']
  trial_ends_at TIMESTAMP
  indexes {
    owner_id
  }
}

// Subscription (simplified for MVP)
Table subscriptions {
  subscription_id UUID [pk]
  business_id UUID [ref: > businesses.business_id, not null]
  plan_type VARCHAR(50) [not null, note: 'trial, solo, salon, premium']
  status VARCHAR(20) [not null, default: 'active', note: 'active, canceled, past_due, trialing']
  start_date TIMESTAMP [not null]
  end_date TIMESTAMP
  trial_end_date TIMESTAMP
  billing_cycle VARCHAR(20) [not null, default: 'monthly', note: 'monthly, annual']
  created_at TIMESTAMP [not null, default: `now()`]
  updated_at TIMESTAMP [not null, default: `now()`]
  indexes {
    business_id
    (business_id, status)
  }
}

// Business Staff Members
Table staff {
  staff_id UUID [pk]
  business_id UUID [ref: > businesses.business_id, not null]
  user_id UUID [ref: > users.user_id, not null]
  position VARCHAR(100) [not null]
  bio TEXT
  color VARCHAR(20) [note: 'Color for calendar display']
  is_active BOOLEAN [not null, default: true]
  created_at TIMESTAMP [not null, default: `now()`]
  updated_at TIMESTAMP [not null, default: `now()`]
  indexes {
    business_id
    user_id
    (business_id, is_active)
  }
}

// Staff Working Hours - Simplifies scheduling logic
Table staff_working_hours {
  hours_id UUID [pk]
  staff_id UUID [ref: > staff.staff_id, not null]
  business_id UUID [ref: > businesses.business_id, not null]
  day_of_week INTEGER [not null, note: '0=Sunday, 1=Monday, ..., 6=Saturday']
  start_time TIME [not null]
  end_time TIME [not null]
  is_working BOOLEAN [not null, default: true]
  created_at TIMESTAMP [not null, default: `now()`]
  updated_at TIMESTAMP [not null, default: `now()`]
  indexes {
    staff_id
    business_id
    (staff_id, day_of_week)
    (business_id, day_of_week)
  }
}

// Clients
Table clients {
  client_id UUID [pk]
  business_id UUID [ref: > businesses.business_id, not null]
  first_name VARCHAR(100) [not null]
  last_name VARCHAR(100) [not null]
  email VARCHAR(255)
  phone VARCHAR(20)
  date_of_birth DATE
  notes TEXT
  client_since DATE [not null, default: `CURRENT_DATE`]
  marketing_consent BOOLEAN [default: false]
  created_at TIMESTAMP [not null, default: `now()`]
  updated_at TIMESTAMP [not null, default: `now()`]
  is_active BOOLEAN [not null, default: true]
  indexes {
    business_id
    (business_id, is_active)
    (business_id, email)
    (business_id, phone)
    (business_id, last_name, first_name)
  }
}

// Service Categories
Table service_categories {
  category_id UUID [pk]
  business_id UUID [ref: > businesses.business_id, not null]
  category_name VARCHAR(100) [not null]
  description TEXT
  color VARCHAR(20) [note: 'Color for display purposes']
  sort_order INTEGER [default: 0]
  is_active BOOLEAN [not null, default: true]
  created_at TIMESTAMP [not null, default: `now()`]
  updated_at TIMESTAMP [not null, default: `now()`]
  indexes {
    business_id
    (business_id, is_active)
  }
}

// Services
Table services {
  service_id UUID [pk]
  business_id UUID [ref: > businesses.business_id, not null]
  category_id UUID [ref: > service_categories.category_id]
  service_name VARCHAR(255) [not null]
  description TEXT
  duration INTEGER [not null, note: 'Duration in minutes']
  price DECIMAL(10, 2) [not null]
  color VARCHAR(20) [note: 'Color for calendar display']
  is_active BOOLEAN [not null, default: true]
  created_at TIMESTAMP [not null, default: `now()`]
  updated_at TIMESTAMP [not null, default: `now()`]
  indexes {
    business_id
    category_id
    (business_id, is_active)
    (business_id, category_id)
  }
}

// Service Staff Assignment - Who can perform what services
Table service_staff_assignment {
  assignment_id UUID [pk]
  business_id UUID [ref: > businesses.business_id, not null]
  service_id UUID [ref: > services.service_id, not null]
  staff_id UUID [ref: > staff.staff_id, not null]
  can_perform BOOLEAN [not null, default: true]
  created_at TIMESTAMP [not null, default: `now()`]
  indexes {
    business_id
    service_id
    staff_id
    (business_id, service_id, staff_id)
    (staff_id, service_id)
  }
}

// Appointments
Table appointments {
  appointment_id UUID [pk]
  business_id UUID [ref: > businesses.business_id, not null]
  client_id UUID [ref: > clients.client_id, not null]
  staff_id UUID [ref: > staff.staff_id, not null]
  service_id UUID [ref: > services.service_id, not null]
  start_time TIMESTAMP [not null]
  end_time TIMESTAMP [not null]
  status VARCHAR(20) [not null, default: 'scheduled', note: 'scheduled, confirmed, completed, canceled, no_show']
  notes TEXT
  created_by UUID [ref: > users.user_id]
  created_at TIMESTAMP [not null, default: `now()`]
  updated_at TIMESTAMP [not null, default: `now()`]
  indexes {
    business_id
    client_id
    staff_id
    service_id
    (business_id, start_time)
    (business_id, client_id)
    (business_id, staff_id, start_time)
    (business_id, status)
    (business_id, start_time, end_time)
  }
}

// Appointment Notes & History - For tracking changes to appointments
Table appointment_history {
  history_id UUID [pk]
  business_id UUID [ref: > businesses.business_id, not null]
  appointment_id UUID [ref: > appointments.appointment_id, not null]
  changed_by UUID [ref: > users.user_id, not null]
  action VARCHAR(50) [not null, note: 'created, rescheduled, canceled, completed, etc.']
  previous_status VARCHAR(20)
  new_status VARCHAR(20)
  previous_start TIMESTAMP
  new_start TIMESTAMP
  previous_end TIMESTAMP
  new_end TIMESTAMP
  previous_staff UUID [ref: > staff.staff_id]
  new_staff UUID [ref: > staff.staff_id]
  previous_service UUID [ref: > services.service_id]
  new_service UUID [ref: > services.service_id]
  notes TEXT
  created_at TIMESTAMP [not null, default: `now()`]
  indexes {
    business_id
    appointment_id
    (business_id, appointment_id)
  }
}

// Client Notes - For tracking important client information
Table client_notes {
  note_id UUID [pk]
  business_id UUID [ref: > businesses.business_id, not null]
  client_id UUID [ref: > clients.client_id, not null]
  created_by UUID [ref: > users.user_id, not null]
  note_text TEXT [not null]
  is_important BOOLEAN [default: false]
  created_at TIMESTAMP [not null, default: `now()`]
  updated_at TIMESTAMP [not null, default: `now()`]
  indexes {
    business_id
    client_id
    (business_id, client_id)
    (business_id, client_id, is_important)
  }
}

// Business Notifications Settings
Table notification_settings {
  setting_id UUID [pk]
  business_id UUID [ref: > businesses.business_id, not null]
  notification_type VARCHAR(50) [not null, note: 'appointment_reminder, appointment_confirmation, etc.']
  is_enabled BOOLEAN [not null, default: true]
  template_subject VARCHAR(255)
  template_content TEXT
  send_time_before INTEGER [note: 'Minutes before event to send notification']
  created_at TIMESTAMP [not null, default: `now()`]
  updated_at TIMESTAMP [not null, default: `now()`]
  indexes {
    business_id
    (business_id, notification_type)
  }
}

// Tenant Settings - For business-specific configurations
Table business_settings {
  setting_id UUID [pk]
  business_id UUID [ref: > businesses.business_id, not null, unique]
  calendar_start_hour INTEGER [not null, default: 8, note: '24-hour format start time']
  calendar_end_hour INTEGER [not null, default: 20, note: '24-hour format end time']
  appointment_buffer_minutes INTEGER [not null, default: 0]
  allow_online_booking BOOLEAN [not null, default: false]
  default_appointment_duration INTEGER [not null, default: 60]
  currency VARCHAR(3) [not null, default: 'EUR']
  date_format VARCHAR(20) [not null, default: 'DD/MM/YYYY']
  time_format VARCHAR(20) [not null, default: 'HH:mm']
  created_at TIMESTAMP [not null, default: `now()`]
  updated_at TIMESTAMP [not null, default: `now()`]
  indexes {
    business_id [unique]
  }
}

// User Sessions - For authentication management
Table user_sessions {
  session_id UUID [pk]
  user_id UUID [ref: > users.user_id, not null]
  business_id UUID [ref: > businesses.business_id]
  token VARCHAR(255) [not null]
  ip_address VARCHAR(50)
  user_agent TEXT
  expires_at TIMESTAMP [not null]
  created_at TIMESTAMP [not null, default: `now()`]
  last_activity TIMESTAMP [not null, default: `now()`]
  indexes {
    user_id
    token
    (user_id, business_id)
  }
}

// Permission Scopes - For multi-tenant role-based access control
Table permission_scopes {
  scope_id UUID [pk]
  user_id UUID [ref: > users.user_id, not null]
  business_id UUID [ref: > businesses.business_id, not null]
  resource VARCHAR(50) [not null, note: 'clients, appointments, services, etc.']
  action VARCHAR(50) [not null, note: 'read, write, delete, manage']
  created_at TIMESTAMP [not null, default: `now()`]
  indexes {
    user_id
    business_id
    (user_id, business_id, resource)
  }
}

// Relationships

Ref: businesses.owner_id > users.user_id
Ref: subscriptions.business_id > businesses.business_id
Ref: staff.business_id > businesses.business_id
Ref: staff.user_id > users.user_id
Ref: staff_working_hours.staff_id > staff.staff_id
Ref: staff_working_hours.business_id > businesses.business_id
Ref: clients.business_id > businesses.business_id
Ref: service_categories.business_id > businesses.business_id
Ref: services.business_id > businesses.business_id
Ref: services.category_id > service_categories.category_id
Ref: service_staff_assignment.business_id > businesses.business_id
Ref: service_staff_assignment.service_id > services.service_id
Ref: service_staff_assignment.staff_id > staff.staff_id
Ref: appointments.business_id > businesses.business_id
Ref: appointments.client_id > clients.client_id
Ref: appointments.staff_id > staff.staff_id
Ref: appointments.service_id > services.service_id
Ref: appointments.created_by > users.user_id
Ref: appointment_history.business_id > businesses.business_id
Ref: appointment_history.appointment_id > appointments.appointment_id
Ref: appointment_history.changed_by > users.user_id
Ref: client_notes.business_id > businesses.business_id
Ref: client_notes.client_id > clients.client_id
Ref: client_notes.created_by > users.user_id
Ref: notification_settings.business_id > businesses.business_id
Ref: business_settings.business_id > businesses.business_id
Ref: user_sessions.user_id > users.user_id
Ref: user_sessions.business_id > businesses.business_id
Ref: permission_scopes.user_id > users.user_id
Ref: permission_scopes.business_id > businesses.business_id
```

## Multi-Tenancy Implementation Considerations

### 1. Business ID as Tenant Identifier

All tables (except `users`) include a `business_id` column that serves as the tenant identifier. This approach:

- Ensures data isolation between tenants
- Simplifies queries by including tenant context in all operations
- Allows efficient indexing strategies

### 2. Row-Level Security Policies

At the database level, we will implement row-level security policies that:

```sql
-- Example PostgreSQL RLS policy (would be applied to all tenant tables)
CREATE POLICY tenant_isolation_policy ON appointments
    USING (business_id = current_setting('app.current_tenant')::uuid);
```

This ensures that even if application logic fails, data cannot be leaked across tenant boundaries.

### 3. Connection Pooling Strategy

For multi-tenancy efficiency:

- Use a middleware layer to set tenant context for each request
- Implement connection pooling that maintains tenant separation
- Consider PgBouncer for PostgreSQL connection management

### 4. Indexing Strategy

All tables include carefully designed indexes:

- Every table has a `business_id` index
- Composite indexes lead with `business_id` for efficient tenant-scoped queries
- Indexes are designed for common query patterns within tenant context

### 5. Tenant Context Management

The application will:

- Extract tenant information from authentication
- Set tenant context at the beginning of each request
- Include tenant context in database queries
- Validate tenant permissions for each operation

## Data Migration Path for Future Growth

This MVP data model is designed to scale with the product. As we move beyond Phase 1:

1. **Advanced Features:** Additional tables can be added without disrupting the core model
2. **Performance Optimization:** As tenant data grows, we can implement table partitioning by tenant
3. **Reporting Capabilities:** Analytical tables can be added later without schema changes to operational tables
4. **Integration Capabilities:** External service connectors can be built on top of this foundation

## Security Considerations for Multi-Tenancy

1. **Authentication Flow:**
   - Users authenticate and receive authorized business contexts
   - JWTs include tenant information for stateless verification
   - Session validation checks tenant permissions

2. **Authorization Model:**
   - Permission scopes table enforces granular access control
   - Each user can have different permissions across businesses
   - System validates both user identity AND tenant context for operations

3. **Data Privacy:**
   - No cross-tenant queries are possible at database level
   - Encryption at rest for sensitive client data
   - Audit logging captures tenant context for all operations

## Implementation Guidelines

1. **Middleware Setup:**
```javascript
// Example Express middleware for tenant context
app.use((req, res, next) => {
  const businessId = req.user?.currentBusinessId;
  if (!businessId) {
    return res.status(403).json({ error: 'No business context' });
  }
  // Set tenant context for this request
  req.tenantId = businessId;
  // Set database session variable
  pool.query('SET app.current_tenant = $1', [businessId]);
  next();
});
```

2. **Query Pattern:**
```javascript
// Example query pattern for tenant-scoped operations
async function getClientsByBusiness(businessId) {
  return db.query(
    'SELECT * FROM clients WHERE business_id = $1',
    [businessId]
  );
}
```

3. **Multi-Tenant Validation:**
```javascript
// Example multi-tenant validation
async function validateTenantAccess(userId, businessId, resource, action) {
  const permission = await db.query(
    'SELECT * FROM permission_scopes WHERE user_id = $1 AND business_id = $2 AND resource = $3 AND action = $4',
    [userId, businessId, resource, action]
  );
  return permission.rowCount > 0;
}
```

## Conclusion

This MVP data model provides a solid foundation for the BeautyBiz platform with multi-tenancy as a core architectural principle. By focusing on the essential entities and relationships needed for the initial launch while incorporating robust tenant isolation, we ensure the platform can scale securely as we add more features in subsequent phases.
