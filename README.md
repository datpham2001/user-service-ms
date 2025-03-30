# User Service Microservice 🚀

A Go microservice designed to handle user authentication, including registration, standard login, and Google OAuth integration.

## ✨ Features

*   User registration with email and password.
*   Password validation (length, uppercase, lowercase, number).
*   Email format validation.
*   Google OAuth 2.0 for login/signup.
*   REST API for authentication endpoints.
*   Configuration management via YAML files and environment variables.
*   Structured logging with Logrus.
*   Custom error handling.
*   Integration with database and cache (specific types determined by configuration).
*   Potential for gRPC interface (basic setup present).

## 🛠️ Tech Stack

*   **Language:** Go <img src="https://raw.githubusercontent.com/devicons/devicon/master/icons/go/go-original.svg" alt="Go" width="20" height="20"/>
*   **Web Framework:** Gin (`gin-gonic/gin`)
*   **Configuration:** Viper (`spf13/viper`)
*   **Logging:** Logrus (`sirupsen/logrus`)
*   **OAuth:** `golang.org/x/oauth2` (Google)
*   **gRPC:** `google.golang.org/grpc`
*   **Database:** (Requires configuration - e.g., PostgreSQL <img src="https://raw.githubusercontent.com/devicons/devicon/master/icons/postgresql/postgresql-original.svg" alt="PostgreSQL" width="20" height="20"/>) - Managed via `pkg/database`
*   **Cache:** (Requires configuration - e.g., Redis <img src="https://raw.githubusercontent.com/devicons/devicon/master/icons/redis/redis-original.svg" alt="Redis" width="20" height="20"/>) - Managed via `pkg/cache`

*(Consider using [Shields.io](https://shields.io/) for more dynamic/professional badges here)*

## ⚙️ Prerequisites

*   Go (version 1.20 or higher recommend)
*   Access to a database instance (e.g., PostgreSQL)
*   Access to a cache instance (e.g., Redis)
*   Google OAuth 2.0 Credentials (Client ID and Client Secret)

## 🔧 Configuration

Configuration is loaded from `config/config.yaml` and can be overridden by environment variables.

1.  **Create `config/config.yaml`:** Start with a basic structure based on the `config.Config` struct (details not fully shown, but likely includes database, cache, server port, etc.).
2.  **Environment Variables:** Set the following environment variables, especially for sensitive data:
    *   `GOOGLE_CLIENT_ID`: Your Google Cloud project OAuth Client ID.
    *   `GOOGLE_CLIENT_SECRET`: Your Google Cloud project OAuth Client Secret.
    *   Database connection details (e.g., `DB_HOST`, `DB_PORT`, `DB_USER`, `DB_PASSWORD`, `DB_NAME`). Viper uses `.` -> `_` replacement, so a config `database.host` becomes `DATABASE_HOST`.
    *   Cache connection details (e.g., `CACHE_ADDR`, `CACHE_PASSWORD`).
    *   `HTTP_PORT` / `GRPC_PORT` (if defined in config struct).

*Note: The Google Redirect URI is currently hardcoded in `internal/service/auth_service.go` as `http://localhost:8080/auth/google/callback`. Ensure this matches your Google Cloud OAuth configuration or update the code to use a configuration value.*

## 📦 Installation

```bash
# Clone the repository
git clone <your-repository-url>
cd user-service-ms

# Install dependencies
go mod tidy
```

## 🔨 Building

```bash
go build -o user-service ./cmd/main.go
```

## ▶️ Running the Service

```bash
# Ensure required environment variables are set or config.yaml is present
./user-service
```

Alternatively, use `go run`:

```bash
go run ./cmd/main.go
```

The service will start, attempting to connect to the configured database and cache, and listen for HTTP (and potentially gRPC) connections on the configured ports (default seems to be 8080 for HTTP based on the Google callback).

## 🔗 API Endpoints

The following REST endpoints are exposed under the `/api/v1/auth` prefix:

*   `POST /api/v1/auth/signup`: Register a new user with email and password.
    *   **Request Body:**
        ```json
        {
          "email": "test@example.com",
          "password": "Password123"
        }
        ```
    *   **Response:** `201 Created` on success.
*   `POST /api/v1/auth/login`: Log in an existing user (Implementation is currently a placeholder).
*   `GET /auth/google/login`: Initiates the Google OAuth flow (Redirects user to Google).
*   `GET /auth/google/callback`: Callback URL for Google OAuth flow after user grants permission. Handles token exchange and user info retrieval.

*(Note: The Google endpoints `/auth/google/...` might need adjustment based on how the HTTP server and routing are fully configured in `cmd/main.go` - the provided snippets focus on the handlers and API definitions)*

## 👋 Contributing

Contributions are welcome! Please follow standard Go practices and ensure code is formatted (`gofmt`) and tested.

## 📜 License

(Specify your license here, e.g., MIT, Apache 2.0, or proprietary)