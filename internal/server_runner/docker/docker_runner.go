package docker_runner

import (
	"context"
	"fmt"
	serverrunner "mcp/internal/server_runner"
	"time"

	docker "github.com/docker/docker/client"
)

var _ serverrunner.ServerRunner = &DockerServerRunner{}

type DockerServerRunner struct {
	docker *docker.Client
}

func NewDockerServerRunner(ctx context.Context, ops ...docker.Opt) (*DockerServerRunner, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	docker, err := docker.NewClientWithOpts(append([]docker.Opt{docker.WithHostFromEnv()}, ops...)...)
	if err != nil {
		return nil, err
	}

	if _, err := docker.Ping(ctx); err != nil {
		return nil, fmt.Errorf("error connecting to docker: %w", err)
	}

	return &DockerServerRunner{
		docker: docker,
	}, nil
}

func (r *DockerServerRunner) Run(ctx context.Context, manifest *serverrunner.RunnableServer) (*serverrunner.RunningServer, error) {
	return nil, nil
}
