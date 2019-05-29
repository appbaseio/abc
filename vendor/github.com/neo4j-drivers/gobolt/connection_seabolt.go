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
	"fmt"
	"time"
	"unsafe"
)

type seaboltConnection struct {
	connector   *seaboltConnector
	cInstance   *C.struct_BoltConnection
	valueSystem *boltValueSystem
}

var newSeaboltConnection = func(connector *seaboltConnector, mode AccessMode) (*seaboltConnection, error) {
	var cMode uint32 = C.BOLT_ACCESS_MODE_WRITE
	if mode == AccessModeRead {
		cMode = C.BOLT_ACCESS_MODE_READ
	}

	cStatus := C.BoltStatus_create()
	defer C.BoltStatus_destroy(cStatus)
	cConnection := C.BoltConnector_acquire(connector.cInstance, C.BoltAccessMode(cMode), cStatus)
	if cConnection == nil {
		state := C.BoltStatus_get_state(cStatus)
		code := C.BoltStatus_get_error(cStatus)
		codeText := C.GoString(C.BoltError_get_string(code))
		context := C.GoString(C.BoltStatus_get_error_context(cStatus))

		return nil, newConnectorError(int(state), int(code), codeText, context, "unable to acquire connection from connector")
	}

	return &seaboltConnection{connector: connector, cInstance: cConnection, valueSystem: connector.valueSystem}, nil
}

func (connection *seaboltConnection) Id() (string, error) {
	return C.GoString(C.BoltConnection_id(connection.cInstance)), nil
}

func (connection *seaboltConnection) RemoteAddress() (string, error) {
	connectedAddress := C.BoltConnection_remote_endpoint(connection.cInstance)
	if connectedAddress == nil {
		return "UNKNOWN", nil
	}

	return fmt.Sprintf("%s:%s", C.GoString(C.BoltAddress_host(connectedAddress)), C.GoString(C.BoltAddress_port(connectedAddress))), nil
}

func (connection *seaboltConnection) Server() (string, error) {
	server := C.BoltConnection_server(connection.cInstance)
	if server == nil {
		return "UNKNOWN", nil
	}

	return C.GoString(server), nil
}

func (connection *seaboltConnection) Begin(bookmarks []string, txTimeout time.Duration, txMetadata map[string]interface{}) (RequestHandle, error) {
	var res C.int32_t

	res = C.BoltConnection_clear_begin(connection.cInstance)
	if res != C.BOLT_SUCCESS {
		return -1, newError(connection, "unable to clear begin message")
	}

	if len(bookmarks) > 0 {
		bookmarksValue, err := connection.valueSystem.valueToConnector(bookmarks)
		if err != nil {
			return -1, connection.valueSystem.genericErrorFactory("unable to convert bookmarks to connector value for begin message: %v", err)
		}
		res := C.BoltConnection_set_begin_bookmarks(connection.cInstance, bookmarksValue)
		C.BoltValue_destroy(bookmarksValue)
		if res != C.BOLT_SUCCESS {
			return -1, newError(connection, "unable to set bookmarks for begin message")
		}
	}

	if txTimeout > 0 {
		timeOut := C.int64_t(txTimeout / time.Millisecond)
		res := C.BoltConnection_set_begin_tx_timeout(connection.cInstance, timeOut)
		if res != C.BOLT_SUCCESS {
			return -1, newError(connection, "unable to set tx timeout for begin message")
		}
	}

	if len(txMetadata) > 0 {
		metadataValue, err := connection.valueSystem.valueToConnector(txMetadata)
		if err != nil {
			return -1, connection.valueSystem.genericErrorFactory("unable to convert tx metadata to connector value for begin message: %v", err)
		}
		res := C.BoltConnection_set_begin_tx_metadata(connection.cInstance, metadataValue)
		C.BoltValue_destroy(metadataValue)
		if res != C.BOLT_SUCCESS {
			return -1, newError(connection, "unable to set tx metadata for begin message")
		}
	}

	res = C.BoltConnection_load_begin_request(connection.cInstance)
	if res != C.BOLT_SUCCESS {
		return -1, newError(connection, "unable to generate begin message")
	}

	return RequestHandle(C.BoltConnection_last_request(connection.cInstance)), nil
}

