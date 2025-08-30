package main

import (
	"RESTCryptoServer/internal/auth"
	"RESTCryptoServer/internal/crypto"
	"RESTCryptoServer/internal/db"
	"RESTCryptoServer/internal/redis"
	"RESTCryptoServer/internal/updater"
	"RESTCryptoServer/monitoring"
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	logLevel := os.Getenv("LOG_LEVEL")
	if logLevel == "" {
		logLevel = "info"
	}
	monitoring.InitLogger(logLevel, os.Getenv("ENV") != "production")

	monitoring.Logger.Info().
		Str("version", "1.0.0").
		Str("go_version", "1.23").
		Msg("Starting Crypto Server")

	userdb, err := db.NewUserDB()
	if err != nil {
		log.Println("error during opening/creation users postgres db: ", err)
		return
	}
	defer userdb.Close()

	cryptodb, err := db.NewCryptoDB()
	if err != nil {
		log.Println("error during opening/creation crypto postgres db: ", err)
		return
	}
	defer cryptodb.Close()

	cache, err := redis.NewRedisClient()
	if err != nil {
		log.Println("error during cache redis db: ", err)
		return
	}
	defer cache.Close()

	monitoring.Logger.Info().Msg("All database connections established")

	authService := auth.NewAuthService(userdb)
	cryptoService := crypto.NewCryptoService(cryptodb, cache)
	updaterService := updater.NewUpdater(cryptoService, 30)

	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Recoverer)
	router.Use(monitoring.TracingMiddleware)
	router.Use(monitoring.MetricsMiddleware)

	router.Use(middleware.Throttle(100))

	router.Use(middleware.SetHeader("Access-Control-Allow-Origin", "*"))
	router.Use(middleware.SetHeader("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS"))
	router.Use(middleware.SetHeader("Access-Control-Allow-Headers", "Authorization, Content-Type"))

	router.Get("/health", monitoring.HealthHandler(userdb, cryptodb, cache))
	router.Get("/ready", monitoring.ReadinessHandler(userdb, cryptodb, cache))
	router.Get("/live", monitoring.LivenessHandler())

	router.Handle("/metrics", promhttp.Handler())

	router.Post("/auth/register", auth.RegisterHandler(authService))
	router.Post("/auth/login", auth.LoginHandler(authService))

	router.Group(func(r chi.Router) {
		r.Use(auth.AuthMiddleware)
		
		r.Get("/crypto", crypto.GETCryptosHandler(cryptoService))
		r.Post("/crypto", crypto.POSTCryptoHandler(cryptoService))
		r.Get("/crypto/{symbol}", crypto.GETCryptoSymbolHandler(cryptoService))
		r.Put("/crypto/{symbol}/refresh", crypto.PUTCryptoSymbolRefreshHandler(cryptoService))
		r.Get("/crypto/{symbol}/history", crypto.GETCryptoHistoryHandler(cryptoService))
		r.Get("/crypto/{symbol}/stats", crypto.GETCryptoStatsHandler(cryptoService))
		r.Delete("/crypto/{symbol}", crypto.DELETECryptoSymbolHandler(cryptoService))
		
		r.Get("/schedule", updater.GETScheduleParamsHandler(updaterService))
		r.Put("/schedule", updater.PUTScheduleParamsHandler(updaterService))
		r.Post("/schedule/trigger", updater.POSTScheduleTriggerHandler(updaterService))
	})

	log.Println("Starting server on http://:8080")
	err = http.ListenAndServe(":8080", router)
	if err != nil {
		log.Fatalln(err)
	}

	srv := &http.Server{
		Addr:         ":8080",
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		monitoring.Logger.Info().Str("addr", srv.Addr).Msg("Starting HTTP server")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			monitoring.Logger.Fatal().Err(err).Msg("Failed to start server")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	
	monitoring.Logger.Info().Msg("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	monitoring.Logger.Info().Msg("Stopping updater service...")
	updaterService.EndUpdating()

	if err := srv.Shutdown(ctx); err != nil {
		monitoring.Logger.Error().Err(err).Msg("Server forced to shutdown")
	}

	monitoring.Logger.Info().Msg("Server exited")
}