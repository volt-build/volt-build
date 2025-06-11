package compiler

import (
	"sync"
)

type anyFuture interface {
	Result() any
	addWaiting(anyFuture)
	Start()
	Done() bool
}

type Future[T any] struct {
	mu      sync.Mutex
	cv      *sync.Cond
	deps    []anyFuture
	waiting []anyFuture
	val     T
	done    bool
	fn      func([]any) T
}

func NewFuture[T any]() *Future[T] {
	f := &Future[T]{}
	f.cv = sync.NewCond(&f.mu)
	return f
}

func NewFutureFromDeps[T any](deps []anyFuture, fn func([]any) T) *Future[T] {
	f := &Future[T]{
		deps: deps,
		fn:   fn,
	}
	f.cv = sync.NewCond(&f.mu)
	return f
}

func (f *Future[T]) Complete(v T) {
	f.mu.Lock()
	defer f.mu.Unlock()
	if !f.done {
		f.val = v
		f.done = true
		f.cv.Broadcast()
		// trigger all waiting futures
		for _, w := range f.waiting {
			w.Start()
		}
	}
}

func (f *Future[T]) Result() any {
	f.mu.Lock()
	defer f.mu.Unlock()
	for !f.done {
		f.cv.Wait()
	}
	return f.val
}

func (f *Future[T]) Done() bool {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.done
}

func (f *Future[T]) addWaiting(wf anyFuture) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.waiting = append(f.waiting, wf)
}

func (f *Future[T]) Start() {
	// Register as waiting on deps
	for _, dep := range f.deps {
		dep.addWaiting(f)
	}

	f.mu.Lock()
	if f.done {
		f.mu.Unlock()
		return
	}
	f.mu.Unlock()

	// Check if all deps are done
	for _, dep := range f.deps {
		if !dep.Done() {
			return // wait until all deps are ready
		}
	}

	go func() {
		args := make([]any, len(f.deps))
		for i, dep := range f.deps {
			args[i] = dep.Result()
		}
		f.Complete(f.fn(args))
	}()
}
