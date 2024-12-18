package util

import (
	"slices"
	"sync"

	"github.com/hashicorp/go-multierror"
)

type MovableDisposer struct {
	// stack of callbacks to run in reverse order
	disposers []func() error
	mu        sync.Mutex
}

func (md *MovableDisposer) Defer(d func()) {
	md.mu.Lock()
	defer md.mu.Unlock()

	md.disposers = append(md.disposers, func() error {
		d()
		return nil
	})
}

func (md *MovableDisposer) DeferWithError(d func() error) {
	md.mu.Lock()
	defer md.mu.Unlock()

	md.disposers = append(md.disposers, d)
}

func (md *MovableDisposer) Dispose() error {
	md.mu.Lock()
	defer md.mu.Unlock()

	var result error

	for i := len(md.disposers) - 1; i >= 0; i-- {
		if err := md.disposers[i](); err != nil {
			result = multierror.Append(result, err)
		}
	}

	return result
}

func (md *MovableDisposer) Move() *MovableDisposer {
	md.mu.Lock()
	defer md.mu.Unlock()

	disposers := slices.Clone(md.disposers)
	md.disposers = nil

	return &MovableDisposer{
		disposers: disposers,
	}
}

func (md *MovableDisposer) MoveTo(other *MovableDisposer) {
	md.mu.Lock()
	defer md.mu.Unlock()

	other.mu.Lock()
	defer other.mu.Unlock()

}
