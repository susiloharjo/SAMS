# SAMS Authentication & RBAC Setup Guide

This guide explains how to set up and use the authentication and Role-Based Access Control (RBAC) system in SAMS.

## Overview

The SAMS application now includes:
- **JWT-based authentication** with access and refresh tokens
- **Role-based access control** with three user roles: `admin`, `manager`, and `user`
- **Protected API endpoints** that require authentication and appropriate permissions
- **Frontend authentication** with login/logout functionality

## User Roles & Permissions

### Admin Role
- **Full access** to all system features
- Can manage users, assets, categories, departments
- Can perform all CRUD operations
- Access to AI features and system administration

### Manager Role
- **Full access** to asset management
- Can create, read, update, and delete assets
- Can manage categories and departments
- Access to AI features
- **Cannot** manage users

### User Role
- **Read-only access** to assets
- Can view asset details and generate QR codes
- **Cannot** create, update, or delete assets
- **Cannot** access management features

## Backend Setup

### 1. Environment Variables

Create a `.env` file in the `backend/` directory with the following variables:

```bash
# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_USER=sams_user
DB_PASSWORD=sams_password
DB_NAME=sams_db
DB_SSL_MODE=disable

# Server Configuration
SERVER_PORT=8080

# JWT Configuration (IMPORTANT: Change these in production!)
JWT_SECRET=your-super-secret-jwt-key-change-this-in-production
JWT_REFRESH_SECRET=your-super-secret-refresh-key-change-this-in-production
JWT_EXPIRATION_HOURS=24
JWT_REFRESH_EXPIRATION_DAYS=7

# API Configuration
NEXT_PUBLIC_API_URL=http://localhost:8080

# Google AI Configuration
GOOGLE_AI_API_KEY=your-google-ai-api-key-here
```

### 2. Install Dependencies

```bash
cd backend
go mod tidy
```

### 3. Create Default Users

Run the user creation script to set up default accounts:

```bash
cd backend
go run cmd/create-users/main.go
```

This will create three default users:
- **Admin**: `admin` / `admin123`
- **Manager**: `manager` / `manager123`
- **User**: `user` / `user123`

### 4. Start the Backend

```bash
cd backend
go run cmd/main.go
```

## Frontend Setup

### 1. Environment Variables

Create a `.env.local` file in the `frontend/` directory:

```bash
NEXT_PUBLIC_API_URL=http://localhost:8080
```

### 2. Install Dependencies

```bash
cd frontend
npm install
```

### 3. Start the Frontend

```bash
cd frontend
npm run dev
```

## API Endpoints

### Public Endpoints (No Authentication Required)
- `POST /api/v1/auth/login` - User login
- `POST /api/v1/auth/refresh` - Refresh access token

### Protected Endpoints

#### Asset Management
- `GET /api/v1/assets` - List assets (all authenticated users)
- `GET /api/v1/assets/:id` - Get asset details (all authenticated users)
- `POST /api/v1/assets` - Create asset (admin, manager only)
- `PUT /api/v1/assets/:id` - Update asset (admin, manager only)
- `DELETE /api/v1/assets/:id` - Delete asset (admin, manager only)

#### Category Management
- `GET /api/v1/categories` - List categories (all authenticated users)
- `POST /api/v1/categories` - Create category (admin, manager only)
- `PUT /api/v1/categories/:id` - Update category (admin, manager only)
- `DELETE /api/v1/categories/:id` - Delete category (admin, manager only)

#### Department Management
- `GET /api/v1/departments` - List departments (all authenticated users)
- `POST /api/v1/departments` - Create department (admin, manager only)
- `PUT /api/v1/departments/:id` - Update department (admin, manager only)
- `DELETE /api/v1/departments/:id` - Delete department (admin, manager only)

#### User Management (Admin Only)
- `GET /api/v1/users` - List users
- `POST /api/v1/users` - Create user
- `PUT /api/v1/users/:id` - Update user
- `DELETE /api/v1/users/:id` - Delete user

#### AI Features (Admin, Manager Only)
- `POST /api/v1/ai/query` - AI-powered queries

## Frontend Usage

### 1. Login

Navigate to `/login` to access the login page. Use the default credentials created by the setup script.

### 2. Protected Routes

The frontend automatically protects routes based on user roles:
- **All authenticated users** can access `/assets`
- **Admin and Manager** can access management features
- **Regular users** see read-only views

### 3. Logout

Users can logout using the logout button, which will clear their session and redirect to the login page.

## Security Features

### JWT Tokens
- **Access tokens** expire after 24 hours (configurable)
- **Refresh tokens** expire after 7 days (configurable)
- Tokens are stored in localStorage (consider httpOnly cookies for production)

### Password Security
- Passwords are hashed using bcrypt with cost factor 10
- No plain text passwords are stored

### Role Validation
- All API endpoints validate user roles before allowing access
- Frontend components check permissions before rendering features

## Production Considerations

### Security
1. **Change default JWT secrets** to strong, unique values
2. **Use HTTPS** in production
3. **Consider httpOnly cookies** instead of localStorage for tokens
4. **Implement rate limiting** for authentication endpoints
5. **Add password complexity requirements**

### Database
1. **Use strong database passwords**
2. **Enable SSL connections**
3. **Regular database backups**
4. **Monitor failed login attempts**

### Monitoring
1. **Log authentication events**
2. **Monitor API usage patterns**
3. **Set up alerts for security events**

## Troubleshooting

### Common Issues

1. **"Invalid credentials" error**
   - Verify the user exists in the database
   - Check if the user account is active

2. **"Insufficient permissions" error**
   - Verify the user has the required role
   - Check if the endpoint requires specific permissions

3. **Token expiration issues**
   - Check JWT expiration settings
   - Verify refresh token functionality

4. **Database connection issues**
   - Verify database credentials
   - Check if PostgreSQL is running
   - Verify database exists and is accessible

### Debug Mode

To enable debug logging, set the following environment variable:
```bash
LOG_LEVEL=debug
```

## Support

For issues or questions about the authentication system:
1. Check the application logs
2. Verify environment variable configuration
3. Ensure database connectivity
4. Test with default user accounts first
