package routes

import (
	"github.com/gofiber/fiber/v2"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"sort"
	"strings"
)

type Route struct {
	Fiber *fiber.App
}

func NewRoute(fiber *fiber.App) *Route {
	return &Route{Fiber: fiber}
}

var Module = fx.Module("routes", fx.Provide(NewRoute), fx.Invoke(registerRoutes))

func registerRoutes(
	r *Route,

	// yourNewController yourModule.IYourNewController, // Replace with your actual controllers

	log *zap.Logger,
) {
	api := r.Fiber.Group("/api")

	// yourNewController.SetupRoutes(api)

	log.Info("🚀 Routes:")
	stack := r.Fiber.Stack()

	// Recopilar todas las rutas API
	var routes []string
	maxMethodLen := 0

	for _, routeGroup := range stack {
		for _, route := range routeGroup {
			if !strings.HasPrefix(route.Path, "/api") || strings.HasPrefix(route.Method, "HEAD") {
				continue
			}
			routes = append(routes, route.Method+" "+route.Path)

			// Encontrar el método más largo para alineación
			if len(route.Method) > maxMethodLen {
				maxMethodLen = len(route.Method)
			}
		}
	}

	// Ordenar rutas alfabéticamente para mejor organización
	sort.Strings(routes)

	// Registrar las rutas con formato
	if len(routes) > 0 {
		log.Info("📋 API Endpoints:")
		for _, route := range routes {
			parts := strings.SplitN(route, " ", 2)
			if len(parts) == 2 {
				method, path := parts[0], parts[1]
				// Padding para alinear los métodos
				paddedMethod := method + strings.Repeat(" ", maxMethodLen-len(method))
				log.Info("  " + paddedMethod + "  │  " + path)
			}
		}
	} else {
		log.Info("No API routes found")
	}
}
