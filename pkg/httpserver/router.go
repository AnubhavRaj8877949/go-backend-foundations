package httpserver

import (
	"log/slog"
	"net/http"

	"example.com/go-backend-foundation/pkg/httpserver/handlers"
	"example.com/go-backend-foundation/pkg/messaging"
	appmiddleware "example.com/go-backend-foundation/pkg/middleware"
	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
	goredis "github.com/redis/go-redis/v9"
)

func NewRouter(
	log *slog.Logger,
	db *pgxpool.Pool,
	redisClient *goredis.Client,
	rabbitMQ *messaging.RabbitMQ,
) http.Handler {
	r := chi.NewRouter()
	healthHandler := handlers.HealthHandler{
		DB:       db,
		Redis:    redisClient,
		RabbitMQ: rabbitMQ,
	}

	// Global middleware
	r.Use(chimiddleware.RequestID)
	r.Use(chimiddleware.RealIP)
	r.Use(chimiddleware.Recoverer)
	r.Use(appmiddleware.RequestLogger(log))

	// Public routes
	r.Get("/health", healthHandler.Health)
	r.Get("/ready", healthHandler.Ready)

	// API v1 routes
	r.Route("/api/v1", func(r chi.Router) {
		r.Get("/hello", handlers.Hello)
	})

	return r
}
