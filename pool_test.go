// Copyright 2018 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package worker

import (
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPool(t *testing.T) {
	gr := runtime.NumGoroutine()

	i := 0

	p := New(2, 4, func(payload interface{}) {
		i += payload.(int)
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

	assert.Equal(t, 5, i)

	// check gorouting leak
	assert.Equal(t, gr, runtime.NumGoroutine())
}
