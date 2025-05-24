-- Migration to remove user role and update staff table for business-context roles
-- This migration supports the architectural change from global user roles to business-context roles

-- ========================================
-- Drop existing role-based indexes and constraints
-- ========================================

-- Drop the role index from users table
DROP INDEX IF EXISTS idx_users_role;

-- ========================================
-- Update staff table structure
-- ========================================

-- First, check if the staff table exists and what structure it has
-- The existing staff table has different fields than our new business role model

-- Add new columns for business roles to the existing staff table
ALTER TABLE public.staff 
    ADD COLUMN IF NOT EXISTS role VARCHAR(20) CHECK (role IN ('owner', 'manager', 'employee', 'assistant')),
    ADD COLUMN IF NOT EXISTS permissions JSONB DEFAULT '{}',
    ADD COLUMN IF NOT EXISTS start_date TIMESTAMP WITH TIME ZONE,
    ADD COLUMN IF NOT EXISTS end_date TIMESTAMP WITH TIME ZONE;

-- Update existing data - set default role based on position if available
UPDATE public.staff 
SET role = CASE 
    WHEN LOWER(position) LIKE '%owner%' OR LOWER(position) LIKE '%manager%' THEN 'manager'
    WHEN LOWER(position) LIKE '%senior%' OR LOWER(position) LIKE '%lead%' THEN 'employee'
    ELSE 'employee'
END
WHERE role IS NULL;

-- Make role column NOT NULL after setting default values
ALTER TABLE public.staff 
    ALTER COLUMN role SET NOT NULL;

-- Set start_date to join_date if available, otherwise created_at
UPDATE public.staff 
SET start_date = COALESCE(join_date, created_at)
WHERE start_date IS NULL;

-- Create indexes for new staff role fields
CREATE INDEX idx_staff_role ON public.staff(role);
CREATE INDEX idx_staff_start_date ON public.staff(start_date);
CREATE INDEX idx_staff_end_date ON public.staff(end_date) WHERE end_date IS NOT NULL;

-- ========================================
-- Remove role column from users table
-- ========================================

-- Remove the role column from users table (this is the main change)
ALTER TABLE public.users DROP COLUMN IF EXISTS role;

-- ========================================
-- Create unique constraint for staff table
-- ========================================

-- Ensure a user can only have one active staff position per business
-- But allow multiple historical positions (with end_date set)
CREATE UNIQUE INDEX idx_staff_unique_active_position 
ON public.staff(business_id, user_id) 
WHERE is_active = true AND end_date IS NULL;

-- ========================================
-- Add comments for new structure
-- ========================================

COMMENT ON COLUMN public.staff.role IS 'Business-context role: owner, manager, employee, or assistant';
COMMENT ON COLUMN public.staff.permissions IS 'JSON object containing specific permissions for this staff member';
COMMENT ON COLUMN public.staff.start_date IS 'Date when the staff member started in this role';
COMMENT ON COLUMN public.staff.end_date IS 'Date when the staff member ended this position (null for active positions)';

-- Update table comment to reflect new purpose
COMMENT ON TABLE public.staff IS 'Staff members with business-context roles and permissions. Users can have different roles in different businesses.';