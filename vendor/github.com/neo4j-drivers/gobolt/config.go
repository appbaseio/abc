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
	"bytes"
	"crypto/x509"
	"encoding/pem"
	"net/url"
	"time"
	"unsafe"
)

// Config holds the available configurations options applicable to the connector
type Config struct {
	Encryption             bool
	TLSCertificates        []*x509.Certificate
	TLSSkipVerify          bool
	TLSSkipVerifyHostname  bool
	MaxPoolSize            int
	MaxConnLifetime        time.Duration
	ConnAcquisitionTimeout time.Duration
	SockConnectTimeout     time.Duration
	SockKeepalive          bool
	ConnectorErrorFactory  func(state, code int, codeText, context, description string) ConnectorError
	DatabaseErrorFactory   func(classification, code, message string) DatabaseError
	GenericErrorFactory    func(format string, args ...interface{}) GenericError
	Log                    Logging
	AddressResolver        URLAddressResolver
	ValueHandlers          []ValueHandler
}

func pemEncodeCerts(certs []*x509.Certificate) (*bytes.Buffer, error) {
	if len(certs) == 0 {
		return nil, nil
	}

	var buf = &bytes.Buffer{}
	for _, cert := range certs {
		if err := pem.Encode(buf, &pem.Block{Type: "CERTIFICATE", Bytes: cert.Raw}); err != nil {
			return nil, err
		}
	}
	return buf, nil
}

func createConfig(key int, uri *url.URL, config *Config, valueSystem *boltValueSystem) (*C.struct_BoltConfig, error) {
	var err error
	var cTrust *C.struct_BoltTrust
	var cSocketOpts *C.struct_BoltSocketOptions
	var cRoutingContext *C.struct_BoltValue
	var cUserAgent *C.char
	var cConfig *C.struct_BoltConfig

	if cTrust, err = createTrust(config); err != nil {
		return nil, valueSystem.genericErrorFactory("unable to create trust settings: %v", err)
	}
	defer C.BoltTrust_destroy(cTrust)

	if cRoutingContext, err = createRoutingContext(uri, valueSystem); err != nil {
		return nil, valueSystem.genericErrorFactory("unable to extract routing context: %v", err)
	}
	defer C.BoltValue_destroy(cRoutingContext)

	cSocketOpts = createSocketOptions(config)
	defer C.BoltSocketOptions_destroy(cSocketOpts)

	cUserAgent = C.CString("Go Driver/1.7")
	defer C.free(unsafe.Pointer(cUserAgent))

	var cLogger = registerLogging(key, config.Log)
	defer C.BoltLog_destroy(cLogger)

	var cResolver = registerResolver(key, config.AddressResolver)
	defer C.BoltAddressResolver_destroy(cResolver)

	cConfig = C.BoltConfig_create()
	C.BoltConfig_set_scheme(cConfig, scheme(uri))
	C.BoltConfig_set_transport(cConfig, transport(config))
	C.BoltConfig_set_trust(cConfig, cTrust)
	C.BoltConfig_set_user_agent(cConfig, cUserAgent)
	C.BoltConfig_set_routing_context(cConfig, cRoutingContext)
	C.BoltConfig_set_address_resolver(cConfig, cResolver)
	C.BoltConfig_set_log(cConfig, cLogger)
	C.BoltConfig_set_max_pool_size(cConfig, C.int32_t(config.MaxPoolSize))
	C.BoltConfig_set_max_connection_life_time(cConfig, C.int32_t(config.MaxConnLifetime/time.Millisecond))
	C.BoltConfig_set_max_connection_acquisition_time(cConfig, C.int32_t(config.ConnAcquisitionTimeout/time.Millisecond))
	C.BoltConfig_set_socket_options(cConfig, cSocketOpts)
	return cConfig, nil
}

func scheme(uri *url.URL) C.BoltScheme {
	var mode C.BoltScheme = C.BOLT_SCHEME_DIRECT
	if uri.Scheme == "bolt+routing" {
		mode = C.BOLT_SCHEME_ROUTING
	}
	if uri.Scheme == "neo4j" {
		mode = C.BOLT_SCHEME_NEO4J
	}

	return mode
}

func transport(config *Config) C.BoltTransport {
	var transport C.BoltTransport = C.BOLT_TRANSPORT_PLAINTEXT
	if config.Encryption {
		transport = C.BOLT_TRANSPORT_ENCRYPTED
	}
	return transport
}

func createSocketOptions(config *Config) *C.struct_BoltSocketOptions {
	var cSocketOpts = C.BoltSocketOptions_create()

	C.BoltSocketOptions_set_connect_timeout(cSocketOpts, C.int32_t(config.SockConnectTimeout/time.Millisecond))
	C.BoltSocketOptions_set_keep_alive(cSocketOpts, 1)
	if !config.SockKeepalive {
		C.BoltSocketOptions_set_keep_alive(cSocketOpts, 0)
	}

	return cSocketOpts
}

func createTrust(config *Config) (*C.struct_BoltTrust, error) {
	var cTrust = C.BoltTrust_create()
	C.BoltTrust_set_certs(cTrust, nil, 0)
	C.BoltTrust_set_skip_verify(cTrust, 0)
	C.BoltTrust_set_skip_verify_hostname(cTrust, 0)

	certsBuf, err := pemEncodeCerts(config.TLSCertificates)
	if err != nil {
		C.BoltTrust_destroy(cTrust)

		return nil, err
	}

	if certsBuf != nil {
		certsBytes := certsBuf.String()
		C.BoltTrust_set_certs(cTrust, C.CString(certsBytes), C.uint64_t(certsBuf.Len()))
	}

	if config.TLSSkipVerify {
		C.BoltTrust_set_skip_verify(cTrust, 1)
	}

	if config.TLSSkipVerifyHostname {
		C.BoltTrust_set_skip_verify_hostname(cTrust, 1)
	}

	return cTrust, nil
}

func createRoutingContext(source *url.URL, valueSystem *boltValueSystem) (*C.struct_BoltValue, error) {
	var err error
	var values url.Values
	var result map[string]string

	if values, err = url.ParseQuery(source.RawQuery); err != nil {
		return nil, valueSystem.genericErrorFactory("unable to parse routing context '%s'", source.RawQuery)
	}

	if len(values) == 0 {
		return nil, nil
	}

	result = make(map[string]string, len(values))
	for key, value := range values {
		if len(value) > 1 {
			return nil, valueSystem.genericErrorFactory("duplicate value specified for '%s' as routing context", key)
		}

		result[key] = value[0]
	}

	return valueSystem.valueToConnector(result)
}
