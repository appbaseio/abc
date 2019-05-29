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
	"time"

	"github.com/stretchr/testify/mock"
)

type mockConnection struct {
	mock.Mock
}

func (m *mockConnection) Id() (string, error) {
	args := m.Called()
	return args.String(0), args.Error(1)
}

func (m *mockConnection) RemoteAddress() (string, error) {
	args := m.Called()
	return args.String(0), args.Error(1)
}

func (m *mockConnection) Server() (string, error) {
	args := m.Called()
	return args.String(0), args.Error(1)
}

func (m *mockConnection) Begin(bookmarks []string, txTimeout time.Duration, txMetadata map[string]interface{}) (RequestHandle, error) {
	args := m.Called(bookmarks, txTimeout, txMetadata)
	return args.Get(0).(RequestHandle), args.Error(1)
}

func (m *mockConnection) Commit() (RequestHandle, error) {
	args := m.Called()
	return args.Get(0).(RequestHandle), args.Error(1)
}

func (m *mockConnection) Rollback() (RequestHandle, error) {
	args := m.Called()
	return args.Get(0).(RequestHandle), args.Error(1)
}

func (m *mockConnection) Run(cypher string, parameters map[string]interface{}, bookmarks []string, txTimeout time.Duration, txMetadata map[string]interface{}) (RequestHandle, error) {
	args := m.Called(cypher, parameters, bookmarks, txTimeout, txMetadata)
	return args.Get(0).(RequestHandle), args.Error(1)
}

func (m *mockConnection) PullAll() (RequestHandle, error) {
	args := m.Called()
	return args.Get(0).(RequestHandle), args.Error(1)
}

func (m *mockConnection) DiscardAll() (RequestHandle, error) {
	args := m.Called()
	return args.Get(0).(RequestHandle), args.Error(1)
}

func (m *mockConnection) Reset() (RequestHandle, error) {
	args := m.Called()
	return args.Get(0).(RequestHandle), args.Error(1)
}

func (m *mockConnection) Flush() error {
	args := m.Called()
	return args.Error(0)
}

func (m *mockConnection) Fetch(request RequestHandle) (FetchType, error) {
	args := m.Called(request)
	return args.Get(0).(FetchType), args.Error(1)
}

func (m *mockConnection) FetchSummary(request RequestHandle) (int, error) {
	args := m.Called(request)
	return args.Int(0), args.Error(1)
}

func (m *mockConnection) LastBookmark() (string, error) {
	args := m.Called()
	return args.String(0), args.Error(1)
}

func (m *mockConnection) Fields() ([]string, error) {
	args := m.Called()
	return args.Get(0).([]string), args.Error(1)
}

func (m *mockConnection) Metadata() (map[string]interface{}, error) {
	args := m.Called()
	return args.Get(0).(map[string]interface{}), args.Error(1)
}

func (m *mockConnection) Data() ([]interface{}, error) {
	args := m.Called()
	return args.Get(0).([]interface{}), args.Error(1)
}

func (m *mockConnection) Close() error {
	args := m.Called()
	return args.Error(0)
}
