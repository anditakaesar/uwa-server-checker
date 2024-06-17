package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/anditakaesar/uwa-server-checker/internal/env"
	"github.com/anditakaesar/uwa-server-checker/internal/initializer"
	internalRouter "github.com/anditakaesar/uwa-server-checker/internal/router"
)

func main() {
	env := env.New()
	router := &internalRouter.Router{
		ServeMux: http.NewServeMux(),
		Env:      env,
	}

	init := initializer.New(router)
	err := init.InitModules()
	if err != nil {
		log.Fatalf("couldn't start modules with err: %v", err)
	}

	defer init.Log.Flush()

	go init.Services.TelebotSvc.Run()

	server := &http.Server{
		Addr:    env.GetAddrPort(),
		Handler: internalRouter.NewHandlerServer(router, env),
	}

	init.Log.Info(fmt.Sprintf("server run on port: %s", env.AppPort()))
	err = server.ListenAndServe()
	if err != nil {
		log.Fatalf("couldn't start server with err: %v", err)
	}
}
