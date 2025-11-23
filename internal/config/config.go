package config

import (
	"fmt"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Environment string `yaml:"env" env:"ENV" env-default:"development"`
	LogLevel    string `yaml:"log_level" env:"LOG_LEVEL" env-default:"info"`

	GRPC struct {
		Host            string        `yaml:"host" env:"GRPC_HOST" env-default:"0.0.0.0"`
		Port            string        `yaml:"port" env:"GRPC_PORT" env-default:"50051"`
		RateLimit       int           `yaml:"rate_limit" env:"GRPC_RATE_LIMIT" env-default:"5"`
		RateLimitWindow time.Duration `yaml:"rate_limit_window" env:"GRPC_RATE_LIMIT_WINDOW" env-default:"1s"`
		Timeout         time.Duration `yaml:"timeout" env:"GRPC_TIMEOUT" env-default:"5s"`
		MaxConnAge      time.Duration `yaml:"max_conn_age" env:"GRPC_MAX_CONN_AGE" env-default:"5m"`
	} `yaml:"grpc"`

	Gateway struct {
		Host            string        `yaml:"host" env:"GATEWAY_HOST" env-default:"0.0.0.0"`
		Port            string        `yaml:"port" env:"GATEWAY_PORT" env-default:"8080"`
		RateLimit       int           `yaml:"rate_limit" env:"GATEWAY_RATE_LIMIT" env-default:"5"`
		RateLimitWindow time.Duration `yaml:"rate_limit_window" env:"GATEWAY_RATE_LIMIT_WINDOW" env-default:"1s"`
		GRPCServerHost  string        `yaml:"grpc_server_host" env:"GATEWAY_GRPC_SERVER_HOST" env-default:"localhost"`
		GRPCServerPort  string        `yaml:"grpc_server_port" env:"GATEWAY_GRPC_SERVER_PORT" env-default:"50051"`
		Timeout         time.Duration `yaml:"timeout" env:"GATEWAY_TIMEOUT" env-default:"10s"`
	} `yaml:"gateway"`

	Swagger struct {
		JSONPath string `yaml:"json_path" env:"SWAGGER_JSON_PATH" env-default:"./pkg/auth/auth.swagger.json"`
	} `yaml:"swagger"`

	Postgres struct {
		DSN             string        `yaml:"dsn" env:"POSTGRES_DSN" env-default:"postgres://postgres:postgres@localhost:5432/auth?sslmode=disable"`
		ReadDSN         string        `yaml:"read_dsn" env:"POSTGRES_READ_DSN" env-default:""`
		MaxOpenConns    int           `yaml:"max_open_conns" env:"POSTGRES_MAX_CONNS" env-default:"25"`
		MinConns        int           `yaml:"min_conns" env:"POSTGRES_MIN_CONNS" env-default:"25"`
		ConnMaxLifetime time.Duration `yaml:"conn_max_lifetime" env:"POSTGRES_CONN_MAX_LIFETIME" env-default:"5m"`
	} `yaml:"postgres"`
}

func New() (*Config, error) {
	cfg := &Config{}

	err := cleanenv.ReadConfig("./config/config.yml", cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

// GetGRPCAddress возвращает полный адрес gRPC сервера
func (c *Config) GetGRPCAddress() string {
	return fmt.Sprintf("%s:%s", c.GRPC.Host, c.GRPC.Port)
}

// GetGatewayAddress возвращает полный адрес HTTP Gateway
func (c *Config) GetGatewayAddress() string {
	return fmt.Sprintf("%s:%s", c.Gateway.Host, c.Gateway.Port)
}

// GetGRPCEndpoint возвращает адрес gRPC сервера для клиентов
func (c *Config) GetGRPCEndpoint() string {
	return fmt.Sprintf("%s:%s", c.Gateway.GRPCServerHost, c.Gateway.GRPCServerPort)
}

// GetPostgresDSN возвращает строку подключения к PostgreSQL
func (c *Config) GetPostgresDSN() string {
	return c.Postgres.DSN
}
