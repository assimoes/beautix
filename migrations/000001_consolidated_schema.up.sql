-- Consolidated database schema for BeautyBiz application
-- This migration represents the complete database state after all migrations up to version 16

-- ========================================
-- Extensions
-- ========================================
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- ========================================
-- Users table
-- ========================================
CREATE TABLE public.users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) NOT NULL UNIQUE,
    clerk_id VARCHAR(255) UNIQUE,
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    phone VARCHAR(50),
    role VARCHAR(20) NOT NULL CHECK (role IN ('admin', 'owner', 'staff', 'user')),
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    created_by UUID,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_by UUID,
    deleted_at TIMESTAMP WITH TIME ZONE,
    deleted_by UUID
);

COMMENT ON TABLE public.users IS 'Users of the system with authentication details and basic profile information';
COMMENT ON COLUMN public.users.created_by IS 'User who created this record';
COMMENT ON COLUMN public.users.updated_by IS 'User who last updated this record';
COMMENT ON COLUMN public.users.deleted_by IS 'User who soft deleted this record';
COMMENT ON COLUMN public.users.deleted_at IS 'Timestamp of soft deletion';

-- Create indexes for users table
CREATE INDEX idx_users_email ON public.users(email);
CREATE INDEX idx_users_clerk_id ON public.users(clerk_id);
CREATE INDEX idx_users_role ON public.users(role);
CREATE INDEX idx_users_is_active ON public.users(is_active) WHERE is_active = TRUE;
CREATE INDEX idx_users_created_by ON public.users(created_by);
CREATE INDEX idx_users_updated_by ON public.users(updated_by);
CREATE INDEX idx_users_deleted_by ON public.users(deleted_by);
CREATE INDEX idx_users_deleted_at ON public.users(deleted_at) WHERE deleted_at IS NULL;

-- ========================================
-- User connected accounts table
-- ========================================
CREATE TABLE public.user_connected_accounts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    provider_type VARCHAR(50) NOT NULL,
    provider_id VARCHAR(255) NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    created_by UUID,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_by UUID,
    deleted_at TIMESTAMP WITH TIME ZONE,
    deleted_by UUID,
    CONSTRAINT fk_user_connected_accounts_user FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE
);

CREATE UNIQUE INDEX idx_user_connected_accounts_provider ON public.user_connected_accounts(user_id, provider_type, provider_id);

-- ========================================
-- Businesses table (renamed from business)
-- ========================================
CREATE TABLE public.businesses (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    name VARCHAR(100) NOT NULL,
    display_name VARCHAR(100),
    business_type VARCHAR(50),
    tax_id VARCHAR(50),
    email VARCHAR(255) NOT NULL,
    website VARCHAR(255),
    logo_url VARCHAR(255),
    cover_photo_url VARCHAR(255),
    is_verified BOOLEAN DEFAULT FALSE,
    social_links JSONB DEFAULT '{}'::jsonb,
    settings JSONB DEFAULT '{}'::jsonb,
    currency VARCHAR(3) DEFAULT 'EUR',
    time_zone VARCHAR(50) NOT NULL DEFAULT 'Europe/Lisbon',
    business_hours JSONB,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    subscription_tier VARCHAR(50) NOT NULL DEFAULT 'free',
    trial_ends_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    created_by UUID,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_by UUID,
    deleted_at TIMESTAMP WITH TIME ZONE,
    deleted_by UUID,
    CONSTRAINT fk_businesses_user FOREIGN KEY (user_id) REFERENCES public.users(id)
);

COMMENT ON TABLE public.businesses IS 'Business entities that serve as tenants in the multi-tenant system';
COMMENT ON COLUMN public.businesses.time_zone IS 'Default timezone for the business';
COMMENT ON COLUMN public.businesses.email IS 'Primary contact email for the business';

-- Create indexes for businesses table
CREATE INDEX idx_businesses_user_id ON public.businesses(user_id);
CREATE INDEX idx_businesses_name ON public.businesses(name);
CREATE INDEX idx_businesses_is_active ON public.businesses(is_active) WHERE is_active = TRUE;
CREATE INDEX idx_businesses_subscription_tier ON public.businesses(subscription_tier);

-- ========================================
-- Business locations table
-- ========================================
CREATE TABLE public.business_locations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    business_id UUID NOT NULL,
    name VARCHAR(100) NOT NULL,
    address VARCHAR(255),
    city VARCHAR(100),
    state VARCHAR(100),
    postal_code VARCHAR(20),
    country VARCHAR(50) NOT NULL DEFAULT 'Portugal',
    phone VARCHAR(50),
    email VARCHAR(255),
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    is_main BOOLEAN NOT NULL DEFAULT FALSE,
    timezone VARCHAR(50) NOT NULL DEFAULT 'Europe/Lisbon',
    settings JSONB DEFAULT '{}'::jsonb,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    created_by UUID,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_by UUID,
    deleted_at TIMESTAMP WITH TIME ZONE,
    deleted_by UUID,
    CONSTRAINT fk_business_locations_business FOREIGN KEY (business_id) REFERENCES public.businesses(id) ON DELETE CASCADE
);

COMMENT ON TABLE public.business_locations IS 'Physical locations for a business (supports multi-location businesses)';

-- Create indexes for business_locations table
CREATE INDEX idx_business_locations_business_id ON public.business_locations(business_id);
CREATE INDEX idx_business_locations_city ON public.business_locations(city);
CREATE INDEX idx_business_locations_postal_code ON public.business_locations(postal_code);
CREATE INDEX idx_business_locations_is_main ON public.business_locations(is_main);

-- Add constraint to ensure only one main location per business
ALTER TABLE public.business_locations
ADD CONSTRAINT unique_main_location_per_business
EXCLUDE USING btree (business_id WITH =) 
WHERE (is_main = true);

-- ========================================
-- Business settings table
-- ========================================
CREATE TABLE public.business_settings (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    business_id UUID NOT NULL UNIQUE,
    calendar_start_hour INTEGER NOT NULL DEFAULT 9,
    calendar_end_hour INTEGER NOT NULL DEFAULT 18,
    appointment_buffer_minutes INTEGER NOT NULL DEFAULT 0,
    allow_online_booking BOOLEAN NOT NULL DEFAULT TRUE,
    default_appointment_duration INTEGER NOT NULL DEFAULT 60,
    currency VARCHAR(3) NOT NULL DEFAULT 'EUR',
    date_format VARCHAR(20) NOT NULL DEFAULT 'DD-MM-YYYY',
    time_format VARCHAR(10) NOT NULL DEFAULT '24h',
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    created_by UUID NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_by UUID,
    deleted_at TIMESTAMP WITH TIME ZONE,
    deleted_by UUID,
    CONSTRAINT fk_business_settings_business FOREIGN KEY (business_id) REFERENCES public.businesses(id) ON DELETE CASCADE,
    CONSTRAINT fk_business_settings_created_by FOREIGN KEY (created_by) REFERENCES public.users(id),
    CONSTRAINT fk_business_settings_updated_by FOREIGN KEY (updated_by) REFERENCES public.users(id),
    CONSTRAINT fk_business_settings_deleted_by FOREIGN KEY (deleted_by) REFERENCES public.users(id),
    CONSTRAINT check_calendar_hours CHECK (calendar_start_hour >= 0 AND calendar_start_hour < 24 AND calendar_end_hour > 0 AND calendar_end_hour <= 24 AND calendar_start_hour < calendar_end_hour)
);

CREATE INDEX idx_business_settings_created_by ON public.business_settings(created_by);
CREATE INDEX idx_business_settings_updated_by ON public.business_settings(updated_by);
CREATE INDEX idx_business_settings_deleted_by ON public.business_settings(deleted_by);
CREATE INDEX idx_business_settings_deleted_at ON public.business_settings(deleted_at) WHERE deleted_at IS NULL;

