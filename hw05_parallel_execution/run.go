package hw05parallelexecution

import (
	"errors"
	"sync"
	"sync/atomic"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	allowedErrors := int64(m)
	ch := make(chan Task, n)
	wg := sync.WaitGroup{}

	for i := 0; i < n; i++ {
		wg.Add(1)
		go processTasks(ch, &allowedErrors, &wg)
	}

	for _, task := range tasks {
		ch <- task
	}
	close(ch)

	wg.Wait()

	if atomic.LoadInt64(&allowedErrors) <= 0 {
		return ErrErrorsLimitExceeded
	}

	return nil
}

func processTasks(ch <-chan Task, allowedErrors *int64, wg *sync.WaitGroup) {
	defer wg.Done()

	for task := range ch {
		if atomic.LoadInt64(allowedErrors) <= 0 {
			continue
		}

		if err := task(); err != nil {
			atomic.AddInt64(allowedErrors, -1)
		}
	}
}
