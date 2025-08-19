package server

import (
	"context"
	"fmt"
	"gitlab.stat4market.com/reelsmarket/fiber-di-server-template/src/bootstrap/config"
	"gitlab.stat4market.com/reelsmarket/fiber-di-server-template/src/bootstrap/probe"
	"gitlab.stat4market.com/reelsmarket/fiber-di-server-template/src/internal/middleware"
	"gitlab.stat4market.com/reelsmarket/fiber-di-server-template/src/internal/utils"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"time"
)

func NewFiberApp(cfg *config.Config) *fiber.App {
	app := fiber.New(fiber.Config{
		ErrorHandler:          middleware.ErrorHandler,
		DisableStartupMessage: true,
	})

	app.Use(middleware.CorsMiddleware())
	if cfg.App.SwaggerEnabled {
		app.Get("/swagger/*", middleware.SwaggerMiddleware())
	}

	app.Get("/health", func(ctx *fiber.Ctx) error {
		return utils.SuccessOnly(ctx)
	})

	return app
}

func StartFiberApp(lifecycle fx.Lifecycle, app *fiber.App, logger *zap.Logger, cfg *config.Config, probeServer *probe.ProbeServer) {
	lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				logger.Info(fmt.Sprintf("🚀 Fiber server starting on :%d", cfg.App.Port))

				if cfg.App.SwaggerEnabled {
					logger.Info("📚 Swagger documentation is available at: /swagger")
				}

				// Mark startup as in progress
				if probeServer != nil {
					probeServer.MarkStartupComplete()
				}

				if err := app.Listen(fmt.Sprintf(":%d", cfg.App.Port)); err != nil {
					logger.Fatal("Fiber server failed to start", zap.Error(err))
				}
			}()

			// Wait a bit to ensure server has started properly, then mark as ready
			go func() {
				if probeServer != nil {
					// Small delay to ensure server is up
					time.Sleep(1 * time.Second)
					probeServer.MarkReady()
					logger.Info("🟢 Service marked as ready for traffic")
				}
			}()

			return nil
		},
		OnStop: func(ctx context.Context) error {
			// Mark service as not ready before shutdown
			if probeServer != nil {
				probeServer.MarkNotReady()
			}

			logger.Info("🛑 Shutting down Fiber server")
			return app.Shutdown()
		},
	})
}

var Module = fx.Options(
	fx.Provide(NewFiberApp),
	fx.Invoke(StartFiberApp),
)
