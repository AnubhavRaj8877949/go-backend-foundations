package handlers

import (
	"net/http"

	"example.com/go-backend-foundation/pkg/messaging"
	"example.com/go-backend-foundation/pkg/response"
	"github.com/jackc/pgx/v5/pgxpool"
	goredis "github.com/redis/go-redis/v9"
)

type HealthHandler struct {
	DB       *pgxpool.Pool
	Redis    *goredis.Client
	RabbitMQ *messaging.RabbitMQ
}

func (h HealthHandler) Health(w http.ResponseWriter, _ *http.Request) {
	response.JSON(w, http.StatusOK, map[string]string{
		"status": "ok",
	})
}

func (h HealthHandler) Ready(w http.ResponseWriter, r *http.Request) {
	checks := make(map[string]string, 3)
	ready := true

	if h.DB == nil || h.DB.Ping(r.Context()) != nil {
		checks["postgres"] = "unavailable"
		ready = false
	} else {
		checks["postgres"] = "ok"
	}

	if h.Redis == nil || h.Redis.Ping(r.Context()).Err() != nil {
		checks["redis"] = "unavailable"
		ready = false
	} else {
		checks["redis"] = "ok"
	}

	if h.RabbitMQ == nil || h.RabbitMQ.Connection == nil || h.RabbitMQ.Connection.IsClosed() {
		checks["rabbitmq"] = "unavailable"
		ready = false
	} else {
		checks["rabbitmq"] = "ok"
	}

	status := http.StatusOK
	result := "ready"
	if !ready {
		status = http.StatusServiceUnavailable
		result = "not ready"
	}

	response.JSON(w, status, map[string]any{
		"status": result,
		"checks": checks,
	})
}
