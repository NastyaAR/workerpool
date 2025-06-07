package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"sync"
)

var (
	ErrSoMuchWorkers  = errors.New("number of workers so much")
	ErrNoWorkers      = errors.New("no workers")
	ErrNilTask        = errors.New("nil task cant be submitted")
	ErrNilTaskChannel = errors.New("nil task channel")
	ErrBadWorkersNums = errors.New("initial workers num cant be less then max workers num")
)

type WorkerPool interface {
	AddWorker() error
	DeleteWorker()
	Submit(r io.Reader) error
	Free()
	Wait()
}

type workerPool struct {
	mtx        sync.Mutex
	wg         *sync.WaitGroup
	numWorkers uint64
	maxWorkers uint64
	tasks      chan io.Reader
}

func NewWorkerPool(
	ctx context.Context,
	initialWorkersNum uint64,
	maxWorkersNum uint64,
	maxTasks uint64,
) (*workerPool, error) {
	tasks := make(chan io.Reader, maxTasks)
	wg := sync.WaitGroup{}

	if initialWorkersNum > maxWorkersNum {
		return nil, ErrBadWorkersNums
	}

	var numWorkers uint64
	for range initialWorkersNum {
		wrk, err := NewWorker(tasks, &wg)
		if err != nil {
			free(tasks, uint64(numWorkers))
			return nil, fmt.Errorf("failed to add worker: %w", err)
		}
		go wrk.Process(ctx)
		numWorkers++
	}

	return &workerPool{
		mtx:        sync.Mutex{},
		wg:         &wg,
		maxWorkers: maxWorkersNum,
		tasks:      tasks,
		numWorkers: initialWorkersNum}, nil
}

func (w *workerPool) AddWorker(ctx context.Context) error {
	w.mtx.Lock()
	defer w.mtx.Unlock()
	if w.numWorkers >= w.maxWorkers {
		return ErrSoMuchWorkers
	}

	newWorker, err := NewWorker(w.tasks, w.wg)
	if err != nil {
		return fmt.Errorf("failed to add worker: %w", err)
	}

	go newWorker.Process(ctx)
	w.numWorkers++

	return nil
}
func (w *workerPool) DeleteWorker() error {
	w.mtx.Lock()
	defer w.mtx.Unlock()

	if w.tasks == nil {
		return ErrNilTaskChannel
	}

	if w.numWorkers == 0 {
		return ErrNoWorkers
	}

	w.tasks <- nil
	w.numWorkers--

	return nil
}

func free(tasks chan io.Reader, numWorkers uint64) {
	if tasks == nil {
		return
	}

	for range numWorkers {
		tasks <- nil
	}
}

func (w *workerPool) Free() {
	w.mtx.Lock()
	defer w.mtx.Unlock()
	w.Wait()

	if w.tasks == nil {
		return
	}

	for range w.numWorkers {
		w.tasks <- nil
		w.numWorkers--
	}
}

func (w *workerPool) Submit(r io.Reader) error {
	if w.tasks == nil {
		return ErrNilTaskChannel
	}

	if r == nil {
		return ErrNilTask
	}

	w.tasks <- r
	w.wg.Add(1)

	return nil
}

func (w *workerPool) Wait() {
	w.wg.Wait()
}
