# Hospital Middleware System

A hospital middleware system built with Go/Gin that allows hospital staff to search patient records. Each staff member can only access patients from their own hospital.

## Tech Stack
- **Go** with **Gin** framework
- **PostgreSQL** for data storage
- **Nginx** as reverse proxy
- **Docker Compose** for containerization

## Quick Start

```bash
# Start all services (PostgreSQL, Go API, Nginx)
docker compose up --build

# The API is available at:
# - http://localhost (via Nginx)
# - http://localhost:8080 (direct to Go API)
```

## API Endpoints

### POST /staff/create
Create a new staff member.
```bash
curl -X POST http://localhost/staff/create \
  -H "Content-Type: application/json" \
  -d '{"username": "john", "password": "secret123", "hospital": "hospital_a"}'
```

### POST /staff/login
Login and get a JWT token.
```bash
curl -X POST http://localhost/staff/login \
  -H "Content-Type: application/json" \
  -d '{"username": "john", "password": "secret123", "hospital": "hospital_a"}'
```

### GET /patient/search
Search patients (requires authentication).
```bash
curl http://localhost/patient/search?national_id=1234567890123 \
  -H "Authorization: Bearer <your-token>"
```

**Optional query parameters:** `national_id`, `passport_id`, `first_name`, `middle_name`, `last_name`, `date_of_birth`, `phone_number`, `email`

## Running Tests
Unit Test for all functions (Create,Login,Search) and cover all cases
```bash
go test ./tests/... -v
```
