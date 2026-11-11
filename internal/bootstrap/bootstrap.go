package bootstrap

import (
	"context"
	"fmt"
	"time"

	"example.com/go-backend-foundation/internal/config"
	"example.com/go-backend-foundation/pkg/database"
	"example.com/go-backend-foundation/pkg/logger"
	"example.com/go-backend-foundation/pkg/messaging"
	redisclient "example.com/go-backend-foundation/pkg/redis"
	"example.com/go-backend-foundation/pkg/retry"
	"github.com/jackc/pgx/v5/pgxpool"
	goredis "github.com/redis/go-redis/v9"
)

func New() (*App, error) {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		return nil, fmt.Errorf("load configuration: %w", err)
	}

	// Initialize logger
	appLogger, err := logger.New(
		cfg.Logging.Level,
		cfg.Logging.Format,
	)
	if err != nil {
		return nil, fmt.Errorf("initialize logger: %w", err)
	}

	appLogger.Info("Initializing application")

	// Create database connection with timeout
	dbCtx, cancel := context.WithTimeout(
		context.Background(),
		10*time.Second,
	)
	defer cancel()

	var dbPool *pgxpool.Pool

	err = retry.Do(
		dbCtx,
		retry.Config{
			Attempts: 4,
			Delay:    time.Second,
		},
		func() error {
			var err error

			dbPool, err = database.NewPostgresPool(
				dbCtx,
				cfg.Database,
			)
			if err != nil {
				appLogger.Warn("PostgreSQL connection attempt failed", "error", err)
			}

			return err
		},
	)
	if err != nil {
		return nil, fmt.Errorf("initialize PostgreSQL after retries: %w", err)
	}

	appLogger.Info("PostgreSQL connected successfully")

	// Redis connection with timeout.
	var redisClient *goredis.Client

	redisCtx, redisCancel := context.WithTimeout(
		context.Background(),
		5*time.Second,
	)
	defer redisCancel()

	redisErr := retry.Do(
		redisCtx,
		retry.Config{
			Attempts: 3,
			Delay:    time.Second,
		},
		func() error {
			var err error

			redisClient, err = redisclient.NewClient(
				redisCtx,
				cfg.Redis,
			)
			if err != nil {
				appLogger.Warn(
					"Redis connection attempt failed",
					"error", err,
				)
			}

			return err
		},
	)

	if redisErr != nil {
		appLogger.Warn(
			"Redis unavailable; continuing without Redis",
			"error", redisErr,
		)
		redisClient = nil
	} else {
		appLogger.Info("Redis connected successfully")
	}
	// RabbitMQ connection with timeout.
	var rabbitMQ *messaging.RabbitMQ

	rabbitCtx, rabbitCancel := context.WithTimeout(
		context.Background(),
		15*time.Second,
	)
	defer rabbitCancel()

	rabbitErr := retry.Do(
		rabbitCtx,
		retry.Config{
			Attempts: 4,
			Delay:    time.Second,
		},
		func() error {
			var err error

			rabbitMQ, err = messaging.NewRabbitMQ(cfg.RabbitMQ.URL)
			if err != nil {
				appLogger.Warn(
					"RabbitMQ connection attempt failed",
					"error", err,
				)
			}

			return err
		},
	)

	if rabbitErr != nil {
		if redisClient != nil {
			_ = redisClient.Close()
		}
		dbPool.Close()

		return nil, fmt.Errorf(
			"initialize RabbitMQ after retries: %w",
			rabbitErr,
		)
	}

	if err := rabbitMQ.DeclareQueue("default"); err != nil {
		_ = rabbitMQ.Close()

		if redisClient != nil {
			_ = redisClient.Close()
		}

		dbPool.Close()

		return nil, fmt.Errorf("declare RabbitMQ queue: %w", err)
	}

	appLogger.Info("RabbitMQ connected successfully")

	appLogger.Info("RabbitMQ connected successfully")

	// Create application
	app := &App{
		Config:   &cfg,
		Logger:   appLogger,
		DB:       dbPool,
		Redis:    redisClient,
		RabbitMQ: rabbitMQ,
	}

	return app, nil
}
