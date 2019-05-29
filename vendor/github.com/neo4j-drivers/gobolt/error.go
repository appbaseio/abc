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
#include "bolt/bolt.h"
*/
import "C"
import (
	"fmt"
	"strings"
)

// BoltError is a marker interface to identify neo4j errors
type BoltError interface {
	BoltError() bool
}

// DatabaseError represents errors returned from the server a FAILURE messages
type DatabaseError interface {
	// Classification returns classification of the error returned from the database
	Classification() string
	// Code returns code of the error returned from the database
	Code() string
	// Message returns message of the error returned from the database
	Message() string
	// Error returns textual representation of the error returned from the database
	Error() string
}

// ConnectorError represents errors that occur on the connector/client side, like network errors, etc.
type ConnectorError interface {
	// State returns the state of the related connection
	State() int
	// Code returns the error code set on the related connection
	Code() int
	// Context returns the error context set by the connector
	Context() string
	// Description returns any additional description set
	Description() string
	// Error returns textual representation of the connector level error
	Error() string
}

// GenericError represents errors which originates from the connector wrapper itself
type GenericError interface {
	// Message returns the underlying error message
	Message() string
	// Error returns textual representation of the generic error
	Error() string
}

type defaultDatabaseError struct {
	classification string
	code           string
	message        string
}

type defaultConnectorError struct {
	state       int
	code        int
	codeText    string
	context     string
	description string
}

type defaultGenericError struct {
	message string
}

func (failure *defaultDatabaseError) BoltError() bool {
	return true
}

func (failure *defaultDatabaseError) Classification() string {
	return failure.classification
}

func (failure *defaultDatabaseError) Code() string {
	return failure.code
}

func (failure *defaultDatabaseError) Message() string {
	return failure.message
}

func (failure *defaultDatabaseError) Error() string {
	return fmt.Sprintf("database returned error [%s]: %s", failure.code, failure.message)
}

func (failure *defaultConnectorError) BoltError() bool {
	return true
}

func (failure *defaultConnectorError) State() int {
	return failure.state
}

func (failure *defaultConnectorError) Code() int {
	return failure.code
}

func (failure *defaultConnectorError) Context() string {
	return failure.context
}

func (failure *defaultConnectorError) Description() string {
	return failure.description
}

func (failure *defaultConnectorError) Error() string {
	if failure.description != "" {
		return fmt.Sprintf("%s: error: [%d] %s, state: %d, context: %s", failure.description, failure.code, failure.codeText, failure.state, failure.context)
	}

	return fmt.Sprintf("error: [%d] %s, state: %d, context: %s", failure.code, failure.codeText, failure.state, failure.context)
}

func (failure *defaultGenericError) BoltError() bool {
	return true
}

func (failure *defaultGenericError) Message() string {
	return failure.message
}

func (failure *defaultGenericError) Error() string {
	return failure.message
}

func newError(connection *seaboltConnection, description string) error {
	cStatus := C.BoltConnection_status(connection.cInstance)
	errorCode := C.BoltStatus_get_error(cStatus)

	if errorCode == C.BOLT_SERVER_FAILURE {
		failure, err := connection.valueSystem.valueAsDictionary(C.BoltConnection_failure(connection.cInstance))
		if err != nil {
			return connection.valueSystem.genericErrorFactory("unable to construct database error: %s", err.Error())
		}

		var ok bool
		var codeInt, messageInt interface{}
		var code, message string

		if codeInt, ok = failure["code"]; !ok {
			return connection.valueSystem.genericErrorFactory("expected 'code' key to be present in map '%v'", failure)
		}
		if code, ok = codeInt.(string); !ok {
			return connection.valueSystem.genericErrorFactory("expected 'code' value to be of type 'string': '%v'", codeInt)
		}

		if messageInt, ok = failure["message"]; !ok {
			return connection.valueSystem.genericErrorFactory("expected 'message' key to be present in map '%v'", failure)
		}
		if message, ok = messageInt.(string); !ok {
			return connection.valueSystem.genericErrorFactory("expected 'message' value to be of type 'string': '%v'", messageInt)
		}

		classification := ""
		if codeParts := strings.Split(code, "."); len(codeParts) >= 2 {
			classification = codeParts[1]
		}

		return connection.valueSystem.databaseErrorFactory(classification, code, message)
	}

	state := C.BoltStatus_get_state(cStatus)
	errorText := C.GoString(C.BoltError_get_string(errorCode))
	context := C.GoString(C.BoltStatus_get_error_context(cStatus))

	return connection.valueSystem.connectorErrorFactory(int(state), int(errorCode), errorText, context, description)
}

func newGenericError(format string, args ...interface{}) GenericError {
	return &defaultGenericError{message: fmt.Sprintf(format, args...)}
}

func newDatabaseError(classification, code, message string) DatabaseError {
	return &defaultDatabaseError{code: code, message: message, classification: classification}
}

func newConnectorError(state int, code int, codeText, context, description string) ConnectorError {
	return &defaultConnectorError{state: state, code: code, codeText: codeText, context: context, description: description}
}

// IsDatabaseError checkes whether given err is a DatabaseError
func IsDatabaseError(err error) bool {
	if _, ok := err.(DatabaseError); !ok {
		return false
	}

	if _, ok := err.(BoltError); !ok {
		return false
	}

	return true
}

