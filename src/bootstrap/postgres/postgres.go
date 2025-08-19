package postgres

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
	"gitlab.stat4market.com/reelsmarket/fiber-di-server-template/src/bootstrap/config"
	"go.uber.org/zap"
	"os"
	"path/filepath"
	"runtime"
)

type Postgres struct {
	Client *sqlx.DB
	Logger *zap.Logger
}

func NewPostgres(cfg *config.Config, logger *zap.Logger) (*Postgres, error) {
	// Use the DSN method from the config package instead of manually constructing the string
	dsn := cfg.Postgres.DSN()
	logger.Log(zap.InfoLevel, "🔌 Connecting to PostgreSQL")
	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		logger.Error("❌ Failed to connect to PostgreSQL", zap.Error(err))
		return nil, err
	}

	client := &Postgres{
		Client: db,
		Logger: logger,
	}

	return client, nil
}

func (p *Postgres) onStart(ctx context.Context) error {
	return p.migrate()
}

func (p *Postgres) migrate() error {
	p.Logger.Info("🪿 Running migrations")
	_, b, _, _ := runtime.Caller(0)
	projectRoot := filepath.Join(filepath.Dir(b), "..", "..", "..")
	migrationsDir := filepath.Join(projectRoot, "migrations")

	if _, err := os.Stat(migrationsDir); os.IsNotExist(err) {
		p.Logger.Error("Migrations directory not found", zap.String("path", migrationsDir))
		return fmt.Errorf("migrations directory not found: %s", migrationsDir)
	}

	if err := goose.SetDialect("postgres"); err != nil {
		p.Logger.Error("Failed to set dialect", zap.Error(err))
		return fmt.Errorf("failed to set dialect: %w", err)
	}

	if err := goose.Up(p.Client.DB, migrationsDir); err != nil {
		p.Logger.Error("Failed to run migrations", zap.Error(err))
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	p.Logger.Info("✅ Migrations completed")

	return nil
}

func (p *Postgres) onStop(ctx context.Context) error {

	err := p.Client.Close()
	if err != nil {
		p.Logger.Error("Failed to close PostgreSQL connection", zap.Error(err))
		return err
	} else {
		p.Logger.Info("🔌 Disconnected from PostgreSQL")
	}

	return nil
}
