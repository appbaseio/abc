/*
 * Copyright (c) 2002-2019 "Neo4j,"
 * Neo4j Sweden AB [http://neo4j.com]
 *
 * This file is part of Neo4j.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package gobolt

import (
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_WorkerPool(t *testing.T) {
	t.Run("newWorkerPool", func(t *testing.T) {
		t.Run("shouldUsePassedInParametersInConstruction", func(t *testing.T) {
			pool := newWorkerPool(2, 10, 1*time.Minute)
			defer pool.close()

			assert.Equal(t, 2, pool.minWorkers)
			assert.Equal(t, 10, pool.maxWorkers)
			assert.Equal(t, 1*time.Minute, pool.keepAlive)
		})

		t.Run("shouldPanicOnNegativeMinWorkers", func(t *testing.T) {
			assert.Panics(t, func() {
				_ = newWorkerPool(-1, 10, 1*time.Minute)
			})
		})

		t.Run("shouldPanicOnZeroMaxWorkers", func(t *testing.T) {
			assert.Panics(t, func() {
				_ = newWorkerPool(0, 0, 1*time.Minute)
			})
		})

		t.Run("shouldPanicOnInvalidMaxWorkers", func(t *testing.T) {
			assert.Panics(t, func() {
				_ = newWorkerPool(10, 5, 1*time.Minute)
			})
		})

		t.Run("shouldStartMinWorkers", func(t *testing.T) {
			pool := newWorkerPool(20, 50, 1*time.Minute)
			defer pool.close()

			assert.Equal(t, int32(20), pool.workerCount)
			assert.NotEmpty(t, pool.workers)
		})
	})

	t.Run("close", func(t *testing.T) {
		t.Run("shouldSetClosedOnClose", func(t *testing.T) {
			pool := newWorkerPool(0, 2, 1*time.Minute)

			assert.False(t, pool.isClosed())

			pool.close()

			assert.True(t, pool.isClosed())
		})

		t.Run("shouldWaitForWorkersOnClose", func(t *testing.T) {
			pool := newWorkerPool(20, 40, 1*time.Minute)

			assert.Equal(t, int32(20), pool.workerCount)
			assert.NotEmpty(t, pool.workers)

			pool.close()

			assert.Equal(t, int32(0), pool.workerCount)
			assert.Empty(t, pool.workers)
		})
	})

	t.Run("workerEntryPoint", func(t *testing.T) {
		t.Run("shouldDieAfterKeepalive", func(t *testing.T) {
			pool := newWorkerPool(4, 8, 3*time.Second)
			defer pool.close()

			assert.Equal(t, int32(4), pool.workerCount)
			assert.NotEmpty(t, pool.workers)

			<-time.After(5 * time.Second)

			assert.Equal(t, int32(0), pool.workerCount)
			assert.Empty(t, pool.workers)
		})
	})

	t.Run("submit", func(t *testing.T) {
		t.Run("shouldSucceedWhenThereIsAvailableWorker", func(t *testing.T) {
			var workExecutions int32
			var workSignal = make(chan bool, 1)
			workFunc := func(stopper <-chan signal) {
				atomic.AddInt32(&workExecutions, 1)
				workSignal <- true
			}

			pool := newWorkerPool(1, 1, 5*time.Minute)
			defer pool.close()

			assert.NoError(t, pool.submit(workFunc))

			<-workSignal

			assert.Equal(t, int32(1), workExecutions)
		})

		t.Run("shouldSucceedBySpawningNewWorker", func(t *testing.T) {
			var workExecutions int32
			var workSignal = make(chan bool, 1)
			workFunc := func(stopper <-chan signal) {
				atomic.AddInt32(&workExecutions, 1)
				workSignal <- true
			}

			pool := newWorkerPool(0, 1, 5*time.Minute)
			defer pool.close()

			assert.NoError(t, pool.submit(workFunc))

			<-workSignal

			assert.Equal(t, int32(1), workExecutions)
		})

		t.Run("shouldFailWhenClosed", func(t *testing.T) {
			var workSignal = make(chan bool, 1)
			workFunc := func(stopper <-chan signal) {
				<-workSignal
			}

			pool := newWorkerPool(1, 5, 5*time.Minute)
			pool.submit(workFunc)

			go func() {
				pool.close()
			}()

			for !pool.isClosed() {
				time.Sleep(500 * time.Millisecond)
			}

			assert.EqualError(t, pool.submit(workFunc), "unable to submit job to a closed worker pool")

			close(workSignal)
		})
	})

}
