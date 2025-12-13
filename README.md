# SubsAPI

SubsAPI is a subscription management API built with Go (Golang), Gin, and GORM, using PostgreSQL. It allows creating, reading, updating, deleting, and calculating subscription costs.

# Features

- CRUD operations for subscriptions
- User-specific subscription cost calculation
- Handles subscription start and end dates (Month-Year format)
- Structured logging with Logrus
- Swagger documentation for all endpoints
- Dockerized for easy setup



# INSTALLATION

1. Clone the repository:

```bash

git clone https://github.com/cyb3rkh4l1d/subsapi.git
cd subsapi

```

2. Prepare environment variables:

```bash

mv env.example .env

```

3. Edit .env with your database credentials, for example:

```ini

APP_PORT=:8080
DB_HOST=postgres
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=admin
DB_NAME=subscriptions_db
DB_SSLMODE=disable
GIN_MODE=release


```

4. Start the application using Docker Compose:

```bash

sudo docker-compose --env-file .env -f deployment/docker-compose.yaml up --build
```

# API Endpoints

```bash
POST   /api/v1/subscriptions/        Create a new subscription
GET    /api/v1/subscriptions/        List all subscriptions
GET    /api/v1/subscriptions/{id}    Get subscription by ID
PUT    /api/v1/subscriptions/{id}    Update subscription by ID
DELETE /api/v1/subscriptions/{id}    Delete subscription by ID
GET    /api/v1/subscriptions/stats?user_id=&service_name=&from=&to=     Calculate total subscription cost for a user
GET    /api/swagger/index.html            Swagger API documentation
```

Visit Swagger Docs endpoints

http://host:port/api/swagger/index.html
