package health

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/anditakaesar/uwa-server-checker/internal"
	"github.com/anditakaesar/uwa-server-checker/internal/env"
	"github.com/anditakaesar/uwa-server-checker/internal/logger"
	"github.com/anditakaesar/uwa-server-checker/internal/router"
)

type Module struct{}

func (s *Module) Start(ctx context.Context, wg *sync.WaitGroup, errCh chan<- error, dep internal.Dependency) {
	log := logger.GetLogInstance()
	wg.Add(1)
	go func() {
		defer wg.Done()
		env := env.New()
		httprouter := &router.Router{
			ServeMux: http.NewServeMux(),
			Env:      env,
		}
		new(httprouter)

		server := &http.Server{
			Addr:    env.GetAddrPort(),
			Handler: router.NewHandlerServer(httprouter),
		}

		// Start server in separate goroutine
		go func() {
			log.Info(fmt.Sprintf("HTTP server starting on %s", env.GetAddrPort()))
			if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				errCh <- fmt.Errorf("HTTP server error: %w", err)
			}
		}()

		// Wait for shutdown signal
		<-ctx.Done()
		log.Info("Shutting down HTTP server...")

		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := server.Shutdown(shutdownCtx); err != nil {
			errCh <- fmt.Errorf("HTTP shutdown error: %w", err)
		}
	}()
}

const modulePrefix string = "health"

func new(router *router.Router) {
	handler := &Handler{}

	router.AddEndpointWithPrefix(modulePrefix, getInfo(handler))
	router.AddEndpointWithPrefix(modulePrefix, getStatus(handler))
}

func getInfo(h *Handler) router.Endpoint {
	return router.Endpoint{
		HTTPMethod: http.MethodGet,
		Path:       "info",
		Handler:    h.getInfo(),
	}
}

func getStatus(h *Handler) router.Endpoint {
	return router.Endpoint{
		HTTPMethod: http.MethodGet,
		Path:       "status",
		Handler:    h.getStatus(),
	}
}
