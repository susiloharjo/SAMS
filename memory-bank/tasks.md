# SAMS Project Tasks - Memory Bank

## Project Overview
**SAMS: Smart Asset Management System**
A focused MVP for a comprehensive web-based Asset Management System designed to manage the tracking, organization, and querying of various assets, following ISO 55001 standards.

## Project Status: MVP Planning
- **Complexity Level**: 2-3 (Moderate to High - Focused MVP)
- **Current Phase**: VAN (Initialization) - Complete
- **Next Phase**: PLAN (Task Planning)
- **Timeline**: 1 Week MVP Development

## MVP Features (1 Week Target)
- **Asset CRUD** - ISO 55001 compliant asset management
- **Dynamic Categories** - User-defined asset categories
- **Location Management** - Map integration with geotagging (lat/lon)
- **QR Code Generation** - Asset identification and tagging
- **LLM Integration** - AI queries for asset information
- **Demo Data** - Government buildings, IT equipment, vehicles

## Asset Fields (ISO 55001 Compliant)
- **Basic Info**: Asset ID, Name, Description, Category
- **Technical**: Type, Model, Serial Number, Manufacturer
- **Financial**: Acquisition Cost, Current Value, Depreciation
- **Operational**: Status, Condition, Criticality
- **Location**: Latitude, Longitude, Address, Building/Room
- **Lifecycle**: Acquisition Date, Expected Life, Maintenance Schedule
- **Compliance**: Certifications, Standards, Audit Information

## Technology Stack (MVP)
- **Frontend**: Next.js (React) with Tailwind CSS
- **Backend**: Go (Golang) with simple REST API
- **Database**: PostgreSQL (single database, simplified schema)
- **Maps**: Google Maps API for geotagging and location display
- **AI**: Google Gemini API (direct integration, no complex LangChain)
- **QR**: Basic QR generation library
- **Containerization**: Docker for development

## Architecture (MVP Simplified)
- Single backend service (no microservices for MVP)
- RESTful API design
- Simple database schema with ISO compliance
- Basic map integration and QR generation
- Direct AI API integration

## Development Workflow (Adjusted for MVP)
This project requires a focused MVP workflow:
1. ✅ **VAN** - Project initialization and scope definition
2. 🔄 **PLAN** - Detailed MVP implementation planning
3. ⏳ **IMPLEMENT** - Rapid MVP development (1 week)
4. ⏳ **REFLECT** - MVP review and lessons learned
5. ⏳ **ARCHIVE** - MVP documentation and future planning

## Current Focus
MVP planning and 1-week development strategy.

## Next Actions
1. Complete MVP planning in PLAN mode
2. Begin rapid implementation for 1-week timeline
3. Focus on core features for government/client demo
