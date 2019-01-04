// Copyright 2018 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package worker

import (
	"runtime"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestPool(t *testing.T) {
	gr := runtime.NumGoroutine()

	var wg sync.WaitGroup

	i := int64(0)

	wg.Add(5)
	p := New(2, 4, func(payload interface{}) {
		defer wg.Done()

		atomic.AddInt64(&i, int64(payload.(int)))
	})

	assert.Equal(t, 2, p.WorkerSize())
	assert.Equal(t, 4, p.QueueSize())
	assert.Equal(t, 0, p.QueueLength())

	p.Start()

	err := p.Enqueue(1)
	assert.NoError(t, err)

	err = p.Enqueue(1)
	assert.NoError(t, err)

	err = p.Enqueue(1)
	assert.NoError(t, err)

	err = p.Enqueue(1)
	assert.NoError(t, err)

	err = p.Enqueue(1)
	assert.NoError(t, err)

	p.Close()

	err = p.Enqueue(1)
	assert.Error(t, err)

	wg.Wait()

	assert.Equal(t, int64(5), i)

	// check gorouting leak
	assert.Equal(t, gr, runtime.NumGoroutine())
}

func TestPoolClose(t *testing.T) {
	//gr := runtime.NumGoroutine()

	i := int64(0)

	p := New(2, 400, func(payload interface{}) {
		atomic.AddInt64(&i, int64(payload.(int)))
	})

	p.Start()

	go func() {
		for index := 0; index < 500000; index++ {
			p.Enqueue(index)
		}
	}()

	go func() {
		for index := 500001; index < 1000000; index++ {
			p.Enqueue(index)
		}
	}()

	timber := time.NewTimer(1 * time.Second)

	<-timber.C

	p.Close()

	assert.True(t, i != 0)
	assert.False(t, i == 100000)

	// check gorouting leak
	// fail with -race
	//assert.Equal(t, gr, runtime.NumGoroutine())
}