-- ========================================
-- Service categories table
-- ========================================
CREATE TABLE public.service_categories (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    business_id UUID NOT NULL,
    name VARCHAR(50) NOT NULL,
    description TEXT,
    display_order INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    created_by UUID NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE,
    updated_by UUID,
    deleted_at TIMESTAMP WITH TIME ZONE,
    deleted_by UUID,
    CONSTRAINT fk_service_categories_business FOREIGN KEY (business_id) REFERENCES public.businesses(id),
    CONSTRAINT fk_service_categories_created_by FOREIGN KEY (created_by) REFERENCES public.users(id),
    CONSTRAINT fk_service_categories_updated_by FOREIGN KEY (updated_by) REFERENCES public.users(id),
    CONSTRAINT fk_service_categories_deleted_by FOREIGN KEY (deleted_by) REFERENCES public.users(id),
    CONSTRAINT uq_category_business_name UNIQUE (business_id, name)
);

CREATE INDEX idx_service_categories_created_by ON public.service_categories(created_by);
CREATE INDEX idx_service_categories_updated_by ON public.service_categories(updated_by);
CREATE INDEX idx_service_categories_deleted_by ON public.service_categories(deleted_by);
CREATE INDEX idx_service_categories_deleted_at ON public.service_categories(deleted_at) WHERE deleted_at IS NULL;

-- ========================================
-- Services table
-- ========================================
CREATE TABLE public.services (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    business_id UUID NOT NULL,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    duration INTEGER NOT NULL, -- in minutes
    price DECIMAL(10,2) NOT NULL,
    category VARCHAR(50),
    color VARCHAR(7), -- Hex color code e.g., #FFFFFF
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    preparation_time INTEGER NOT NULL DEFAULT 0, -- in minutes
    cleanup_time INTEGER NOT NULL DEFAULT 0, -- in minutes
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    created_by UUID NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE,
    updated_by UUID,
    deleted_at TIMESTAMP WITH TIME ZONE,
    deleted_by UUID,
    CONSTRAINT fk_services_business FOREIGN KEY (business_id) REFERENCES public.businesses(id),
    CONSTRAINT fk_services_created_by FOREIGN KEY (created_by) REFERENCES public.users(id),
    CONSTRAINT fk_services_updated_by FOREIGN KEY (updated_by) REFERENCES public.users(id),
    CONSTRAINT fk_services_deleted_by FOREIGN KEY (deleted_by) REFERENCES public.users(id)
);

COMMENT ON TABLE public.services IS 'Services offered by businesses with pricing, duration, and categorization';

-- Create indexes for services table
CREATE INDEX idx_services_business_id ON public.services(business_id);
CREATE INDEX idx_services_category ON public.services(category);
CREATE INDEX idx_services_is_active ON public.services(is_active) WHERE is_active = TRUE;

-- ========================================
-- Staff table
-- ========================================
CREATE TABLE public.staff (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    business_id UUID NOT NULL,
    user_id UUID NOT NULL,
    position VARCHAR(100) NOT NULL,
    bio TEXT,
    specialty_areas TEXT[] DEFAULT '{}',
    profile_image_url TEXT,
    working_hours JSONB,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    employment_type VARCHAR(50) NOT NULL,
    join_date TIMESTAMP WITH TIME ZONE NOT NULL,
    end_date TIMESTAMP WITH TIME ZONE,
    commission_rate DECIMAL(5,2),
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    created_by UUID NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE,
    updated_by UUID,
    deleted_at TIMESTAMP WITH TIME ZONE,
    deleted_by UUID,
    CONSTRAINT fk_staff_business FOREIGN KEY (business_id) REFERENCES public.businesses(id),
    CONSTRAINT fk_staff_user FOREIGN KEY (user_id) REFERENCES public.users(id),
    CONSTRAINT fk_staff_created_by FOREIGN KEY (created_by) REFERENCES public.users(id),
    CONSTRAINT fk_staff_updated_by FOREIGN KEY (updated_by) REFERENCES public.users(id),
    CONSTRAINT fk_staff_deleted_by FOREIGN KEY (deleted_by) REFERENCES public.users(id)
);

COMMENT ON TABLE public.staff IS 'Staff members who work at businesses and provide services to clients';

-- Create indexes for staff table
CREATE INDEX idx_staff_business_id ON public.staff(business_id);
CREATE INDEX idx_staff_user_id ON public.staff(user_id);
CREATE INDEX idx_staff_is_active ON public.staff(is_active) WHERE is_active = TRUE;
CREATE INDEX idx_staff_employment_type ON public.staff(employment_type);

-- ========================================
-- Service assignment table
-- ========================================
CREATE TABLE public.service_assignment (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    business_id UUID NOT NULL,
    staff_id UUID NOT NULL,
    service_id UUID NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    created_by UUID NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE,
    updated_by UUID,
    deleted_at TIMESTAMP WITH TIME ZONE,
    deleted_by UUID,
    CONSTRAINT fk_service_assignment_business FOREIGN KEY (business_id) REFERENCES public.businesses(id),
    CONSTRAINT fk_service_assignment_staff FOREIGN KEY (staff_id) REFERENCES public.staff(id),
    CONSTRAINT fk_service_assignment_service FOREIGN KEY (service_id) REFERENCES public.services(id),
    CONSTRAINT fk_service_assignment_created_by FOREIGN KEY (created_by) REFERENCES public.users(id),
    CONSTRAINT fk_service_assignment_updated_by FOREIGN KEY (updated_by) REFERENCES public.users(id),
    CONSTRAINT fk_service_assignment_deleted_by FOREIGN KEY (deleted_by) REFERENCES public.users(id),
    CONSTRAINT uq_staff_service UNIQUE (staff_id, service_id)
);

COMMENT ON TABLE public.service_assignment IS 'Assignments of services to staff members, indicating which staff can perform which services';

-- Create indexes for service_assignment table
CREATE INDEX idx_service_assignment_business_id ON public.service_assignment(business_id);
CREATE INDEX idx_service_assignment_staff_id ON public.service_assignment(staff_id);
CREATE INDEX idx_service_assignment_service_id ON public.service_assignment(service_id);
CREATE INDEX idx_service_assignment_is_active ON public.service_assignment(is_active) WHERE is_active = TRUE;

-- ========================================
-- Availability exception table
-- ========================================
CREATE TABLE public.availability_exception (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    business_id UUID NOT NULL,
    staff_id UUID NOT NULL,
    exception_type VARCHAR(50) NOT NULL, -- 'time_off', 'holiday', 'custom_hours'
    start_time TIMESTAMP WITH TIME ZONE NOT NULL,
    end_time TIMESTAMP WITH TIME ZONE NOT NULL,
    is_full_day BOOLEAN NOT NULL DEFAULT FALSE,
    is_recurring BOOLEAN NOT NULL DEFAULT FALSE,
    recurrence_rule TEXT, -- iCalendar RRULE format
    notes TEXT,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    created_by UUID NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE,
    updated_by UUID,
    deleted_at TIMESTAMP WITH TIME ZONE,
    deleted_by UUID,
    CONSTRAINT fk_availability_exception_business FOREIGN KEY (business_id) REFERENCES public.businesses(id),
    CONSTRAINT fk_availability_exception_staff FOREIGN KEY (staff_id) REFERENCES public.staff(id),
    CONSTRAINT fk_availability_exception_created_by FOREIGN KEY (created_by) REFERENCES public.users(id),
    CONSTRAINT fk_availability_exception_updated_by FOREIGN KEY (updated_by) REFERENCES public.users(id),
    CONSTRAINT fk_availability_exception_deleted_by FOREIGN KEY (deleted_by) REFERENCES public.users(id),
    CONSTRAINT check_dates CHECK (start_time <= end_time)
);

COMMENT ON TABLE public.availability_exception IS 'Exceptions to staff regular working hours, such as time off, holidays, or custom working hours';

-- Create indexes for availability_exception table
CREATE INDEX idx_availability_exception_business_id ON public.availability_exception(business_id);
CREATE INDEX idx_availability_exception_staff_id ON public.availability_exception(staff_id);
CREATE INDEX idx_availability_exception_date_range ON public.availability_exception(start_time, end_time);
CREATE INDEX idx_availability_exception_type ON public.availability_exception(exception_type);