func (connection *seaboltConnection) Commit() (RequestHandle, error) {
	res := C.BoltConnection_load_commit_request(connection.cInstance)
	if res != C.BOLT_SUCCESS {
		return -1, newError(connection, "unable to generate commit message")
	}

	return RequestHandle(C.BoltConnection_last_request(connection.cInstance)), nil
}

func (connection *seaboltConnection) Rollback() (RequestHandle, error) {
	res := C.BoltConnection_load_rollback_request(connection.cInstance)
	if res != C.BOLT_SUCCESS {
		return -1, newError(connection, "unable to generate rollback message")
	}

	return RequestHandle(C.BoltConnection_last_request(connection.cInstance)), nil
}

func (connection *seaboltConnection) Run(cypher string, params map[string]interface{}, bookmarks []string, txTimeout time.Duration, txMetadata map[string]interface{}) (RequestHandle, error) {
	var res C.int32_t

	res = C.BoltConnection_clear_run(connection.cInstance)
	if res != C.BOLT_SUCCESS {
		return -1, newError(connection, "unable to clear run message")
	}

	cypherStr := C.CString(cypher)
	res = C.BoltConnection_set_run_cypher(connection.cInstance, cypherStr, C.uint64_t(len(cypher)), C.int32_t(len(params)))
	C.free(unsafe.Pointer(cypherStr))
	if res != C.BOLT_SUCCESS {
		return -1, newError(connection, "unable to set cypher statement")
	}

	var index C.int32_t
	for paramName, paramValue := range params {
		paramNameLen := C.uint64_t(len(paramName))
		paramNameStr := C.CString(paramName)

		boltValue := C.BoltConnection_set_run_cypher_parameter(connection.cInstance, index, paramNameStr, paramNameLen)
		C.free(unsafe.Pointer(paramNameStr))
		if boltValue == nil {
			return -1, newError(connection, "unable to retrieve reference to cypher statement parameter value")
		}

		if err := connection.valueSystem.valueAsConnector(boltValue, paramValue); err != nil {
			return -1, connection.valueSystem.genericErrorFactory("unable to convert parameter %q to connector value for run message: %v", paramName, err)
		}

		index++
	}

	if len(bookmarks) > 0 {
		bookmarksValue, err := connection.valueSystem.valueToConnector(bookmarks)
		if err != nil {
			return -1, connection.valueSystem.genericErrorFactory("unable to convert bookmarks to connector value for run message: %v", err)
		}
		res := C.BoltConnection_set_run_bookmarks(connection.cInstance, bookmarksValue)
		C.BoltValue_destroy(bookmarksValue)
		if res != C.BOLT_SUCCESS {
			return -1, newError(connection, "unable to set bookmarks for run message")
		}
	}

	if txTimeout > 0 {
		timeOut := C.int64_t(txTimeout / time.Millisecond)
		res := C.BoltConnection_set_run_tx_timeout(connection.cInstance, timeOut)
		if res != C.BOLT_SUCCESS {
			return -1, newError(connection, "unable to set tx timeout for run message")
		}
	}

	if len(txMetadata) > 0 {
		metadataValue, err := connection.valueSystem.valueToConnector(txMetadata)
		if err != nil {
			return -1, connection.valueSystem.genericErrorFactory("unable to convert tx metadata to connector value for run message: %v", err)
		}
		res := C.BoltConnection_set_run_tx_metadata(connection.cInstance, metadataValue)
		C.BoltValue_destroy(metadataValue)
		if res != C.BOLT_SUCCESS {
			return -1, newError(connection, "unable to set tx metadata for run message")
		}
	}

	res = C.BoltConnection_load_run_request(connection.cInstance)
	if res != C.BOLT_SUCCESS {
		return -1, newError(connection, "unable to generate run message")
	}

	return RequestHandle(C.BoltConnection_last_request(connection.cInstance)), nil
}

