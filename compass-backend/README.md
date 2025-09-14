# Compass Backend

A project management backend system built with Golang, featuring JWT authentication, role-based access control, and comprehensive project management capabilities.

## Features

- **Authentication**: JWT-based authentication with access and refresh tokens
- **User Management**: Admin-only user creation and management
- **Project Management**: Create, update, and track project status
- **Specifications**: Version-controlled project specifications
- **RFIs**: Request for Information system with yes/no responses
- **Clean Architecture**: Well-structured codebase following clean architecture principles
- **Docker Support**: Easy deployment with Docker and docker-compose

## Tech Stack

- **Language**: Go 1.21
- **Framework**: Gin Web Framework
- **Database**: MySQL/MariaDB
- **ORM**: GORM
- **Authentication**: JWT
- **Containerization**: Docker

## Getting Started

### Prerequisites

- Go 1.21 or higher
- MySQL 8.0 or higher (or Docker)
- Docker and Docker Compose (optional)

### Installation

1. Clone the repository:
```bash
cd compass-backend
```

2. Copy the environment file:
```bash
cp .env.example .env
```

3. Update the `.env` file with your configuration.

### Running with Docker

The easiest way to run the application:

```bash
docker-compose up -d
```

This will start both the MySQL database and the backend application.

### Running Locally

1. Install dependencies:
```bash
go mod download
```

2. Make sure MySQL is running and create the database:
```sql
CREATE DATABASE compass_db;
```

3. Run the application:
```bash
go run cmd/api/main.go
```

## API Endpoints

### Authentication
- `POST /auth/signin` - Sign in with email and password
- `POST /auth/refresh` - Refresh access token
- `POST /api/auth/logout` - Logout (requires auth)

### User Management (Admin only)
- `POST /api/users` - Create new user
- `PATCH /api/users/:id/status` - Update user status
- `GET /api/users` - List all users

### Projects
- `POST /api/projects` - Create new project
- `GET /api/projects` - List all projects
- `GET /api/projects/:id` - Get project details
- `PATCH /api/projects/:id/status` - Update project status
- `DELETE /api/projects/:id` - Delete project (Admin only)

### Project Specifications
- `POST /api/projects/:id/specifications` - Create/update specification
- `GET /api/projects/:id/specifications` - List all specification versions

### RFIs
- `POST /api/projects/:id/rfis` - Create new RFI
- `GET /api/projects/:id/rfis` - List project RFIs
- `PATCH /api/rfis/:id/answer` - Answer an RFI

### Health Check
- `GET /health` - Application health status

## Default Admin Account

On first run, the system creates a default admin account:
- Email: admin@compass.com
- Password: AdminPassword123!

**Important**: Change these credentials immediately in production!

## Project Structure

```
compass-backend/
├── cmd/
│   └── api/
│       └── main.go           # Application entry point
├── config/
│   └── config.go            # Configuration management
├── db/
│   ├── connection.go        # Database connection
│   └── seed.go             # Database seeding
├── internal/
│   ├── controllers/        # HTTP handlers
│   ├── middleware/         # Middleware functions
│   ├── models/            # Domain models
│   ├── repositories/      # Data access layer
│   ├── routes/           # Route definitions
│   ├── services/         # Business logic
│   └── utils/           # Utility functions
├── docker-compose.yml    # Docker Compose configuration
├── Dockerfile           # Docker image definition
├── go.mod              # Go module file
└── README.md          # This file
```

## Security Features

- Password hashing with bcrypt
- JWT token-based authentication
- Role-based access control (Admin/User)
- Request logging
- CORS support (can be added)

## Development

### Running Tests

```bash
go test ./...
```

### Building for Production

```bash
go build -o main ./cmd/api
```

## Environment Variables

See `.env.example` for all available configuration options.

## License

This project is proprietary software.