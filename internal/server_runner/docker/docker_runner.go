package docker_runner

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"mcp/internal/jsonrpc"
	"mcp/internal/mcp"
	serverrunner "mcp/internal/server_runner"
	"mcp/internal/util"
	"slices"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/api/types/network"
	docker "github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/sourcegraph/jsonrpc2"
	"golang.org/x/sync/errgroup"
)

var (
	DEFAULT_MEMORY_LIMIT_MB      int = 64
	SERVER_START_TIMEOUT_SECONDS int = 30
	SERVER_STOP_TIMEOUT_SECONDS      = 15
)

var _ serverrunner.ServerStarter = &DockerServerRunner{}

type DockerServerOptions struct {
	// TODO: Cache dir for making npx faster
}

type DockerServerRunner struct {
	docker *docker.Client
	logger *slog.Logger
}

func NewDockerServerRunner(ctx context.Context, logger *slog.Logger, ops DockerServerOptions) (*DockerServerRunner, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	docker, err := docker.NewClientWithOpts(docker.WithHostFromEnv())
	if err != nil {
		return nil, err
	}

	if _, err := docker.Ping(ctx); err != nil {
		return nil, fmt.Errorf("error connecting to docker: %w", err)
	}

	return &DockerServerRunner{
		docker: docker,
		logger: logger,
	}, nil
}

func (r *DockerServerRunner) Close() error {
	return r.docker.Close()
}

// Create creates a new server instance from the given manifest.
//
// The server is not started until Start is called, which blocks for the
// duration of the server's execution.
func (r *DockerServerRunner) Create(ctx context.Context, manifest serverrunner.ServerDescription) (serverrunner.ServerInstance, error) {
	var err error

	runtime, err := serverrunner.ParseRuntime(manifest.Runtime)
	if err != nil {
		return nil, fmt.Errorf("error parsing runtime: %w", err)
	}

	config := container.Config{
		StdinOnce:    true,
		StopTimeout:  &SERVER_STOP_TIMEOUT_SECONDS,
		AttachStdin:  true,
		AttachStdout: true,
		AttachStderr: true,
		OpenStdin:    true,
		User:         "node",

		Cmd: append([]string{manifest.Command}, slices.Clone(manifest.Args)...),
		Env: envMapToSlice(map[string]string{"NPM_CONFIG_CACHE": "/tmp/.npm", "GITHUB_PERSONAL_ACCESS_TOKEN": "FAKE 123"}),
	}

	initTrue := true
	memorySwappinessZero := int64(0)
	memoryLimitMB := manifest.MemoryLimitMB

	if memoryLimitMB == 0 {
		memoryLimitMB = DEFAULT_MEMORY_LIMIT_MB
	}

	hostConfig := container.HostConfig{
		AutoRemove: false,
		Init:       &initTrue,
		Resources: container.Resources{
			Memory:           int64(memoryLimitMB * 1024 * 1024),
			MemorySwap:       0,
			MemorySwappiness: &memorySwappinessZero,
		},
		ReadonlyRootfs: true,
		DNS:            []string{"8.8.8.8"},
		Tmpfs:          map[string]string{"/tmp": "rw"},
	}
	networkingConfig := network.NetworkingConfig{}

	switch runtime.Name {
	case "node":
		config.Image = "node"
		if runtime.Version != "" {
			config.Image += ":" + runtime.Version
		}

	case "python":
		config.Image = "python:" + runtime.Version
		if runtime.Version != "" {
			config.Image += ":" + runtime.Version
		}
	}

	pull, err := r.docker.ImagePull(ctx, config.Image, image.PullOptions{})
	if err != nil {
		return nil, fmt.Errorf("error pulling image: %w", err)
	}

	pullBuf, err := io.ReadAll(pull)
	if err != nil {
		return nil, fmt.Errorf("error reading image pull response: %w", err)
	}

	r.logger.Info("pulled image", "image", config.Image, "output", string(pullBuf))

	dsi := &DockerServerInstance{
		docker:           r.docker,
		logger:           r.logger.With("cmd", manifest.Command, "args", manifest.Args),
		containerConfig:  config,
		hostConfig:       hostConfig,
		networkingConfig: networkingConfig,
	}

	return dsi, nil
}

type DockerServerInstance struct {
	docker *docker.Client
	logger *slog.Logger

	containerConfig  container.Config
	hostConfig       container.HostConfig
	networkingConfig network.NetworkingConfig
}

