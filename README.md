# JWT Auth API Demo

This service provides a secure and efficient RESTful API for managing **consumer data**, built with **Go (Gin framework)** and **PostgreSQL**. It implements **JWT-based authentication and authorization**, ensuring that only authenticated clients can access protected endpoints.

---


## ✨ Features

This application delivers a simple yet robust JWT-based authentication flow, with PostgreSQL as the primary datastore. Below is an overview of the core features:

### 🔐 JWT Authentication

- Full **JWT authentication** system: 
  - Protects API endpoints using **JSON Web Tokens (JWT)**.
  - All secured requests must include a valid JWT in the `Authorization` header (`Bearer <token>`).
  - The token is verified before any processing happens, and requests with missing or invalid tokens receive a `401 Unauthorized` response.
  - Ensure only authenticated users can interact with the transaction API.
  - Helps trace and attribute transactions to specific authenticated consumers or services.

- **Authentication Endpoints**:
  - `POST /auth/login` — Authenticates user with `username` and `password`, returns, returns:
    - `AccessToken`
    - `RefreshToken`
    - `ExpirationDate`
    - `TokenType`
  - `POST /auth/refresh-token` — Accepts a valid `RefreshToken` and issues a new `AccessToken`.

- **RSA key pairs** are used to sign and verify tokens (more secure than symmetric secrets)
  - Stored in `/keys` directory: `privateKey.pem` and `publicKey.pem`
  - Keys are generated using `OpenSSL`


### 🛡️ Security & Middleware

The service is designed with security and extensibility in mind, using several middlewares:

- **Authorization Middleware**:
  - Validates JWT
  - Enforces Role-Based Access Control (RBAC)

- **Security Headers Middleware**:
  - CORS
  - Secure HTTP headers (e.g., `X-Frame-Options`, `X-Content-Type-Options`, etc.)


### 🗄️ Logging

- Uses `github.com/sirupsen/logrus` for structured logging
- Integrates with `gopkg.in/natefinch/lumberjack.v2` for automatic log rotation based on size and age
- Logs are separated by level: **info**, **request**, **warn**, **error**, **fatal**, and **panic**


---

## 🧭 Business Process Flow

The following flow illustrates how clients interact with the API for **authentication** and **consumer management**. The system enforces **JWT-based authorization** and **role-based access control (RBAC)** to protect resources.

```pgsql
┌──────────────────────────────────────────────┐
│           [1] User Logs In                   │
│----------------------------------------------│
│ - POST /auth/login                           │
│ - Body: { username, password }               │
└──────────────────────────────────────────────┘
              │
              ▼
┌──────────────────────────────────────────────┐
│       [2] JWT Token Generation               │
│----------------------------------------------│
│ - If credentials valid:                      │
│   → issue AccessToken & RefreshToken         │
│ - If invalid: respond with 401               │
└──────────────────────────────────────────────┘
              │
              ▼
┌──────────────────────────────────────────────┐
│ [3] Access Protected Consumer API Endpoint   │
│----------------------------------------------│
│ - Include AccessToken in header:             │
│   Authorization: Bearer <access_token>       │
└──────────────────────────────────────────────┘
              │
              ▼
┌──────────────────────────────────────────────┐
│ [4] Middleware: JWT & Role Authorization     │
│----------------------------------------------│
│ - Validate JWT signature & expiry            │
│ - Check user's role (e.g., ROLE_ADMIN)       │
│ - If unauthorized: respond 403 or 401        │
└──────────────────────────────────────────────┘
              │
              ▼
┌──────────────────────────────────────────────┐
│    [5] Handler Executes Requested Action     │
│----------------------------------------------│
│ - GET /consumers → list all (ADMIN/USER)     │
│ - GET /consumers/:id → detail (ADMIN/USER)   │
│ - GET /consumers/active|inactive|suspended   │
│ - POST /consumers → create (ADMIN only)      │
│ - PATCH /consumers/:id → update status       │
└──────────────────────────────────────────────┘

```
---


## 🤖 Tech Stack

This project leverages a modern and robust set of technologies to ensure performance, security, and maintainability. Below is an overview of the core tools and libraries used in the development:

