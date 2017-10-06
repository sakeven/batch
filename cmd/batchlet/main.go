package main

import (
	"net/http"
	"time"

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

type jobLister struct {
	client *http.Client
}

func (jl *jobLister) List() []*api.Job {
	// jl.client.Get("/jobs")
	return nil
}

// ContainerRuntime is a type of container runtime
type ContainerRuntime struct {
	*client.Client
}

// Batchlet a
type Batchlet struct {
	jobLister        JobLister
	containerRuntime *ContainerRuntime
}

// NewBatchlet creates a new batchlet instance
func NewBatchlet(jobLister JobLister, containerRuntime *ContainerRuntime) *Batchlet {
	return &Batchlet{
		jobLister:        jobLister,
		containerRuntime: containerRuntime,
	}
}

// Run forever
func (b *Batchlet) Run() {
	for {
		for _, job := range b.jobLister.List() {
			if job.Spec.NodeName == "" {
				err := b.Bind(job)
				if err != nil {
					log.Errorf("Failed to bind job %s", job.Name)
					continue
				}
				err = b.Create(job)
				if err != nil {
					log.Errorf("Failed to create job %s", job.Name)
				}
			}
		}
		time.Sleep(time.Second)
	}
}

// Bind binds a job on node
func (b *Batchlet) Bind(job *api.Job) error {
	return nil
}

// Create creates a job
func (b *Batchlet) Create(job *api.Job) error {
	for _, container := range job.Spec.Containers {
		containerConfig := &containertype.Config{
			Image: container.Image,
		}

		hostConfig := &containertype.HostConfig{}
		dockerContainer, err := b.containerRuntime.ContainerCreate(nil, containerConfig, hostConfig, nil, container.Name)
		if err != nil {
			log.Errorf("Failed to create container %s", container.Name)
			return err
		}

		err = b.containerRuntime.ContainerStart(nil, dockerContainer.ID, dockertypes.ContainerStartOptions{})
		if err != nil {
			log.Errorf("Failed to start container %s", dockerContainer.ID)
			return err
		}
	}
	return nil
}

func main() {
	dockerClient, err := client.NewEnvClient()
	if err != nil {
		panic("can't connect to docker")
	}

	joblister := &jobLister{}
	batchlet := NewBatchlet(joblister, &ContainerRuntime{dockerClient})

	log.Info("Running Batchlet")
	batchlet.Run()
}
