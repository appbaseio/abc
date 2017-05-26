package all

import (
	// Initialize all adapters by importing this package
	_ "github.com/aviaryan/abc/adaptor/elasticsearch"
	_ "github.com/aviaryan/abc/adaptor/file"
	_ "github.com/aviaryan/abc/adaptor/mongodb"
	_ "github.com/aviaryan/abc/adaptor/postgres"
	_ "github.com/aviaryan/abc/adaptor/rabbitmq"
	_ "github.com/aviaryan/abc/adaptor/rethinkdb"
)
