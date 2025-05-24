# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Changed
- Consolidated 16 separate migration files into a single comprehensive migration (000001_consolidated_schema)
- Simplified database migration management by creating a fresh starting point for the schema
- Removed all old migration files (000001 through 000016) to reduce complexity
- Database schema now represents the complete state after all migrations up to version 16
- Removed database-level overlap check triggers (appointment and resource booking) - these will be handled in the application layer

### Technical Details
- The consolidated migration includes:
  - All table definitions with their final structure
  - All indexes and constraints
  - Essential triggers and functions (updated_at timestamp management, inventory stock updates)
  - Complete audit fields (created_at, created_by, updated_at, updated_by, deleted_at, deleted_by) on all tables
  - Proper foreign key relationships
  - PostgreSQL extensions (uuid-ossp, pgcrypto)
  - Comments on tables and columns for documentation
  - Business rule configuration tables (appointment_booking_rules) for application-level validation

### Migration Instructions
- For new installations: Run the single migration file 000001_consolidated_schema.up.sql
- For existing installations: Ensure all migrations up to 000016 have been applied before using this consolidated version