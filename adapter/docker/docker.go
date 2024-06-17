package docker

import (
	"github.com/anditakaesar/uwa-server-checker/dto"
	"github.com/anditakaesar/uwa-server-checker/internal/logger"
	docker "github.com/fsouza/go-dockerclient"
)

type Interface interface {
	GetContainerList() ([]dto.Container, error)
}

type Docker struct {
	Client *docker.Client
	Log    logger.Interface
}

func New(log logger.Interface) (Interface, error) {
	client, err := docker.NewClientFromEnv()
	if err != nil {
		panic(err)
	}

	return &Docker{
		Client: client,
		Log:    log,
	}, nil
}

func (d *Docker) GetContainerList() ([]dto.Container, error) {
	containers, err := d.Client.ListContainers(docker.ListContainersOptions{
		All: true,
	})
	if err != nil {
		return nil, err
	}

	response := make([]dto.Container, 0, len(containers))
	for _, container := range containers {
		response = append(response, dto.Container{
			ID:     container.ID[:12],
			State:  container.State,
			Status: container.Status,
			Names:  container.Names,
		})
	}

	return response, nil
}
