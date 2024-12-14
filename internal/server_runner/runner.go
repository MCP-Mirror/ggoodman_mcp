package serverrunner

import "context"

type ServerRunner interface {
	Run(ctx context.Context, manifest *RunnableServer) (*RunningServer, error)
}
