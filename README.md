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

## 🛡️ Security Features

### **Production Security Measures**
- **Environment-Based Debug Access**: Debug endpoints (`/debug/memory`) are automatically disabled in production
- **Input Validation**: Whitelist validation for all debug actions
- **Secure Error Handling**: Generic error messages prevent information disclosure
- **No File System Exposure**: Limited system information access

### **Security Best Practices for Production**
```bash
# Set production environment
export ENVIRONMENT=prod
export APP_VERSION=1.0.0

# Debug endpoints will be automatically disabled
# Only essential health check endpoints available
# Enhanced security measures active
```

### **Security Headers to Add in Production**
```yaml
# In your Kubernetes ingress or load balancer
annotations:
  nginx.ingress.kubernetes.io/configuration-snippet: |
    add_header X-Content-Type-Options nosniff;
    add_header X-Frame-Options DENY;
    add_header X-XSS-Protection "1; mode=block";
    add_header Strict-Transport-Security "max-age=31536000; includeSubDomains";
```

## 🚀 What it does

- **Production-ready health checks** optimized for Kubernetes (tested on GKE, designed for any cluster)
- **Graceful shutdown handling** with 30-second timeout for zero-downtime deployments
- **Environment-aware responses** for different deployment stages
- **Memory monitoring** with system and Go runtime checks
- **Local testing capabilities** for simulating warning conditions
- **JSON API** with environment info, version, hostname, and timestamp

> **Platform Compatibility**: This service uses standard Kubernetes health check endpoints (`/health`, `/ready`) and Docker containers, making it compatible with any Kubernetes cluster (GKE, EKS, AKS, minikube, etc.) and any cloud provider. Currently tested on Google Kubernetes Engine (GKE).

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
- **`/debug/memory`** - Debug endpoint for testing memory warnings

## 🛑 Graceful Shutdown

The service implements graceful shutdown handling for clean termination in containerized environments:

### **How it works:**
- Listens for `SIGTERM` and `SIGINT` signals (Kubernetes pod termination, Ctrl+C)
- Stops accepting new requests immediately
- Allows existing requests to complete within **30-second timeout**
- Logs shutdown process for monitoring and debugging

### **Kubernetes Integration:**
```yaml
# In your deployment manifest
spec:
  template:
    spec:
      terminationGracePeriodSeconds: 30  # Matches app timeout
      containers:
      - name: health-check-service
        # ... other config
```

### **Benefits:**
- **Zero dropped connections** during pod restarts/scaling
- **Clean deployment rollouts** without service interruption
- **Load balancer friendly** - proper connection cleanup
- **Production reliability** - prevents abrupt termination errors

### **Shutdown Sequence:**
1. Signal received (SIGTERM from Kubernetes)
2. Server stops accepting new requests
3. Existing requests complete (up to 30 seconds)
4. Clean shutdown with status logging
5. Container exits with proper code

## 🧪 Testing Health Check Warnings

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

## 💡 How to Use This Service

### **1. As a Starting Template**
```bash
# Copy to start a new microservice
cp -r app/ my-new-service/
cd my-new-service
docker-compose up -d  # Works immediately!
```

### **2. In Kubernetes Clusters**
```yaml
# Use with any Helm chart or manifest
apiVersion: apps/v1
kind: Deployment
metadata:
  name: health-check-service
spec:
  template:
    spec:
      containers:
      - name: app
        image: your-registry/health-check-service:latest
        livenessProbe:
          httpGet:
            path: /health
            port: 8080
        readinessProbe:
          httpGet:
            path: /ready
            port: 8080
```

### **3. For Load Balancer Health Checks**
```bash
# Configure load balancer to check /health endpoint
# Returns 200 for healthy, 503 for degraded
curl -I http://your-service/health
```

### **4. In CI/CD Pipelines**
```bash
# Test health before deployment
curl -f http://localhost:8080/health || exit 1
curl -f http://localhost:8080/ready || exit 1
```

### **5. For Monitoring Systems**
```bash
# Prometheus can scrape /health for metrics
# Grafana can use the JSON response for dashboards
# AlertManager can trigger on 503 responses
```

### **6. In Development Teams**
- **New developers** get working health checks immediately
- **QA teams** can test warning conditions locally
- **DevOps** have production-ready monitoring endpoints
- **Architects** can use as a reference implementation

### **7. For Different Deployment Methods**
- **Docker Compose** - Local development
- **Kubernetes** - Any cluster (tested on GKE, should work on EKS, AKS, minikube, etc.)
- **Cloud Run** - Serverless container deployment
- **VM/Server** - Traditional deployment
- **Edge Computing** - Lightweight health monitoring

## ⚙️ Environment Variables

- **`ENVIRONMENT`** - Set to `dev`, `test`, `stage`, or `prod` (defaults to `dev`)
- **`APP_VERSION`** - Application version (defaults to `dev`)
- **`PORT`** - Server port (defaults to `8080`)
- **`LOG_LEVEL`** - Logging level (defaults to `debug`)

## 🐳 Docker

- **Docker Compose** - Runs app on port 8080
- **Health checks** - Built-in container health monitoring

## 🧪 Quick Start

```bash
# Start the app
cd app
docker-compose up -d

# Test endpoints
curl http://localhost:8080/
curl http://localhost:8080/health
curl http://localhost:8080/ready

# Test memory simulation
curl "http://localhost:8080/debug/memory?action=allocate"
curl -s http://localhost:8080/health | jq '.'
```

## 🎯 Production Features

- **Kubernetes-ready** health and readiness probes
- **Graceful shutdown handling** with 30-second timeout for clean container termination
- **Memory leak detection** with configurable thresholds
- **System resource monitoring** for production environments
- **Fast response times** optimized for high-frequency health checks
- **Clean JSON responses** for monitoring systems

## 🔧 Prerequisites

- Docker and Docker Compose
- Go 1.24+ (for local development)

## 🚀 Deployment

This service is designed to work with any deployment method:

1. **Local Development** - Use docker-compose
2. **Kubernetes** - Deploy with any Helm chart or manifest
3. **Cloud Run** - Deploy as a container (Google Cloud)
4. **VM/Server** - Run as a binary

## 📚 Documentation

- [App Documentation](app/README.md) - Detailed setup, testing, and production features

## 🔄 Reusability

This repository is designed to be **copied and reused** in other projects:

### **Git Clone Method (Recommended):**
```bash
# Clone the template repository
git clone https://github.com/userravs/go-health-check-service.git my-new-service
cd my-new-service

# Remove git history and start fresh
rm -rf .git
git init
git add .
git commit -m "Initial commit: health-check-service implementation"

# Start the service
docker-compose up -d  # Works immediately!
```

### **Clean Copy Method (Alternative):**
```bash
# Copy source files (excluding git history)
cp -r app/* my-new-service/
cd my-new-service

# Initialize as new git repository
git init
git add .
git commit -m "Initial commit: health-check-service template"

# Start the service
docker-compose up -d  # Works immediately!
```

### **What to Update After Copying:**
1. **`go.mod`** - Change module name to match new project
2. **`docker-compose.yml`** - Update service names if needed
3. **`README.md`** - Update project-specific information
4. **Environment variables** - Adjust for new project needs

## 🤝 Contributing

Found a bug or have an idea? Contributions are welcome!

- **Report bugs** using the bug report template
- **Suggest features** using the feature request template
- **Keep it simple** - this is a focused utility app
- **Be respectful** in all interactions

Perfect for teams that want a **production-ready health check service** without building from scratch!
