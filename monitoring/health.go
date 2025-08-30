package monitoring

import (
	"encoding/json"
	"net/http"
	"time"
	"RESTCryptoServer/internal/db"
	"RESTCryptoServer/internal/redis"
)

type HealthStatus struct {
	Status    string                 `json:"status"`
	Timestamp time.Time              `json:"timestamp"`
	Uptime    string                 `json:"uptime"`
	Checks    map[string]CheckResult `json:"checks"`
}

type CheckResult struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
}

var startTime = time.Now()

func HealthHandler(userDB *db.UserDB, cryptoDB *db.CryptoDB, cache *redis.RedisClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		health := HealthStatus{
			Timestamp: time.Now(),
			Uptime:    time.Since(startTime).String(),
			Checks:    make(map[string]CheckResult),
		}
		
		if err := userDB.Ping(); err != nil {
			health.Checks["postgres_users"] = CheckResult{
				Status:  "unhealthy",
				Message: err.Error(),
			}
		} else {
			health.Checks["postgres_users"] = CheckResult{Status: "healthy"}
		}
		
		if err := cryptoDB.Ping(); err != nil {
			health.Checks["postgres_crypto"] = CheckResult{
				Status:  "unhealthy",
				Message: err.Error(),
			}
		} else {
			health.Checks["postgres_crypto"] = CheckResult{Status: "healthy"}
		}
		
		if err := cache.Ping(); err != nil {
			health.Checks["redis"] = CheckResult{
				Status:  "unhealthy",
				Message: err.Error(),
			}
		} else {
			health.Checks["redis"] = CheckResult{Status: "healthy"}
		}
		
		health.Status = "healthy"
		statusCode := http.StatusOK
		
		for _, check := range health.Checks {
			if check.Status == "unhealthy" {
				health.Status = "unhealthy"
				statusCode = http.StatusServiceUnavailable
				break
			}
		}
		
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(statusCode)
		json.NewEncoder(w).Encode(health)
	}
}

func ReadinessHandler(userDB *db.UserDB, cryptoDB *db.CryptoDB, cache *redis.RedisClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := userDB.Ping(); err != nil {
			http.Error(w, "Database not ready", http.StatusServiceUnavailable)
			return
		}
		
		if err := cryptoDB.Ping(); err != nil {
			http.Error(w, "Crypto database not ready", http.StatusServiceUnavailable)
			return
		}
		
		if err := cache.Ping(); err != nil {
			http.Error(w, "Cache not ready", http.StatusServiceUnavailable)
			return
		}
		
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}
}

func LivenessHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}
}