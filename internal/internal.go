package internal

import (
	"context"
	"sync"
)

// Module, interface that define a module to start on separate go routine
type Module interface {
	Start(ctx context.Context, wg *sync.WaitGroup, errCh chan<- error, dep Dependency)
}

type Dependency struct {
	Adapter Adapter
}
