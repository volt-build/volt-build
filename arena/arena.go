package arena

import (
	"sync"
)

type Arena[T any] struct {
	mu      sync.Mutex
	Futures []*Future[T]
	Done    bool
	Chained *Arena[T] // use in super big projects
}

func NewArena[T any](initialFuture *Future[T]) *Arena[T] {
	if initialFuture != nil && len(initialFuture.waiters) != 0 {
		panic("initial future must have 0 waiting futures")
	}
	futures := make([]*Future[T], 0)
	if initialFuture != nil {
		futures = append(futures, initialFuture)
	}
	return &Arena[T]{
		Futures: futures,
		Done:    false,
	}
}

func (a *Arena[T]) AddFuture(f *Future[T]) {
	a.mu.Lock()
	defer a.mu.Unlock()
	if a.Done {
		panic("cannot add future to completed arena")
	}
	a.Futures = append(a.Futures, f)
}

func (a *Arena[T]) Chain() *Arena[T] {
	a.mu.Lock()
	defer a.mu.Unlock()
	if a.Done {
		panic("cannot chain completed arena")
	}
	newArena := &Arena[T]{
		Futures: make([]*Future[T], 0),
		Done:    false,
	}
	a.Chained = newArena
	return newArena
}

func (a *Arena[T]) AwaitAll() []T {
	a.mu.Lock()
	if a.Done {
		a.mu.Unlock()
		return a.collectResults()
	}
	futures := append([]*Future[T]{}, a.Futures...)
	a.Done = true
	a.mu.Unlock()

	results := make([]T, 0, len(futures))
	for _, f := range futures {
		results = append(results, f.Await())
	}

	if a.Chained != nil {
		results = append(results, a.Chained.AwaitAll()...)
	}
	return results
}

func (a *Arena[T]) collectResults() []T {
	results := make([]T, 0, len(a.Futures))
	for _, f := range a.Futures {
		results = append(results, f.result)
	}
	if a.Chained != nil {
		results = append(results, a.Chained.collectResults()...)
	}
	return results
}

func (a *Arena[T]) Destroy() error {
	for _, future := range a.Futures {
		future.destroy()
	}
	return nil
}
