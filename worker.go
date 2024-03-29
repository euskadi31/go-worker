// Copyright 2018 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package worker

import (
	"sync"
)

type worker struct {
	queue    <-chan interface{}
	handler  Handler
	wg       *sync.WaitGroup
	shutdown chan struct{}
}

func newWorker(queue <-chan interface{}, wg *sync.WaitGroup, handler Handler) *worker {
	return &worker{
		queue:    queue,
		handler:  handler,
		wg:       wg,
		shutdown: make(chan struct{}, 1),
	}
}

func (w *worker) run() {
	go func() {
		defer w.wg.Done()

		for {
			select {
			case item := <-w.queue:
				if item != nil {
					w.handler(item)
				}
			case <-w.shutdown:
				return
			}
		}
	}()
}

func (w *worker) close() {
	w.shutdown <- struct{}{}
}
