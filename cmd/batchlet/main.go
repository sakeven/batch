package main

import (
	"github.com/sakeven/batch/pkg/api"

	dockertypes "github.com/docker/docker/api/types"
	containertype "github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	log "github.com/sirupsen/logrus"
)

// JobLister is an interface to list jobs
type JobLister interface {
	List() []*api.Job
}

func create(job *api.Job) {
	for _, container := range job.Spec.Containers {
		containerConfig := &containertype.Config{
			Image: container.Image,
		}

		hostConfig := &containertype.HostConfig{}
		dockerContainer, err := dockerClient.ContainerCreate(nil, containerConfig, hostConfig, nil, container.Name)
		if err != nil {
			log.Errorf("failed to create container %s", container.Name)
			return
		}

		err = dockerClient.ContainerStart(nil, dockerContainer.ID, dockertypes.ContainerStartOptions{})
		if err != nil {
			log.Errorf("failed to start container %s", dockerContainer.ID)
		}
	}
}

var dockerClient *client.Client

func main() {
	var err error
	dockerClient, err = client.NewEnvClient()
	if err != nil {
		panic("can't connect to docker")
	}
}
