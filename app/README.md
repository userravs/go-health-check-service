# Health Check Service

[![License: BSD-2-Clause](https://img.shields.io/badge/License-BSD%202--Clause-blue.svg)](https://opensource.org/licenses/BSD-2-Clause)
[![Go](https://img.shields.io/badge/Go-1.24+-blue.svg)](https://golang.org/)
[![Kubernetes](https://img.shields.io/badge/Kubernetes-1.24+-blue.svg)](https://kubernetes.io/)
[![Docker](https://img.shields.io/badge/Docker-20.0+-blue.svg)](https://www.docker.com/)
[![Platform](https://img.shields.io/badge/Platform-Linux%20%7C%20Alpine-blue.svg)](https://alpinelinux.org/)
[![Architecture](https://img.shields.io/badge/Architecture-amd64-blue.svg)](https://golang.org/)
[![HTTP](https://img.shields.io/badge/HTTP-REST%20API-blue.svg)](https://en.wikipedia.org/wiki/REST)
[![Security](https://img.shields.io/badge/Security-Production%20Hardened-green.svg)](https://github.com/userravs/go-health-check-service#security-features)
[![Version](https://img.shields.io/badge/Version-v1.0.0-green.svg)](https://github.com/userravs/go-health-check-service/releases/tag/v1.0.0)

A production-ready Go microservice with comprehensive health checks, designed for Kubernetes deployment and local development.

## 🚀 Features

- **Health endpoints** - `/health` and `/ready` for monitoring
- **JSON API** - Returns environment info, version, hostname, and timestamp
- **Production-ready** - Kubernetes-optimized health checks with minimal overhead (tested on GKE)
- **Memory monitoring** - System and Go runtime memory checks
- **Debug endpoints** - For testing health scenarios locally (non-production only)

## 🛡️ Security Features

### **Production Security Measures**
- **Environment-Based Debug Access**: Debug endpoints (`/debug/memory`) are automatically disabled in production
- **Input Validation**: Whitelist validation for all debug actions
- **Secure Error Handling**: Generic error messages prevent information disclosure
- **Production Hardening**: Debug endpoints completely removed when `ENVIRONMENT=prod`

### **Environment Behavior**
| Environment | Debug Access | Security Level |
|-------------|--------------|----------------|
| `dev` | ✅ Full access | 🔓 Development mode |
| `test` | ✅ Full access | 🔓 Testing mode |
| `stage` | ✅ Full access | 🔓 Staging mode |
| `prod` | ❌ Disabled | 🔒 Production hardened |

> **⚠️ Security Note**: In production (`ENVIRONMENT=prod`), debug endpoints are completely disabled for security. Only essential health check endpoints (`/`, `/health`, `/ready`) are available.

## 🌍 Environment Behavior

| Environment | Message | Emoji |
|-------------|---------|-------|
| `dev` | Development environment - safe for debugging! | 🛠️ |
| `test` | Test environment - safe for validation! | 🧬 |
| `stage` | Stage environment - safe for testing! | 🧪 |
| `prod` | Live environment - handle with care! | 🚀 |

## 📡 API Endpoints

- **`/`** - Main endpoint with environment info
- **`/health`** - Health check (status: "healthy" or "degraded")
- **`/ready`** - Readiness check (status: "ready" or "not ready")
- **`/debug/memory`** - Debug endpoint for testing memory warnings *(disabled in production)*

## 🧪 Testing Health Check Warnings

> **🔒 Security Note**: Debug endpoints are only available in non-production environments (`dev`, `test`, `stage`). In production (`ENVIRONMENT=prod`), these endpoints are completely disabled for security.

### **Simulate Memory Issues Locally**

```bash
# Check current memory status
curl "http://localhost:8080/debug/memory?action=status"

# Allocate 150MB to trigger warning
curl "http://localhost:8080/debug/memory?action=allocate"

# Check health endpoint (should show "degraded")
curl -s http://localhost:8080/health | jq '.'

# Free memory and return to normal
curl "http://localhost:8080/debug/memory?action=free"

# Check health endpoint (should show "healthy")
curl -s http://localhost:8080/health | jq '.'
```

### **Expected Responses**

**Normal (Healthy):**
```json
{
  "status": "healthy",
  "timestamp": "2025-08-28T22:16:56.312486885Z"
}
```

**Warning (Degraded):**
```json
{
  "status": "degraded",
  "details": {
    "go_memory": "WARNING: 159 MB"
  },
  "timestamp": "2025-08-28T22:17:15.486873192Z"
}
```

## ⚙️ Environment Variables

- **`ENVIRONMENT`** - Set to `dev`, `test`, `stage`, or `prod` (defaults to `dev`)
- **`APP_VERSION`** - Application version (defaults to `dev`)
- **`PORT`** - Server port (defaults to `8080`)
- **`LOG_LEVEL`** - Logging level (defaults to `debug`)

## 🐳 Docker

- **Docker Compose** - Runs app on port 8080
- **Health checks** - Built-in container health monitoring

## 🧪 Quick Test

```bash
# Start the app
docker-compose up

# Test endpoints
curl http://localhost:8080/
curl http://localhost:8080/health
curl http://localhost:8080/ready

# Test memory simulation
curl "http://localhost:8080/debug/memory?action=allocate"
curl -s http://localhost:8080/health | jq '.'
```

## 🚀 Production Deployment

### **Security Configuration**
```bash
# Set production environment
export ENVIRONMENT=prod
export APP_VERSION=1.0.0

# Debug endpoints will be automatically disabled
# Only essential health check endpoints available
# Enhanced security measures active
```

### **Kubernetes Health Checks**
```yaml
# Production-ready health check configuration
livenessProbe:
  httpGet:
    path: /health
    port: 8080
  initialDelaySeconds: 5
  periodSeconds: 10
  timeoutSeconds: 3
  failureThreshold: 3

readinessProbe:
  httpGet:
    path: /ready
    port: 8080
  initialDelaySeconds: 2
  periodSeconds: 5
  timeoutSeconds: 3
  failureThreshold: 3
```

## 🔄 Reusability

This service can be easily copied to other projects:

### **Git Clone Method (Recommended):**
```bash
# Clone the template repository
git clone https://github.com/userravs/go-health-check-service.git my-new-service
cd my-new-service
rm -rf .git
git init
git add .
git commit -m "Initial commit: health-check-service implementation"
docker-compose up -d
```

### **Clean Copy Method:**
```bash
# Copy source files
cp -r app/* my-new-service/
cd my-new-service
git init
git add .
git commit -m "Initial commit: health-check-service template"
docker-compose up -d
```
