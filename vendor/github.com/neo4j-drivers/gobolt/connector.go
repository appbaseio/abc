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
)

// AccessMode is used by the routing driver to decide if a transaction should be routed to a write server
// or a read server in a cluster. When running a transaction, a write transaction requires a server that
// supports writes. A read transaction, on the other hand, requires a server that supports read operations.
// This classification is key for routing driver to route transactions to a cluster correctly.
type AccessMode int

const (
	// AccessModeWrite makes the driver return a session towards a write server
	AccessModeWrite AccessMode = 0
	// AccessModeRead makes the driver return a session towards a follower or a read-replica
	AccessModeRead AccessMode = 1
)

// Connector represents an initialised seabolt connector
type Connector interface {
	Acquire(mode AccessMode) (Connection, error)
	Close() error
}

func NewConnector(uri *url.URL, authToken map[string]interface{}, config *Config) (Connector, error) {
	if workersEnabled() {
		return newWorkerConnector(uri, authToken, config)
	} else {
		return newSeaboltConnector(uri, authToken, config)
	}
}
