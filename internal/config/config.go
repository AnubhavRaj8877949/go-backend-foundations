package config

type AppConfig struct {
	Name        string
	Environment string
}
type HTTPConfig struct {
	Host string
	Port int
}
type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Database string
}
type RedisConfig struct {
	Address  string
	Password string
	Database int
}

type RabbitMQConfig struct {
	URL string
}

type LoggingConfig struct {
	Level  string
	Format string
}
type Config struct {
	App      AppConfig
	HTTP     HTTPConfig
	Database DatabaseConfig
	Redis    RedisConfig
	RabbitMQ RabbitMQConfig
	Logging  LoggingConfig
}