-- ========================================
-- Staff performance table
-- ========================================
CREATE TABLE public.staff_performance (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    business_id UUID NOT NULL,
    staff_id UUID NOT NULL,
    period VARCHAR(20) NOT NULL, -- 'daily', 'weekly', 'monthly', 'yearly'
    start_date TIMESTAMP WITH TIME ZONE NOT NULL,
    end_date TIMESTAMP WITH TIME ZONE NOT NULL,
    total_appointments INTEGER NOT NULL DEFAULT 0,
    completed_appointments INTEGER NOT NULL DEFAULT 0,
    canceled_appointments INTEGER NOT NULL DEFAULT 0,
    no_show_appointments INTEGER NOT NULL DEFAULT 0,
    total_revenue DECIMAL(10,2) NOT NULL DEFAULT 0,
    average_rating DECIMAL(3,2),
    client_retention_rate DECIMAL(5,2),
    new_clients INTEGER NOT NULL DEFAULT 0,
    return_clients INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    created_by UUID NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_by UUID,
    deleted_at TIMESTAMP WITH TIME ZONE,
    deleted_by UUID,
    CONSTRAINT fk_staff_performance_business FOREIGN KEY (business_id) REFERENCES public.businesses(id),
    CONSTRAINT fk_staff_performance_staff FOREIGN KEY (staff_id) REFERENCES public.staff(id),
    CONSTRAINT fk_staff_performance_created_by FOREIGN KEY (created_by) REFERENCES public.users(id),
    CONSTRAINT fk_staff_performance_updated_by FOREIGN KEY (updated_by) REFERENCES public.users(id),
    CONSTRAINT fk_staff_performance_deleted_by FOREIGN KEY (deleted_by) REFERENCES public.users(id),
    CONSTRAINT uq_staff_performance_period UNIQUE (staff_id, period, start_date),
    CONSTRAINT check_performance_dates CHECK (start_date <= end_date)
);

COMMENT ON TABLE public.staff_performance IS 'Performance metrics for staff members, including appointment statistics, revenue, and client retention';

-- Create indexes for staff_performance table
CREATE INDEX idx_staff_performance_business_id ON public.staff_performance(business_id);
CREATE INDEX idx_staff_performance_staff_id ON public.staff_performance(staff_id);
CREATE INDEX idx_staff_performance_period ON public.staff_performance(period, start_date);
CREATE INDEX idx_staff_performance_date_range ON public.staff_performance(start_date, end_date);
CREATE INDEX idx_staff_performance_created_by ON public.staff_performance(created_by);
CREATE INDEX idx_staff_performance_updated_by ON public.staff_performance(updated_by);
CREATE INDEX idx_staff_performance_deleted_by ON public.staff_performance(deleted_by);
CREATE INDEX idx_staff_performance_deleted_at ON public.staff_performance(deleted_at) WHERE deleted_at IS NULL;

-- ========================================
-- Clients table
-- ========================================
CREATE TABLE public.clients (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    business_id UUID NOT NULL,
    user_id UUID,
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    email VARCHAR(255),
    phone VARCHAR(50),
    date_of_birth DATE,
    address_line1 VARCHAR(255),
    city VARCHAR(100),
    postal_code VARCHAR(20),
    country VARCHAR(50),
    notes TEXT,
    allergies TEXT,
    health_conditions TEXT,
    referral_source VARCHAR(100),
    accepts_marketing BOOLEAN NOT NULL DEFAULT FALSE,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    created_by UUID NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE,
    updated_by UUID,
    deleted_at TIMESTAMP WITH TIME ZONE,
    deleted_by UUID,
    CONSTRAINT fk_clients_business FOREIGN KEY (business_id) REFERENCES public.businesses(id),
    CONSTRAINT fk_clients_user FOREIGN KEY (user_id) REFERENCES public.users(id),
    CONSTRAINT fk_clients_created_by FOREIGN KEY (created_by) REFERENCES public.users(id),
    CONSTRAINT fk_clients_updated_by FOREIGN KEY (updated_by) REFERENCES public.users(id),
    CONSTRAINT fk_clients_deleted_by FOREIGN KEY (deleted_by) REFERENCES public.users(id)
);

COMMENT ON TABLE public.clients IS 'Client profiles with contact information, preferences, and health details';

-- Create indexes for clients table
CREATE INDEX idx_clients_business_id ON public.clients(business_id);
CREATE INDEX idx_clients_user_id ON public.clients(user_id);
CREATE INDEX idx_clients_email ON public.clients(email);
CREATE INDEX idx_clients_phone ON public.clients(phone);
CREATE INDEX idx_clients_is_active ON public.clients(is_active) WHERE is_active = TRUE;

-- ========================================
-- Appointments table
-- ========================================
CREATE TABLE public.appointments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    business_id UUID NOT NULL,
    client_id UUID NOT NULL,
    staff_id UUID NOT NULL,
    start_time TIMESTAMP WITH TIME ZONE NOT NULL,
    end_time TIMESTAMP WITH TIME ZONE NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'scheduled', -- 'scheduled', 'confirmed', 'completed', 'cancelled', 'no-show'
    notes TEXT,
    estimated_price DECIMAL(10,2),
    actual_price DECIMAL(10,2),
    payment_method VARCHAR(20), -- 'cash', 'card', 'transfer', 'other'
    payment_status VARCHAR(20) DEFAULT 'pending', -- 'pending', 'paid', 'partial', 'refunded'
    client_confirmed BOOLEAN DEFAULT FALSE,
    staff_confirmed BOOLEAN DEFAULT FALSE,
    cancellation_reason TEXT,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    created_by UUID NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE,
    updated_by UUID,
    deleted_at TIMESTAMP WITH TIME ZONE,
    deleted_by UUID,
    CONSTRAINT fk_appointments_business FOREIGN KEY (business_id) REFERENCES public.businesses(id),
    CONSTRAINT fk_appointments_client FOREIGN KEY (client_id) REFERENCES public.clients(id),
    CONSTRAINT fk_appointments_staff FOREIGN KEY (staff_id) REFERENCES public.staff(id),
    CONSTRAINT fk_appointments_created_by FOREIGN KEY (created_by) REFERENCES public.users(id),
    CONSTRAINT fk_appointments_updated_by FOREIGN KEY (updated_by) REFERENCES public.users(id),
    CONSTRAINT fk_appointments_deleted_by FOREIGN KEY (deleted_by) REFERENCES public.users(id),
    CONSTRAINT check_appointment_times CHECK (start_time < end_time),
    CONSTRAINT chk_appointment_time_order CHECK (end_time > start_time)
);

COMMENT ON TABLE public.appointments IS 'Appointments scheduled between clients and staff for specific services';

-- Create indexes for appointments table
CREATE INDEX idx_appointments_business_id ON public.appointments(business_id);
CREATE INDEX idx_appointments_client_id ON public.appointments(client_id);
CREATE INDEX idx_appointments_staff_id ON public.appointments(staff_id);
CREATE INDEX idx_appointments_status ON public.appointments(status);
CREATE INDEX idx_appointments_start_time ON public.appointments(start_time);
CREATE INDEX idx_appointments_date_range ON public.appointments(start_time, end_time);
CREATE INDEX idx_appointments_payment_status ON public.appointments(payment_status);

-- ========================================
-- Appointment services table
-- ========================================
CREATE TABLE public.appointment_services (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    appointment_id UUID NOT NULL,
    service_id UUID NOT NULL,
    staff_id UUID NOT NULL,
    duration INTEGER NOT NULL, -- in minutes
    price DECIMAL(10,2) NOT NULL,
    notes TEXT,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    created_by UUID NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE,
    updated_by UUID,
    deleted_at TIMESTAMP WITH TIME ZONE,
    deleted_by UUID,
    CONSTRAINT fk_appointment_services_appointment FOREIGN KEY (appointment_id) REFERENCES public.appointments(id) ON DELETE CASCADE,
    CONSTRAINT fk_appointment_services_service FOREIGN KEY (service_id) REFERENCES public.services(id),
    CONSTRAINT fk_appointment_services_staff FOREIGN KEY (staff_id) REFERENCES public.staff(id),
    CONSTRAINT fk_appointment_services_created_by FOREIGN KEY (created_by) REFERENCES public.users(id),
    CONSTRAINT fk_appointment_services_updated_by FOREIGN KEY (updated_by) REFERENCES public.users(id),
    CONSTRAINT fk_appointment_services_deleted_by FOREIGN KEY (deleted_by) REFERENCES public.users(id)
);

