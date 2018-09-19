// Copyright 2018 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package worker

import (
	"errors"
	"sync"
	"sync/atomic"
	"time"
)

// Errors messages
var (
	ErrPoolClosed = errors.New("pool is closed")
)

// Pool interface
type Pool interface {
	WorkerSize() int
	QueueSize() int
	QueueLength() int
	Start()
	Enqueue(payload interface{}) error
	Close()
}

// Handler of worker
type Handler func(payload interface{})

type pool struct {
	workerSize int
	queueSize  int
	handler    Handler
	queue      chan interface{}
	closed     int32
	workers    []*worker
	wg         sync.WaitGroup
}

// New pool worker
func New(workerSize int, queueSize int, handler Handler) Pool {
	return &pool{
		workerSize: workerSize,
		queueSize:  queueSize,
		handler:    handler,
		queue:      make(chan interface{}, queueSize),
		workers:    make([]*worker, 0, workerSize),
	}
}

func (p *pool) Start() {
	p.wg.Add(p.workerSize)

	for i := 0; i < p.workerSize; i++ {
		wkr := newWorker(p.queue, &p.wg, p.handler)

		wkr.run()

		p.workers = append(p.workers, wkr)
	}
}

// WorkerSize retruns number of worker
func (p *pool) WorkerSize() int {
	return p.workerSize
}

func (p *pool) QueueSize() int {
	return p.queueSize
}

func (p *pool) QueueLength() int {
	return len(p.queue)
}

func (p *pool) Enqueue(payload interface{}) error {
	if atomic.LoadInt32(&p.closed) != 0 {
		return ErrPoolClosed
	}

	p.queue <- payload

	return nil
}

func (p *pool) clean() {
	for _, w := range p.workers {
		w.close()
	}
}

func (p *pool) Close() {
	atomic.StoreInt32(&p.closed, int32(1))

	ticker := time.NewTicker(100 * time.Millisecond)

	for range ticker.C {
		if len(p.queue) == 0 {
			p.clean()

			ticker.Stop()

			break
		}
	}

	p.wg.Wait()

	p.workers = make([]*worker, 0, p.workerSize)

	close(p.queue)
}
