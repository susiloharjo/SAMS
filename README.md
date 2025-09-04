# SAMS: Smart Asset Management System

## Project Overview
SAMS is a comprehensive web-based Asset Management System designed to manage the tracking, organization, and querying of various assets, including fixed assets (e.g., IT equipment, vehicles, machinery) and non-fixed assets (e.g., land, buildings).

## 🚀 **Current Status: MVP Development (Day 1-2)**

### **What's Been Implemented**
- ✅ **Docker Compose Setup** - PostgreSQL, Redis, pgAdmin
- ✅ **Database Schema** - ISO 55001 compliant asset management
- ✅ **Sample Data** - Government/company demo assets
- 🔄 **Backend Setup** - Go project structure (in progress)
- ⏳ **Frontend Setup** - Next.js application (pending)

### **Next Steps**
1. **Start Backend Development** - Go API with Fiber framework
2. **Create Frontend** - Next.js with Tailwind CSS
3. **Implement Core Features** - Asset CRUD, categories, location
4. **Add Map Integration** - Google Maps with geotagging
5. **QR Code Generation** - Asset identification
6. **AI Integration** - Gemini API for asset queries

## 🐳 **Quick Start with Docker**

### **1. Start the Database Environment**
```bash
# Start all services
docker-compose up -d

# Check status
docker-compose ps

# View logs
docker-compose logs -f postgres
```

### **2. Access Services**
- **Frontend**: `http://localhost:3000`
- **Backend API**: `http://localhost:8081`
- **PostgreSQL**: `localhost:5433` (configured via environment variables)
- **pgAdmin**: `http://localhost:5051` (configured via environment variables)
- **Redis**: `localhost:6380`

### **3. Environment Setup**
Before starting services, copy `env.example` to `.env` and configure your environment variables:
```bash
cp env.example .env
# Edit .env with your actual values
```

### **4. Demo Credentials**
The system comes with pre-configured demo accounts for testing:

| Role | Username | Password | Access Level |
|------|----------|----------|--------------|
| **Admin** | `admin` | `user.1001` | Full system access, user management |
| **Manager** | `manager` | `manager123` | Asset CRUD, categories, departments, AI assistant |
| **User** | `user` | `user123` | View-only access to assets and reports |

### **5. Stop Services**
```bash
docker-compose down
```

## Memory Bank System Setup

This project uses the Memory Bank system for structured development workflow management. The Memory Bank system has been initialized and configured for this project.

### Memory Bank Files
- `memory-bank/tasks.md` - Main project task tracking
- `memory-bank/activeContext.md` - Current project focus and context
- `memory-bank/progress.md` - Implementation progress tracking

### Custom Modes Setup

To use the Memory Bank system, you need to set up custom modes in Cursor:

1. **Open Cursor** and click on the mode selector in the chat panel
2. **Select "Add custom mode"**
3. **Configure each mode** with the following settings:

#### 🔍 VAN Mode (Initialization)
- **Name**: 🔍 VAN
- **Tools**: Enable "Codebase Search", "Read File", "Terminal", "List Directory"
- **Advanced options**: Copy content from `custom_modes/van_instructions.md`

#### 📋 PLAN Mode (Task Planning)
- **Name**: 📋 PLAN
- **Tools**: Enable "Codebase Search", "Read File", "Terminal", "List Directory"
- **Advanced options**: Copy content from `custom_modes/plan_instructions.md`

#### 🎨 CREATIVE Mode (Design Decisions)
- **Name**: 🎨 CREATIVE
- **Tools**: Enable "Codebase Search", "Read File", "Terminal", "List Directory", "Edit File"
- **Advanced options**: Copy content from `custom_modes/creative_instructions.md`

#### ⚒️ IMPLEMENT Mode (Code Implementation)
- **Name**: ⚒️ IMPLEMENT
- **Tools**: Enable all tools
- **Advanced options**: Copy content from `custom_modes/implement_instructions.md`

#### 🔍 REFLECT Mode (Review)
- **Name**: 🔍 REFLECT
- **Tools**: Enable "Codebase Search", "Read File", "Terminal", "Terminal", "List Directory"
- **Advanced options**: Copy content from `custom_modes/reflect_archive_instructions.md` (REFLECT section)

#### 📚 ARCHIVE Mode (Documentation)
- **Name**: 📚 ARCHIVE
- **Tools**: Enable "Codebase Search", "Read File", "Terminal", "List Directory", "Edit File"
- **Advanced options**: Copy content from `custom_modes/reflect_archive_instructions.md` (ARCHIVE section)

### Usage

1. **Start with VAN Mode**: Switch to VAN mode and type "VAN" to initialize the project
2. **Follow the Workflow**: Based on complexity assessment, follow the recommended workflow:
   - **Level 1**: VAN → IMPLEMENT
   - **Level 2**: VAN → PLAN → IMPLEMENT → REFLECT
   - **Level 3-4**: VAN → PLAN → CREATIVE → IMPLEMENT → REFLECT → ARCHIVE
3. **Use QA Functions**: Type "QA" in any mode for technical validation

### Current Status

- **Complexity Level**: 2-3 (Moderate to High - Focused MVP)
- **Current Phase**: IMPLEMENT (MVP Development)
- **Timeline**: 1 Week MVP Development
- **Workflow**: Simplified MVP workflow (VAN → PLAN → IMPLEMENT → REFLECT → ARCHIVE)

### Next Steps

1. Set up custom modes in Cursor using the instruction files
2. Start backend development with Go and Fiber
3. Create frontend with Next.js and Tailwind CSS
4. Implement core MVP features for government/client demo

## Project Features

- Asset Registry (Master Data Asset)
- QR Code / Barcode Integration
- Lifecycle Management
- Location & Ownership Tracking
- Maintenance Management
- Document Management
- Reporting & Analytics
- User & Role Management (RBAC)
- Audit Trail & Compliance
- AI Assistant / Smart Search

## Technology Stack

- **Frontend**: Next.js (React) with Tailwind CSS
- **Backend**: Go (Golang) with Fiber or Echo
- **Database**: PostgreSQL + Qdrant (vector database)
- **AI Integration**: Google Gemini API via LangChain
- **Containerization**: Docker + Docker Compose

## Architecture

- Microservices approach with modular services
- RESTful API design
- Vector search capabilities for AI-driven queries
- Cloud storage for document management

For more information about the Memory Bank system, see the documentation in the `.cursor/rules/isolation_rules/` directory.
