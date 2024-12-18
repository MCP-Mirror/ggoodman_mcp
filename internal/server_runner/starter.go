package serverrunner

import (
	"context"
)

type ServerDescription struct {
	Runtime string
	Command string
	Args    []string
	Env     map[string]string

	MemoryLimitMB int
}

type StartedServer interface {
	Close() error
}

type ServerInstance interface {
	Run(ctx context.Context) error
}

type ServerStarter interface {
	Close() error
	Create(ctx context.Context, manifest ServerDescription) (ServerInstance, error)
}
