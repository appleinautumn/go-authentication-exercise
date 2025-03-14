# Go Authentication Exercise

A RESTful API service built with Go that provides user authentication and management. This project demonstrates clean architecture principles with a domain-driven design approach and includes:

- User authentication (login/signup) with JWT token
- User management with pagination support
- Middleware for protected routes
- PostgreSQL database integration

## Requirements

This project is tested with:

- Go 1.22
- Postgres 16

## Installation

Clone the project

```bash
git clone git@github.com:appleinautumn/go-authentication-exercise.git
```

Go to the project directory

```bash
cd go-authentication-exercise
```

This service contains a `.env.example` file that defines environment variables you need to set. Copy and set the variables to a new `.env` file.

```bash
cp .env.example .env
```

Start the app

```bash
go run main.go
```

## Database

If you have not created the database, please create one before going to the next step.

We're using [golang-migrate](https://github.com/golang-migrate/migrate) for the migration.

### Without Docker

Install the package

```bash
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
```

Run the migration (change the value accordingly)

```bash
migrate -path=migrations -database "postgres://postgres:password@127.0.0.1:5432/database?sslmode=disable" up
```

To rollback

```bash
migrate -path=migrations -database "postgres://postgres:password@127.0.0.1:5432/database?sslmode=disable" down 1
```

### With Docker

Run the migration (change the value accordingly)

```bash
docker run -v "$(pwd)"/migrations:/migrations --network host migrate/migrate -path=/migrations/ -database "postgres://postgres:password@127.0.0.1:5432/database?sslmode=disable" up
```

To rollback

```bash
docker run -v "$(pwd)"/migrations:/migrations --network host migrate/migrate -path=/migrations/ -database "postgres://postgres:password@127.0.0.1:5432/database?sslmode=disable" down 1
```

## Deployment

Build the binary

```bash
go build -o go-auth
```

Run it

```bash
./go-auth
```

## API Endpoints

| Method | Endpoint     | Description                                     | Authentication |
| ------ | ------------ | ----------------------------------------------- | -------------- |
| GET    | /            | Root endpoint (health check)                    | No             |
| POST   | /auth/signup | Create a new user account                       | No             |
| POST   | /auth/login  | Authenticate and receive JWT token              | No             |
| GET    | /user/list   | List users with pagination (page & limit query) | Yes (JWT)      |

Example requests can be found in the `requests.http` file, which can be used with REST client extensions in various IDEs.

## Project Structure

The codebase follows a domain-driven design with clear separation of concerns and standard Go project layout:

```
├── internal/           # Application-specific code
│   ├── auth/           # Authentication domain
│   │   ├── handler/    # HTTP request handlers
│   │   ├── request/    # Request validation
│   │   └── service/    # Business logic
│   ├── middleware/     # HTTP middleware components
│   ├── user/           # User domain
│   │   ├── entity/     # Data models
│   │   ├── handler/    # HTTP request handlers
│   │   ├── repository/ # Data access layer
│   │   └── service/    # Business logic
│   └── util/           # Shared utilities
├── migrations/         # Database migrations
└── main.go             # Application entry point
```

The project uses the `internal` package pattern to indicate code that is private to this application and not intended for reuse by other packages.
