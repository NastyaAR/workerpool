package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"sync"

	"github.com/google/uuid"
)

var (
	ErrSoMuchWorkers = errors.New("number of workers so much")
)

type Workerpool interface {
	AddWorker() error
	DeleteWorker()
	Submit(r io.Reader) error
	Free()
	Wait()
	Close()
}

type workerPool struct {
	mtx sync.Mutex
	wg  *sync.WaitGroup

	numWorkers uint64
	maxWorkers uint64

	tasks chan io.Reader
}

func NewWorkerPool(
	ctx context.Context,
	initialWorkersNum uint64,
	maxWorkersNum uint64,
	maxTasks uint64,
) *workerPool {
	tasks := make(chan io.Reader, maxTasks)
	wg := sync.WaitGroup{}

	for range initialWorkersNum {
		uid, _ := uuid.NewV7()
		wrk := NewWorker(uid, tasks, &wg)
		go wrk.Process(ctx)
		wg.Add(1)
	}

	return &workerPool{
		mtx:        sync.Mutex{},
		wg:         &wg,
		maxWorkers: maxWorkersNum,
		tasks:      tasks,
		numWorkers: initialWorkersNum,
	}
}

func (w *workerPool) AddWorker(ctx context.Context) error {
	w.mtx.Lock()
	defer w.mtx.Unlock()

	if w.numWorkers >= w.maxWorkers {
		return ErrSoMuchWorkers
	}

	uid, err := uuid.NewV7()
	if err != nil {
		return fmt.Errorf("failed to gen uuid: %w", err)
	}

	newWorker := NewWorker(uid, w.tasks, w.wg)
	w.numWorkers++

	go newWorker.Process(ctx)
	w.wg.Add(1)

	return nil
}

func (w *workerPool) DeleteWorker() {
	w.mtx.Lock()
	defer w.mtx.Unlock()

	if w.numWorkers == 0 {
		return
	}

	w.tasks <- nil
	w.numWorkers--
}

func (w *workerPool) Free() {
	w.mtx.Lock()
	defer w.mtx.Unlock()

	if w.tasks == nil {
		return
	}

	for range w.numWorkers {
		w.tasks <- nil
		w.numWorkers--
	}

	close(w.tasks)
	w.tasks = nil
}

func (w *workerPool) Submit(r io.Reader) error {
	w.tasks <- r
	return nil
}

func (w *workerPool) Wait() {
	fmt.Println("start wait in wait")
	w.wg.Wait()
	fmt.Println("stop wait in wait")
}

func (w *workerPool) Close() error {
	w.Free()
	fmt.Println("start wait in close")
	w.Wait()
	fmt.Println("stop wait in close")
	return nil
}
