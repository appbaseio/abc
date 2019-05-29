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
	"net/url"
	"os"
	"strconv"
	"time"
)

type workerConnector struct {
	config      Config
	delegate    *seaboltConnector
	pool        *workerPool
	closeSignal chan signal
}

func (w *workerConnector) Acquire(mode AccessMode) (Connection, error) {
	return newWorkerConnection(w, mode)
}

func (w *workerConnector) Close() error {
	var err error

	var done = make(chan bool, 1)
	if poolErr := w.pool.submit(func(stopper <-chan signal) {
		err = w.delegate.Close()
		done <- true
	}); poolErr != nil {
		err = poolErr
		done <- true
	}

	<-done

	if err != nil {
		return err
	}

	w.pool.close()
	close(w.closeSignal)

	return nil
}

func newWorkerConnector(url *url.URL, authToken map[string]interface{}, config *Config) (Connector, error) {
	var err error
	var connector *seaboltConnector
	var pool = newWorkerPool(minWorkers(config), maxWorkers(config), keepAlive(config))

	var configOverride = *config
	configOverride.ConnAcquisitionTimeout = 0

	var done = make(chan bool, 1)
	if poolErr := pool.submit(func(stopper <-chan signal) {
		connector, err = newSeaboltConnector(url, authToken, &configOverride)
		done <- true
	}); poolErr != nil {
		err = poolErr
		done <- true
	}

	// wait for connector creation to complete
	<-done

	if err != nil {
		defer pool.close()
		return nil, err
	}

	return &workerConnector{
		config:      *config,
		delegate:    connector,
		pool:        pool,
		closeSignal: make(chan signal, config.MaxPoolSize),
	}, nil
}

func workersEnabled() bool {
	var workersEnabled = true
	if val, ok := os.LookupEnv("BOLTWORKERS"); ok {
		if parsed, err := strconv.ParseBool(val); err == nil {
			workersEnabled = parsed
		}
	}
	return workersEnabled
}

func maxWorkers(config *Config) int {
	var workersMax = int(float64(config.MaxPoolSize) * float64(1.2))
	if val, ok := os.LookupEnv("BOLTWORKERSMAX"); ok {
		if parsed, err := strconv.ParseInt(val, 10, 32); err == nil {
			workersMax = int(parsed)
		}
	}
	return workersMax
}

func minWorkers(config *Config) int {
	var workersMin = 0
	if val, ok := os.LookupEnv("BOLTWORKERSMIN"); ok {
		if parsed, err := strconv.ParseInt(val, 10, 32); err == nil {
			workersMin = int(parsed)
		}
	}
	return workersMin
}

func keepAlive(config *Config) time.Duration {
	var workersKeepAlive = 5 * time.Minute
	if val, ok := os.LookupEnv("BOLTWORKERSKEEPALIVE"); ok {
		if parsed, err := time.ParseDuration(val); err == nil {
			workersKeepAlive = parsed
		}
	}
	return workersKeepAlive
}