COMMENT ON TABLE public.appointment_services IS 'Additional services included in appointments for bundled services';

-- Create indexes for appointment_services table
CREATE INDEX idx_appointment_services_appointment_id ON public.appointment_services(appointment_id);
CREATE INDEX idx_appointment_services_service_id ON public.appointment_services(service_id);
CREATE INDEX idx_appointment_services_staff_id ON public.appointment_services(staff_id);
CREATE INDEX idx_appointment_services_created_by ON public.appointment_services(created_by);
CREATE INDEX idx_appointment_services_updated_by ON public.appointment_services(updated_by);
CREATE INDEX idx_appointment_services_deleted_by ON public.appointment_services(deleted_by);
CREATE INDEX idx_appointment_services_deleted_at ON public.appointment_services(deleted_at) WHERE deleted_at IS NULL;

-- ========================================
-- Appointment notes table
-- ========================================
CREATE TABLE public.appointment_notes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    appointment_id UUID NOT NULL,
    note_text TEXT NOT NULL,
    is_private BOOLEAN NOT NULL DEFAULT FALSE,
    products_used JSONB,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    created_by UUID NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE,
    updated_by UUID,
    deleted_at TIMESTAMP WITH TIME ZONE,
    deleted_by UUID,
    CONSTRAINT fk_appointment_notes_appointment FOREIGN KEY (appointment_id) REFERENCES public.appointments(id) ON DELETE CASCADE,
    CONSTRAINT fk_appointment_notes_created_by FOREIGN KEY (created_by) REFERENCES public.users(id),
    CONSTRAINT fk_appointment_notes_updated_by FOREIGN KEY (updated_by) REFERENCES public.users(id),
    CONSTRAINT fk_appointment_notes_deleted_by FOREIGN KEY (deleted_by) REFERENCES public.users(id)
);

COMMENT ON TABLE public.appointment_notes IS 'Treatment notes and records for appointments';

-- Create indexes for appointment_notes table
CREATE INDEX idx_appointment_notes_appointment_id ON public.appointment_notes(appointment_id);
CREATE INDEX idx_appointment_notes_created_by ON public.appointment_notes(created_by);
CREATE INDEX idx_appointment_notes_deleted_at ON public.appointment_notes(deleted_at) WHERE deleted_at IS NULL;

-- ========================================
-- Appointment reminders table
-- ========================================
CREATE TABLE public.appointment_reminders (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    appointment_id UUID NOT NULL,
    reminder_type VARCHAR(20) NOT NULL, -- 'email', 'sms', 'push', 'whatsapp'
    scheduled_time TIMESTAMP WITH TIME ZONE NOT NULL,
    sent_time TIMESTAMP WITH TIME ZONE,
    status VARCHAR(20) NOT NULL DEFAULT 'pending', -- 'pending', 'sent', 'failed', 'cancelled'
    error_message TEXT,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    created_by UUID NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE,
    updated_by UUID,
    deleted_at TIMESTAMP WITH TIME ZONE,
    deleted_by UUID,
    CONSTRAINT fk_appointment_reminders_appointment FOREIGN KEY (appointment_id) REFERENCES public.appointments(id) ON DELETE CASCADE,
    CONSTRAINT fk_appointment_reminders_created_by FOREIGN KEY (created_by) REFERENCES public.users(id),
    CONSTRAINT fk_appointment_reminders_updated_by FOREIGN KEY (updated_by) REFERENCES public.users(id),
    CONSTRAINT fk_appointment_reminders_deleted_by FOREIGN KEY (deleted_by) REFERENCES public.users(id)
);

COMMENT ON TABLE public.appointment_reminders IS 'Scheduled and sent reminders for upcoming appointments';

-- Create indexes for appointment_reminders table
CREATE INDEX idx_appointment_reminders_appointment_id ON public.appointment_reminders(appointment_id);
CREATE INDEX idx_appointment_reminders_scheduled_time ON public.appointment_reminders(scheduled_time);
CREATE INDEX idx_appointment_reminders_status ON public.appointment_reminders(status);
CREATE INDEX idx_appointment_reminders_type ON public.appointment_reminders(reminder_type);
CREATE INDEX idx_appointment_reminders_created_by ON public.appointment_reminders(created_by);
CREATE INDEX idx_appointment_reminders_updated_by ON public.appointment_reminders(updated_by);
CREATE INDEX idx_appointment_reminders_deleted_by ON public.appointment_reminders(deleted_by);
CREATE INDEX idx_appointment_reminders_deleted_at ON public.appointment_reminders(deleted_at) WHERE deleted_at IS NULL;

-- ========================================
-- Service completions table
-- ========================================
CREATE TABLE public.service_completions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    appointment_id UUID NOT NULL,
    price_charged DECIMAL(10,2) NOT NULL,
    payment_method VARCHAR(20) NOT NULL, -- 'cash', 'card', 'transfer', 'other'
    provider_confirmed BOOLEAN NOT NULL DEFAULT FALSE,
    client_confirmed BOOLEAN NOT NULL DEFAULT FALSE,
    completion_date TIMESTAMP WITH TIME ZONE,
    actual_duration INTEGER, -- in minutes
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    created_by UUID,
    updated_at TIMESTAMP WITH TIME ZONE,
    updated_by UUID,
    deleted_at TIMESTAMP WITH TIME ZONE,
    deleted_by UUID,
    CONSTRAINT fk_service_completions_appointment FOREIGN KEY (appointment_id) REFERENCES public.appointments(id) ON DELETE CASCADE,
    CONSTRAINT fk_service_completions_created_by FOREIGN KEY (created_by) REFERENCES public.users(id),
    CONSTRAINT fk_service_completions_updated_by FOREIGN KEY (updated_by) REFERENCES public.users(id),
    CONSTRAINT fk_service_completions_deleted_by FOREIGN KEY (deleted_by) REFERENCES public.users(id)
);

COMMENT ON TABLE public.service_completions IS 'Records of completed services with financial tracking information';

-- Create indexes for service_completions table
CREATE INDEX idx_service_completions_appointment_id ON public.service_completions(appointment_id);
CREATE INDEX idx_service_completions_payment_method ON public.service_completions(payment_method);
CREATE INDEX idx_service_completions_completion_date ON public.service_completions(completion_date);

-- ========================================
-- Service ratings table
-- ========================================
CREATE TABLE public.service_ratings (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    appointment_id UUID NOT NULL,
    client_id UUID NOT NULL,
    staff_id UUID NOT NULL,
    service_id UUID NOT NULL,
    rating INTEGER NOT NULL CHECK (rating BETWEEN 1 AND 5),
    feedback TEXT,
    is_published BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    created_by UUID NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE,
    updated_by UUID,
    deleted_at TIMESTAMP WITH TIME ZONE,
    deleted_by UUID,
    CONSTRAINT fk_service_ratings_appointment FOREIGN KEY (appointment_id) REFERENCES public.appointments(id) ON DELETE CASCADE,
    CONSTRAINT fk_service_ratings_client FOREIGN KEY (client_id) REFERENCES public.clients(id),
    CONSTRAINT fk_service_ratings_staff FOREIGN KEY (staff_id) REFERENCES public.staff(id),
    CONSTRAINT fk_service_ratings_service FOREIGN KEY (service_id) REFERENCES public.services(id),
    CONSTRAINT fk_service_ratings_created_by FOREIGN KEY (created_by) REFERENCES public.users(id),
    CONSTRAINT fk_service_ratings_updated_by FOREIGN KEY (updated_by) REFERENCES public.users(id),
    CONSTRAINT fk_service_ratings_deleted_by FOREIGN KEY (deleted_by) REFERENCES public.users(id),
    CONSTRAINT uq_rating_appointment UNIQUE (appointment_id),
    CONSTRAINT chk_service_rating_range CHECK (rating >= 1 AND rating <= 5)
);

