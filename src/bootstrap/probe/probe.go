package probe

import (
	"context"
	"fmt"
	"gitlab.stat4market.com/reelsmarket/fiber-di-server-template/src/bootstrap/config"
	"go.uber.org/zap"
	"net/http"
	"sync/atomic"
	"time"
)

// ProbeServer represents a minimal HTTP server for Kubernetes probes and health checks
type ProbeServer struct {
	Server       *http.Server
	Logger       *zap.Logger
	readyFlag    atomic.Bool
	startupFlag  atomic.Bool
	shutdownFlag atomic.Bool
}

// NewProbeServer creates a new probe server
func NewProbeServer(cfg *config.Config, logger *zap.Logger) *ProbeServer {
	p := &ProbeServer{
		Logger: logger.Named("probe"),
	}

	// Create a simple HTTP server for health checks
	mux := http.NewServeMux()

	// Liveness probe - always returns 200 if the server is running
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("OK"))
	})

	// Readiness probe - returns 200 only if the application is ready to serve traffic
	mux.HandleFunc("/readyz", func(w http.ResponseWriter, r *http.Request) {
		if p.readyFlag.Load() {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("OK"))
		} else {
			w.WriteHeader(http.StatusServiceUnavailable)
			_, _ = w.Write([]byte("Service not ready yet"))
		}
	})

	// Startup probe - returns 200 only after the application has completed its initial startup
	mux.HandleFunc("/startupz", func(w http.ResponseWriter, r *http.Request) {
		if p.startupFlag.Load() {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("OK"))
		} else {
			w.WriteHeader(http.StatusServiceUnavailable)
			_, _ = w.Write([]byte("Service starting up"))
		}
	})

	p.Server = &http.Server{
		Addr:              fmt.Sprintf(":%d", cfg.Probe.Port),
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}

	return p
}

// Start starts the probe server
func (p *ProbeServer) Start(ctx context.Context) error {
	p.Logger.Info(fmt.Sprintf("🔍 Starting probe server on port %s", p.Server.Addr))

	go func() {
		if err := p.Server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			p.Logger.Error("Probe server error", zap.Error(err))
		}
	}()

	// Reset flags
	p.readyFlag.Store(false)
	p.startupFlag.Store(false)
	p.shutdownFlag.Store(false)

	return nil
}

// Stop stops the probe server
func (p *ProbeServer) Stop(ctx context.Context) error {
	// Mark as not ready and shutting down
	p.readyFlag.Store(false)
	p.shutdownFlag.Store(true)

	p.Logger.Info("🛑 Stopping probe server")
	stopCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := p.Server.Shutdown(stopCtx); err != nil {
		p.Logger.Error("Error shutting down probe server", zap.Error(err))
		return err
	}

	return nil
}

// MarkReady marks the service as ready
func (p *ProbeServer) MarkReady() {
	p.readyFlag.Store(true)
	p.Logger.Info("✅ Service marked as ready")
}

// MarkNotReady marks the service as not ready
func (p *ProbeServer) MarkNotReady() {
	p.readyFlag.Store(false)
	p.Logger.Info("❌ Service marked as not ready")
}

// MarkStartupComplete marks the startup as complete
func (p *ProbeServer) MarkStartupComplete() {
	p.startupFlag.Store(true)
	p.Logger.Info("✅ Service startup completed")
}

// IsReady returns whether the service is ready
func (p *ProbeServer) IsReady() bool {
	return p.readyFlag.Load()
}

// IsStartupComplete returns whether the startup has completed
func (p *ProbeServer) IsStartupComplete() bool {
	return p.startupFlag.Load()
}

// IsShuttingDown returns whether the service is shutting down
func (p *ProbeServer) IsShuttingDown() bool {
	return p.shutdownFlag.Load()
}
