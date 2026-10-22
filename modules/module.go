package modules

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
	"github.com/anditakaesar/uwa-server-checker/modules/health"
	"github.com/anditakaesar/uwa-server-checker/modules/telebot"
)

type Module struct {
	Telebot *telebot.Telebot
}

type ModDep struct {
	Adapter internal.Adapter
}

func BuildModules(dep ModDep) (*Module, error) {
	log := logger.GetLogInstance()
	botObj, err := telebot.New(telebot.Dependency{
		Docker: dep.Adapter.Docker,
	})
	if err != nil {
		log.Error("failed to initialize bot object", err)
		return nil, err
	}

	_ = botObj.SendMessage("Power Surge detected! Check system uptime!")

	return &Module{
		Telebot: botObj,
	}, nil
}

func (m *Module) StartTelebot(ctx context.Context, wg *sync.WaitGroup, errCh chan<- error) {
	log := logger.GetLogInstance()
	wg.Add(1)
	botObj := m.Telebot
	go func() {
		defer wg.Done()

		// Start server in separate goroutine
		go func() {
			if err := botObj.Run(); err != nil && err != http.ErrServerClosed {
				errCh <- fmt.Errorf("TelebotModule start error: %w", err)
			}
		}()

		// Wait for shutdown signal
		<-ctx.Done()
		log.Info("Shutting down TelebotModule...")
		if err := botObj.Shutdown(); err != nil {
			errCh <- fmt.Errorf("TelebotModule shutdown error: %w", err)
		}
	}()
}

func (m *Module) StartApi(ctx context.Context, wg *sync.WaitGroup, errCh chan<- error) {
	log := logger.GetLogInstance()
	wg.Add(1)
	go func() {
		defer wg.Done()
		env := env.New()
		httprouter := &router.Router{
			ServeMux: http.NewServeMux(),
			Env:      env,
			Telebot:  m.Telebot,
		}
		health.NewModule(httprouter, m.Telebot)

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
