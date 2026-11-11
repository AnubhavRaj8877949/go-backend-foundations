package bootstrap

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"example.com/go-backend-foundation/internal/config"
	"example.com/go-backend-foundation/pkg/httpserver"
	"example.com/go-backend-foundation/pkg/messaging"
	"github.com/jackc/pgx/v5/pgxpool"
	goredis "github.com/redis/go-redis/v9"
)

type App struct {
	Config   *config.Config
	Logger   *slog.Logger
	DB       *pgxpool.Pool
	Redis    *goredis.Client
	RabbitMQ *messaging.RabbitMQ
}

func (a *App) Run(ctx context.Context) error {
	addr := fmt.Sprintf(
		"%s:%d",
		a.Config.HTTP.Host,
		a.Config.HTTP.Port,
	)

	server := httpserver.NewServer(addr, a.Logger, a.DB, a.Redis, a.RabbitMQ)
	serverErrors := make(chan error, 1)

	//server ko separate goroutine mein run karta hai, taaki main goroutine shutdown signal bhi sun sake.
	go func() {
		a.Logger.Info("HTTP server starting", "address", server.Addr)
		serverErrors <- server.ListenAndServe() // ye line server ko run karta hai aur agar error aata hai to serverErrors channel mein bhejta hai.
	}()

	select {
	case err := <-serverErrors: // agar serverErrors channel mein error aata hai to ye case execute hoga.
		if !errors.Is(err, http.ErrServerClosed) {
			return fmt.Errorf("HTTP server: %w", err)
		}

	case <-ctx.Done():
		a.Logger.Info("Shutdown signal received")

		shutdownCtx, cancel := context.WithTimeout(
			context.Background(), // ye line ek naya context create karta hai jo 10 second ke liye valid rahega.
			10*time.Second,       // server ko gracefully shutdown karne ke liye 10 second ka timeout set karta hai.
		)
		defer cancel()

		if err := server.Shutdown(shutdownCtx); err != nil { // server ko gracefully shutdown karne ki koshish karta hai, agar 10 second ke andar shutdown nahi hota to error return karega.
			return fmt.Errorf("HTTP server shutdown: %w", err)
		}

		a.Logger.Info("HTTP server stopped gracefully")
	}

	return nil
}
