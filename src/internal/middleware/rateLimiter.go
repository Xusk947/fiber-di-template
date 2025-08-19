package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"gitlab.stat4market.com/reelsmarket/fiber-di-server-template/src/bootstrap/config"
	"go.uber.org/zap"
)

func RateLimiterMiddleware(cfg *config.Config, logger *zap.Logger) fiber.Handler {
	logger.Info("🚀 Rate limiter middleware enabled")
	logger.Info("Config", zap.Any("RateLimit", cfg.App.RateLimit), zap.Any("RateWindow", cfg.App.RateWindow))

	return limiter.New(limiter.Config{
		Max:        cfg.App.RateLimit,
		Expiration: cfg.App.RateWindow,
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.IP()
		},
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"error": "Rate limit exceeded",
			})
		},
	})
}
