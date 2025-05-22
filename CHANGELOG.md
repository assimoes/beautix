# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- Database connection setup with GORM
- Domain models for business entities
- Infrastructure layer with database connection handling
- Enhanced Make commands for database management
- Data model implementation based on BeautyBiz MVP requirements
- Staff domain model implementation with related entities (AvailabilityException, ServiceAssignment, StaffPerformance)
- Unit tests for Staff models and associated entities
- Unit tests for service assignments, availability exceptions, and staff performance metrics
- Repository unit tests using mockery for Staff, ServiceAssignment, AvailabilityException, and StaffPerformance
- Repository implementations using GORM for Staff, ServiceAssignment, AvailabilityException, and StaffPerformance
- Integration tests for Staff, ServiceAssignment, AvailabilityException, and StaffPerformance repositories
- Transaction-based testing framework for truly idempotent integration tests
- Repository implementations for Client, Service, ServiceCategory, and Appointment with transaction-based testing
- Extended test helpers to support new repository implementations
- True test isolation for model tests using transaction-based approach

### Changed
- Restructured code to follow Clean Architecture principles
- Updated Makefile with improved database commands
- Modified README with updated setup instructions
- Simplified main.go to focus on database setup

### Removed
- GraphQL server implementation (will be reimplemented later)
- Mock services (replaced with actual database models)