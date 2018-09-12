// Copyright 2018 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package worker

import (
	"runtime"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWorker(t *testing.T) {
	gr := runtime.NumGoroutine()

	var wg sync.WaitGroup

	queue := make(chan interface{}, 2)

	handler := func(payload interface{}) {
		assert.Equal(t, "foo", payload.(string))
	}

	wg.Add(1)

	w := newWorker(queue, &wg, handler)

	w.run()

	queue <- "foo"

	close(queue)

	wg.Wait()

	// check gorouting leak
	assert.Equal(t, gr, runtime.NumGoroutine())
}
