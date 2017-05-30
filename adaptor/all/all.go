package all

import (
	// Initialize all adapters by importing this package
	_ "github.com/appbaseio/abc/adaptor/elasticsearch"
	_ "github.com/appbaseio/abc/adaptor/file"
	_ "github.com/appbaseio/abc/adaptor/mongodb"
	_ "github.com/appbaseio/abc/adaptor/postgres"
	_ "github.com/appbaseio/abc/adaptor/rabbitmq"
	_ "github.com/appbaseio/abc/adaptor/rethinkdb"
)
