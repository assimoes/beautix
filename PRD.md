1. Executive Summary
BeautyBiz is a comprehensive provider-centric platform designed specifically for aestheticians in Portugal. This platform empowers beauty professionals to manage and grow their businesses through integrated tools for client management, scheduling, business operations, and client acquisition.
The updated MVP focuses on critical features for independent aestheticians and small salons in Portugal, including appointment scheduling, client relationship management, service tracking, business performance monitoring, and client acquisition tools. By prioritizing provider needs while adding client-facing features, BeautyBiz creates a sustainable ecosystem that helps beauty professionals thrive.
This PRD outlines the key features, target users, monetization strategy, and implementation roadmap for the BeautyBiz MVP, with a focus on the Portuguese market's specific needs.
2. Product Vision
2.1 Vision Statement
To become the essential digital partner for beauty professionals in Portugal by providing a comprehensive business management and client acquisition platform that empowers aestheticians to build successful, sustainable practices.
2.2 Mission
BeautyBiz helps aestheticians focus on their craft by simplifying the business side of beauty through intuitive tools, client acquisition features, and performance tracking capabilities.
2.3 Key Objectives

Simplify Business Management: Reduce administrative burden through streamlined scheduling and client management
Enhance Client Relationships: Build stronger client connections through communication tools and loyalty programs
Drive Client Acquisition: Help providers attract new clients through discovery and promotion features
Track Business Performance: Provide simple tools for monitoring financial performance regardless of payment method
Ensure Regulatory Compliance: Keep providers updated on relevant regulations and standards

3. Target Users
[Target users section remains mostly unchanged - Sofia, Miguel, Carolina, and Beatriz personas are still relevant]
4. MVP Features
4.1 Core Business Management
4.1.1 Calendar & Scheduling

Smart Scheduling System

Multi-provider calendar view
Service-specific duration settings
Buffer time configuration
Color-coding by service type
Recurring appointment setup


Appointment Management

Appointment creation, modification, and cancellation
Service bundling capabilities
Client booking history access
Waiting list functionality
Conflict detection and resolution
Holiday and time-off management


Public Calendar

Client-facing calendar showing available slots
Direct appointment request capability
Service selection with estimated pricing
Appointment request notifications
Provider approval workflow



4.1.2 Client Management

Client Profiles

Comprehensive contact information
Treatment history and preferences
Allergy and contraindication tracking
Service preferences and notes
Average spend and visit frequency


Client Communication

In-app messaging system between provider and client
Appointment reminders (push notifications, SMS, email, WhatsApp)
Customizable message templates
Two-way chat functionality
Special offers and announcements


Treatment Documentation

Before/after photo storage
Treatment notes and protocols
Product usage tracking
Service completion confirmation system



4.2 Client Acquisition & Retention
4.2.1 Provider Discovery

Search & Discovery

Location-based provider search
Service type filtering
Rating and review system
Provider profile pages with services and pricing
Featured providers section


Provider Promotion

Subscription-based promotion tiers
Algorithm boosting (X times per month based on subscription)
Featured placement in search results
Special badges for premium providers
Enhanced profile visibility



4.2.2 Mini Website Builder

Provider Profile Creation

Customizable templates
Logo and branding upload
Service showcase with images
Team member profiles
Location and contact information
Gallery for work samples
Customer testimonials section



4.2.3 Loyalty & Campaigns

Loyalty Program Builder

Customizable loyalty schemes (visit-based, spending-based)
Configurable rewards (discounts, free services)
Client progress tracking
Automatic reward notifications
Digital loyalty cards


Campaign Management

Seasonal promotion creation
Target audience selection
Limited-time offer creation
Campaign performance tracking
Automated client communications



4.3 Business Operations
4.3.1 Service Management

Service Catalog

Customizable service menu
Service categorization
Duration and pricing configuration
Service description and preparation notes



4.3.2 Performance Tracking

Service Completion System

Estimated price calculation at booking
Service completion confirmation by both parties
Post-service rating prompts
Push notifications for confirmation requests
Automatic triggers based on scheduled end time


Financial Tracking

Simple payment recording (cash, card, transfer)
Revenue tracking by service type
Client spending history
Daily/weekly/monthly performance views
Receipt generation for completed services
Tax reporting preparation


Business Analytics

