package config

import (
	"context"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/sethvargo/go-envconfig"
	"sync"
	"time"
)

// Config represents the application configuration
type Config struct {
	App        AppConfig
	Postgres   PostgresConfig
	ClickHouse ClickHouseConfig
	Redis      RedisConfig
	Kafka      KafkaConfig
	Probe      ProbeConfig
	Logger     LoggerConfig
}

// AppConfig contains application-specific configuration
type AppConfig struct {
	Env            string        `env:"ENV,default=development"`
	Port           int           `env:"PORT,default=8080"`
	SwaggerEnabled bool          `env:"SWAGGER_ENABLED,default=true"`
	RateLimit      int           `env:"RATE_LIMIT,default=100"`
	RateWindow     time.Duration `env:"RATE_WINDOW,default=1m"`
}

// PostgresConfig contains database connection configuration
type PostgresConfig struct {
	Host     string `env:"POSTGRES_HOST,default=localhost"`
	Port     int    `env:"POSTGRES_PORT,default=5432"`
	User     string `env:"POSTGRES_USER,default=postgres"`
	Password string `env:"POSTGRES_PASSWORD,default=postgres"`
	DBName   string `env:"POSTGRES_DBNAME,default=reelsmarket"`
	SSLMode  string `env:"POSTGRES_SSL_MODE,default=disable"`
}

// ClickHouseConfig contains ClickHouse connection configuration
type ClickHouseConfig struct {
	Host     string `env:"CLICKHOUSE_HOST,default=localhost"`
	Port     int    `env:"CLICKHOUSE_PORT,default=9000"`
	User     string `env:"CLICKHOUSE_USER,default=default"`
	Password string `env:"CLICKHOUSE_PASSWORD,default="`
	DBName   string `env:"CLICKHOUSE_DBNAME,default=default"`
}

// RedisConfig contains Redis connection configuration
type RedisConfig struct {
	Host           string `env:"REDIS_HOST,default=localhost"`
	Port           int    `env:"REDIS_PORT,default=6379"`
	Password       string `env:"REDIS_PASSWORD,default="`
	Username       string `env:"REDIS_USERNAME,default="`
	DB             int    `env:"REDIS_DB,default=0"`
	PoolSize       int    `env:"REDIS_POOL_SIZE,default=10"`
	MinIdleConns   int    `env:"REDIS_MIN_IDLE_CONNS,default=2"`
	TimeoutSeconds int    `env:"REDIS_TIMEOUT_SECONDS,default=5"`
}

// KafkaConfig contains Kafka connection configuration
type KafkaConfig struct {
	Brokers               []string `env:"KAFKA_BROKERS,default=localhost:9092"`
	Username              string   `env:"KAFKA_USERNAME,default="`
	Password              string   `env:"KAFKA_PASSWORD,default="`
	TimeoutSeconds        int      `env:"KAFKA_TIMEOUT_SECONDS,default=10"`
	Async                 bool     `env:"KAFKA_ASYNC,default=true"`
	BatchSize             int      `env:"KAFKA_BATCH_SIZE,default=100"`
	BatchTimeoutMs        int      `env:"KAFKA_BATCH_TIMEOUT_MS,default=1000"`
	AllowAutoTopicCreation bool     `env:"KAFKA_ALLOW_AUTO_TOPIC_CREATION,default=true"`
	MinBytes              int      `env:"KAFKA_MIN_BYTES,default=10240"` // 10KB
	MaxBytes              int      `env:"KAFKA_MAX_BYTES,default=10485760"` // 10MB
	CommitIntervalMs      int      `env:"KAFKA_COMMIT_INTERVAL_MS,default=1000"`
}

// ProbeConfig contains configuration for the health probe server
type ProbeConfig struct {
	Port              int  `env:"PROBE_PORT,default=8081"`
	EnableProbeServer bool `env:"ENABLE_PROBE_SERVER,default=true"`
}

// DSN returns the PostgreSQL connection string
func (p PostgresConfig) DSN() string {
	fmt.Println(p.Host, p.Port, p.User, p.Password, p.DBName, p.SSLMode)
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		p.Host, p.Port, p.User, p.Password, p.DBName, p.SSLMode)
}

// DSN returns the ClickHouse connection string
func (c ClickHouseConfig) DSN() string {
	return fmt.Sprintf("tcp://%s:%d?database=%s&username=%s&password=%s",
		c.Host, c.Port, c.DBName, c.User, c.Password)
}

// LoggerConfig contains logging configuration
type LoggerConfig struct {
	Level string `mapstructure:"level" default:"info"`
}

var (
	cfg  Config
	once sync.Once
)

// LoadEnv loads configuration from .env file
func LoadEnv() *Config {
	once.Do(func() {
		_ = godotenv.Load()

		if err := envconfig.Process(context.Background(), &cfg); err != nil {
			panic(err)
		}
	})

	return &cfg
}
