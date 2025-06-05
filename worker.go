package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"sync"

	"github.com/google/uuid"
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

func NewWorker(id uuid.UUID, tasks chan io.Reader, wg *sync.WaitGroup) *worker {
	return &worker{
		id:    id,
		tasks: tasks,
		wg:    wg,
	}
}

func (w *worker) Process(ctx context.Context) {
	for {
		select {
		case msg := <-w.tasks:
			if msg == nil {
				w.wg.Done()
				return
			}
			text := make([]byte, TextBuffSize)
			_, err := msg.Read(text)
			if err != nil {
				log.Printf("Worker #%s: read error: %v", w.id.String(), err)
				continue
			}
			fmt.Printf("Worker #%s: %s\n", w.id.String(), text)
		case <-ctx.Done():
			return
		}
	}
}
