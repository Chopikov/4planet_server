# 4Planet Backend

A production-ready Go backend that implements the REST API for the 4Planet tree planting platform. This backend handles user authentication, donations, payments, and project management with a focus on sustainability and transparency.

## Features

- **Authentication**: Cookie-based sessions with email verification
- **Payments**: CloudPayments integration with webhook support
- **Database**: PostgreSQL with GORM ORM
- **Admin Interface**: QOR Admin for data management
- **API Documentation**: OpenAPI 3.0.3 spec with Swagger UI
- **Email**: SMTP support with development mock (MailHog)
- **Security**: Password hashing, session management, CORS

## Tech Stack

- **Language**: Go 1.22+
- **Framework**: Gin (HTTP router)
- **ORM**: GORM with PostgreSQL driver
- **Database**: PostgreSQL 15
- **Admin**: QOR Admin
- **Authentication**: Custom session-based auth
- **Payments**: CloudPayments integration
- **Email**: SMTP with MailHog for development

## Prerequisites

- Go 1.22 or later
- Docker and Docker Compose
- PostgreSQL (if running locally)
- Make (optional, for convenience)

## Quick Start

### 1. Clone and Setup

```bash
git clone <repository-url>
cd TreeProject
```

### 2. Environment Configuration

Create a `.env` file in the root directory:

```bash
# App Configuration
APP_BASE_URL=http://localhost:8080
APP_COOKIE_NAME=session_id
APP_COOKIE_DOMAIN=localhost
APP_COOKIE_SECURE=false

# Database
DB_DSN=postgres://postgres:postgres@localhost:5432/planet?sslmode=disable

# SMTP (optional, uses MailHog in development)
SMTP_HOST=
SMTP_PORT=587
SMTP_USER=
SMTP_PASSWORD=
SMTP_FROM=noreply@4planet.local

# CloudPayments (optional for development)
CLOUDPAYMENTS_PUBLIC_ID=
CLOUDPAYMENTS_SECRET=

# Logging
LOG_LEVEL=debug

# Admin credentials
ADMIN_USERNAME=admin
ADMIN_PASSWORD=admin
```

### 3. Development Setup

```bash
# Install dependencies
make deps

# Start development environment
make dev-setup
```

This will:
- Download Go dependencies
- Start PostgreSQL and MailHog
- Create database schema
- Seed initial data
- Start the API server

### 4. Access Services

- **API**: http://localhost:8080
- **Admin Interface**: http://localhost:8080/admin (admin/admin)
- **API Documentation**: http://localhost:8080/docs
- **OpenAPI Spec**: http://localhost:8080/openapi.yaml
- **MailHog**: http://localhost:8025 (SMTP testing)

## API Endpoints

### Authentication
- `POST /v1/auth/register` - User registration
- `POST /v1/auth/login` - User login
- `POST /v1/auth/logout` - User logout
- `GET /v1/auth/me` - Current user info
- `POST /v1/auth/verify-email/confirm` - Email verification
- `POST /v1/auth/password/forgot` - Password reset request
- `POST /v1/auth/password/reset` - Password reset

### Payments & Donations
- `POST /v1/payments/intents` - Create payment intent
- `POST /v1/subscriptions/intents` - Create subscription intent
- `GET /v1/donations` - List user donations

### Projects & Media
- `GET /v1/projects` - List projects
- `GET /v1/projects/{id}` - Get project details
- `GET /v1/projects/{id}/media` - Get project media

### Achievements & Leaderboard
- `GET /v1/me/achievements` - User achievements (authenticated)
- `GET /v1/achievements` - All available achievements catalog (authenticated)
- `GET /v1/badges` - All available achievements catalog (public)
- `GET /v1/leaderboard` - Top users by trees

### Sharing
- `POST /v1/shares/profile` - Create profile share
- `POST /v1/shares/donation/{id}` - Create donation share
- `GET /v1/shares/resolve/{slug}` - Resolve share link

### Webhooks
- `POST /webhooks/{provider}` - Payment provider webhooks

## Database Schema

The application uses PostgreSQL with the following key tables:

- **users** - User accounts and authentication
- **sessions** - User sessions for cookie auth
- **payments** - Payment transactions
- **donations** - Tree planting donations
- **projects** - Tree planting projects
- **achievements** - User achievements and badges
- **tree_prices** - Tree prices by currency

## Development

### Running Locally

```bash
# Build and run
make run

# Or run directly
go run ./cmd/api
```

### Testing

```bash
# Run all tests
make test

# Run specific test
go test ./pkg/auth -v
```

### Code Quality

```bash
# Format code
make fmt

# Lint code
make lint

# Install development tools
make install-tools
```

### Database Operations

```bash
# Run migrations
make migrate

# Seed database
make seed

# View logs
make docker-logs
```

## Docker

### Build Image

```bash
make docker-build
```

### Run Services

```bash
# Start all services
make docker-run

# Stop services
make docker-stop

# View logs
make docker-logs
```

## Production Deployment

### 1. Build Production Binary

```bash
make prod-build
```

### 2. Environment Variables

Set production environment variables:

```bash
export APP_BASE_URL=https://api.yourdomain.com
export APP_COOKIE_SECURE=true
export APP_COOKIE_DOMAIN=.yourdomain.com
export DB_DSN=postgres://user:pass@db:5432/planet?sslmode=require
export SMTP_HOST=smtp.yourdomain.com
export SMTP_USER=your-smtp-user
export SMTP_PASSWORD=your-smtp-password
export CLOUDPAYMENTS_PUBLIC_ID=your-public-id
export CLOUDPAYMENTS_SECRET=your-secret
export LOG_LEVEL=info
```

### 3. Run Application

```bash
./bin/api
```

## Documentation

- **`README.md`** - This file, project overview and setup
- **`MIGRATIONS.md`** - Complete guide for database migrations and GORM integration
- **`openapi.yaml`** - API specification (v2.1.0)
- **`migrations/`** - Database migration files generated from GORM models

## Project Structure

```
.
├── cmd/api/                 # Main application entry point
├── internal/                # Internal packages
│   ├── config/             # Configuration management
│   ├── database/           # Database connection and setup
│   ├── handlers/           # HTTP request handlers
│   ├── middleware/         # HTTP middleware
│   └── models/             # GORM data models
├── pkg/                    # Public packages
│   ├── auth/               # Authentication service
│   ├── mailer/             # Email service
│   └── payments/           # Payment processing
├── migrations/             # Database migrations (GORM-generated)
├── scripts/                # Database seeding scripts
├── web/admin/              # Admin interface assets
├── docker-compose.yml      # Development environment
├── Dockerfile              # Production container
├── Makefile                # Development commands
├── go.mod                  # Go module definition
└── openapi.yaml            # API specification
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests for new functionality
5. Ensure all tests pass
6. Submit a pull request

## Testing

The application includes unit tests for core functionality:

- Authentication (register/login/verify/reset)
- Payment webhook idempotency
- Donation creation and user counter updates
- Achievement auto-awarding

Run tests with:

```bash
make test
```

## Security

- Passwords are hashed using bcrypt
- Sessions use secure, HttpOnly cookies
- CORS is configured for cross-origin requests
- Input validation on all endpoints
- SQL injection protection via GORM

## Monitoring

- Request ID tracking for debugging
- Structured logging with logrus
- Health check endpoint at `/health`
- Database connection monitoring

## License

[Add your license information here]

## Support

For questions and support, please [create an issue](link-to-issues) or contact the development team.