COMMENT ON TABLE public.service_ratings IS 'Client ratings and feedback for services received';

-- Create indexes for service_ratings table
CREATE INDEX idx_service_ratings_appointment_id ON public.service_ratings(appointment_id);
CREATE INDEX idx_service_ratings_client_id ON public.service_ratings(client_id);
CREATE INDEX idx_service_ratings_staff_id ON public.service_ratings(staff_id);
CREATE INDEX idx_service_ratings_service_id ON public.service_ratings(service_id);
CREATE INDEX idx_service_ratings_rating ON public.service_ratings(rating);
CREATE INDEX idx_service_ratings_created_by ON public.service_ratings(created_by);
CREATE INDEX idx_service_ratings_updated_by ON public.service_ratings(updated_by);
CREATE INDEX idx_service_ratings_deleted_by ON public.service_ratings(deleted_by);
CREATE INDEX idx_service_ratings_deleted_at ON public.service_ratings(deleted_at) WHERE deleted_at IS NULL;

-- ========================================
-- Waiting list table
-- ========================================
CREATE TABLE public.waiting_list (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    business_id UUID NOT NULL,
    client_id UUID NOT NULL,
    service_id UUID NOT NULL,
    staff_id UUID,
    preferred_date DATE,
    preferred_time_range JSONB, -- Start and end time preferences
    flexibility_level VARCHAR(20) NOT NULL DEFAULT 'medium', -- 'low', 'medium', 'high'
    status VARCHAR(20) NOT NULL DEFAULT 'waiting', -- 'waiting', 'contacted', 'scheduled', 'expired', 'cancelled'
    notes TEXT,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    created_by UUID NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE,
    updated_by UUID,
    deleted_at TIMESTAMP WITH TIME ZONE,
    deleted_by UUID,
    CONSTRAINT fk_waiting_list_business FOREIGN KEY (business_id) REFERENCES public.businesses(id),
    CONSTRAINT fk_waiting_list_client FOREIGN KEY (client_id) REFERENCES public.clients(id),
    CONSTRAINT fk_waiting_list_service FOREIGN KEY (service_id) REFERENCES public.services(id),
    CONSTRAINT fk_waiting_list_staff FOREIGN KEY (staff_id) REFERENCES public.staff(id),
    CONSTRAINT fk_waiting_list_created_by FOREIGN KEY (created_by) REFERENCES public.users(id),
    CONSTRAINT fk_waiting_list_updated_by FOREIGN KEY (updated_by) REFERENCES public.users(id),
    CONSTRAINT fk_waiting_list_deleted_by FOREIGN KEY (deleted_by) REFERENCES public.users(id)
);

COMMENT ON TABLE public.waiting_list IS 'Waiting list for clients requesting appointments when no slots are available';

-- Create indexes for waiting_list table
CREATE INDEX idx_waiting_list_business_id ON public.waiting_list(business_id);
CREATE INDEX idx_waiting_list_client_id ON public.waiting_list(client_id);
CREATE INDEX idx_waiting_list_service_id ON public.waiting_list(service_id);
CREATE INDEX idx_waiting_list_staff_id ON public.waiting_list(staff_id);
CREATE INDEX idx_waiting_list_preferred_date ON public.waiting_list(preferred_date);
CREATE INDEX idx_waiting_list_status ON public.waiting_list(status);

-- ========================================
-- Loyalty programs table
-- ========================================
CREATE TABLE public.loyalty_programs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    business_id UUID NOT NULL,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    program_type VARCHAR(50) NOT NULL, -- 'visit', 'spend', 'service', 'tier'
    rules JSONB NOT NULL, -- Configuration rules based on program type
    reward_type VARCHAR(50) NOT NULL, -- 'percentage', 'fixed', 'free_service', 'upgrade', 'product'
    reward_value JSONB NOT NULL, -- Details of the reward
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    start_date TIMESTAMP WITH TIME ZONE,
    end_date TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    created_by UUID NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE,
    updated_by UUID,
    deleted_at TIMESTAMP WITH TIME ZONE,
    deleted_by UUID,
    CONSTRAINT fk_loyalty_programs_business FOREIGN KEY (business_id) REFERENCES public.businesses(id),
    CONSTRAINT fk_loyalty_programs_created_by FOREIGN KEY (created_by) REFERENCES public.users(id),
    CONSTRAINT fk_loyalty_programs_updated_by FOREIGN KEY (updated_by) REFERENCES public.users(id),
    CONSTRAINT fk_loyalty_programs_deleted_by FOREIGN KEY (deleted_by) REFERENCES public.users(id)
);

COMMENT ON TABLE public.loyalty_programs IS 'Loyalty program configurations created by businesses';

-- Create indexes for loyalty_programs table
CREATE INDEX idx_loyalty_programs_business_id ON public.loyalty_programs(business_id);
CREATE INDEX idx_loyalty_programs_program_type ON public.loyalty_programs(program_type);
CREATE INDEX idx_loyalty_programs_is_active ON public.loyalty_programs(is_active) WHERE is_active = TRUE;
CREATE INDEX idx_loyalty_programs_date_range ON public.loyalty_programs(start_date, end_date);

-- ========================================
-- Client loyalty memberships table
-- ========================================
CREATE TABLE public.client_loyalty_memberships (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    program_id UUID NOT NULL,
    client_id UUID NOT NULL,
    current_points INTEGER NOT NULL DEFAULT 0,
    visits_count INTEGER NOT NULL DEFAULT 0,
    total_spent DECIMAL(10,2) NOT NULL DEFAULT 0,
    tier_level VARCHAR(20),
    progress JSONB, -- Progress data specific to program type
    join_date TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    expiry_date TIMESTAMP WITH TIME ZONE,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    created_by UUID NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE,
    updated_by UUID,
    deleted_at TIMESTAMP WITH TIME ZONE,
    deleted_by UUID,
    CONSTRAINT fk_client_loyalty_memberships_program FOREIGN KEY (program_id) REFERENCES public.loyalty_programs(id) ON DELETE CASCADE,
    CONSTRAINT fk_client_loyalty_memberships_client FOREIGN KEY (client_id) REFERENCES public.clients(id) ON DELETE CASCADE,
    CONSTRAINT fk_client_loyalty_memberships_created_by FOREIGN KEY (created_by) REFERENCES public.users(id),
    CONSTRAINT fk_client_loyalty_memberships_updated_by FOREIGN KEY (updated_by) REFERENCES public.users(id),
    CONSTRAINT fk_client_loyalty_memberships_deleted_by FOREIGN KEY (deleted_by) REFERENCES public.users(id),
    CONSTRAINT uq_client_program UNIQUE (client_id, program_id),
    CONSTRAINT chk_loyalty_points_non_negative CHECK (current_points >= 0)
);

COMMENT ON TABLE public.client_loyalty_memberships IS 'Client membership in loyalty programs with progress tracking';

-- Create indexes for client_loyalty_memberships table
CREATE INDEX idx_client_loyalty_memberships_program_id ON public.client_loyalty_memberships(program_id);
CREATE INDEX idx_client_loyalty_memberships_client_id ON public.client_loyalty_memberships(client_id);
CREATE INDEX idx_client_loyalty_memberships_is_active ON public.client_loyalty_memberships(is_active) WHERE is_active = TRUE;
CREATE INDEX idx_client_loyalty_memberships_created_by ON public.client_loyalty_memberships(created_by);
CREATE INDEX idx_client_loyalty_memberships_updated_by ON public.client_loyalty_memberships(updated_by);
CREATE INDEX idx_client_loyalty_memberships_deleted_by ON public.client_loyalty_memberships(deleted_by);
CREATE INDEX idx_client_loyalty_memberships_deleted_at ON public.client_loyalty_memberships(deleted_at) WHERE deleted_at IS NULL;