| **Component**             | **Description**                                                                             |
|---------------------------|---------------------------------------------------------------------------------------------|
| **Language**              | Go (Golang), a statically typed, compiled language known for concurrency and efficiency     |
| **Web Framework**         | Gin, a fast and minimalist HTTP web framework for Go                                        |
| **ORM**                   | GORM, an ORM library for Go supporting SQL and migrations                                   |
| **Database**              | PostgreSQL, a powerful open-source relational database system                               |
| **JWT Signing**           | RSA asymmetric key pairs generated via OpenSSL, used to securely sign and verify JWT tokens |
| **Logging**               | Logrus for structured logging, combined with Lumberjack for log rotation                    |
| **Validation**            | `go-playground/validator.v9` for input validation and data integrity enforcement            |

---

## 🧱 Architecture Overview

This project follows a **modular** and **maintainable** architecture inspired by **Clean Architecture** principles. Each domain feature (e.g., **entity**, **handler**, **repository**, **service**) is organized into self-contained modules with clear separation of concerns.

```bash
📁 go-jwt-auth-demo/
├── 📂cert/                                 # Stores self-signed TLS certificates used for local development (e.g., for HTTPS or JWT signing verification)
├── 📂cmd/                                  # Contains the application's entry point.
├── 📂config/
│   └── 📂database/                         # Config for PostgreSQL (DSN, pool settings, migration, etc.)
├── 📂docker/                               # Docker-related configuration for building and running services
│   ├── 📂app/                              # Contains Dockerfile to build the main Go application image
│   └── 📂postgres/                         # Contains PostgreSQL container configuration
├── 📂internal/                             # Core domain logic and business use cases, organized by module
│   ├── 📂entity/                           # Data models/entities representing business concepts like Transaction, Consumer
│   ├── 📂handler/                          # HTTP handlers (controllers) that parse requests and return responses
│   ├── 📂repository/                       # Data access layer, communicating with DB or cache
│   └── 📂service/                          # Business logic layer orchestrating operations between handlers and repositories
├── 📂keys/                                 # Contains RSA public/private keys used for signing and verifying JWT tokens
├── 📂logs/                                 # Application log files (error, request, info) written and rotated using Logrus + Lumberjack
├── 📂pkg/                                  # Reusable utility and middleware packages shared across modules
│   ├── 📂contextdata/                      # Stores and retrieves contextual data like User Information
│   ├── 📂customtype/                       # Defines custom types, enums, constants used throughout the application
│   ├── 📂diagnostics/                      # Health check endpoints, metrics, and diagnostics handlers for monitoring
│   ├── 📂logger/                           # Centralized log initialization and configuration
│   ├── 📂middleware/                       # Request processing middleware
│   │   ├── 📂authorization/                # JWT validation and Role-Based Access Control (RBAC)
│   │   ├── 📂headers/                      # Manages request headers like CORS, security, request ID
│   │   └── 📂logging/                      # Logs incoming requests
│   └── 📂util/                             # General utility functions and helpers
│       ├── 📂http-util/                    # Utilities for common HTTP tasks (e.g., write JSON, status helpers)
│       ├── 📂jwt-util/                     # Token generation, parsing, and validation logic
│       └── 📂validation-util/              # Common input validators (e.g., UUID, numeric range)
├── 📂routes/                               # Route definitions, groups APIs, and applies middleware per route scope
└── 📂tests/                                # Contains unit or integration tests for business logic
```

---

## 🛠️ Installation & Setup  

Follow the instructions below to get the project up and running in your local development environment. You may run it natively or via Docker depending on your preference.  

### ✅ Prerequisites

Make sure the following tools are installed on your system:

