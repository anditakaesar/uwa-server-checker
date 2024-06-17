package initializer

import (
	"net/http"

	"github.com/anditakaesar/uwa-server-checker/adapter/docker"
	"github.com/anditakaesar/uwa-server-checker/internal/env"
	"github.com/anditakaesar/uwa-server-checker/internal/logger"
	"github.com/anditakaesar/uwa-server-checker/internal/router"
	"github.com/anditakaesar/uwa-server-checker/modules/beelink"
	"github.com/anditakaesar/uwa-server-checker/modules/telebot"
)

type Adapters struct {
	Docker docker.Interface
}

type Services struct {
	TelebotSvc *telebot.Telebot
}

type Initializer struct {
	Router   *router.Router
	Log      logger.Interface
	Services Services
	Adapters Adapters
}

func New(r *router.Router) *Initializer {
	log := NewLogger()
	r.Log = log

	initializer := &Initializer{
		Router:   r,
		Log:      log,
		Adapters: Adapters{},
	}

	err := initializer.InitializeAdapters()
	if err != nil {
		panic(err)
	}

	err = initializer.InitializeServices()
	if err != nil {
		panic(err)
	}

	return initializer
}

func NewLogger() logger.Interface {
	newClient := http.Client{}
	logglyCore := logger.NewLogglyZapCore(logger.NewLogglyLogWriter(
		logger.LogglyLogWriterDependency{
			HttpClient:    &newClient,
			BaseUrl:       env.LogglyBaseUrl(),
			CustomerToken: env.LogglyToken(),
			Tag:           env.LogglyTag(),
		},
	))

	internalLogger := logger.BuildNewLogger(logglyCore)
	return internalLogger
}

func (i *Initializer) InitModules() error {
	beelink.New(i.Router)

	return nil
}

func (i *Initializer) InitializeAdapters() error {
	dockerAdapter, err := docker.New(i.Log)
	if err != nil {
		i.Log.Error("failed to initialize docker adapter with", err)
		return err
	}

	i.Adapters.Docker = dockerAdapter
	return nil
}

func (i *Initializer) InitializeServices() error {
	botObj, err := telebot.New(i.Log, i.Adapters.Docker)
	if err != nil {
		i.Log.Error("failed to initialize docker adapter with", err)
		return err
	}

	i.Services.TelebotSvc = botObj
	return nil
}
