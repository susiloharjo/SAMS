#!/bin/bash

echo "=== SAMS Deployment Test ==="
echo "This script will test the deployment on a different server"
echo ""

# Check if .env file exists
if [ ! -f .env ]; then
    echo "❌ .env file not found. Please create it with the required environment variables."
    exit 1
fi

echo "✅ .env file found"

# Source environment variables
source .env

echo "📋 Environment variables:"
echo "  DB_HOST: $DB_HOST"
echo "  DB_PORT: $DB_PORT"
echo "  DB_NAME: $DB_NAME"
echo "  DB_USER: $DB_USER"
echo ""

# Clean up any existing containers
echo "🧹 Cleaning up existing containers..."
docker-compose down -v

# Remove any existing volumes to ensure clean start
echo "🗑️  Removing existing volumes..."
docker volume rm sams_postgres_data 2>/dev/null || true

# Build and start services
echo "🔨 Building and starting services..."
docker-compose up --build -d

# Wait for services to be ready
echo "⏳ Waiting for services to be ready..."
sleep 30

# Check if services are running
echo "🔍 Checking service status..."
docker-compose ps

# Test database connection
echo "🗄️  Testing database connection..."
docker exec sams-postgres psql -U $DB_USER -d $DB_NAME -c "SELECT version();" 2>/dev/null && echo "✅ Database connection successful" || echo "❌ Database connection failed"

# Test UUID extension
echo "🔧 Testing UUID extension..."
docker exec sams-postgres psql -U $DB_USER -d $DB_NAME -c "SELECT uuid_generate_v4();" 2>/dev/null && echo "✅ UUID extension working" || echo "❌ UUID extension failed"

# Test backend health
echo "🏥 Testing backend health..."
sleep 10
curl -s http://localhost:8081/health 2>/dev/null && echo "✅ Backend health check passed" || echo "❌ Backend health check failed"

# Test MCP server
echo "🤖 Testing MCP server..."
curl -s http://localhost:8082/ 2>/dev/null && echo "✅ MCP server responding" || echo "❌ MCP server not responding"

echo ""
echo "=== Deployment Test Complete ==="
echo "Check the logs with: docker-compose logs -f"
