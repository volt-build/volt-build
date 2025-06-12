package arena

type Arena[T any] struct {
	Futures []*Future[T]
	Done    bool
	Chained *Arena[T]
}

func NewArena[T any](initialFuture *Future[T]) *Arena[T] {
	if len(initialFuture.waiters) != 0 {
		// (example) because a ProgramNode must not depend on another Node since that is the biggest node in the hierarchy
		panic("initial future must have 0 waiting futures")
	}
	futures := []*Future[T]{}
	futures[0] = initialFuture
	arena := &Arena[T]{
		Futures: futures,
		Done:    false,
	}
	return arena
}
