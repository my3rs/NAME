# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a blog/content management system written in Go using the Iris web framework. It provides REST API endpoints for managing blog posts, pages, users, comments, and file attachments with JWT authentication.

## Development Commands

### Build and Run
```bash
# Install dependencies
go mod tidy

# Build the application
go build -o NAME main.go

# Run in development
go run main.go

# Run with Docker
docker build -t name .
docker run -p 8000:8000 name
```

### Testing
```bash
# Run all tests
go test ./...

# Run specific test packages
go test ./test/model/
go test ./test/service/
```

## Configuration

Configuration uses TOML files with environment variable overrides. Copy `conf/example.toml` to `name.toml` in the project root and modify as needed.

**Configuration precedence:**
1. Environment variables (prefixed with `NAME_`)
2. `./name.toml`
3. `./config/name.toml` 
4. `./bin/config/name.toml`
5. `$HOME/.name/name.toml`

**Key settings:**
- `PORT`: Server port (default: 8000)
- `MODE`: "development" or "production"
- `DATA_PATH`: Storage location for uploads and database
- Database: Supports SQLite (development) and PostgreSQL (production)
- JWT: Requires RSA private/public key pair for token signing

## Architecture

The codebase follows clean architecture with these layers:

- **main.go**: Application entry point with Iris server setup
- **conf/**: Configuration management using Viper
- **route/**: API route definitions and middleware setup
- **controller/**: HTTP request handlers for API endpoints
- **service/**: Business logic layer
- **model/**: Database entities and GORM models
- **database/**: Database connection and initialization
- **middleware/**: HTTP middleware (JWT, logging, etc.)
- **utils/**: Utility functions
- **customerror/**: Custom error types

## API Structure

RESTful API with the following main endpoints:
- `/api/v1/auth`: Authentication (login, register)
- `/api/v1/posts`: Blog posts (public read, protected write)
- `/api/v1/pages`: Static pages (protected)
- `/api/v1/categories`: Content categories
- `/api/v1/tags`: Content tags  
- `/api/v1/comments`: Comments system
- `/api/v1/users`: User management
- `/api/v1/attachments`: File uploads (protected)
- `/api/v1/settings`: Application settings
- `/api/v1/meta`: Blog metadata and statistics

## Database Schema

**Content Types:**
- `post`: Blog articles
- `page`: Static pages
- `digu`: Micro-posts/short updates

**Main Entities:**
- **Content**: Blog posts/pages with markdown support
- **User**: User accounts with role-based access
- **Category**: Content categorization
- **Tag**: Content tagging (many-to-many relationship)
- **Comment**: Comments on content
- **Attachment**: File uploads
- **Setting**: Application configuration

## Key Technologies

- **Iris v12**: Web framework
- **GORM**: ORM for database operations
- **Viper**: Configuration management
- **JWT**: Authentication with custom fork (`github.com/my3rs/jwt`)
- **Blackfriday v2**: Markdown to HTML conversion
- **Bluemonday**: HTML sanitization
- **Testify**: Testing framework

## Development Notes

- Uses SQLite for development, PostgreSQL for production
- JWT tokens signed with RSA keys (generate or use existing in `bin/config/`)
- File uploads stored in `DATA_PATH/uploads/`
- HTML content is sanitized for security
- Supports Chinese language interface elements
- CORS middleware configured for cross-origin requests
- Role-based access control (admin/user roles)