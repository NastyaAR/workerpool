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
	ErrGenWorkerID = errors.New("failed to gen uuid")
)

const (
	TextBuffSize = 1024
)

type Worker interface {
	Process(ctx context.Context)
}

type worker struct {
	id    uuid.UUID
	tasks chan io.Reader
	wg    *sync.WaitGroup
}

func NewWorker(tasks chan io.Reader, wg *sync.WaitGroup) (*worker, error) {
	uid, err := uuid.NewV7()
	if err != nil {
		return nil, fmt.Errorf("failed to create worker: %w", ErrGenWorkerID)
	}

	return &worker{
		id: uid, tasks: tasks,
		wg: wg}, nil
}

func (w *worker) Process(ctx context.Context) {
	for {
		select {
		case msg := <-w.tasks:
			if msg == nil {
				return
			}
			text := make([]byte, TextBuffSize)
			_, err := msg.Read(text)
			if err != nil {
				continue
			}
			fmt.Printf("Worker #%s: %s\n", w.id.String(), text)
			w.wg.Done()
		case <-ctx.Done():
			return
		}
	}
}