-- ========================================
-- Loyalty transactions table
-- ========================================
CREATE TABLE public.loyalty_transactions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    membership_id UUID NOT NULL,
    appointment_id UUID,
    transaction_type VARCHAR(20) NOT NULL, -- 'earn', 'redeem', 'adjust', 'expire'
    points INTEGER NOT NULL,
    description TEXT,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    created_by UUID NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE,
    updated_by UUID,
    deleted_at TIMESTAMP WITH TIME ZONE,
    deleted_by UUID,
    CONSTRAINT fk_loyalty_transactions_membership FOREIGN KEY (membership_id) REFERENCES public.client_loyalty_memberships(id) ON DELETE CASCADE,
    CONSTRAINT fk_loyalty_transactions_appointment FOREIGN KEY (appointment_id) REFERENCES public.appointments(id) ON DELETE SET NULL,
    CONSTRAINT fk_loyalty_transactions_created_by FOREIGN KEY (created_by) REFERENCES public.users(id),
    CONSTRAINT fk_loyalty_transactions_updated_by FOREIGN KEY (updated_by) REFERENCES public.users(id),
    CONSTRAINT fk_loyalty_transactions_deleted_by FOREIGN KEY (deleted_by) REFERENCES public.users(id)
);

COMMENT ON TABLE public.loyalty_transactions IS 'Transactions for earning and redeeming loyalty program rewards';

-- Create indexes for loyalty_transactions table
CREATE INDEX idx_loyalty_transactions_membership_id ON public.loyalty_transactions(membership_id);
CREATE INDEX idx_loyalty_transactions_appointment_id ON public.loyalty_transactions(appointment_id);
CREATE INDEX idx_loyalty_transactions_transaction_type ON public.loyalty_transactions(transaction_type);
CREATE INDEX idx_loyalty_transactions_created_at ON public.loyalty_transactions(created_at);

-- ========================================
-- Campaigns table
-- ========================================
CREATE TABLE public.campaigns (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    business_id UUID NOT NULL,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    campaign_type VARCHAR(50) NOT NULL, -- 'promotion', 'seasonal', 'reactivation', 'birthday'
    target_audience JSONB, -- Criteria for selecting clients
    offer_type VARCHAR(50) NOT NULL, -- 'discount', 'free_service', 'bundle', 'gift'
    offer_details JSONB NOT NULL, -- Details of the offer
    start_date TIMESTAMP WITH TIME ZONE NOT NULL,
    end_date TIMESTAMP WITH TIME ZONE NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    message_template TEXT,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    created_by UUID NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE,
    updated_by UUID,
    deleted_at TIMESTAMP WITH TIME ZONE,
    deleted_by UUID,
    CONSTRAINT fk_campaigns_business FOREIGN KEY (business_id) REFERENCES public.businesses(id),
    CONSTRAINT fk_campaigns_created_by FOREIGN KEY (created_by) REFERENCES public.users(id),
    CONSTRAINT fk_campaigns_updated_by FOREIGN KEY (updated_by) REFERENCES public.users(id),
    CONSTRAINT fk_campaigns_deleted_by FOREIGN KEY (deleted_by) REFERENCES public.users(id),
    CONSTRAINT check_campaign_dates CHECK (start_date <= end_date),
    CONSTRAINT chk_campaign_date_order CHECK (end_date > start_date)
);

COMMENT ON TABLE public.campaigns IS 'Marketing campaigns created by businesses to promote services or offers';

-- Create indexes for campaigns table
CREATE INDEX idx_campaigns_business_id ON public.campaigns(business_id);
CREATE INDEX idx_campaigns_campaign_type ON public.campaigns(campaign_type);
CREATE INDEX idx_campaigns_is_active ON public.campaigns(is_active) WHERE is_active = TRUE;
CREATE INDEX idx_campaigns_date_range ON public.campaigns(start_date, end_date);

-- ========================================
-- Campaign clients table
-- ========================================
CREATE TABLE public.campaign_clients (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    campaign_id UUID NOT NULL,
    client_id UUID NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'pending', -- 'pending', 'sent', 'opened', 'clicked', 'converted', 'unsubscribed'
    sent_at TIMESTAMP WITH TIME ZONE,
    opened_at TIMESTAMP WITH TIME ZONE,
    clicked_at TIMESTAMP WITH TIME ZONE,
    converted_at TIMESTAMP WITH TIME ZONE,
    conversion_value DECIMAL(10,2),
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    created_by UUID NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE,
    updated_by UUID,
    deleted_at TIMESTAMP WITH TIME ZONE,
    deleted_by UUID,
    CONSTRAINT fk_campaign_clients_campaign FOREIGN KEY (campaign_id) REFERENCES public.campaigns(id) ON DELETE CASCADE,
    CONSTRAINT fk_campaign_clients_client FOREIGN KEY (client_id) REFERENCES public.clients(id) ON DELETE CASCADE,
    CONSTRAINT fk_campaign_clients_created_by FOREIGN KEY (created_by) REFERENCES public.users(id),
    CONSTRAINT fk_campaign_clients_updated_by FOREIGN KEY (updated_by) REFERENCES public.users(id),
    CONSTRAINT fk_campaign_clients_deleted_by FOREIGN KEY (deleted_by) REFERENCES public.users(id),
    CONSTRAINT uq_campaign_client UNIQUE (campaign_id, client_id)
);

COMMENT ON TABLE public.campaign_clients IS 'Client targeting and response tracking for marketing campaigns';

-- Create indexes for campaign_clients table
CREATE INDEX idx_campaign_clients_campaign_id ON public.campaign_clients(campaign_id);
CREATE INDEX idx_campaign_clients_client_id ON public.campaign_clients(client_id);
CREATE INDEX idx_campaign_clients_status ON public.campaign_clients(status);
CREATE INDEX idx_campaign_clients_created_by ON public.campaign_clients(created_by);
CREATE INDEX idx_campaign_clients_updated_by ON public.campaign_clients(updated_by);
CREATE INDEX idx_campaign_clients_deleted_by ON public.campaign_clients(deleted_by);
CREATE INDEX idx_campaign_clients_deleted_at ON public.campaign_clients(deleted_at) WHERE deleted_at IS NULL;

-- ========================================
-- Campaign messages table
-- ========================================
CREATE TABLE public.campaign_messages (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    campaign_id UUID NOT NULL,
    client_id UUID NOT NULL,
    message_type VARCHAR(20) NOT NULL, -- 'email', 'sms', 'push', 'whatsapp'
    message_content TEXT NOT NULL,
    scheduled_time TIMESTAMP WITH TIME ZONE NOT NULL,
    sent_time TIMESTAMP WITH TIME ZONE,
    status VARCHAR(20) NOT NULL DEFAULT 'pending', -- 'pending', 'sent', 'failed', 'cancelled'
    error_message TEXT,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    created_by UUID NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE,
    updated_by UUID,
    deleted_at TIMESTAMP WITH TIME ZONE,
    deleted_by UUID,
    CONSTRAINT fk_campaign_messages_campaign FOREIGN KEY (campaign_id) REFERENCES public.campaigns(id) ON DELETE CASCADE,
    CONSTRAINT fk_campaign_messages_client FOREIGN KEY (client_id) REFERENCES public.clients(id) ON DELETE CASCADE,
    CONSTRAINT fk_campaign_messages_created_by FOREIGN KEY (created_by) REFERENCES public.users(id),
    CONSTRAINT fk_campaign_messages_updated_by FOREIGN KEY (updated_by) REFERENCES public.users(id),
    CONSTRAINT fk_campaign_messages_deleted_by FOREIGN KEY (deleted_by) REFERENCES public.users(id)
);

COMMENT ON TABLE public.campaign_messages IS 'Messages sent to clients as part of marketing campaigns';

-- Create indexes for campaign_messages table
CREATE INDEX idx_campaign_messages_campaign_id ON public.campaign_messages(campaign_id);
CREATE INDEX idx_campaign_messages_client_id ON public.campaign_messages(client_id);
CREATE INDEX idx_campaign_messages_scheduled_time ON public.campaign_messages(scheduled_time);
CREATE INDEX idx_campaign_messages_status ON public.campaign_messages(status);
CREATE INDEX idx_campaign_messages_type ON public.campaign_messages(message_type);
CREATE INDEX idx_campaign_messages_created_by ON public.campaign_messages(created_by);
CREATE INDEX idx_campaign_messages_updated_by ON public.campaign_messages(updated_by);
CREATE INDEX idx_campaign_messages_deleted_by ON public.campaign_messages(deleted_by);
CREATE INDEX idx_campaign_messages_deleted_at ON public.campaign_messages(deleted_at) WHERE deleted_at IS NULL;

