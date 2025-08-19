package redis

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"gitlab.stat4market.com/reelsmarket/fiber-di-server-template/src/bootstrap/config"
	"go.uber.org/zap"
	"time"
)

// Redis struct for managing Redis connections
type Redis struct {
	Client *redis.Client
	Logger *zap.Logger
}

// NewRedis creates and returns a new Redis connection
func NewRedis(cfg *config.Config, logger *zap.Logger) (*Redis, error) {
	logger.Info("🔌 Connecting to Redis")

	options := &redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Redis.Host, cfg.Redis.Port),
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
		Username: cfg.Redis.Username,

		// Connection pool settings
		PoolSize:     cfg.Redis.PoolSize,
		MinIdleConns: cfg.Redis.MinIdleConns,
		DialTimeout:  time.Duration(cfg.Redis.TimeoutSeconds) * time.Second,
		ReadTimeout:  time.Duration(cfg.Redis.TimeoutSeconds) * time.Second,
		WriteTimeout: time.Duration(cfg.Redis.TimeoutSeconds) * time.Second,
		PoolTimeout:  time.Duration(cfg.Redis.TimeoutSeconds) * time.Second,
	}

	client := redis.NewClient(options)

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if _, err := client.Ping(ctx).Result(); err != nil {
		logger.Error("❌ Failed to connect to Redis", zap.Error(err))
		return nil, err
	}

	logger.Info("✅ Connected to Redis")
	return &Redis{
		Client: client,
		Logger: logger,
	}, nil
}

// Get retrieves a value for the given key
func (r *Redis) Get(ctx context.Context, key string) (string, error) {
	return r.Client.Get(ctx, key).Result()
}

// Set stores a value for the given key with optional expiration
func (r *Redis) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return r.Client.Set(ctx, key, value, expiration).Err()
}

// Del removes the specified keys
func (r *Redis) Del(ctx context.Context, keys ...string) error {
	return r.Client.Del(ctx, keys...).Err()
}

// HSet sets field in the hash stored at key to value
func (r *Redis) HSet(ctx context.Context, key string, values ...interface{}) error {
	return r.Client.HSet(ctx, key, values...).Err()
}

// HGet returns the value associated with field in the hash stored at key
func (r *Redis) HGet(ctx context.Context, key string, field string) (string, error) {
	return r.Client.HGet(ctx, key, field).Result()
}

// HGetAll returns all fields and values of the hash stored at key
func (r *Redis) HGetAll(ctx context.Context, key string) (map[string]string, error) {
	return r.Client.HGetAll(ctx, key).Result()
}

// onStart is called during application startup
func (r *Redis) onStart(ctx context.Context) error {
	r.Logger.Info("✅ Redis initialized")
	return nil
}

// onStop is called during application shutdown
func (r *Redis) onStop(ctx context.Context) error {
	r.Logger.Info("🔌 Disconnecting from Redis")
	if err := r.Client.Close(); err != nil {
		r.Logger.Error("❌ Failed to close Redis connection", zap.Error(err))
		return err
	}
	
	r.Logger.Info("✅ Disconnected from Redis")
	return nil
}
