# IMP Assessment

A RESTful API service built with Go that provides user authentication and management. This project demonstrates clean architecture principles with a domain-driven design approach and includes:

- User authentication (login/signup) with JWT token
- User management with pagination support
- Middleware for protected routes
- PostgreSQL database integration

## Requirements

This project is developed with:

- Go 1.20

- Postgres 15

## Installation

Clone the project

```bash
git clone git@github.com:appleinautumn/imp-assessment.git
```

Go to the project directory

```bash
cd imp-assessment
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

If you have not created the database for IMP Assessment, please create one before going to the next step.

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
go build -o imp
```

Run it

```bash
./imp
```

## API Endpoints

| Method | Endpoint     | Description                                     | Authentication |
| ------ | ------------ | ----------------------------------------------- | -------------- |
| GET    | /            | Root endpoint (health check)                    | No             |
| POST   | /auth/signup | Create a new user account                       | No             |
| POST   | /auth/login  | Authenticate and receive JWT token              | No             |
| GET    | /user/list   | List users with pagination (page & limit query) | Yes (JWT)      |

Example requests can be found in the `requests.http` file, which can be used with REST client extensions in various IDEs.
