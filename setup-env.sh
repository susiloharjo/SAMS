#!/bin/bash

# SAMS Environment Setup Script
echo "🔧 Setting up SAMS environment configuration..."

# Check if .env already exists
if [ -f ".env" ]; then
    echo "⚠️  .env file already exists!"
    read -p "Do you want to overwrite it? (y/N): " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        echo "❌ Setup cancelled. .env file unchanged."
        exit 1
    fi
fi

echo ""
echo "🚀 Choose your deployment environment:"
echo "1) Docker (recommended for development)"
echo "2) Local development (outside Docker)"
echo "3) Custom configuration"
echo ""
read -p "Enter your choice (1-3): " -n 1 -r
echo

case $REPLY in
    1)
        echo "🐳 Setting up Docker environment..."
        if [ -f "env.docker" ]; then
            cp env.docker .env
            echo "✅ Created .env file for Docker deployment"
        else
            echo "❌ env.docker not found!"
            exit 1
        fi
        ;;
    2)
        echo "💻 Setting up local development environment..."
        if [ -f "env.local" ]; then
            cp env.local .env
            echo "✅ Created .env file for local development"
        else
            echo "❌ env.local not found!"
            exit 1
        fi
        ;;
    3)
        echo "⚙️  Setting up custom configuration..."
        if [ -f "env.example" ]; then
            cp env.example .env
            echo "✅ Created .env file from env.example"
            echo "📝 Edit .env file with your custom values"
        else
            echo "❌ env.example not found!"
            exit 1
        fi
        ;;
    *)
        echo "❌ Invalid choice. Exiting."
        exit 1
        ;;
esac

echo ""
echo "📝 Environment file created successfully!"
echo ""
echo "🔑 IMPORTANT: Edit .env file with your actual values:"
echo "   - Database credentials"
echo "   - API keys (Google Maps, Gemini)"
echo "   - JWT secrets"
echo "   - Other configuration values"
echo ""
echo "📖 Next steps:"
echo "   1. Edit .env with your actual values"
echo "   2. Run 'docker-compose up -d' to start services"
echo "   3. Run 'make setup' for complete setup"
echo ""
echo "⚠️  Remember: .env is in .gitignore and will NOT be committed to git"
echo ""
echo "🔍 Current .env configuration:"
echo "   DB_HOST: $(grep '^DB_HOST=' .env | cut -d'=' -f2)"
echo "   DB_PORT: $(grep '^DB_PORT=' .env | cut -d'=' -f2)"
echo "   REDIS_HOST: $(grep '^REDIS_HOST=' .env | cut -d'=' -f2)"
echo "   REDIS_PORT: $(grep '^REDIS_PORT=' .env | cut -d'=' -f2)"