| **Tool**                                                      | **Description**                           |
|---------------------------------------------------------------|-------------------------------------------|
| [Go](https://go.dev/dl/)                                      | Go programming language (v1.20+)          |
| [Make](https://www.gnu.org/software/make/)                    | Build automation tool (`make`)            |
| [PostgreSQL](https://www.postgresql.org/)                     | Relational database system (v14+)         |
| [Docker](https://www.docker.com/)                             | Containerization platform (optional)      |

### 🔁 Clone the Project  

Clone the repository:  

```bash
git clone https://github.com/yoanesber/Go-JWT-Auth-Demo.git
cd Go-JWT-Auth-Demo
```

### ⚙️ Configure `.env` File  

Set up your **database** and **JWT configuration** in `.env` file. Create a `.env` file at the project root directory:  

```properties
# Application configuration
ENV=PRODUCTION
API_VERSION=1.0
PORT=1000
IS_SSL=TRUE
SSL_KEYS=./cert/mycert.key
SSL_CERT=./cert/mycert.cer

# Database configuration
DB_HOST=localhost
DB_PORT=5432
DB_USER=appuser
DB_PASS=app@123
DB_NAME=consumer_service
DB_SCHEMA=public
# Options: disable, require, verify-ca, verify-full
DB_SSL_MODE=disable
DB_TIMEZONE=Asia/Jakarta
DB_MIGRATE=TRUE
DB_SEED=TRUE
DB_SEED_FILE=import.sql
# Set to INFO for development and staging, SILENT for production
DB_LOG=SILENT

# JWT configuration
JWT_SECRET=a-string-secret-at-least-256-bits-long
# 2 days
JWT_EXPIRATION_HOUR=48
JWT_ISSUER=your_jwt_issuer
JWT_AUDIENCE=your_jwt_audience
# 30 days
JWT_REFRESH_TOKEN_EXPIRATION_HOUR=720
JWT_PRIVATE_KEY_PATH=./keys/privateKey.pem
JWT_PUBLIC_KEY_PATH=./keys/publicKey.pem
# RS256 or HS256
JWT_ALGORITHM=RS256
# Bearer or JWT
TOKEN_TYPE=Bearer

```

- **🔐 Notes**:  
  - `IS_SSL=TRUE`: Enable this if you want your app to run over `HTTPS`. Make sure to run `generate-certificate.sh` to generate **self-signed certificates** and place them in the `./cert/` directory (e.g., `mycert.key`, `mycert.cer`).
  - `JWT_ALGORITHM=RS256`: Set this if you're using **asymmetric JWT signing**. Be sure to run `generate-jwt-key.sh` to generate **RSA key pairs** and place `privateKey.pem` and `publicKey.pem` in the `./keys/` directory.
  - Make sure your paths (`./cert/`, `./keys/`) exist and are accessible by the application during runtime.
  - `DB_TIMEZONE=Asia/Jakarta`: Adjust this value to your local timezone (e.g., `America/New_York`, etc.).
  - `DB_MIGRATE=TRUE`: Set to `TRUE` to automatically run `GORM` migrations for all entity definitions on app startup.
  - `DB_SEED=TRUE` & `DB_SEED_FILE=import.sql`: Use these settings if you want to insert predefined data into the database using the SQL file provided.
  - `DB_USER=appuser`, `DB_PASS=app@123`: It's strongly recommended to create a dedicated database user instead of using the default postgres superuser.

### 🔑 Generate RSA Key for JWT (If Using `RS256`)  

If you are using `JWT_ALGORITHM=RS256`, generate the **RSA key** pair for **JWT signing** by running this file:  
```bash
./generate-jwt-key.sh
```

- **Notes**:  
  - On **Linux/macOS**: Run the script directly
  - On **Windows**: Use **WSL** to execute the `.sh` script

This will generate:
  - `./keys/privateKey.pem`
  - `./keys/publicKey.pem`


These files will be referenced by your `.env`:
```properties
JWT_PRIVATE_KEY_PATH=./keys/privateKey.pem
JWT_PUBLIC_KEY_PATH=./keys/publicKey.pem
JWT_ALGORITHM=RS256
```

### 🔐 Generate Certificate for HTTPS (Optional)  

If `IS_SSL=TRUE` in your `.env`, generate the certificate files by running this file:  
```bash
./generate-certificate.sh
```

- **Notes**:  
  - On **Linux/macOS**: Run the script directly
  - On **Windows**: Use **WSL** to execute the `.sh` script

This will generate:
  - `./cert/mycert.key`
  - `./cert/mycert.cer`


Ensure these are correctly referenced in your `.env`:
```properties
IS_SSL=TRUE
SSL_KEYS=./cert/mycert.key
SSL_CERT=./cert/mycert.cer
```

### 👤 Create Dedicated PostgreSQL User (Recommended)

For security reasons, it's recommended to avoid using the default postgres superuser. Use the following SQL script to create a dedicated user (`appuser`) and assign permissions:

```sql
-- Create appuser and database
CREATE USER appuser WITH PASSWORD 'app@123';

-- Allow user to connect to database
GRANT CONNECT, TEMP, CREATE ON DATABASE consumer_service TO appuser;

-- Grant permissions on public schema
GRANT USAGE, CREATE ON SCHEMA public TO appuser;

-- Grant all permissions on existing tables
GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA public TO appuser;

-- Grant all permissions on sequences (if using SERIAL/BIGSERIAL ids)
GRANT USAGE, SELECT, UPDATE ON ALL SEQUENCES IN SCHEMA public TO appuser;

-- Ensure future tables/sequences will be accessible too
ALTER DEFAULT PRIVILEGES IN SCHEMA public
GRANT SELECT, INSERT, UPDATE, DELETE ON TABLES TO appuser;

-- Ensure future sequences will be accessible too
ALTER DEFAULT PRIVILEGES IN SCHEMA public
GRANT USAGE, SELECT, UPDATE ON SEQUENCES TO appuser;
```

Update your `.env` accordingly:
```properties
DB_USER=appuser
DB_PASS=app@123
```

---


## 🚀 Running the Application  

This section provides step-by-step instructions to run the application either **locally** or via **Docker containers**.

- **Notes**:  
  - All commands are defined in the `Makefile`.
  - To run using `make`, ensure that `make` is installed on your system.
  - To run the application in containers, make sure `Docker` is installed and running.
  - Ensure you have `Go` installed on your system

### 📦 Install Dependencies

Make sure all Go modules are properly installed:  

```bash
make tidy
```

### 🧪 Run Unit Tests

```bash
make test
```

### 🔧 Run Locally (Non-containerized)

Ensure PostgreSQL are running locally, then:

```bash
make run
```

### 🐳 Run Using Docker

To build and run all services (PostgreSQL, Go app):

```bash
make docker-up
```

To stop and remove all containers:

```bash
make docker-down
```

- **Notes**:  
  - Before running the application inside Docker, make sure to update your environment variables `.env`
    - Change `DB_HOST=localhost` to `DB_HOST=postgres-server`.

### 🟢 Application is Running

Now your application is accessible at:
```bash
http://localhost:1000
```

or 

```bash
https://localhost:1000 (if SSL is enabled)
```

---

## 🧪 Testing Scenarios  

### 🔐 Login API

**Endpoint**: `POST https://localhost:1000/auth/login`

#### ✅ Scenario 1: Successful Login

**Request**:

```json
{
  "username": "admin",
  "password": "P@ssw0rd"
}
```

**Response**:

```json
{
  "message": "Login successful",
  "error": null,
  "path": "/auth/login",
  "status": 200,
  "data": {
    "accessToken": "<JWT>",
    "refreshToken": "<UUID>",
    "expirationDate": "2025-05-25T12:58:00Z",
    "tokenType": "Bearer"
  },
  "timestamp": "2025-05-23T12:58:00Z"
}
```

#### ❌ Scenario 2: Invalid Credentials

**Request with invalid user**:
```json
{
  "username": "invalid_user",
  "password": "P@ssw0rd"
}
```

**Response**:
```json
{
  "message": "Failed to login",
  "error": "user with the given username not found",
  "path": "/auth/login",
  "status": 401,
  "data": null,
  "timestamp": "2025-05-23T15:18:23Z"
}
```

**Request with invalid password**:
```json
{
  "username": "admin",
  "password": "invalid_password"
}
```

**Response**:
```json
{
    "message": "Failed to login",
    "error": "invalid password",
    "path": "/auth/login",
    "status": 401,
    "data": null,
    "timestamp": "2025-05-23T15:51:39.288150079Z"
}
```

#### 🚫 Scenario 3: Disabled User

Precondition:
```sql
UPDATE users SET is_enabled = false WHERE id = 2;
```

**Request**:
```json
{
  "username": "userone",
  "password": "P@ssw0rd"
}
```

**Response**:
```json
{
  "message": "Failed to login",
  "error": "user is not enabled",
  "path": "/auth/login",
  "status": 401,
  "data": null,
  "timestamp": "2025-05-23T15:19:24Z"
}
```


### 🔄 Refresh Token API

**Endpoint**: `POST https://localhost:1000/auth/refresh-token`

#### ✅ Scenario 1: Successful Refresh Token

**Request**:
```json
{
  "refreshToken": "<valid_refresh_token>"
}
```

**Response**:
```json
{
  "message": "Token refreshed successfully",
  "error": null,
  "path": "/auth/refresh-token",
  "status": 200,
  "data": {
    "accessToken": "<JWT>",
    "refreshToken": "<new_UUID>",
    "expirationDate": "2025-05-25T15:23:51Z",
    "tokenType": "Bearer"
  },
  "timestamp": "2025-05-23T15:23:51Z"
}
```

#### ❌ Scenario 2: Invalid Refresh Token

**Request**:
```json
{
  "refreshToken": "<invalid_refresh_token>"
}
```

**Response**:
```json
{
  "message": "Failed to refresh token",
  "error": "record not found",
  "path": "/auth/refresh-token",
  "status": 401,
  "data": null,
  "timestamp": "2025-05-23T15:24:47Z"
}
```

#### 🔁 Scenario 3: Expired Refresh Token (Auto Regenerate)

**Request**:
```json
{
  "refreshToken": "<expired_refresh_token>"
}
```

**Response**:
```json
{
  "message": "Token refreshed successfully",
  "error": null,
  "path": "/auth/refresh-token",
  "status": 200,
  "data": {
    "accessToken": "<new_JWT>",
    "refreshToken": "<new_UUID>",
    "expirationDate": "2025-05-25T15:29:02Z",
    "tokenType": "Bearer"
  },
  "timestamp": "2025-05-23T15:29:02Z"
}
```

### 👨‍👩‍👧‍👦 Consumer API

All requests below must include a valid JWT token in the `Authorization` header:
```http
Authorization: Bearer <valid_token>
```

#### Scenario 1: Create Consumer

**Endpoint**: 
```http
POST https://localhost:1000/api/v1/consumers
```

**Request**:
```json
{
    "fullname": "Austin Libertus",
    "username": "auslibertus",
    "email": "austin.libertus@example.com",
    "phone": "+628997452753",
    "address": "Jl. Anggrek No. 4, Jakarta",
    "birthDate": "1990-03-05"
}
```

**Response**:
```json
{
    "message": "Consumer created successfully",
    "error": null,
    "path": "/api/v1/consumers",
    "status": 201,
    "data": {
        "id": "4c6c42bc-3b82-4f34-9eaf-c4dcfb246ec0",
        "fullname": "Austin Libertus",
        "username": "auslibertus",
        "email": "austin.libertus@example.com",
        "phone": "628997452753",
        "address": "Jl. Anggrek No. 4, Jakarta",
        "birthDate": "1990-03-05",
        "status": "inactive",
        "createdAt": "2025-06-18T11:42:13.165068Z",
        "updatedAt": "2025-06-18T11:42:13.165068Z"
    },
    "timestamp": "2025-06-18T11:42:13.171205664Z"
}
```

#### Scenario 2: Update Consumer Status

**Endpoint**: 
```http
PATCH https://localhost:1000/api/v1/consumers/4c6c42bc-3b82-4f34-9eaf-c4dcfb246ec0?status=active
```

**Response**:
```json
{
    "message": "Consumer status updated successfully",
    "error": null,
    "path": "/api/v1/consumers/4c6c42bc-3b82-4f34-9eaf-c4dcfb246ec0",
    "status": 200,
    "data": {
        "id": "4c6c42bc-3b82-4f34-9eaf-c4dcfb246ec0",
        "fullname": "Austin Libertus",
        "username": "auslibertus",
        "email": "austin.libertus@example.com",
        "phone": "628997452753",
        "address": "Jl. Anggrek No. 4, Jakarta",
        "birthDate": "1990-03-05",
        "status": "active",
        "createdAt": "2025-06-18T11:42:13.165068Z",
        "updatedAt": "2025-06-18T11:44:52.059458364Z"
    },
    "timestamp": "2025-06-18T11:44:52.062880937Z"
}
```

#### Scenario 3: Get All Consumers

**Endpoint**: 
```http
GET https://localhost:1000/api/v1/consumers?page=1&limit=10
```

**Response**:
```json
{
    "message": "All consumers retrieved successfully",
    "error": null,
    "path": "/api/v1/consumers",
    "status": 200,
    "data": [
        {
            "id": "74fe86f3-6324-42c2-97b4-fa3225461299",
            "fullname": "John Doe",
            "username": "johndoe",
            "email": "john.doe@example.com",
            "phone": "6281234567890",
            "address": "Jl. Merdeka No. 123, Jakarta",
            "birthDate": "1990-05-10",
            "status": "active",
            "createdAt": "2025-06-18T11:40:56.66591Z",
            "updatedAt": "2025-06-18T11:40:56.66591Z"
        }
        ...
    ],
    "timestamp": "2025-06-18T13:11:24.539972654Z"
}
```
