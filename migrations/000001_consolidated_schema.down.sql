-- Drop all tables and functions in reverse order of creation

-- Drop all triggers first
DROP TRIGGER IF EXISTS update_inventory_stock_trigger ON public.inventory_transactions;

-- Drop all functions
DROP FUNCTION IF EXISTS update_inventory_stock();
DROP FUNCTION IF EXISTS update_updated_at_column();

-- Drop tables in reverse order of dependencies
DROP TABLE IF EXISTS public.appointment_booking_rules CASCADE;
DROP TABLE IF EXISTS public.providers CASCADE;
DROP TABLE IF EXISTS public.resource_bookings CASCADE;
DROP TABLE IF EXISTS public.resource_availability CASCADE;
DROP TABLE IF EXISTS public.resources CASCADE;
DROP TABLE IF EXISTS public.inventory_transactions CASCADE;
DROP TABLE IF EXISTS public.inventory_products CASCADE;
DROP TABLE IF EXISTS public.campaign_messages CASCADE;
DROP TABLE IF EXISTS public.campaign_clients CASCADE;
DROP TABLE IF EXISTS public.campaigns CASCADE;
DROP TABLE IF EXISTS public.loyalty_transactions CASCADE;
DROP TABLE IF EXISTS public.client_loyalty_memberships CASCADE;
DROP TABLE IF EXISTS public.loyalty_programs CASCADE;
DROP TABLE IF EXISTS public.waiting_list CASCADE;
DROP TABLE IF EXISTS public.service_ratings CASCADE;
DROP TABLE IF EXISTS public.service_completions CASCADE;
DROP TABLE IF EXISTS public.appointment_reminders CASCADE;
DROP TABLE IF EXISTS public.appointment_notes CASCADE;
DROP TABLE IF EXISTS public.appointment_services CASCADE;
DROP TABLE IF EXISTS public.appointments CASCADE;
DROP TABLE IF EXISTS public.clients CASCADE;
DROP TABLE IF EXISTS public.staff_performance CASCADE;
DROP TABLE IF EXISTS public.availability_exception CASCADE;
DROP TABLE IF EXISTS public.service_assignment CASCADE;
DROP TABLE IF EXISTS public.staff CASCADE;
DROP TABLE IF EXISTS public.services CASCADE;
DROP TABLE IF EXISTS public.service_categories CASCADE;
DROP TABLE IF EXISTS public.business_settings CASCADE;
DROP TABLE IF EXISTS public.business_locations CASCADE;
DROP TABLE IF EXISTS public.businesses CASCADE;
DROP TABLE IF EXISTS public.user_connected_accounts CASCADE;
DROP TABLE IF EXISTS public.users CASCADE;

-- Drop extensions
DROP EXTENSION IF EXISTS "pgcrypto";
DROP EXTENSION IF EXISTS "uuid-ossp";