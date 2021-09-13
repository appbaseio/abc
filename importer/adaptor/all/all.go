package all

import (
	// Initialize all adapters by importing this package
	_ "github.com/appbaseio/abc/importer/adaptor/csv"
	_ "github.com/appbaseio/abc/importer/adaptor/elasticsearch"
	_ "github.com/appbaseio/abc/importer/adaptor/firestore"
	_ "github.com/appbaseio/abc/importer/adaptor/json"
	_ "github.com/appbaseio/abc/importer/adaptor/jsonl"
	_ "github.com/appbaseio/abc/importer/adaptor/kafka"
	_ "github.com/appbaseio/abc/importer/adaptor/mongodb"
	_ "github.com/appbaseio/abc/importer/adaptor/mssql"
	_ "github.com/appbaseio/abc/importer/adaptor/mysql"
	_ "github.com/appbaseio/abc/importer/adaptor/postgres"
	_ "github.com/appbaseio/abc/importer/adaptor/rabbitmq"
	_ "github.com/appbaseio/abc/importer/adaptor/redis"
	_ "github.com/appbaseio/abc/importer/adaptor/rethinkdb"
	_ "github.com/appbaseio/abc/importer/adaptor/sqlite"
)
