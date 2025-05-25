package taskpool

import (
	"sync"
)

type TaskPool struct {
	wg        sync.WaitGroup
	semaphore *semaphore
}

func New(taskLimit int) *TaskPool {
	return &TaskPool{
		semaphore: newSemaphore(taskLimit),
	}
}

func (t *TaskPool) Close() {
	t.wg.Wait()
}

func (t *TaskPool) AddTask(task func()) {
	t.wg.Add(1)

	go func() {
		defer t.wg.Done()
		t.semaphore.Acquire()
		defer t.semaphore.Release()

		task()
	}()
}
