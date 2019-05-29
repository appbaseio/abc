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

/*
#include <stdlib.h>

#include "bolt/bolt.h"
*/
import "C"
import (
	"errors"
	"net/url"
	"reflect"
	"unsafe"
)

type seaboltConnector struct {
	key int

	uri       *url.URL
	authToken map[string]interface{}
	config    Config

	cAddress    *C.BoltAddress
	cInstance   *C.BoltConnector
	valueSystem *boltValueSystem
}

func (conn *seaboltConnector) Close() error {
	if conn.cInstance != nil {
		C.BoltConnector_destroy(conn.cInstance)
		conn.cInstance = nil
	}

	if conn.cAddress != nil {
		C.BoltAddress_destroy(conn.cAddress)
		conn.cAddress = nil
	}

	unregisterLogging(conn.key)
	unregisterResolver(conn.key)

	shutdownLibrary()

	return nil
}

func (conn *seaboltConnector) Acquire(mode AccessMode) (Connection, error) {
	return newSeaboltConnection(conn, mode)
}

// NewConnector returns a new connector instance with given parameters
func newSeaboltConnector(uri *url.URL, authToken map[string]interface{}, config *Config) (*seaboltConnector, error) {
	var err error
	var key int
	var cAddress *C.struct_BoltAddress
	var valueSystem *boltValueSystem
	var cAuthToken *C.struct_BoltValue
	var cConfig *C.struct_BoltConfig

	if uri == nil {
		return nil, errors.New("provided uri should not be nil")
	}

	if config == nil {
		config = &Config{
			Encryption:  true,
			MaxPoolSize: 100,
		}
	}

	valueSystem = createValueSystem(config)
	cAddress = createAddress(uri)
	key = startupLibrary()

	if cAuthToken, err = valueSystem.valueToConnector(authToken); err != nil {
		return nil, valueSystem.genericErrorFactory("unable to convert authentication token: %v", err)
	}
	defer C.BoltValue_destroy(cAuthToken)

	if cConfig, err = createConfig(key, uri, config, valueSystem); err != nil {
		return nil, err
	}
	defer C.BoltConfig_destroy(cConfig)

	cInstance := C.BoltConnector_create(cAddress, cAuthToken, cConfig)
	conn := &seaboltConnector{
		key:         key,
		uri:         uri,
		authToken:   authToken,
		config:      *config,
		cAddress:    cAddress,
		valueSystem: valueSystem,
		cInstance:   cInstance,
	}

	return conn, nil
}

func createValueSystem(config *Config) *boltValueSystem {
	valueHandlersBySignature := make(map[int16]ValueHandler, len(config.ValueHandlers))
	valueHandlersByType := make(map[reflect.Type]ValueHandler, len(config.ValueHandlers))
	for _, handler := range config.ValueHandlers {
		for _, readSignature := range handler.ReadableStructs() {
			valueHandlersBySignature[readSignature] = handler
		}

		for _, writeType := range handler.WritableTypes() {
			valueHandlersByType[writeType] = handler
		}
	}

	databaseErrorFactory := newDatabaseError
	connectorErrorFactory := newConnectorError
	genericErrorFactory := newGenericError
	if config.DatabaseErrorFactory != nil {
		databaseErrorFactory = config.DatabaseErrorFactory
	}
	if config.ConnectorErrorFactory != nil {
		connectorErrorFactory = config.ConnectorErrorFactory
	}
	if config.GenericErrorFactory != nil {
		genericErrorFactory = config.GenericErrorFactory
	}

	return &boltValueSystem{
		valueHandlers:            config.ValueHandlers,
		valueHandlersBySignature: valueHandlersBySignature,
		valueHandlersByType:      valueHandlersByType,
		connectorErrorFactory:    connectorErrorFactory,
		databaseErrorFactory:     databaseErrorFactory,
		genericErrorFactory:      genericErrorFactory,
	}
}

func createAddress(uri *url.URL) *C.struct_BoltAddress {
	var hostname = C.CString(uri.Hostname())
	var port = C.CString(uri.Port())
	defer C.free(unsafe.Pointer(hostname))
	defer C.free(unsafe.Pointer(port))

	return C.BoltAddress_create(hostname, port)
}