func (connection *seaboltConnection) PullAll() (RequestHandle, error) {
	res := C.BoltConnection_load_pull_request(connection.cInstance, -1)
	if res != C.BOLT_SUCCESS {
		return -1, newError(connection, "unable to generate pullall message")
	}
	return RequestHandle(C.BoltConnection_last_request(connection.cInstance)), nil
}

func (connection *seaboltConnection) DiscardAll() (RequestHandle, error) {
	res := C.BoltConnection_load_discard_request(connection.cInstance, -1)
	if res != C.BOLT_SUCCESS {
		return -1, newError(connection, "unable to generate discardall message")
	}
	return RequestHandle(C.BoltConnection_last_request(connection.cInstance)), nil
}

func (connection *seaboltConnection) assertReadyState() error {
	cStatus := C.BoltConnection_status(connection.cInstance)

	if C.BoltStatus_get_state(cStatus) != C.BOLT_CONNECTION_STATE_READY {
		return newError(connection, "unexpected connection state")
	}

	return nil
}

func (connection *seaboltConnection) Flush() error {
	res := C.BoltConnection_send(connection.cInstance)
	if res < 0 {
		return newError(connection, "unable to flush")
	}

	return connection.assertReadyState()
}

func (connection *seaboltConnection) Fetch(request RequestHandle) (FetchType, error) {
	res := C.BoltConnection_fetch(connection.cInstance, C.BoltRequest(request))

	if err := connection.assertReadyState(); err != nil {
		return FetchTypeError, err
	}

	return FetchType(res), nil
}

func (connection *seaboltConnection) FetchSummary(request RequestHandle) (int, error) {
	res := C.BoltConnection_fetch_summary(connection.cInstance, C.BoltRequest(request))
	if res < 0 {
		return -1, newError(connection, "unable to fetch summary")
	}

	err := connection.assertReadyState()
	if err != nil {
		return -1, err
	}

	return int(res), nil
}

func (connection *seaboltConnection) LastBookmark() (string, error) {
	bookmark := C.BoltConnection_last_bookmark(connection.cInstance)
	if bookmark != nil {
		return C.GoString(bookmark), nil
	}

	return "", nil
}

func (connection *seaboltConnection) Fields() ([]string, error) {
	fields, err := connection.valueSystem.valueAsGo(C.BoltConnection_field_names(connection.cInstance))
	if err != nil {
		return nil, err
	}

	if fields != nil {
		fieldsAsList := fields.([]interface{})
		fieldsAsStr := make([]string, len(fieldsAsList))
		for i := range fieldsAsList {
			fieldsAsStr[i] = fieldsAsList[i].(string)
		}
		return fieldsAsStr, nil
	}

	return nil, connection.valueSystem.genericErrorFactory("field names not available")
}

func (connection *seaboltConnection) Metadata() (map[string]interface{}, error) {
	metadata, err := connection.valueSystem.valueAsGo(C.BoltConnection_metadata(connection.cInstance))
	if err != nil {
		return nil, err
	}

	if metadataAsGenericMap, ok := metadata.(map[string]interface{}); ok {
		return metadataAsGenericMap, nil
	}

	return nil, connection.valueSystem.genericErrorFactory("metadata is not of expected type")
}

func (connection *seaboltConnection) Data() ([]interface{}, error) {
	fields, err := connection.valueSystem.valueAsGo(C.BoltConnection_field_values(connection.cInstance))
	if err != nil {
		return nil, err
	}

	return fields.([]interface{}), nil
}

func (connection *seaboltConnection) Reset() (RequestHandle, error) {
	res := C.BoltConnection_load_reset_request(connection.cInstance)
	if res != C.BOLT_SUCCESS {
		return -1, newError(connection, "unable to generate reset message")
	}
	return RequestHandle(C.BoltConnection_last_request(connection.cInstance)), nil
}

func (connection *seaboltConnection) Close() error {
	C.BoltConnector_release(connection.connector.cInstance, connection.cInstance)
	return nil
}