-- ========================================
-- Inventory products table
-- ========================================
CREATE TABLE public.inventory_products (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    business_id UUID NOT NULL,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    brand VARCHAR(100),
    category VARCHAR(50),
    sku VARCHAR(50),
    barcode VARCHAR(50),
    unit VARCHAR(20) NOT NULL DEFAULT 'item',
    cost_price DECIMAL(10,2),
    retail_price DECIMAL(10,2),
    stock_quantity DECIMAL(10,2) NOT NULL DEFAULT 0,
    reorder_level DECIMAL(10,2),
    quantity_on_hand DECIMAL(10,2) NOT NULL DEFAULT 0,
    minimum_stock_level DECIMAL(10,2) NOT NULL DEFAULT 0,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    product_type VARCHAR(20) NOT NULL DEFAULT 'professional', -- 'professional', 'retail', 'both'
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    created_by UUID NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE,
    updated_by UUID,
    deleted_at TIMESTAMP WITH TIME ZONE,
    deleted_by UUID,
    CONSTRAINT fk_inventory_products_business FOREIGN KEY (business_id) REFERENCES public.businesses(id),
    CONSTRAINT fk_inventory_products_created_by FOREIGN KEY (created_by) REFERENCES public.users(id),
    CONSTRAINT fk_inventory_products_updated_by FOREIGN KEY (updated_by) REFERENCES public.users(id),
    CONSTRAINT fk_inventory_products_deleted_by FOREIGN KEY (deleted_by) REFERENCES public.users(id),
    CONSTRAINT chk_inventory_quantity_non_negative CHECK (quantity_on_hand >= 0 AND minimum_stock_level >= 0)
);

COMMENT ON TABLE public.inventory_products IS 'Professional and retail products used in services or sold to clients';

-- Create indexes for inventory_products table
CREATE INDEX idx_inventory_products_business_id ON public.inventory_products(business_id);
CREATE INDEX idx_inventory_products_category ON public.inventory_products(category);
CREATE INDEX idx_inventory_products_name ON public.inventory_products(name);
CREATE INDEX idx_inventory_products_sku ON public.inventory_products(sku);
CREATE INDEX idx_inventory_products_barcode ON public.inventory_products(barcode);
CREATE INDEX idx_inventory_products_is_active ON public.inventory_products(is_active) WHERE is_active = TRUE;
CREATE INDEX idx_inventory_products_product_type ON public.inventory_products(product_type);

-- ========================================
-- Inventory transactions table
-- ========================================
CREATE TABLE public.inventory_transactions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    business_id UUID NOT NULL,
    product_id UUID NOT NULL,
    location_id UUID,
    transaction_type VARCHAR(20) NOT NULL, -- 'purchase', 'sale', 'service_use', 'adjustment', 'return'
    quantity DECIMAL(10,2) NOT NULL,
    unit_price DECIMAL(10,2),
    appointment_id UUID,
    notes TEXT,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    created_by UUID NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE,
    updated_by UUID,
    deleted_at TIMESTAMP WITH TIME ZONE,
    deleted_by UUID,
    CONSTRAINT fk_inventory_transactions_business FOREIGN KEY (business_id) REFERENCES public.businesses(id),
    CONSTRAINT fk_inventory_transactions_product FOREIGN KEY (product_id) REFERENCES public.inventory_products(id),
    CONSTRAINT fk_inventory_transactions_appointment FOREIGN KEY (appointment_id) REFERENCES public.appointments(id) ON DELETE SET NULL,
    CONSTRAINT fk_inventory_transactions_created_by FOREIGN KEY (created_by) REFERENCES public.users(id),
    CONSTRAINT fk_inventory_transactions_updated_by FOREIGN KEY (updated_by) REFERENCES public.users(id),
    CONSTRAINT fk_inventory_transactions_deleted_by FOREIGN KEY (deleted_by) REFERENCES public.users(id)
);

COMMENT ON TABLE public.inventory_transactions IS 'Inventory movement records for tracking stock changes';

-- Create indexes for inventory_transactions table
CREATE INDEX idx_inventory_transactions_business_id ON public.inventory_transactions(business_id);
CREATE INDEX idx_inventory_transactions_product_id ON public.inventory_transactions(product_id);
CREATE INDEX idx_inventory_transactions_location_id ON public.inventory_transactions(location_id);
CREATE INDEX idx_inventory_transactions_transaction_type ON public.inventory_transactions(transaction_type);
CREATE INDEX idx_inventory_transactions_created_at ON public.inventory_transactions(created_at);
CREATE INDEX idx_inventory_transactions_appointment_id ON public.inventory_transactions(appointment_id);

-- ========================================
-- Resources table
-- ========================================
CREATE TABLE public.resources (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    business_id UUID NOT NULL,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    resource_type VARCHAR(50) NOT NULL, -- 'room', 'equipment', 'station', 'vehicle', 'other'
    capacity INTEGER,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    created_by UUID NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE,
    updated_by UUID,
    deleted_at TIMESTAMP WITH TIME ZONE,
    deleted_by UUID,
    CONSTRAINT fk_resources_business FOREIGN KEY (business_id) REFERENCES public.businesses(id),
    CONSTRAINT fk_resources_created_by FOREIGN KEY (created_by) REFERENCES public.users(id),
    CONSTRAINT fk_resources_updated_by FOREIGN KEY (updated_by) REFERENCES public.users(id),
    CONSTRAINT fk_resources_deleted_by FOREIGN KEY (deleted_by) REFERENCES public.users(id)
);

COMMENT ON TABLE public.resources IS 'Physical resources like rooms and equipment that can be booked for services';

-- Create indexes for resources table
CREATE INDEX idx_resources_business_id ON public.resources(business_id);
CREATE INDEX idx_resources_resource_type ON public.resources(resource_type);
CREATE INDEX idx_resources_is_active ON public.resources(is_active) WHERE is_active = TRUE;

-- ========================================
-- Resource availability table
-- ========================================
CREATE TABLE public.resource_availability (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    resource_id UUID NOT NULL,
    day_of_week INTEGER NOT NULL CHECK (day_of_week BETWEEN 0 AND 6), -- 0=Sunday, 6=Saturday
    start_time TIME NOT NULL,
    end_time TIME NOT NULL,
    is_available BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    created_by UUID NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE,
    updated_by UUID,
    deleted_at TIMESTAMP WITH TIME ZONE,
    deleted_by UUID,
    CONSTRAINT fk_resource_availability_resource FOREIGN KEY (resource_id) REFERENCES public.resources(id) ON DELETE CASCADE,
    CONSTRAINT fk_resource_availability_created_by FOREIGN KEY (created_by) REFERENCES public.users(id),
    CONSTRAINT fk_resource_availability_updated_by FOREIGN KEY (updated_by) REFERENCES public.users(id),
    CONSTRAINT fk_resource_availability_deleted_by FOREIGN KEY (deleted_by) REFERENCES public.users(id),
    CONSTRAINT check_resource_times CHECK (start_time < end_time)
);

COMMENT ON TABLE public.resource_availability IS 'Regular availability patterns for resources by day of week';

-- Create indexes for resource_availability table
CREATE INDEX idx_resource_availability_resource_id ON public.resource_availability(resource_id);
CREATE INDEX idx_resource_availability_day_of_week ON public.resource_availability(day_of_week);
CREATE INDEX idx_resource_availability_created_by ON public.resource_availability(created_by);
CREATE INDEX idx_resource_availability_updated_by ON public.resource_availability(updated_by);
CREATE INDEX idx_resource_availability_deleted_by ON public.resource_availability(deleted_by);
CREATE INDEX idx_resource_availability_deleted_at ON public.resource_availability(deleted_at) WHERE deleted_at IS NULL;

