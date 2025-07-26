package initializer

import (
	"github.com/anditakaesar/uwa-server-checker/adapter/docker"
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
	Services Services
	Adapters Adapters
}

func New(r *router.Router) *Initializer {
	initializer := &Initializer{
		Router:   r,
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

func (i *Initializer) InitModules() error {
	beelink.New(i.Router)

	return nil
}

func (i *Initializer) InitializeAdapters() error {
	dockerAdapter, err := docker.New()
	if err != nil {
		logger.GetLogInstance().Error("failed to initialize docker adapter with", err)
		return err
	}

	i.Adapters.Docker = dockerAdapter
	return nil
}

func (i *Initializer) InitializeServices() error {
	log := logger.GetLogInstance()
	botObj, err := telebot.New(telebot.Dependency{
		Docker: i.Adapters.Docker,
	})
	if err != nil {
		log.Error("failed to initialize docker adapter with", err)
		return err
	}

	i.Services.TelebotSvc = botObj
	return nil
}
