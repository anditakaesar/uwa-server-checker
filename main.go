package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/anditakaesar/uwa-server-checker/internal"
	"github.com/anditakaesar/uwa-server-checker/internal/logger"
	"github.com/anditakaesar/uwa-server-checker/modules/health"
	"github.com/anditakaesar/uwa-server-checker/modules/telebot"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	log := logger.GetLogInstance()
	defer log.Flush()

	adapter, err := internal.InitializeAdapters()
	if err != nil {
		os.Exit(1)
	}

	modules := []internal.Module{
		&health.Module{},
		&telebot.Module{},
	}

	var wg sync.WaitGroup
	errCh := make(chan error, len(modules))
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	moduleDep := internal.Dependency{
		Adapter: adapter,
	}

	for _, module := range modules {
		module.Start(ctx, &wg, errCh, moduleDep)
	}

	select {
	case sig := <-sigCh:
		log.Info(fmt.Sprintf("Received signal: %v", sig))
		cancel()
	case err := <-errCh:
		log.Error("Module Error", err)
		cancel()
	}

	wg.Wait()
	log.Info("Application shutdown complete")
}
