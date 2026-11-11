package bootstrap

func (a *App) Close() {
	if a.Redis != nil {
		if err := a.Redis.Close(); err != nil {
			a.Logger.Error("failed to close Redis client", "error", err)
		} else {
			a.Logger.Info("Redis client closed")
		}
	}
	if a.RabbitMQ != nil {
		if err := a.RabbitMQ.Close(); err != nil {
			a.Logger.Error("failed to close RabbitMQ connection", "error", err)
		} else {
			a.Logger.Info("RabbitMQ connection closed")
		}
	}

	if a.DB != nil {
		a.DB.Close()
		a.Logger.Info("PostgreSQL connection pool closed")
	}
}
