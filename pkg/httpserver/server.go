package httpserver

import (
	"log/slog"
	"net/http"
	"time"

	"example.com/go-backend-foundation/pkg/messaging"
	"github.com/jackc/pgx/v5/pgxpool"
	goredis "github.com/redis/go-redis/v9"
)

func NewServer(
	addr string,
	log *slog.Logger,
	db *pgxpool.Pool,
	redisClient *goredis.Client,
	rabbitMQ *messaging.RabbitMQ,
) *http.Server {
	return &http.Server{
		Addr:              addr,
		Handler:           NewRouter(log, db, redisClient, rabbitMQ),
		ReadHeaderTimeout: 5 * time.Second,
		IdleTimeout:       60 * time.Second,
	}
}
