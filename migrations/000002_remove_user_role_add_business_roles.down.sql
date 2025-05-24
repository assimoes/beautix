-- Rollback migration: restore user role and revert staff table changes
-- This rollback restores the previous structure with global user roles

-- ========================================
-- Restore role column to users table
-- ========================================

-- Add the role column back to users table
ALTER TABLE public.users 
    ADD COLUMN role VARCHAR(20) CHECK (role IN ('admin', 'owner', 'staff', 'user'));

-- Set default roles for existing users based on their staff roles
-- This is a best-effort migration and may need manual adjustment
UPDATE public.users 
SET role = CASE 
    WHEN id IN (
        SELECT DISTINCT s.user_id 
        FROM public.staff s 
        WHERE s.role = 'owner' AND s.is_active = true
    ) THEN 'owner'
    WHEN id IN (
        SELECT DISTINCT s.user_id 
        FROM public.staff s 
        WHERE s.role IN ('manager', 'employee', 'assistant') AND s.is_active = true
    ) THEN 'staff'
    ELSE 'user'
END;

-- Make role column NOT NULL and set default
ALTER TABLE public.users 
    ALTER COLUMN role SET DEFAULT 'user',
    ALTER COLUMN role SET NOT NULL;

-- Recreate the role index
CREATE INDEX idx_users_role ON public.users(role);

-- ========================================
-- Revert staff table structure
-- ========================================

-- Drop the new unique constraint
DROP INDEX IF EXISTS idx_staff_unique_active_position;

-- Drop new indexes
DROP INDEX IF EXISTS idx_staff_role;
DROP INDEX IF EXISTS idx_staff_start_date;
DROP INDEX IF EXISTS idx_staff_end_date;

-- Remove the new business role columns
ALTER TABLE public.staff 
    DROP COLUMN IF EXISTS role,
    DROP COLUMN IF EXISTS permissions,
    DROP COLUMN IF EXISTS start_date,
    DROP COLUMN IF EXISTS end_date;

-- Restore original table comment
COMMENT ON TABLE public.staff IS 'Staff members who work at businesses and provide services to clients';

-- ========================================
-- Clean up comments
-- ========================================

-- Remove comments for columns that no longer exist
COMMENT ON COLUMN public.staff.role IS NULL;
COMMENT ON COLUMN public.staff.permissions IS NULL;
COMMENT ON COLUMN public.staff.start_date IS NULL;
COMMENT ON COLUMN public.staff.end_date IS NULL;