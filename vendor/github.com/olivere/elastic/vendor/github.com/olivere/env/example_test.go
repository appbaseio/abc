package env_test

import (
	"flag"
	"fmt"

	"github.com/olivere/env"
)

func Example() {
	var (
		// Parse addr from flag, use HTTP_ADDR and ADDR env vars as fallback
		addr = flag.String("addr", env.String("127.0.0.1:3000", "HTTP_ADDR", "ADDR"), "Bind to this address")
	)
	flag.Parse()

	fmt.Println(*addr)
	// Output: 127.0.0.1:3000
}
