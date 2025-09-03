# SAMS Security Checklist

## 🔐 Secrets Management

This document outlines all the secrets and sensitive configuration that need to be properly managed in your SAMS deployment.

## ⚠️ Critical Secrets (Must be configured)

### Database Credentials
- `DB_USER` - Database username
- `DB_PASSWORD` - Database password  
- `DB_NAME` - Database name

### API Keys
- `NEXT_PUBLIC_GOOGLE_MAPS_API_KEY` - Google Maps API key for location services
- `GEMINI_API_KEY` - Google Gemini AI API key for AI assistant features

### Authentication & Security
- `JWT_SECRET` - Secret key for JWT token signing (use a strong, random string)
- `JWT_EXPIRY` - JWT token expiration time

### pgAdmin Access
- `PGADMIN_EMAIL` - pgAdmin login email
- `PGADMIN_PASSWORD` - pgAdmin login password

## 🚨 Security Best Practices

### 1. Environment Variables
- ✅ All secrets are now stored in `.env` file
- ✅ `.env` is in `.gitignore` and will NOT be committed
- ✅ No hardcoded secrets in source code

### 2. Database Security
- ✅ Use strong, unique passwords for database users
- ✅ Consider using environment-specific databases
- ✅ Enable SSL connections in production

### 3. API Key Security
- ✅ Never commit API keys to version control
- ✅ Use least-privilege API keys
- ✅ Rotate API keys regularly
- ✅ Monitor API key usage

### 4. JWT Security
- ✅ Use cryptographically strong random strings for JWT_SECRET
- ✅ Set appropriate expiration times
- ✅ Consider using asymmetric keys for production

## 🔧 Setup Instructions

### 1. Create Environment File
```bash
# Run the setup script
./setup-env.sh

# Or manually copy and edit
cp env.example .env
```

### 2. Configure Secrets
Edit `.env` file with your actual values:
```bash
# Database
DB_USER=your_db_user
DB_PASSWORD=your_strong_password
DB_NAME=your_db_name

# API Keys
NEXT_PUBLIC_GOOGLE_MAPS_API_KEY=your_google_maps_key
GEMINI_API_KEY=your_gemini_key

# JWT
JWT_SECRET=your_very_long_random_string_here
JWT_EXPIRY=24h

# pgAdmin
PGADMIN_EMAIL=your_email@domain.com
PGADMIN_PASSWORD=your_strong_password
```

### 3. Generate Strong Secrets
```bash
# Generate JWT secret (64 characters)
openssl rand -base64 48

# Generate database password (32 characters)
openssl rand -base64 24
```

## 🚫 What NOT to do

- ❌ Never commit `.env` files to git
- ❌ Never share API keys in public repositories
- ❌ Never use default passwords in production
- ❌ Never log secrets to console or files
- ❌ Never hardcode secrets in source code

## 🔍 Verification

After setup, verify:
1. `.env` file exists and contains your values
2. `.env` is listed in `.gitignore`
3. Docker containers start without errors
4. Database connections work
5. API endpoints respond correctly

## 📞 Support

If you encounter security-related issues:
1. Check that all environment variables are set
2. Verify `.env` file permissions (should be 600)
3. Ensure no secrets are logged in application output
4. Review Docker container environment variables

## 🔄 Updates

- **Last Updated**: $(date)
- **Version**: 1.0
- **Status**: ✅ All hardcoded secrets removed
