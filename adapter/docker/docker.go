package docker

import (
	"github.com/anditakaesar/uwa-server-checker/dto"
	docker "github.com/fsouza/go-dockerclient"
)

const DefaultTimeout uint = 5

type Interface interface {
	GetContainerListWithPaging(size, page int) (*dto.ContainerWithPaging, error)
	GetContainersByName(name string) ([]dto.Container, error)
	StartContainer(containerID string) error
	StopContainer(containerID string) error
}

type Docker struct {
	Client *docker.Client
}

func New() (Interface, error) {
	client, err := docker.NewClientFromEnv()
	if err != nil {
		panic(err)
	}

	return &Docker{
		Client: client,
	}, nil
}

func (d *Docker) GetContainersByName(name string) ([]dto.Container, error) {
	containers, err := d.Client.ListContainers(docker.ListContainersOptions{
		All: true,
		Filters: map[string][]string{
			"name": []string{name},
		},
	})
	if err != nil {
		return nil, err
	}

	resultList := make([]dto.Container, len(containers))
	for i, container := range containers {
		resultList[i] = dto.Container{
			ID:     container.ID[:12],
			State:  container.State,
			Status: container.Status,
			Names:  container.Names,
		}
	}

	return resultList, nil
}

func (d *Docker) GetContainerListWithPaging(size, page int) (*dto.ContainerWithPaging, error) {
	result := dto.ContainerWithPaging{
		Size: size,
		Page: page,
	}

	containers, err := d.Client.ListContainers(docker.ListContainersOptions{
		All: true,
	})
	if err != nil {
		return nil, err
	}

	total := len(containers)
	result.Total = total
	maxPage := total / size
	if total%size > 0 {
		maxPage += 1
	}

	if page < maxPage {
		containers = containers[(page-1)*size : page*size]
	} else {
		containers = containers[total-1:]
	}

	resultList := make([]dto.Container, len(containers))
	for i, container := range containers {
		resultList[i] = dto.Container{
			ID:     container.ID[:12],
			State:  container.State,
			Status: container.Status,
			Names:  container.Names,
		}
	}
	result.List = resultList
	result.HasPrev = page-1 > 0
	result.HasNext = page+1 < maxPage

	return &result, nil
}

func (d *Docker) StartContainer(containerID string) error {
	return d.Client.StartContainer(containerID, nil)
}

func (d *Docker) StopContainer(containerID string) error {
	return d.Client.StopContainer(containerID, DefaultTimeout)
}
