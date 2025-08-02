package internal

import (
	"github.com/anditakaesar/uwa-server-checker/adapter/docker"
	"github.com/anditakaesar/uwa-server-checker/internal/logger"
)

type Adapter struct {
	Docker docker.Interface
}

func InitializeAdapters() (Adapter, error) {
	log := logger.GetLogInstance()
	dockerAdp, err := docker.New()
	if err != nil {
		log.Error("Error initialize docker", err)
		return Adapter{}, err
	}

	return Adapter{
		Docker: dockerAdp,
	}, nil
}
