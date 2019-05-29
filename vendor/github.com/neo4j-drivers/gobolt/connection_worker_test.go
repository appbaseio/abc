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
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_WorkerConnection(t *testing.T) {
	newMockedConnection := func() (*workerConnection, *mockConnection, func()) {
		pool := newWorkerPool(1, 1, 1*time.Minute)
		delegate := new(mockConnection)
		connection := &workerConnection{
			pool:     pool,
			delegate: delegate,
			active:   0,
		}

		return connection, delegate, func() {
			pool.close()
		}
	}

	newMockedConnectionInUse := func() (*workerConnection, *mockConnection, func()) {
		conn, mocked, cleanup := newMockedConnection()

		conn.active = 1

		return conn, mocked, cleanup
	}

	t.Run("newWorkerConnection", func(t *testing.T) {
		var valueSystem = createValueSystem(&Config{})
		var someOtherFailure = fmt.Errorf("some other failure")
		var poolFullError = newPoolFullError(valueSystem)
		var connectionAcquisitionTimedOutError = newConnectionAcquisitionTimedOutError(valueSystem)

		t.Run("shouldInvokeNewSeaboltConnection", func(t *testing.T) {
			var cases = []struct {
				name    string
				timeout time.Duration
			}{
				{"AcquisitionTimeout=0", 0},
				{"AcquisitionTimeout=5s", 5 * time.Second},
			}

			for _, testCase := range cases {
				t.Run(testCase.name, func(t *testing.T) {
					var originalNewSeaboltConnection = newSeaboltConnection
					defer func() {
						newSeaboltConnection = originalNewSeaboltConnection
					}()

					var newSeaboltConnectionCount = 0
					newSeaboltConnection = func(connector *seaboltConnector, mode AccessMode) (*seaboltConnection, error) {
						newSeaboltConnectionCount++
						return &seaboltConnection{}, nil
					}

					var connector = &workerConnector{
						config: Config{ConnAcquisitionTimeout: testCase.timeout},
						pool:   newWorkerPool(1, 1, 1*time.Minute),
					}
					defer connector.pool.close()

					connection, err := newWorkerConnection(connector, AccessModeRead)

					assert.NoError(t, err)
					assert.NotNil(t, connection)
					assert.Equal(t, 1, newSeaboltConnectionCount)
				})
			}
		})

		t.Run("shouldReturnErrorFromNewSeaboltConnection", func(t *testing.T) {
			var originalNewSeaboltConnection = newSeaboltConnection
			defer func() {
				newSeaboltConnection = originalNewSeaboltConnection
			}()

			var newSeaboltConnectionCount = 0
			newSeaboltConnection = func(connector *seaboltConnector, mode AccessMode) (*seaboltConnection, error) {
				newSeaboltConnectionCount++
				return nil, someOtherFailure
			}

			var connector = &workerConnector{
				config: Config{ConnAcquisitionTimeout: 0},
				pool:   newWorkerPool(1, 1, 1*time.Minute),
			}
			defer connector.pool.close()

			connection, err := newWorkerConnection(connector, AccessModeRead)

			assert.EqualError(t, err, "some other failure")
			assert.Nil(t, connection)
			assert.Equal(t, 1, newSeaboltConnectionCount)
		})

		t.Run("shouldReturnPoolFullErrorWhenAcquisitionTimeoutIsZero", func(t *testing.T) {
			var originalNewSeaboltConnection = newSeaboltConnection
			defer func() {
				newSeaboltConnection = originalNewSeaboltConnection
			}()

			var newSeaboltConnectionCount = 0
			newSeaboltConnection = func(connector *seaboltConnector, mode AccessMode) (*seaboltConnection, error) {
				newSeaboltConnectionCount++
				return nil, poolFullError
			}

			var connector = &workerConnector{
				config: Config{ConnAcquisitionTimeout: 0},
				pool:   newWorkerPool(1, 1, 1*time.Minute),
			}
			defer connector.pool.close()

			connection, err := newWorkerConnection(connector, AccessModeRead)

			assert.EqualError(t, err, poolFullError.Error())
			assert.Nil(t, connection)
			assert.Equal(t, 1, newSeaboltConnectionCount)
		})

		t.Run("shouldInvokeWaitClosedWhenPoolIsFullAndSucceedWhenWaitSucceeds", func(t *testing.T) {
			var originalNewSeaboltConnection = newSeaboltConnection
			var originalWaitClosed = waitClosed
			defer func() {
				newSeaboltConnection = originalNewSeaboltConnection
				waitClosed = originalWaitClosed
			}()

			var newSeaboltConnectionCount = 0
			newSeaboltConnection = func(connector *seaboltConnector, mode AccessMode) (*seaboltConnection, error) {
				newSeaboltConnectionCount++
				if newSeaboltConnectionCount > 1 {
					return &seaboltConnection{}, nil
				}
				return nil, poolFullError
			}

			var waitClosedCount = 0
			waitClosed = func(w *workerConnection, timeout time.Duration) bool {
				waitClosedCount++
				return true
			}

			var connector = &workerConnector{
				config: Config{ConnAcquisitionTimeout: 5 * time.Second},
				pool:   newWorkerPool(1, 1, 1*time.Minute),
			}
			defer connector.pool.close()

			connection, err := newWorkerConnection(connector, AccessModeRead)

			assert.NoError(t, err)
			assert.NotNil(t, connection)
			assert.Equal(t, 1, waitClosedCount)
			assert.Equal(t, 2, newSeaboltConnectionCount)
		})

		t.Run("shouldInvokeWaitClosedWhenPoolIsFullAndFailWhenWaitFails", func(t *testing.T) {
			var originalNewSeaboltConnection = newSeaboltConnection
			var originalWaitClosed = waitClosed
			defer func() {
				newSeaboltConnection = originalNewSeaboltConnection
				waitClosed = originalWaitClosed
			}()

			var newSeaboltConnectionCount = 0
			newSeaboltConnection = func(connector *seaboltConnector, mode AccessMode) (*seaboltConnection, error) {
				newSeaboltConnectionCount++
				return nil, poolFullError
			}

			var waitClosedCount = 0
			waitClosed = func(w *workerConnection, timeout time.Duration) bool {
				waitClosedCount++
				return false
			}

			var connector = &workerConnector{
				config: Config{ConnAcquisitionTimeout: 5 * time.Second},
				pool:   newWorkerPool(1, 1, 1*time.Minute),
				delegate: &seaboltConnector{
					valueSystem: valueSystem,
				},
			}
			defer connector.pool.close()

			connection, err := newWorkerConnection(connector, AccessModeRead)

			assert.EqualError(t, err, connectionAcquisitionTimedOutError.Error())
			assert.Nil(t, connection)
			assert.Equal(t, 1, waitClosedCount)
			assert.Equal(t, 1, newSeaboltConnectionCount)
		})
	})

	t.Run("shouldInvokeSignalClosedOnClose", func(t *testing.T) {
		var originalSignalClosed = signalClosed
		defer func() {
			signalClosed = originalSignalClosed
		}()

		var signalClosedCount = 0
		signalClosed = func(w *workerConnection) {
			signalClosedCount++
		}

		conn, delegate, cleanup := newMockedConnection()
		defer cleanup()

		delegate.On("Close").Return(nil)

		assert.NoError(t, conn.Close())
		assert.Equal(t, 1, signalClosedCount)
	})

	t.Run("shouldSurfacePoolFullErrorWhenAcquisitionTimeoutIsZero", func(t *testing.T) {

	})

	t.Run("shouldInterceptPoolFullErrorWhenAcquisitionTimeoutIsNotZero", func(t *testing.T) {

	})

	t.Run("shouldInvokeDelegate", func(t *testing.T) {
		failure := fmt.Errorf("some error")
		handle := RequestHandle(500)

		t.Run("Id", func(t *testing.T) {
			conn, delegate, cleanup := newMockedConnection()
			defer cleanup()

			delegate.On("Id").Return("123", failure)

			id, err := conn.Id()
			assert.Equal(t, "123", id)
			assert.Equal(t, failure, err)

			delegate.AssertExpectations(t)
		})

		t.Run("RemoteAddress", func(t *testing.T) {
			conn, delegate, cleanup := newMockedConnection()
			defer cleanup()

			delegate.On("RemoteAddress").Return("localhost:7687", failure)

			remoteAddress, err := conn.RemoteAddress()
			assert.Equal(t, "localhost:7687", remoteAddress)
			assert.Equal(t, failure, err)

			delegate.AssertExpectations(t)
		})

		t.Run("Server", func(t *testing.T) {
			conn, delegate, cleanup := newMockedConnection()
			defer cleanup()

			delegate.On("Server").Return("Neo4j/3.5.0", failure)

			server, err := conn.Server()
			assert.Equal(t, "Neo4j/3.5.0", server)
			assert.Equal(t, failure, err)

			delegate.AssertExpectations(t)
		})

		t.Run("Begin", func(t *testing.T) {
			conn, delegate, cleanup := newMockedConnection()
			defer cleanup()

			bookmarks := []string{"1", "2", "3"}
			txTimeout := 5 * time.Minute
			txMetadata := map[string]interface{}{"a": 1, "b": true, "c": "yes"}
			delegate.On("Begin", bookmarks, txTimeout, txMetadata).Return(handle, failure)

			beginHandle, err := conn.Begin(bookmarks, txTimeout, txMetadata)
			assert.Equal(t, handle, beginHandle)
			assert.Equal(t, failure, err)

			delegate.AssertExpectations(t)
		})

		t.Run("Commit", func(t *testing.T) {
			conn, delegate, cleanup := newMockedConnection()
			defer cleanup()

			delegate.On("Commit").Return(handle, failure)

			commitHandle, err := conn.Commit()
			assert.Equal(t, handle, commitHandle)
			assert.Equal(t, failure, err)

			delegate.AssertExpectations(t)
		})

		t.Run("Rollback", func(t *testing.T) {
			conn, delegate, cleanup := newMockedConnection()
			defer cleanup()

			delegate.On("Rollback").Return(handle, failure)

			rollbackHandle, err := conn.Rollback()
			assert.Equal(t, handle, rollbackHandle)
			assert.Equal(t, failure, err)

			delegate.AssertExpectations(t)
		})

		t.Run("Run", func(t *testing.T) {
			conn, delegate, cleanup := newMockedConnection()
			defer cleanup()

			cypher := "CREATE (n {id: $x})"
			parameters := map[string]interface{}{"id": 5000}
			bookmarks := []string{"1", "2", "3"}
			txTimeout := 5 * time.Minute
			txMetadata := map[string]interface{}{"a": 1, "b": true, "c": "yes"}
			delegate.On("Run", cypher, parameters, bookmarks, txTimeout, txMetadata).Return(handle, failure)

			runHandle, err := conn.Run(cypher, parameters, bookmarks, txTimeout, txMetadata)
			assert.Equal(t, handle, runHandle)
			assert.Equal(t, failure, err)

			delegate.AssertExpectations(t)
		})

		t.Run("PullAll", func(t *testing.T) {
			conn, delegate, cleanup := newMockedConnection()
			defer cleanup()

			delegate.On("PullAll").Return(handle, failure)

			pullAllHandle, err := conn.PullAll()
			assert.Equal(t, handle, pullAllHandle)
			assert.Equal(t, failure, err)

			delegate.AssertExpectations(t)
		})

		t.Run("DiscardAll", func(t *testing.T) {
			conn, delegate, cleanup := newMockedConnection()
			defer cleanup()

			delegate.On("DiscardAll").Return(handle, failure)

			discardAllHandle, err := conn.DiscardAll()
			assert.Equal(t, handle, discardAllHandle)
			assert.Equal(t, failure, err)

			delegate.AssertExpectations(t)
		})

		t.Run("Reset", func(t *testing.T) {
			conn, delegate, cleanup := newMockedConnection()
			defer cleanup()

			delegate.On("Reset").Return(handle, failure)

			resetHandle, err := conn.Reset()
			assert.Equal(t, handle, resetHandle)
			assert.Equal(t, failure, err)

			delegate.AssertExpectations(t)
		})

		t.Run("Flush", func(t *testing.T) {
			conn, delegate, cleanup := newMockedConnection()
			defer cleanup()

			delegate.On("Flush").Return(failure)

			err := conn.Flush()
			assert.Equal(t, failure, err)

			delegate.AssertExpectations(t)
		})

		t.Run("Fetch", func(t *testing.T) {
			conn, delegate, cleanup := newMockedConnection()
			defer cleanup()

			delegate.On("Fetch", handle).Return(FetchTypeRecord, failure)

			fetched, err := conn.Fetch(handle)
			assert.Equal(t, FetchTypeRecord, fetched)
			assert.Equal(t, failure, err)

			delegate.AssertExpectations(t)
		})

		t.Run("FetchSummary", func(t *testing.T) {
			conn, delegate, cleanup := newMockedConnection()
			defer cleanup()

			delegate.On("FetchSummary", handle).Return(50, failure)

			records, err := conn.FetchSummary(handle)
			assert.Equal(t, 50, records)
			assert.Equal(t, failure, err)

			delegate.AssertExpectations(t)
		})

		t.Run("LastBookmark", func(t *testing.T) {
			conn, delegate, cleanup := newMockedConnection()
			defer cleanup()

			delegate.On("LastBookmark").Return("bookmark:1234", failure)

			bookmark, err := conn.LastBookmark()
			assert.Equal(t, "bookmark:1234", bookmark)
			assert.Equal(t, failure, err)

			delegate.AssertExpectations(t)
		})

		t.Run("Fields", func(t *testing.T) {
			conn, delegate, cleanup := newMockedConnection()
			defer cleanup()

			fields := []string{"x", "y", "z"}
			delegate.On("Fields").Return(fields, failure)

			fieldsReturned, err := conn.Fields()
			assert.Equal(t, fields, fieldsReturned)
			assert.Equal(t, failure, err)

			delegate.AssertExpectations(t)
		})

		t.Run("Metadata", func(t *testing.T) {
			conn, delegate, cleanup := newMockedConnection()
			defer cleanup()

			metadata := map[string]interface{}{"x": 1, "y": "a", "z": false}
			delegate.On("Metadata").Return(metadata, failure)

			metadataReturned, err := conn.Metadata()
			assert.Equal(t, metadata, metadataReturned)
			assert.Equal(t, failure, err)

			delegate.AssertExpectations(t)
		})

		t.Run("Data", func(t *testing.T) {
			conn, delegate, cleanup := newMockedConnection()
			defer cleanup()

			data := []interface{}{"1", 2, false}
			delegate.On("Data").Return(data, failure)

			dataReturned, err := conn.Data()
			assert.Equal(t, data, dataReturned)
			assert.Equal(t, failure, err)

			delegate.AssertExpectations(t)
		})

		t.Run("Close", func(t *testing.T) {
			conn, delegate, cleanup := newMockedConnection()
			defer cleanup()

			delegate.On("Close").Return(failure)

			err := conn.Close()
			assert.Equal(t, failure, err)

			delegate.AssertExpectations(t)
		})

	})

	t.Run("shouldPropagateWorkerError", func(t *testing.T) {
		errText := "a connection is not thread-safe and thus should not be used concurrently"
		t.Run("Id", func(t *testing.T) {
			conn, _, cleanup := newMockedConnectionInUse()
			defer cleanup()

			_, err := conn.Id()
			assert.EqualError(t, err, errText)
		})

		t.Run("RemoteAddress", func(t *testing.T) {
			conn, _, cleanup := newMockedConnectionInUse()
			defer cleanup()

			_, err := conn.RemoteAddress()
			assert.EqualError(t, err, errText)
		})

		t.Run("Server", func(t *testing.T) {
			conn, _, cleanup := newMockedConnectionInUse()
			defer cleanup()

			_, err := conn.Server()
			assert.EqualError(t, err, errText)
		})

		t.Run("Begin", func(t *testing.T) {
			conn, _, cleanup := newMockedConnectionInUse()
			defer cleanup()

			_, err := conn.Begin([]string{}, 1*time.Minute, nil)
			assert.EqualError(t, err, errText)
		})

		t.Run("Commit", func(t *testing.T) {
			conn, _, cleanup := newMockedConnectionInUse()
			defer cleanup()

			_, err := conn.Commit()
			assert.EqualError(t, err, errText)
		})

		t.Run("Rollback", func(t *testing.T) {
			conn, _, cleanup := newMockedConnectionInUse()
			defer cleanup()

			_, err := conn.Rollback()
			assert.EqualError(t, err, errText)
		})

		t.Run("Run", func(t *testing.T) {
			conn, _, cleanup := newMockedConnectionInUse()
			defer cleanup()

			_, err := conn.Run("RETURN 1", nil, nil, 1*time.Second, nil)
			assert.EqualError(t, err, errText)
		})

		t.Run("PullAll", func(t *testing.T) {
			conn, _, cleanup := newMockedConnectionInUse()
			defer cleanup()

			_, err := conn.PullAll()
			assert.EqualError(t, err, errText)
		})

		t.Run("DiscardAll", func(t *testing.T) {
			conn, _, cleanup := newMockedConnectionInUse()
			defer cleanup()

			_, err := conn.DiscardAll()
			assert.EqualError(t, err, errText)
		})

		t.Run("Reset", func(t *testing.T) {
			conn, _, cleanup := newMockedConnectionInUse()
			defer cleanup()

			_, err := conn.Reset()
			assert.EqualError(t, err, errText)
		})

		t.Run("Flush", func(t *testing.T) {
			conn, _, cleanup := newMockedConnectionInUse()
			defer cleanup()

			err := conn.Flush()
			assert.EqualError(t, err, errText)
		})

		t.Run("Fetch", func(t *testing.T) {
			conn, _, cleanup := newMockedConnectionInUse()
			defer cleanup()

			_, err := conn.Fetch(RequestHandle(1))
			assert.EqualError(t, err, errText)
		})

		t.Run("FetchSummary", func(t *testing.T) {
			conn, _, cleanup := newMockedConnectionInUse()
			defer cleanup()

			_, err := conn.FetchSummary(RequestHandle(1))
			assert.EqualError(t, err, errText)
		})

		t.Run("LastBookmark", func(t *testing.T) {
			conn, _, cleanup := newMockedConnectionInUse()
			defer cleanup()

			_, err := conn.LastBookmark()
			assert.EqualError(t, err, errText)
		})

		t.Run("Fields", func(t *testing.T) {
			conn, _, cleanup := newMockedConnectionInUse()
			defer cleanup()

			_, err := conn.Fields()
			assert.EqualError(t, err, errText)
		})

		t.Run("Metadata", func(t *testing.T) {
			conn, _, cleanup := newMockedConnectionInUse()
			defer cleanup()

			_, err := conn.Metadata()
			assert.EqualError(t, err, errText)
		})

		t.Run("Data", func(t *testing.T) {
			conn, _, cleanup := newMockedConnectionInUse()
			defer cleanup()

			_, err := conn.Data()
			assert.EqualError(t, err, errText)
		})

		t.Run("Close", func(t *testing.T) {
			conn, _, cleanup := newMockedConnectionInUse()
			defer cleanup()

			err := conn.Close()
			assert.EqualError(t, err, errText)
		})

	})

	t.Run("queueJob", func(t *testing.T) {
		t.Run("shouldSetReceivingToOneWhenExecuting", func(t *testing.T) {
			var startEvent = make(chan bool, 1)
			var waitEvent = make(chan bool, 1)
			var blockingJob = func() {
				startEvent <- true
				<-waitEvent
			}

			conn, _, cleanup := newMockedConnection()
			defer cleanup()
			defer close(waitEvent)
			defer close(startEvent)

			go conn.queueJob(blockingJob)

			<-startEvent

			assert.Equal(t, int32(1), conn.active)
		})

		t.Run("shouldSetReceivingToZeroWhenExecutionIsComplete", func(t *testing.T) {
			conn, _, cleanup := newMockedConnection()
			defer cleanup()

			conn.queueJob(func() {})

			assert.Equal(t, int32(0), conn.active)
		})

		t.Run("shouldCheckForConcurrentAccess", func(t *testing.T) {
			var startEvent = make(chan bool, 1)
			var waitEvent = make(chan bool, 1)
			var blockingJob = func() {
				startEvent <- true
				<-waitEvent
			}

			conn, _, cleanup := newMockedConnection()
			defer cleanup()
			defer close(waitEvent)
			defer close(startEvent)

			go conn.queueJob(blockingJob)

			<-startEvent

			err := conn.queueJob(blockingJob)

			assert.EqualError(t, err, "a connection is not thread-safe and thus should not be used concurrently")
		})
	})
}
