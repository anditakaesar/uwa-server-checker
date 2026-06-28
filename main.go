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
	"github.com/anditakaesar/uwa-server-checker/modules"
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

	apps, err := modules.BuildModules(modules.ModDep{
		Adapter: adapter,
	})
	if err != nil {
		os.Exit(1)
	}

	var wg sync.WaitGroup
	errCh := make(chan error, 2)
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	apps.StartTelebot(ctx, &wg, errCh)
	apps.StartApi(ctx, &wg, errCh)

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
