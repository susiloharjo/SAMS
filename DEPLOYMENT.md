# SAMS Deployment Guide

## Quick Fix for UUID Extension Error

If you encounter the error:
```
ERROR: function uuid_generate_v4() does not exist (SQLSTATE 42883)
```

This means PostgreSQL is missing the UUID extension. Here's how to fix it:

### Option 1: Use the Updated Docker Setup (Recommended)

The repository now includes fixes for this issue:

1. **Clean deployment:**
   ```bash
   # Stop and remove existing containers
   docker-compose down -v
   
   # Remove volumes to ensure clean start
   docker volume rm sams_postgres_data 2>/dev/null || true
   
   # Start with the fixed setup
   docker-compose up --build -d
   ```

2. **Run the test script:**
   ```bash
   ./test-deployment.sh
   ```

### Option 2: Manual Fix (If Option 1 doesn't work)

1. **Connect to PostgreSQL container:**
   ```bash
   docker exec -it sams-postgres psql -U sams_user -d sams_db
   ```

2. **Enable UUID extension:**
   ```sql
   CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
   ```

3. **Restart backend:**
   ```bash
   docker-compose restart sams-backend
   ```

### What Was Fixed

1. **UUID Extension Priority**: Added `enable-uuid-extension.sql` that runs before the main init script
2. **Database Wait Script**: Backend now waits for PostgreSQL and UUID extension to be ready
3. **Proper Dependencies**: Services now wait for each other to be healthy before starting
4. **Health Checks**: Added health checks to ensure services are ready

### Files Added/Modified

- `enable-uuid-extension.sql` - Ensures UUID extension is available
- `wait-for-db.sh` - Backend waits for database to be ready
- `backend/Dockerfile` - Updated to include wait script and PostgreSQL client
- `docker-compose.yml` - Updated with proper service dependencies
- `test-deployment.sh` - Test script to verify deployment

### Environment Variables Required

Make sure your `.env` file contains:
```env
DB_HOST=sams-postgres
DB_PORT=5432
DB_NAME=sams_db
DB_USER=sams_user
DB_PASSWORD=your_password_here
REDIS_HOST=sams-redis
REDIS_PORT=6379
GEMINI_API_KEY=your_gemini_api_key
```

### Troubleshooting

1. **Check logs:**
   ```bash
   docker-compose logs -f sams-backend
   docker-compose logs -f sams-postgres
   ```

2. **Test database connection:**
   ```bash
   docker exec sams-postgres psql -U sams_user -d sams_db -c "SELECT uuid_generate_v4();"
   ```

3. **Test backend health:**
   ```bash
   curl http://localhost:8081/health
   ```

4. **Test MCP server:**
   ```bash
   curl http://localhost:8082/
   ```

The deployment should now work correctly on any server with Docker and Docker Compose installed.