// IsConnectorError checkes whether given err is a ConnectorError
func IsConnectorError(err error) bool {
	if _, ok := err.(ConnectorError); !ok {
		return false
	}

	if _, ok := err.(BoltError); !ok {
		return false
	}

	return true
}

// IsGenericError checkes whether given err is a GenericError
func IsGenericError(err error) bool {
	if _, ok := err.(GenericError); !ok {
		return false
	}

	if _, ok := err.(BoltError); !ok {
		return false
	}

	return true
}

// IsTransientError checks whether given err is a transient error
func IsTransientError(err error) bool {
	if _, ok := err.(BoltError); !ok {
		return false
	}

	if dbErr, ok := err.(DatabaseError); ok {
		if dbErr.Classification() == "TransientError" {
			switch dbErr.Code() {
			case "Neo.TransientError.Transaction.Terminated":
				fallthrough
			case "Neo.TransientError.Transaction.LockClientStopped":
				return false
			}

			return true
		}
	}

	return false
}

// IsWriteError checks whether given err can be classified as a write error
func IsWriteError(err error) bool {
	if _, ok := err.(BoltError); !ok {
		return false
	}

	if dbErr, ok := err.(DatabaseError); ok {
		switch dbErr.Code() {
		case "Neo.ClientError.Cluster.NotALeader":
			fallthrough
		case "Neo.ClientError.General.ForbiddenOnReadOnlyDatabase":
			return true
		}
	}

	return false
}

// IsServiceUnavailable checkes whether given err represents a service unavailable status
func IsServiceUnavailable(err error) bool {
	if _, ok := err.(BoltError); !ok {
		return false
	}

	if connErr, ok := err.(ConnectorError); ok {
		switch connErr.Code() {
		case C.BOLT_INTERRUPTED:
			fallthrough
		case C.BOLT_CONNECTION_RESET:
			fallthrough
		case C.BOLT_NO_VALID_ADDRESS:
			fallthrough
		case C.BOLT_TIMED_OUT:
			fallthrough
		case C.BOLT_CONNECTION_REFUSED:
			fallthrough
		case C.BOLT_NETWORK_UNREACHABLE:
			fallthrough
		case C.BOLT_TLS_ERROR:
			fallthrough
		case C.BOLT_END_OF_TRANSMISSION:
			fallthrough
		case C.BOLT_POOL_FULL:
			fallthrough
		case C.BOLT_ADDRESS_NOT_RESOLVED:
			fallthrough
		case C.BOLT_ROUTING_UNABLE_TO_RETRIEVE_ROUTING_TABLE:
			fallthrough
		case C.BOLT_ROUTING_UNABLE_TO_REFRESH_ROUTING_TABLE:
			fallthrough
		case C.BOLT_ROUTING_NO_SERVERS_TO_SELECT:
			return true
		}
	}

	return false
}

func IsSecurityError(err error) bool {
	if _, ok := err.(BoltError); !ok {
		return false
	}

	if connErr, ok := err.(ConnectorError); ok {
		return connErr.Code() == C.BOLT_TLS_ERROR
	}

	return IsAuthenticationError(err)
}

func IsAuthenticationError(err error) bool {
	if _, ok := err.(BoltError); !ok {
		return false
	}

	if connErr, ok := err.(ConnectorError); ok {
		return connErr.Code() == C.BOLT_PERMISSION_DENIED
	}

	if dbErr, ok := err.(DatabaseError); ok {
		return dbErr.Code() == "Neo.ClientError.Security.Unauthorized"
	}

	return false
}

func IsClientError(err error) bool {
	if _, ok := err.(BoltError); !ok {
		return false
	}

	if dbErr, ok := err.(DatabaseError); ok {
		if dbErr.Classification() == "ClientError" {
			return dbErr.Code() != "Neo.ClientError.Security.Unauthorized"
		}

		return false
	}

	return IsGenericError(err)
}

func IsSessionExpired(err error) bool {
	if _, ok := err.(BoltError); !ok {
		return false
	}

	if connErr, ok := err.(ConnectorError); ok {
		return connErr.Code() == C.BOLT_ROUTING_NO_SERVERS_TO_SELECT
	}

	return false
}

func isPoolFullError(err error) bool {
	if connectorError, ok := err.(ConnectorError); ok {
		return connectorError.Code() == C.BOLT_POOL_FULL
	}

	return false
}

func newConnectionAcquisitionTimedOutError(valueSystem *boltValueSystem) error {
	return valueSystem.connectorErrorFactory(C.BOLT_CONNECTION_STATE_DISCONNECTED, C.BOLT_POOL_ACQUISITION_TIMED_OUT, C.GoString(C.BoltError_get_string(C.BOLT_POOL_ACQUISITION_TIMED_OUT)), "", "")

}

func newPoolFullError(valueSystem *boltValueSystem) error {
	return valueSystem.connectorErrorFactory(C.BOLT_CONNECTION_STATE_DISCONNECTED, C.BOLT_POOL_FULL, C.GoString(C.BoltError_get_string(C.BOLT_POOL_FULL)), "", "")

}