Revenue reporting (recorded inside or outside platform)
Client retention metrics
Service popularity analysis
New vs. returning client ratio
Appointment fulfillment rate
Cancellation analysis



5. Technical Requirements
5.1 Platform Architecture

Web Application

Responsive design for desktop/tablet/mobile
Progressive Web App capabilities
Modern, intuitive UI with beauty industry aesthetic


Native Mobile Applications

iOS application (iPhone, iPad)
Android application
Push notification support
Offline functionality for core features



5.2 Security Requirements

Data Protection

End-to-end encryption for client-provider communications
GDPR compliance
Client data protection
Secure image storage
Role-based access controls



5.3 Integration Requirements

Communication Tools

Push notification system
SMS providers (integration with Portuguese carriers)
WhatsApp Business API
Email service providers



6. User Experience
6.1 Key User Journeys
6.1.1 Provider Onboarding Journey

Registration and account creation
Business profile and mini-website setup
Service menu configuration
Staff addition (if applicable)
Calendar configuration
Subscription tier selection

6.1.2 Client Discovery Journey

App download/website visit
Location-based provider search
Filter by service type/availability
Provider profile review
Appointment request
Service confirmation and rating

6.1.3 Loyalty Program Setup Journey

Program type selection
Reward configuration
Eligibility criteria setting
Client enrollment
Progress tracking
Reward redemption

6.1.4 Service Completion Journey

Client arrives for appointment
Service delivery
End-of-service trigger (time-based or manual)
Provider marks service as complete
Client receives confirmation request
Both parties confirm completion and rate experience
Provider records payment received (outside platform)
System updates financial tracking

6.2 User Interface Requirements

Provider App Key Screens

Dashboard with daily overview
Calendar (day, week, month views)
Client profile view
Treatment documentation interface
Performance analytics dashboard
Messaging center


Client App Key Screens

Provider discovery page
Provider profile view
Booking interface
Appointment management
Messaging center
Loyalty program progress



7. Critical Features - Detailed Specifications
7.1 Service Completion & Financial Tracking
The system will implement a frictionless approach to tracking business performance, especially for services paid in cash or other methods outside the platform:

Pre-Service Pricing:

When booking is confirmed, system calculates and displays expected service price
Both provider and client have visibility on expected cost
Option for provider to adjust final price based on actual services provided


Service Completion Triggers:

Automatic trigger: 15 minutes after scheduled end time
Manual provider trigger: "Service Completed" button
Manual client trigger: "Rate Your Experience" option
All triggers generate push notifications to the other party


Dual Confirmation Process:

Provider confirms service completion and can record payment details
Client confirms service received and can provide rating
System marks service as complete when both confirmations received
Reminder notifications if confirmation pending after 24 hours


Simple Payment Recording:

Provider indicates payment method (cash, card, transfer, other)
Records actual amount received
Option to note tips or additional charges
Digital receipt generation for client
No actual payment processing within platform


Performance Analytics:

Daily and weekly summaries of completed services
Revenue breakdowns by service type
Reconciliation between scheduled and completed services
Identification of no-shows or cancellations
Export functionality for accounting purposes



This approach provides valuable business insights without requiring the platform to handle actual payments, creating a lightweight solution that respects how most aestheticians currently operate in Portugal.
7.2 Loyalty Program Builder
The loyalty program builder will enable providers to create customized incentive programs:

Program Types:

Visit-based (e.g., every 10th visit free)
Spending-based (e.g., 10% off after €500 spent)
Service-specific (e.g., buy 5 facials, get 1 free)
Tiered rewards (Bronze, Silver, Gold levels)


Reward Configuration:

Percentage discounts
Fixed amount discounts
Free services
Service upgrades
Product gifts


Program Management:

Automatic progress tracking
Client notification of milestones
Expiration date settings
Visibility controls for clients
Performance analytics



This feature will help providers increase client retention and average lifetime value while providing clients with tangible benefits for their loyalty.
8. Implementation Considerations

Development Priorities:

Core scheduling and client management features first
Service completion and financial tracking system
Provider discovery and public calendars
Loyalty program and campaign features


Critical Success Factors:

Intuitive, frictionless provider experience
Reliable notification system
Accurate financial tracking without payment processing
Strong data privacy and security measures


Key Metrics:

Provider adoption rate
Client acquisition through platform
Appointment request to confirmation ratio
Service completion confirmation rate
Feature usage (loyalty programs, campaigns)