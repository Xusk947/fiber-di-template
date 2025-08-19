package clickhouse

import (
	"context"
	"fmt"
	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	"gitlab.stat4market.com/reelsmarket/fiber-di-server-template/src/bootstrap/config"
	"go.uber.org/zap"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

// ClickHouse struct for managing ClickHouse database connections
type ClickHouse struct {
	Client driver.Conn
	Logger *zap.Logger
}

// NewClickHouse creates and returns a new ClickHouse connection
func NewClickHouse(cfg *config.Config, logger *zap.Logger) (*ClickHouse, error) {
	logger.Info("🔌 Connecting to ClickHouse")

	// Configuring ClickHouse connection options
	conn, err := clickhouse.Open(&clickhouse.Options{
		Addr: []string{fmt.Sprintf("%s:%d", cfg.ClickHouse.Host, cfg.ClickHouse.Port)},
		Auth: clickhouse.Auth{
			Database: cfg.ClickHouse.DBName,
			Username: cfg.ClickHouse.User,
			Password: cfg.ClickHouse.Password,
		},
		Settings: clickhouse.Settings{
			"max_execution_time": 60,
		},
		DialTimeout:     time.Second * 10,
		MaxOpenConns:    50,
		MaxIdleConns:    10,
		ConnMaxLifetime: time.Hour,
		Debug:           cfg.App.Env != "production",
	})

	if err != nil {
		logger.Error("❌ Failed to connect to ClickHouse", zap.Error(err))
		return nil, err
	}

	// Test connection
	if err := conn.Ping(context.Background()); err != nil {
		logger.Error("❌ Failed to ping ClickHouse", zap.Error(err))
		return nil, err
	}

	client := &ClickHouse{
		Client: conn,
		Logger: logger,
	}

	return client, nil
}

// onStart is called during application startup
func (c *ClickHouse) onStart(ctx context.Context) error {
	return c.migrate()
}

// migrate runs ClickHouse migrations
func (c *ClickHouse) migrate() error {
	c.Logger.Info("🪿 Running ClickHouse migrations")
	_, b, _, _ := runtime.Caller(0)
	projectRoot := filepath.Join(filepath.Dir(b), "..", "..", "..")
	migrationsDir := filepath.Join(projectRoot, "migrations", "clickhouse")

	if _, err := os.Stat(migrationsDir); os.IsNotExist(err) {
		c.Logger.Error("ClickHouse migrations directory not found", zap.String("path", migrationsDir))
		return fmt.Errorf("clickhouse migrations directory not found: %s", migrationsDir)
	}

	// Read and execute .sql files in the migrations directory
	// Note: Unlike PostgreSQL with goose, ClickHouse doesn't have a standard migration tool
	// This is a simple implementation; you may want to use a proper migration tool

	files, err := os.ReadDir(migrationsDir)
	if err != nil {
		c.Logger.Error("Failed to read migrations directory", zap.Error(err))
		return fmt.Errorf("failed to read migrations directory: %w", err)
	}

	for _, file := range files {
		if filepath.Ext(file.Name()) != ".sql" {
			continue
		}

		filePath := filepath.Join(migrationsDir, file.Name())
		content, err := os.ReadFile(filePath)
		if err != nil {
			c.Logger.Error("Failed to read migration file", zap.String("file", file.Name()), zap.Error(err))
			continue
		}

		query := string(content)
		err = c.Client.Exec(context.Background(), query)
		if err != nil {
			c.Logger.Error("Failed to execute migration", zap.String("file", file.Name()), zap.Error(err))
			return fmt.Errorf("failed to execute migration %s: %w", file.Name(), err)
		}

		c.Logger.Info("Migration executed successfully", zap.String("file", file.Name()))
	}

	c.Logger.Info("✅ ClickHouse migrations completed")
	return nil
}

// onStop is called during application shutdown
func (c *ClickHouse) onStop(ctx context.Context) error {
	err := c.Client.Close()
	if err != nil {
		c.Logger.Error("Failed to close ClickHouse connection", zap.Error(err))
		return err
	}

	c.Logger.Info("🔌 Disconnected from ClickHouse")
	return nil
}
