# Gmoapi - Movie Management API

RESTful API movies management with comprehensive user authentication, role-based authorization, rate limiting, and email notifications.

## 🚀 Live Demo

- **API Base URL**: https://gmoapi.ucokman.web.id/v1/
- **Swagger Documentation**: https://gmoapi.ucokman.web.id/swagger/index.html

## ✨ Features

- **Full CRUD Operations** for movies with pagination, filtering, and sorting
- **User Management** with registration, email verification, and password reset
- **Token-based Authentication** using Bearer tokens
- **Role-based Access Control (RBAC)** with granular permissions
- **Rate Limiting** (2 req/s, burst: 4) to prevent API abuse
- **Email Notifications** for account activation and password reset
- **CORS Support** for cross-origin requests
- **Optimistic Locking** to prevent concurrent modification conflicts
- **Graceful Shutdown** with background task completion

## 🔌 API Endpoints

### Health Check

- `GET /v1/` - Check API health and version

### Movies

- `GET /v1/movies` - List all movies (with filtering, pagination, sorting)
- `GET /v1/movies/:id` - Get movie by ID

- `POST /v1/movies` - Create a new movie (require movies:write permissions)
- `PATCH /v1/movies/:id` - Update movie (require movies:write permissions)
- `DELETE /v1/movies/:id` - Delete movie (require movies:write permissions)

### Users

- `POST /v1/users/register` - Register new user
- `PUT /v1/users/activated` - Activate user account
- `PUT /v1/users/password` - Reset user password

### Tokens

- `POST /v1/tokens/authentication` - Generate authentication token
- `POST /v1/tokens/password-reset` - Request password reset token
- `POST /v1/tokens/activation` - Resend activation token

### Metrics

- `GET /debug/vars` - View application metrics (requires metrics:read permission)

## 🔐 Authentication

Most endpoints require authentication. To authenticate:

1. **Register a new account**:

```bash
curl -X POST https://gmoapi.ucokman.web.id/v1/users/register \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John Doe",
    "email": "john@example.com",
    "password": "your-secure-password"
  }'
```

2. **Activate your account** using the token sent to your email:

```bash
curl -X PUT https://gmoapi.ucokman.web.id/v1/users/activated \
  -H "Content-Type: application/json" \
  -d '{"token": "YOUR_ACTIVATION_TOKEN"}'
```

3. **Get an authentication token**:

```bash
curl -X POST https://gmoapi.ucokman.web.id/v1/tokens/authentication \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john@example.com",
    "password": "your-secure-password"
  }'
```

4. **Use the token** in subsequent requests:

```bash
curl -X GET https://gmoapi.ucokman.web.id/v1/movies \
  -H "Authorization: Bearer YOUR_TOKEN_HERE"
```

Here's a fixed and improved version of your README:

## 🏁 Run Locally

### 🧰 Prerequisites

Make sure you have the following installed:

- **Linux OS**
- **Make**
- **Docker**
- **Go** `v1.21` or higher
- **SMTP server credentials** (for email notifications)
- **Air** – for live reloading: https://github.com/air-verse/air

### ⚙️ Installation

1. **Clone the repository**

   ```bash
   git clone https://github.com/ucok-man/gmoapi.git
   cd gmoapi
   ```

2. **Install Go dependencies**

   ```bash
   go mod download
   ```

3. **Start the database**

   ```bash
   # See database credentials in the docker-compose file
   # PostgreSQL will run on port 5433
   docker-compose up -d
   ```

4. **Create a `.env` file**

   Follow the Configuration section for reference.

5. **Run the application in development mode**

   ```bash
   make dev
   ```

### 🧑‍💻 Database Setup for Mutations & Metrics

To perform data mutations (create, update, delete) or view metrics, follow these steps:

1. **Connect to the PostgreSQL container**

   ```bash
   docker compose exec gmoapi_postgres psql -U gmoapi -d gmoapi
   ```

2. **Create an admin user**

   ```sql
   INSERT INTO users (name, email, password_hash, activated)
   VALUES (
       '<username>',
       '<email>',
       crypt('<password>', gen_salt('bf', 12))::bytea,
       true
   );
   ```

   💡 **Note:** The password is hashed using PostgreSQL's `crypt()` with a bcrypt salt, stored as `bytea`.

3. **Assign all permissions**

   ```sql
   INSERT INTO users_permissions
   SELECT u.id, p.id
   FROM users u, permissions p
   WHERE u.email = '<email>'
   AND p.code IN ('movies:read', 'movies:write', 'metrics:read');
   ```

### 📜 Available Commands

To view all supported `make` commands:

```bash
make help
```

## ⚙️ Configuration

The application can be configured using environment variables or command-line flags. Environment variables take precedence over default values, and command-line flags override everything.

### Environment Variables

All environment variables use the `GMOAPI_` prefix:

```bash
# Server Configuration
GMOAPI_HOST=localhost
GMOAPI_PORT=4000
GMOAPI_ENV=development  # development|staging|production

# Database Configuration
GMOAPI_DB_DSN=postgres://gmoapi:pa55word@localhost:5433/gmoapi?sslmode=disable
GMOAPI_DB_MAX_OPEN_CONN=25
GMOAPI_DB_MAX_IDLE_CONN=25
GMOAPI_DB_MAX_IDLE_TIME=15m

# Rate Limiter Configuration
GMOAPI_LIMITER_RPS=2
GMOAPI_LIMITER_BURST=4
GMOAPI_LIMITER_ENABLED=true

# SMTP Configuration
GMOAPI_SMTP_HOST=smtp.example.com
GMOAPI_SMTP_PORT=587
GMOAPI_SMTP_USERNAME=your-username
GMOAPI_SMTP_PASSWORD=your-password
GMOAPI_SMTP_SENDER="Gmoapi <no-reply@gmoapi.ucokman.web.id>"

# CORS Configuration
GMOAPI_CORS_TRUSTED_ORIGINS=https://example.com,https://app.example.com
```

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## 📄 License

This project is licensed under the Apache License 2.0 - see the [LICENSE](http://www.apache.org/licenses/LICENSE-2.0.html) file for details.
