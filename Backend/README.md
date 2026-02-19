# DevOps API - Beego Backend

A Role-Based Access Control (RBAC) API built with Beego framework and PostgreSQL.

## Features

- JWT-based authentication
- Role-based access control (Admin & User roles)
- PostgreSQL database with **golang-migrate** for schema versioning
- RESTful API endpoints
- Password hashing with bcrypt
- Database migrations (like Alembic/Prisma)

## Prerequisites

- Go 1.23+ installed
- Docker installed
- PostgreSQL container running

## Database Setup

The PostgreSQL container is already configured with:

- Host: localhost:5432
- Database: devops_db
- Username: devops
- Password: devops123

Start the container with Docker Compose:

```bash
docker-compose up -d
```

Or if already running individually:

```bash
docker start postgres-devops
```

### Database Migrations

This project uses **golang-migrate** for database schema versioning. Migrations are automatically applied on application startup.

**Quick Start:**
```powershell
# Apply all pending migrations
.\migrate.ps1 up

# Rollback the last migration
.\migrate.ps1 down

# Create a new migration
.\migrate.ps1 create add_new_feature

# Check current version
.\migrate.ps1 version
```

📖 **See [MIGRATIONS.md](MIGRATIONS.md) for complete migration documentation**

## Installation

1. Navigate to the project directory:

```bash
cd Backend/devops-api
```

2. Install dependencies:

```bash
go mod tidy
```

3. Run the application:

```bash
go run main.go
```

The server will start on `http://localhost:8080`

**Note:** Database migrations run automatically on startup, ensuring your schema is always up-to-date.

## API Endpoints

### Authentication (Public)

#### 1. Login

```http
POST /api/auth/login
Content-Type: application/json

{
  "email": "admin@example.com",
  "password": "admin123"
}

Response:
{
  "access": "eyJhbGc...",
  "refresh": "eyJhbGc...",
  "user": {
    "id": 1,
    "email": "admin@example.com",
    "name": "Admin User",
    "role": "admin"
  }
}
```

#### 2. Refresh Token

```http
POST /api/auth/refresh
Content-Type: application/json

{
  "refresh": "eyJhbGc..."
}

Response:
{
  "access": "eyJhbGc..."
}
```

#### 3. Get Current User

```http
GET /api/auth/me
Authorization: Bearer <access_token>

Response:
{
  "id": 1,
  "email": "admin@example.com",
  "name": "Admin User",
  "role": "admin"
}
```

### Admin Endpoints (Requires Admin Role)

#### 4. Get Statistics

```http
GET /api/admin/stats
Authorization: Bearer <access_token>

Response:
{
  "total_users": 10,
  "active_users": 8
}
```

### User Endpoints (Requires User Role)

#### 5. Get User Summary

```http
GET /api/user/summary
Authorization: Bearer <access_token>

Response:
{
  "last_login": "2026-02-19T15:30:00Z",
  "message": "Welcome, John Doe!"
}
```

## Creating Test Users

Connect to the PostgreSQL database and create test users:

```sql
-- Admin user (password: admin123)
INSERT INTO accounts_user (email, password, name, role, is_active, is_staff, is_superuser)
VALUES ('admin@example.com', '$2a$10$YourHashedPasswordHere', 'Admin User', 'admin', true, true, false);

-- Regular user (password: user123)
INSERT INTO accounts_user (email, password, name, role, is_active, is_staff, is_superuser)
VALUES ('user@example.com', '$2a$10$YourHashedPasswordHere', 'Regular User', 'user', true, false, false);
```

Or use the helper script to create users with properly hashed passwords.

## Configuration

Configuration is stored in `conf/app.conf`:

```ini
# Server
httpport = 8080
runmode = dev

# Database
dbhost = localhost
dbport = 5432
dbuser = devops
dbpass = devops123
dbname = devops_db

# JWT
jwtsecret = your-secret-key-change-this-in-production
jwtexpiration = 3600
```

## Project Structure

```
devops-api/
├── conf/              # Configuration files
│   └── app.conf      # App settings, DB connection, JWT config
├── controllers/       # API controllers
│   ├── auth.go       # Authentication endpoints
│   ├── admin.go      # Admin endpoints
│   └── user.go       # User endpoints
├── database/          # Database connection & migrations
│   └── db.go         # DB init + auto-migration runner
├── middleware/        # Authentication middleware
│   └── auth.go       # JWT verification & RBAC
├── migrations/        # Database migration files
│   ├── 000001_create_users_table.up.sql
│   └── 000001_create_users_table.down.sql
├── models/            # Data models
│   └── user.go       # User model & DB queries
├── routers/           # Route definitions
│   └── router.go     # API route configuration
├── scripts/           # Helper scripts
│   └── create_test_users.go
├── utils/             # Utility functions
│   └── jwt.go        # JWT token generation/parsing
├── migrate.ps1        # Migration management script (Windows)
├── Makefile           # Migration management (Linux/Mac)
├── MIGRATIONS.md      # Migration documentation
├── main.go            # Application entry point
└── README.md          # This file
```

## Database Migrations

### Overview

This project uses **golang-migrate** for database schema versioning:

- ✅ Version controlled schema changes
- ✅ Automatic migrations on startup
- ✅ Rollback support
- ✅ Production ready (like Alembic/Prisma)

### Quick Commands

```powershell
# Apply pending migrations
.\migrate.ps1 up

# Rollback last migration
.\migrate.ps1 down

# Create new migration
.\migrate.ps1 create add_profile_fields

# Check version
.\migrate.ps1 version
```

### Migration Files

Located in `migrations/` directory:
- `*.up.sql` - Apply changes
- `*.down.sql` - Rollback changes

**Example - Creating a new migration:**

```powershell
.\migrate.ps1 create add_user_avatar

# Creates:
# migrations/000002_add_user_avatar.up.sql
# migrations/000002_add_user_avatar.down.sql
```

Edit the files:

```sql
-- 000002_add_user_avatar.up.sql
ALTER TABLE accounts_user ADD COLUMN avatar_url VARCHAR(500);

-- 000002_add_user_avatar.down.sql
ALTER TABLE accounts_user DROP COLUMN IF EXISTS avatar_url;
```

Apply it:
```powershell
.\migrate.ps1 up
```

📖 **Full migration guide:** [MIGRATIONS.md](MIGRATIONS.md)

## Error Codes

- 200: Success
- 400: Bad Request
- 401: Unauthorized (invalid credentials or token)
- 403: Forbidden (insufficient permissions)
- 404: Not Found
- 500: Internal Server Error

## Security Notes

- Change the `jwtsecret` in production
- Access tokens expire in 1 hour (configurable)
- Refresh tokens expire in 7 days
- Passwords are hashed using bcrypt
- CORS is enabled for development

## Testing with cURL

```bash
# Login
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@example.com","password":"admin123"}'

# Get current user
curl -X GET http://localhost:8080/api/auth/me \
  -H "Authorization: Bearer <your_access_token>"

# Get admin stats
curl -X GET http://localhost:8080/api/admin/stats \
  -H "Authorization: Bearer <admin_access_token>"
```