func (dsi *DockerServerInstance) Run(ctx context.Context) error {
	defer dsi.docker.Close()

	cr, err := dsi.docker.ContainerCreate(ctx, &dsi.containerConfig, &dsi.hostConfig, &dsi.networkingConfig, nil, "")
	if err != nil {
		return fmt.Errorf("error creating container: %w", err)
	}
	// defer dsi.docker.ContainerRemove(ctx, cr.ID, container.RemoveOptions{
	// 	Force: true,
	// })

	waitCh, errCh := dsi.docker.ContainerWait(ctx, cr.ID, container.WaitConditionNotRunning)

	if err := dsi.docker.ContainerStart(ctx, cr.ID, container.StartOptions{}); err != nil {
		return fmt.Errorf("error starting container: %w", err)
	}

	dsi.logger.Info("started container", "id", cr.ID)

	// defer dsi.docker.ContainerStop(ctx, cr.ID, container.StopOptions{
	// 	Timeout: &SERVER_STOP_TIMEOUT_SECONDS,
	// })

	g, ctx := errgroup.WithContext(ctx)

	stdoutR, stdoutW := io.Pipe()
	_, stderrW := io.Pipe()

	// Grab stdin and stdout
	attachResp, err := dsi.docker.ContainerAttach(ctx, cr.ID, container.AttachOptions{
		Stream: true,
		Stdin:  true,
		Stdout: true,
		Stderr: true,
	})
	if err != nil {
		return fmt.Errorf("error attaching to container: %w", err)
	}
	defer attachResp.Close()

	handler := jsonrpc2.AsyncHandler(jsonrpc2.HandlerWithError(dsi.handleRequest).SuppressErrClosed())
	jsonRPCLogger := jsonrpc2.LogMessages(jsonrpc.NewJSONRPCLogger(dsi.logger))

	stream := jsonrpc2.NewPlainObjectStream(util.NewReaderWriterCloser(stdoutR, attachResp.Conn))
	defer stream.Close()

	conn := jsonrpc2.NewConn(ctx, stream, handler, jsonRPCLogger)
	defer conn.Close()

	g.Go(func() error {
		select {
		case <-ctx.Done():
			return nil
		case waitRes := <-waitCh:
			dsi.logger.Debug("container exited", "status", waitRes.StatusCode)
			return fmt.Errorf("container exited with code: %d", waitRes.StatusCode)
		case err := <-errCh:
			dsi.logger.Debug("error waiting for container", "err", err)
			return fmt.Errorf("error waiting for container: %w", err)
		case <-conn.DisconnectNotify():
			return fmt.Errorf("connection closed")
		}
	})

	g.Go(func() error {
		if _, err := stdcopy.StdCopy(stdoutW, stderrW, attachResp.Reader); err != nil && err != io.EOF {
			return fmt.Errorf("error copying stdio: %w", err)
		}
		return nil
	})

	g.Go(func() error {
		var initResult mcp.InitializeResult
		if err := conn.Call(ctx, "initialize", &mcp.InitializeRequest{
			ProtocolVersion: mcp.MCP_PROTOCOL_VERSION,
			Capabilities: mcp.ClientCapabilities{
				Roots: &mcp.ListChangesCapability{
					ListChanged: func() *bool { b := true; return &b }(),
				},
				Sampling: &mcp.SamplingCapability{},
			},
		}, &initResult); err != nil {
			return fmt.Errorf("error initializing server: %w", err)
		}

		if err := conn.Notify(ctx, "notifications/initialized", &mcp.InitializedNotification{}); err != nil {
			return fmt.Errorf("error sending initialized notification: %w", err)
		}

		return nil
	})

	return g.Wait()
}

func (dsi *DockerServerInstance) handleRequest(ctx context.Context, conn *jsonrpc2.Conn, req *jsonrpc2.Request) (interface{}, error) {
	return nil, &jsonrpc2.Error{
		Code:    jsonrpc2.CodeMethodNotFound,
		Message: fmt.Sprintf("method %q not found", req.Method),
	}
}

func envMapToSlice(env map[string]string) []string {
	envSlice := make([]string, 0, len(env))
	for k, v := range env {
		envSlice = append(envSlice, fmt.Sprintf("%s=%s", k, v))
	}
	return envSlice
}