-- ========================================
-- Resource bookings table
-- ========================================
CREATE TABLE public.resource_bookings (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    resource_id UUID NOT NULL,
    appointment_id UUID,
    staff_id UUID,
    start_time TIMESTAMP WITH TIME ZONE NOT NULL,
    end_time TIMESTAMP WITH TIME ZONE NOT NULL,
    booking_type VARCHAR(20) NOT NULL DEFAULT 'appointment', -- 'appointment', 'maintenance', 'block', 'other'
    notes TEXT,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    created_by UUID NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE,
    updated_by UUID,
    deleted_at TIMESTAMP WITH TIME ZONE,
    deleted_by UUID,
    CONSTRAINT fk_resource_bookings_resource FOREIGN KEY (resource_id) REFERENCES public.resources(id),
    CONSTRAINT fk_resource_bookings_appointment FOREIGN KEY (appointment_id) REFERENCES public.appointments(id) ON DELETE SET NULL,
    CONSTRAINT fk_resource_bookings_staff FOREIGN KEY (staff_id) REFERENCES public.staff(id),
    CONSTRAINT fk_resource_bookings_created_by FOREIGN KEY (created_by) REFERENCES public.users(id),
    CONSTRAINT fk_resource_bookings_updated_by FOREIGN KEY (updated_by) REFERENCES public.users(id),
    CONSTRAINT fk_resource_bookings_deleted_by FOREIGN KEY (deleted_by) REFERENCES public.users(id),
    CONSTRAINT check_booking_times CHECK (start_time < end_time),
    CONSTRAINT chk_resource_booking_time_order CHECK (end_time > start_time)
);

COMMENT ON TABLE public.resource_bookings IS 'Bookings of resources for appointments or other purposes';

-- Create indexes for resource_bookings table
CREATE INDEX idx_resource_bookings_resource_id ON public.resource_bookings(resource_id);
CREATE INDEX idx_resource_bookings_appointment_id ON public.resource_bookings(appointment_id);
CREATE INDEX idx_resource_bookings_staff_id ON public.resource_bookings(staff_id);
CREATE INDEX idx_resource_bookings_time_range ON public.resource_bookings(start_time, end_time);
CREATE INDEX idx_resource_bookings_booking_type ON public.resource_bookings(booking_type);

-- ========================================
-- Providers table (for discovery)
-- ========================================
CREATE TABLE public.providers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    business_name VARCHAR(100) NOT NULL,
    description TEXT,
    address VARCHAR(255),
    city VARCHAR(100),
    postal_code VARCHAR(20),
    country VARCHAR(50) NOT NULL,
    website VARCHAR(255),
    logo_url VARCHAR(255),
    subscription_tier VARCHAR(50) DEFAULT 'basic',
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    created_by UUID,
    updated_at TIMESTAMP WITH TIME ZONE,
    updated_by UUID,
    deleted_at TIMESTAMP WITH TIME ZONE,
    deleted_by UUID,
    CONSTRAINT fk_providers_user FOREIGN KEY (user_id) REFERENCES public.users(id),
    CONSTRAINT fk_providers_created_by FOREIGN KEY (created_by) REFERENCES public.users(id),
    CONSTRAINT fk_providers_updated_by FOREIGN KEY (updated_by) REFERENCES public.users(id),
    CONSTRAINT fk_providers_deleted_by FOREIGN KEY (deleted_by) REFERENCES public.users(id)
);

COMMENT ON TABLE public.providers IS 'Beauty service providers with business information for discovery and promotion';

-- Create indexes for providers table
CREATE INDEX idx_providers_user_id ON public.providers(user_id);
CREATE INDEX idx_providers_city ON public.providers(city);
CREATE INDEX idx_providers_country ON public.providers(country);
CREATE INDEX idx_providers_subscription_tier ON public.providers(subscription_tier);

-- ========================================
-- Appointment booking rules table
-- ========================================
CREATE TABLE public.appointment_booking_rules (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    business_id UUID NOT NULL,
    buffer_time_minutes INTEGER NOT NULL DEFAULT 0, -- Time between appointments
    max_advance_booking_days INTEGER, -- How far in advance can book
    min_advance_booking_hours INTEGER, -- Minimum notice required
    allow_double_booking BOOLEAN NOT NULL DEFAULT FALSE,
    allow_overlap_different_staff BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    created_by UUID NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE,
    updated_by UUID,
    deleted_at TIMESTAMP WITH TIME ZONE,
    deleted_by UUID,
    CONSTRAINT fk_appointment_booking_rules_business FOREIGN KEY (business_id) REFERENCES public.businesses(id),
    CONSTRAINT fk_appointment_booking_rules_created_by FOREIGN KEY (created_by) REFERENCES public.users(id),
    CONSTRAINT fk_appointment_booking_rules_updated_by FOREIGN KEY (updated_by) REFERENCES public.users(id),
    CONSTRAINT fk_appointment_booking_rules_deleted_by FOREIGN KEY (deleted_by) REFERENCES public.users(id)
);

COMMENT ON TABLE public.appointment_booking_rules IS 'Business-specific rules for appointment booking logic';

CREATE INDEX idx_appointment_booking_rules_business_id ON public.appointment_booking_rules(business_id);

-- ========================================
-- Functions and Triggers
-- ========================================

-- Create function to automatically update updated_at timestamp
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Create triggers for each table with updated_at column
DO $$
DECLARE
    tables CURSOR FOR
        SELECT tablename
        FROM pg_tables
        WHERE schemaname = 'public'
        AND EXISTS (
            SELECT 1 FROM information_schema.columns 
            WHERE table_schema = 'public' 
            AND table_name = pg_tables.tablename 
            AND column_name = 'updated_at'
        );
BEGIN
    FOR table_record IN tables LOOP
        EXECUTE format('
            CREATE TRIGGER set_updated_at
            BEFORE UPDATE ON public.%I
            FOR EACH ROW
            EXECUTE FUNCTION update_updated_at_column()', 
            table_record.tablename
        );
    END LOOP;
END;
$$ LANGUAGE plpgsql;

-- Note: Appointment overlap checks are handled in the application layer

-- Create function to update inventory stock on transaction
CREATE OR REPLACE FUNCTION update_inventory_stock()
RETURNS TRIGGER AS $$
BEGIN
    -- Update product stock quantity based on transaction
    IF NEW.transaction_type = 'purchase' OR NEW.transaction_type = 'return' THEN
        -- Add to stock
        UPDATE public.inventory_products
        SET stock_quantity = stock_quantity + NEW.quantity,
            updated_at = NOW()
        WHERE id = NEW.product_id;
    ELSIF NEW.transaction_type = 'sale' OR NEW.transaction_type = 'service_use' THEN
        -- Remove from stock
        UPDATE public.inventory_products
        SET stock_quantity = stock_quantity - NEW.quantity,
            updated_at = NOW()
        WHERE id = NEW.product_id;
    ELSIF NEW.transaction_type = 'adjustment' THEN
        -- Direct adjustment
        UPDATE public.inventory_products
        SET stock_quantity = stock_quantity + NEW.quantity, -- Can be negative
            updated_at = NOW()
        WHERE id = NEW.product_id;
    END IF;
    
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Create trigger for inventory stock update
CREATE TRIGGER update_inventory_stock_trigger
AFTER INSERT ON public.inventory_transactions
FOR EACH ROW
EXECUTE FUNCTION update_inventory_stock();

-- Note: Resource booking overlap checks are handled in the application layer

-- ========================================
-- Final comments
-- ========================================
COMMENT ON SCHEMA public IS 'Schema standardized for Go/GORM compatibility - primary keys use id, foreign keys reference id columns';

-- Add missing constraint on service completions for actual_duration
ALTER TABLE public.service_completions
    ADD CONSTRAINT chk_service_completion_duration 
    CHECK (actual_duration > 0 OR actual_duration IS NULL);

-- ========================================
-- Add foreign key constraints that reference themselves
-- ========================================
ALTER TABLE public.users
    ADD CONSTRAINT fk_users_created_by FOREIGN KEY (created_by) REFERENCES public.users(id),
    ADD CONSTRAINT fk_users_updated_by FOREIGN KEY (updated_by) REFERENCES public.users(id),
    ADD CONSTRAINT fk_users_deleted_by FOREIGN KEY (deleted_by) REFERENCES public.users(id);