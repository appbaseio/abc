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

extern void go_seabolt_server_address_resolver_cb(void* state, struct BoltAddress *address, struct BoltAddressSet *resolved);
*/
import "C"
import (
	"fmt"
	"net/url"
	"sync"
	"unsafe"
)

// URLAddressResolver is the callback function signature that provided custom resolver
// should adhere to
type URLAddressResolver func(address *url.URL) []*url.URL

//export go_seabolt_server_address_resolver_cb
func go_seabolt_server_address_resolver_cb(state unsafe.Pointer, address *C.struct_BoltAddress, resolved *C.struct_BoltAddressSet) {
	resolver := lookupResolver(state)
	if resolver != nil {
		resolvedAddresses := resolver(&url.URL{Host: fmt.Sprintf("%s:%s", C.GoString(C.BoltAddress_host(address)), C.GoString(C.BoltAddress_port(address)))})

		for _, addr := range resolvedAddresses {
			cHost := C.CString(addr.Hostname())
			cPort := C.CString(addr.Port())
			cAddress := C.BoltAddress_create(cHost, cPort)

			C.BoltAddressSet_add(resolved, cAddress)

			C.BoltAddress_destroy(cAddress)
			C.free(unsafe.Pointer(cHost))
			C.free(unsafe.Pointer(cPort))
		}
	}
}

var mapResolver sync.Map

func registerResolver(key int, resolver URLAddressResolver) *C.struct_BoltAddressResolver {
	if resolver == nil {
		return nil
	}

	mapResolver.Store(key, resolver)

	boltResolver := C.BoltAddressResolver_create(unsafe.Pointer(&key), C.address_resolver_func(C.go_seabolt_server_address_resolver_cb))
	return boltResolver
}

func lookupResolver(key unsafe.Pointer) URLAddressResolver {
	if resolver, ok := mapResolver.Load(*(*int)(key)); ok {
		return resolver.(URLAddressResolver)
	}

	return nil
}

func unregisterResolver(key int) {
	mapResolver.Delete(key)
}
