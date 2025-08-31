package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"
)

// Response structure for JSON responses
type Response struct {
	Message     string    `json:"message"`
	Environment string    `json:"environment"`
	Version     string    `json:"version"`
	Hostname    string    `json:"hostname"`
	Timestamp   time.Time `json:"timestamp"`
}

// HealthResponse for health checks
type HealthResponse struct {
	Status    string            `json:"status"`
	Details   map[string]string `json:"details,omitempty"`
	Timestamp time.Time         `json:"timestamp"`
}

// Global app state
var (
	appInitialized = false
	startTime      = time.Now()
)

// System memory information
type SystemMemory struct {
	total        uint64
	available    uint64
	used         uint64
	usagePercent float64
}

// Get system memory information from /proc/meminfo (Linux)
func getSystemMemory() (*SystemMemory, error) {
	// Try Linux /proc/meminfo
	if memInfo, err := os.ReadFile("/proc/meminfo"); err == nil {
		return parseMemInfo(string(memInfo))
	}

	return nil, fmt.Errorf("unable to get system memory info")
}

// Parse Linux /proc/meminfo
func parseMemInfo(content string) (*SystemMemory, error) {
	lines := strings.Split(content, "\n")
	var total, available uint64

	for _, line := range lines {
		if strings.HasPrefix(line, "MemTotal:") {
			fields := strings.Fields(line)
			if len(fields) >= 2 {
				if val, err := strconv.ParseUint(fields[1], 10, 64); err == nil {
					total = val * 1024 // Convert KB to bytes
				}
			}
		}
		if strings.HasPrefix(line, "MemAvailable:") {
			fields := strings.Fields(line)
			if len(fields) >= 2 {
				if val, err := strconv.ParseUint(fields[1], 10, 64); err == nil {
					available = val * 1024 // Convert KB to bytes
				}
			}
		}
	}

	if total == 0 || available == 0 {
		return nil, fmt.Errorf("invalid memory info")
	}

	used := total - available
	usagePercent := float64(used) / float64(total) * 100

	return &SystemMemory{
		total:        total,
		available:    available,
		used:         used,
		usagePercent: usagePercent,
	}, nil
}

func main() {
	// Get configuration from environment variables
	port := getEnv("PORT", "8080")
	environment := getEnv("ENVIRONMENT", "dev") // Default to dev
	version := getEnv("APP_VERSION", "1.0.0")

	// Validate environment
	if environment != "dev" && environment != "test" && environment != "stage" && environment != "prod" {
		log.Printf("⚠️  Invalid environment '%s', defaulting to 'dev'\n", environment)
		environment = "dev"
	}

	// Simulate app initialization
	go func() {
		time.Sleep(2 * time.Second) // Simulate startup time
		appInitialized = true
		log.Println("✅ Application initialized and ready")
	}()

	// Setup routes
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		homeHandler(w, r, environment, version)
	})
	http.HandleFunc("/health", healthHandler)
	http.HandleFunc("/ready", readyHandler)
	
	// Debug endpoint only available in non-production environments
	if environment != "prod" {
		http.HandleFunc("/debug/memory", debugMemoryHandler)
		log.Printf("🔧 Debug endpoints enabled for non-production environment: '%s'", environment)
	} else {
		log.Printf("🚫 Debug endpoints disabled in production for security: '%s'", environment)
	}

	// Start server
	log.Printf("🚀 Server starting on port %s in %s environment (version: %s)\n", port, environment, version)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}

func homeHandler(w http.ResponseWriter, r *http.Request, environment, version string) {
	hostname, _ := os.Hostname()

	// Different messages for each environment
	var message string
	var emoji string

	if environment == "prod" {
		message = "Hello from PROD! Live environment - handle with care!"
		emoji = "🚀"
	} else if environment == "stage" {
		message = "Hello from STAGE! Stage environment - safe for testing!"
		emoji = "🧪"
	} else if environment == "test" {
		message = "Hello from TEST! Test environment - safe for validation!"
		emoji = "🧬"
	} else if environment == "dev" {
		message = "Hello from DEV! Development environment - safe for debugging!"
		emoji = "🛠️"
	}

	response := Response{
		Message:     fmt.Sprintf("%s %s", emoji, message),
		Environment: environment,
		Version:     version,
		Hostname:    hostname,
		Timestamp:   time.Now(),
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Error encoding response: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	// Production-optimized health checks for GKE
	checks := make(map[string]string)

	// Only critical Go runtime check (memory leak detection)
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	// Alert only if Go memory usage is concerning (> 100MB)
	if m.Sys/1024/1024 > 100 {
		checks["go_memory"] = fmt.Sprintf("WARNING: %d MB", m.Sys/1024/1024)
	}

	// Critical system memory check (only if > 80%)
	if systemMem, err := getSystemMemory(); err == nil {
		if systemMem.usagePercent > 80 {
			checks["system_memory"] = fmt.Sprintf("WARNING: %.1f%%", systemMem.usagePercent)
		}
	}

	// Determine overall health
	status := "healthy"
	httpStatus := http.StatusOK

	// Only unhealthy if there are warnings
	if len(checks) > 0 {
		status = "degraded"
		httpStatus = http.StatusServiceUnavailable
	}

	response := HealthResponse{
		Status:    status,
		Details:   checks,
		Timestamp: time.Now(),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpStatus)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Error encoding response: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

func readyHandler(w http.ResponseWriter, r *http.Request) {
	// Simple readiness check for GKE
	if !appInitialized {
		response := HealthResponse{
			Status:    "not ready",
			Details:   map[string]string{"reason": "initializing"},
			Timestamp: time.Now(),
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusServiceUnavailable) // 503
		if err := json.NewEncoder(w).Encode(response); err != nil {
			log.Printf("Error encoding response: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	// App is ready - minimal response
	response := HealthResponse{
		Status:    "ready",
		Details:   map[string]string{},
		Timestamp: time.Now(),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK) // 200
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Error encoding response: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

// Global variable to hold allocated memory for testing
var debugMemory []byte

func debugMemoryHandler(w http.ResponseWriter, r *http.Request) {
	action := r.URL.Query().Get("action")

	switch action {
	case "allocate":
		// Allocate 150MB to trigger warning
		debugMemory = make([]byte, 150*1024*1024) // 150MB
		log.Printf("Allocated 150MB of memory for testing")
		fmt.Fprintf(w, "Allocated 150MB of memory. Check /health for warning.\n")

	case "free":
		// Free the allocated memory
		debugMemory = nil
		runtime.GC() // Force garbage collection
		log.Printf("Freed debug memory")
		fmt.Fprintf(w, "Freed debug memory. Check /health for clean status.\n")

	case "status":
		var m runtime.MemStats
		runtime.ReadMemStats(&m)
		fmt.Fprintf(w, "Current Go memory usage: %d MB\n", m.Sys/1024/1024)
		if debugMemory != nil {
			fmt.Fprintf(w, "Debug memory allocated: %d MB\n", len(debugMemory)/1024/1024)
		} else {
			fmt.Fprintf(w, "No debug memory allocated\n")
		}

	default:
		fmt.Fprintf(w, "Debug memory endpoint. Use ?action=allocate|free|status\n")
		fmt.Fprintf(w, "Examples:\n")
		fmt.Fprintf(w, "  /debug/memory?action=allocate  - Allocate 150MB\n")
		fmt.Fprintf(w, "  /debug/memory?action=free      - Free memory\n")
		fmt.Fprintf(w, "  /debug/memory?action=status    - Show current usage\n")
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
