// Package executor holds an execution manager for the language itself, to make it parallel and fast
package executor

import (
	"errors"
	"sync"
)

// Context interface for managing task execution context
type Context[T any] interface {
	Previous() T
	Advance(*T)
	AddTask(t *Task[T]) error
	Exec(taskFunc TaskFunc[T]) T
}

// DefaultContext implements the Context interface
type DefaultContext[T any] struct {
	Zero       *T           // zeroed variant to reduce allocation
	Tasks      [16]*Task[T] // list of tasks
	curTaskIdx uint8        // current task index
	Cur        *T           // current T being processed
	Prev       *T           // previous T being processed
	mu         sync.Mutex   // mutex for synchronization
	cond       *sync.Cond   // condition variable for signaling
}

// NewDefaultContext initializes a new DefaultContext
func NewDefaultContext[T any](zero *T) *DefaultContext[T] {
	dc := &DefaultContext[T]{Zero: zero}
	dc.cond = sync.NewCond(&dc.mu)
	return dc
}

// Previous returns the previous T being processed
func (dc *DefaultContext[T]) Previous() T {
	return *dc.Prev
}

// AddTask adds a task to the context's scope
func (dc *DefaultContext[T]) AddTask(t *Task[T]) error {
	dc.mu.Lock()
	defer dc.mu.Unlock()

	for i, cur := range dc.Tasks {
		if cur == nil {
			dc.Tasks[i] = t
			return nil
		}
	}
	return errors.New("overflow of tasks registered on a single context")
}

// Advance moves to the next task
func (dc *DefaultContext[T]) Advance(cur *T) {
	dc.mu.Lock()
	defer dc.mu.Unlock()
	dc.Prev = cur
	dc.Cur = dc.Zero
	dc.curTaskIdx++
	dc.cond.Broadcast()
}

// Exec executes a task function within the context
func (dc *DefaultContext[T]) Exec(taskFunc TaskFunc[T]) T {
	return taskFunc(dc)
}

// TaskFunc defines a function type for tasks
type TaskFunc[T any] func(ctx Context[T]) T

// Task represents a unit of work
type Task[T any] struct {
	ctx  Context[T] // Use Context[T] instead of *Context[T]
	Val  *T
	Done bool
	fun  TaskFunc[T]
	cv   *sync.Cond  // condition variable for signaling task completion
	mu   *sync.Mutex // mutex for synchronization
}

// NewTask initializes a new Task
func NewTask[T any](ctx Context[T], fun TaskFunc[T]) *Task[T] {
	task := &Task[T]{ctx: ctx, fun: fun, Done: false, mu: &sync.Mutex{}}
	task.cv = sync.NewCond(task.mu)
	return task
}

// Execute runs the task and signals completion
func (t *Task[T]) Execute() {
	t.mu.Lock()
	t.Done = true
	val := t.fun(t.ctx) // Use t.ctx directly
	t.Val = &val
	t.cv.Broadcast() // Notify that the task is done
	t.mu.Unlock()
}

// Executor manages the execution of tasks
type Executor[T any] struct {
	ctx *DefaultContext[T]
	wg  sync.WaitGroup
}

// NewExecutor initializes a new Executor
func NewExecutor[T any](zero *T) *Executor[T] {
	return &Executor[T]{ctx: NewDefaultContext(zero)}
}

// AddTask adds a task to the executor's context
func (e *Executor[T]) AddTask(t *Task[T]) error {
	return e.ctx.AddTask(t)
}

// Run starts executing tasks in parallel
func (e *Executor[T]) Run() {
	for _, task := range e.ctx.Tasks {
		if task != nil {
			e.wg.Add(1)
			go func(t *Task[T]) {
				task.Execute()
			}(task)
		}
	}
}
