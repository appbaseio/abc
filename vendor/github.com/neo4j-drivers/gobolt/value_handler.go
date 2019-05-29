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
	"reflect"
)

// ValueHandler is the interface that custom value handlers should implement to
// support reading/writing struct types into custom types
type ValueHandler interface {
	ReadableStructs() []int16
	WritableTypes() []reflect.Type
	Read(signature int16, values []interface{}) (interface{}, error)
	Write(value interface{}) (int16, []interface{}, error)
}

// ValueHandlerError is the special error that ValueHandlers should return in
// case of unexpected cases
type ValueHandlerError struct {
	message string
}

// NewValueHandlerError constructs a new ValueHandlerError
func NewValueHandlerError(format string, args ...interface{}) *ValueHandlerError {
	return &ValueHandlerError{
		message: fmt.Sprintf(format, args...),
	}
}

// Error returns textual representation of the value handler error
func (ns *ValueHandlerError) Error() string {
	return ns.message
}
