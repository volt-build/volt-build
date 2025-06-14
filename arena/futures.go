package arena

import "sync"

type Future[T any] struct {
	mu      sync.Mutex
	cv      *sync.Cond
	done    bool
	result  T
	deps    []*Future[T]
	waiters []*Future[T]
	fn      func([]any) T
}

func NewFuture[T any](fn func([]any) T, deps ...*Future[T]) *Future[T] {
	f := &Future[T]{
		deps:    deps,
		fn:      fn,
		waiters: make([]*Future[T], 0),
	}
	f.cv = sync.NewCond(&f.mu)
	for _, dep := range deps {
		if dep != nil {
			dep.mu.Lock()
			dep.waiters = append(dep.waiters, f)
			dep.mu.Unlock()
		}
	}
	f.maybeStart()
	return f
}

func (f *Future[T]) maybeStart() {
	f.mu.Lock()
	if f.done || f.fn == nil {
		f.mu.Unlock()
		return
	}
	for _, dep := range f.deps {
		if dep != nil {
			if !dep.IsDone() {
				f.mu.Unlock()
				return
			}
		}
	}
	f.mu.Unlock()
	go f.evaluate()
}

func (f *Future[T]) IsDone() bool {
	return f.done
}

func (f *Future[T]) evaluate() {
	args := make([]any, len(f.deps))
	for i, dep := range f.deps {
		if dep != nil {
			args[i] = dep.Await()
		}
	}
	result := f.fn(args)
	f.mu.Lock()
	f.result = result
	f.done = true
	f.cv.Broadcast()
	waiters := append([]*Future[T]{}, f.waiters...)
	f.mu.Unlock()

	for _, w := range waiters {
		w.maybeStart()
	}
}

func (f *Future[T]) Await() T {
	f.mu.Lock()
	for !f.done {
		f.cv.Wait()
	}
	val := f.result
	f.mu.Unlock()
	return val
}
