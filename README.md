# Go Google Auth

A full-stack web application for Google OAuth authentication with comprehensive Cookie/Session management testing capabilities.

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Version](https://img.shields.io/badge/Go-1.21%2B-blue)](https://go.dev/)
[![Vue Version](https://img.shields.io/badge/Vue-3.5-green)](https://vuejs.org/)
[![Node Version](https://img.shields.io/badge/Node-22.15.0%2B-brightgreen)](https://nodejs.org/)

## 📋 Table of Contents

- [Overview](#overview)
- [Features](#features)
- [Technology Stack](#technology-stack)
- [Project Structure](#project-structure)
- [Getting Started](#getting-started)
- [Usage](#usage)
- [API Documentation](#api-documentation)
- [Development](#development)
- [Docker Commands](#docker-commands)
- [Environment Variables](#environment-variables)
- [Troubleshooting](#troubleshooting)
- [Contributing](#contributing)
- [License](#license)

## 🎯 Overview

This project demonstrates a modern full-stack application architecture with:
- **Backend**: Go-based REST API using Gin framework
- **Frontend**: Vue.js 3 with TypeScript and Vite
- **Database**: DynamoDB (AWS DynamoDB Local for development)
- **Authentication**: Google OAuth 2.0 (planned)
- **Session Management**: Cookie-based session handling with testing interface

## ✨ Features

### Current Features
- ✅ Cookie/Session testing interface
- ✅ Set-Cookie header validation
- ✅ Cookie transmission verification
- ✅ Real-time cookie display and management
- ✅ CORS configuration for cross-origin requests
- ✅ Docker-based development environment
- ✅ Hot reload for both frontend and backend

### Planned Features
- 🚧 Google OAuth 2.0 authentication
- 🚧 User session management
- 🚧 Protected routes and authorization
- 🚧 AWS deployment with CDK

## 🛠 Technology Stack

### Backend
| Technology | Version | Purpose |
|------------|---------|---------|
| **Go** | 1.21+ | Primary backend language |
| **Gin** | Latest | Web framework |
| **DynamoDB** | Local | NoSQL database |
| **Air** | Latest | Hot reload for development |
| **Docker** | Latest | Containerization |

### Frontend
| Technology | Version | Purpose |
|------------|---------|---------|
| **Vue.js** | 3.5+ | Frontend framework |
| **TypeScript** | 5.8+ | Type-safe JavaScript |
| **Vite** | 6.0+ | Build tool and dev server |
| **Vue Router** | 4.5+ | Client-side routing |
| **Vue DevTools** | 8.0+ | Development tools |

### Infrastructure
| Technology | Purpose |
|------------|---------|
| **Docker Compose** | Multi-container orchestration |
| **AWS CDK** | Infrastructure as Code (TypeScript) |
| **GitHub Actions** | CI/CD pipeline |

## 📁 Project Structure

```
go-google-auth/
│
├── backend/                    # Go backend application
│   ├── cmd/
│   │   └── api/
│   │       └── main.go        # Application entry point
│   ├── internal/
│   │   ├── config/            # Configuration management
│   │   └── handlers/          # HTTP request handlers
│   │       ├── cookie.go      # Cookie management handlers
│   │       ├── health.go      # Health check endpoints
│   │       └── hello.go       # Example endpoints
│   ├── dockers/
│   │   ├── Dockerfile.local   # Development Docker image
│   │   └── Dockerfile.prod    # Production Docker image
│   ├── compose.yml            # Backend Docker Compose config
│   ├── .air.toml             # Air hot reload configuration
│   ├── go.mod                # Go module definition
│   └── Makefile              # Backend build commands
│
├── frontend/                  # Vue.js frontend application
│   ├── vue-app/
│   │   ├── src/
│   │   │   ├── components/   # Reusable Vue components
│   │   │   │   └── AppHeader.vue
│   │   │   ├── views/        # Page-level components
│   │   │   │   ├── HomeView.vue      # Cookie testing page
│   │   │   │   └── AboutView.vue     # About page
│   │   │   ├── router/       # Vue Router configuration
│   │   │   │   └── index.ts
│   │   │   ├── assets/       # Static assets
│   │   │   │   ├── main.css
│   │   │   │   └── styles/   # Organized CSS
│   │   │   │       ├── base/
│   │   │   │       ├── components/
│   │   │   │       ├── layouts/
│   │   │   │       ├── pages/
│   │   │   │       └── utilities/
│   │   │   └── main.ts       # Application entry point
│   │   ├── index.html
│   │   ├── vite.config.ts
│   │   ├── tsconfig.json
│   │   └── package.json
│   ├── Dockerfile            # Frontend Docker image
│   ├── compose.yml           # Frontend Docker Compose config
│   └── docker-entrypoint.sh  # Container startup script
│
├── iac/                       # Infrastructure as Code
│   ├── bin/
│   │   ├── backend.ts        # Backend infrastructure stack
│   │   └── network.ts        # Network infrastructure stack
│   ├── cdk.json
│   └── package.json
│
├── .github/
│   └── workflows/            # GitHub Actions CI/CD
│       ├── backend.yml
│       ├── backend-setup.yml
│       └── network.yml
│
├── Makefile                  # Root-level commands
├── README.md                 # This file
├── .cursorrules             # Development guidelines
├── .gitignore
└── LICENSE
```

## 🚀 Getting Started

### Prerequisites

Ensure you have the following installed:

- **Docker Desktop** (latest version)
- **Docker Compose** (v2.0+)
- **Git**

Optional (for local development without Docker):
- **Go** (1.21+)
- **Node.js** (22.15.0+)
- **npm** (10.9+)

### Installation

1. **Clone the repository**
   ```bash
   git clone https://github.com/yuki5155/go-google-auth.git
   cd go-google-auth
   ```

2. **Start all services using Docker Compose**
   ```bash
   # Start backend
   cd backend && docker compose up -d
   
   # Start frontend (in a new terminal)
   cd frontend && docker compose up -d
   ```

3. **Verify services are running**
   ```bash
   docker ps | grep go-google-auth
   ```

   You should see three containers running:
   - `go-google-auth-frontend` (Port 5173)
   - `go-google-auth-app` (Port 8080)
   - `go-google-auth-dynamodb` (Port 8000)

4. **Access the application**
   - Frontend: http://localhost:5173
   - Backend API: http://localhost:8080
   - DynamoDB Local: http://localhost:8000

## 💡 Usage

### Cookie/Session Testing

The application provides an interactive interface for testing cookie behavior:

1. **Open the application**
   - Navigate to http://localhost:5173 in your browser

2. **Test Set-Cookie**
   - Click the **"Test Set-Cookie"** button
   - The backend will send a `Set-Cookie` header
   - Check the response in the UI

3. **View Current Cookies**
   - The **"Current Cookies"** section displays all cookies stored in your browser
   - Click **"Refresh"** to update the display

4. **Test Cookie Sending**
   - Click **"Check Cookie Sending"** to verify cookies are sent to the backend
   - The backend will confirm receipt and display cookie contents

5. **Clear Cookies**
   - Use the **"Clear All Cookies"** button to remove test cookies

### Backend URL Configuration

You can change the backend API URL in the frontend interface if your backend is running on a different port or host.

## 📚 API Documentation

### Health Endpoints

#### `GET /health`
Health check endpoint for monitoring.

**Response:**
```json
{
  "status": "ok"
}
```

#### `GET /health/ready`
Readiness probe for Kubernetes/ECS deployments.

### Cookie Testing Endpoints

#### `GET /api/set-cookie`
Sets a test cookie in the response.

**Headers:**
- `Access-Control-Allow-Origin`: Frontend URL
- `Access-Control-Allow-Credentials`: true

**Response:**
```json
{
  "message": "Cookie has been set",
  "cookieName": "test_session",
  "cookieValue": "session_value_12345",
  "timestamp": "2025-12-14T10:09:50Z"
}
```

**Set-Cookie Header:**
```
test_session=session_value_12345; Path=/; Max-Age=3600; HttpOnly
```

#### `GET /api/check-cookie`
Verifies if cookies are received from the client.

**Required Headers:**
- `Cookie`: Must include cookies from browser

**Response:**
```json
{
  "cookieReceived": true,
  "cookies": "test_session=session_value_12345",
  "testSession": "session_value_12345",
  "timestamp": "2025-12-14T10:10:00Z"
}
```

### Future Endpoints (Planned)

- `GET /auth/google` - Initiate Google OAuth flow
- `GET /auth/google/callback` - OAuth callback handler
- `GET /auth/logout` - User logout
- `GET /api/user` - Get authenticated user info

## 🔧 Development

### Backend Development

#### Running Locally (with Air hot reload)
```bash
cd backend

# Install dependencies
go mod download

# Run with hot reload
make dev

# Run tests
make test

# Build binary
make build
```

#### Project Structure
```go
// Handler pattern
type Handler struct {
    Path string
}

func NewHandler() *Handler {
    return &Handler{
        Path: "/api/endpoint",
    }
}

func (h *Handler) Handle(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{
        "message": "success",
    })
}
```

### Frontend Development

#### Running Locally
```bash
cd frontend/vue-app

# Install dependencies
npm install

# Start dev server
npm run dev

# Build for production
npm run build

# Run linter
npm run lint

# Type check
npm run type-check
```

#### Component Structure
```vue
<script setup lang="ts">
import { ref } from 'vue'

interface Props {
  title: string
}

const props = defineProps<Props>()
const count = ref(0)

const increment = () => {
  count.value++
}
</script>

<template>
  <div>
    <h1>{{ props.title }}</h1>
    <button @click="increment">Count: {{ count }}</button>
  </div>
</template>

<style scoped>
/* Component-specific styles */
</style>
```

## 🐳 Docker Commands

### Basic Operations

```bash
# Start all services
cd backend && docker compose up -d
cd frontend && docker compose up -d

# Stop all services
cd backend && docker compose down
cd frontend && docker compose down

# Stop with volume cleanup
docker compose down -v

# View logs
cd backend && docker compose logs -f
cd frontend && docker compose logs -f

# Rebuild containers
cd frontend && docker compose build --no-cache
```

### Container Management

```bash
# List running containers
docker ps

# List all containers
docker ps -a

# Inspect a container
docker inspect go-google-auth-frontend

# Execute command in container
docker exec -it go-google-auth-frontend bash

# Remove stopped containers
docker container prune

# View container logs
docker logs go-google-auth-app
```

### Volume Management

```bash
# List volumes
docker volume ls

# Remove unused volumes
docker volume prune

# Inspect volume
docker volume inspect frontend_node_modules
```

## 🔐 Environment Variables

### Backend Configuration

Create a `.env` file in the `backend/` directory or set environment variables:

```env
# Application
GO_ENV=development
PORT=8080

# AWS
AWS_REGION=ap-northeast-1
DYNAMODB_ENDPOINT=http://dynamodb:8000
AWS_ACCESS_KEY_ID=dummy
AWS_SECRET_ACCESS_KEY=dummy

# Google OAuth (required for authentication)
GOOGLE_CLIENT_ID=your_google_client_id
GOOGLE_CLIENT_SECRET=your_google_client_secret
GOOGLE_REDIRECT_URL=http://localhost:8080/auth/google/callback
```

### Frontend Configuration

Configuration in `frontend/vue-app/vite.config.ts`:

```typescript
export default defineConfig({
  server: {
    host: '0.0.0.0',
    port: 5173,
    watch: {
      usePolling: true  // For Docker compatibility
    }
  }
})
```

## 🐛 Troubleshooting

### Frontend Not Starting

**Symptom:** Container exits immediately or logs show npm errors

**Solution:**
```bash
cd frontend
docker compose down -v           # Remove volumes
docker compose build --no-cache  # Rebuild image
docker compose up -d             # Start container
```

### Permission Errors in Frontend

**Symptom:** `EACCES: permission denied` errors

**Solution:** Already handled in the Dockerfile with proper user permissions. If issues persist:
```bash
cd frontend
docker compose down -v
docker compose up -d
```

### CORS Errors

**Symptom:** Browser console shows CORS policy errors

**Solution:**
1. Verify backend CORS configuration includes `http://localhost:5173`
2. Ensure frontend uses `credentials: 'include'` in fetch requests
3. Check `Access-Control-Allow-Credentials` header is `true`

### Cookie Not Being Set

**Symptom:** Set-Cookie header sent but cookie not stored in browser

**Checklist:**
- ✅ Backend returns `Access-Control-Allow-Credentials: true`
- ✅ Frontend uses `credentials: 'include'` in fetch
- ✅ Frontend and backend are on allowed domains
- ✅ Cookie flags are appropriate (HttpOnly, Secure in production)

### Port Already in Use

**Symptom:** Cannot bind to port (8080, 5173, or 8000)

**Solution:**
```bash
# Find process using the port
lsof -i :8080

# Kill the process
kill -9 <PID>

# Or stop Docker containers
docker ps | grep go-google-auth
docker stop <CONTAINER_ID>
```

### Docker Build Fails

**Symptom:** Build context transfer errors or build failures

**Solution:**
```bash
# Clean Docker system
docker system prune -a

# Rebuild from scratch
cd frontend
docker compose build --no-cache --pull
```

## 🤝 Contributing

Contributions are welcome! Please follow these guidelines:

1. **Fork the repository**
2. **Create a feature branch**
   ```bash
   git checkout -b feature/amazing-feature
   ```
3. **Commit your changes**
   ```bash
   git commit -m 'Add amazing feature'
   ```
4. **Push to the branch**
   ```bash
   git push origin feature/amazing-feature
   ```
5. **Open a Pull Request**

### Code Style

- **Backend (Go)**: Use `gofmt` and follow [Effective Go](https://go.dev/doc/effective_go)
- **Frontend (Vue/TS)**: Use Prettier and ESLint configurations provided

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 👤 Author

**yuki5155**

## 🔗 Resources

- [Go Documentation](https://go.dev/doc/)
- [Gin Framework](https://gin-gonic.com/)
- [Vue.js 3 Documentation](https://vuejs.org/)
- [Vite Documentation](https://vitejs.dev/)
- [Docker Documentation](https://docs.docker.com/)
- [AWS DynamoDB Local](https://docs.aws.amazon.com/amazondynamodb/latest/developerguide/DynamoDBLocal.html)

## 📊 Project Status

This project is currently in active development. Core features are functional, with Google OAuth integration planned for the next release.

---

**Note:** This is a development project. For production deployment, ensure proper security configurations, use HTTPS, and enable appropriate cookie security flags.
