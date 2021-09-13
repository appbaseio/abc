package all

import (
	// ensures init functions get called
	// _ "github.com/appbaseio/abc/importer/adaptor/elasticsearch/clients/v1"
	// _ "github.com/appbaseio/abc/importer/adaptor/elasticsearch/clients/v2"
	_ "github.com/appbaseio/abc/importer/adaptor/elasticsearch/clients/v5"
	_ "github.com/appbaseio/abc/importer/adaptor/elasticsearch/clients/v6"
	_ "github.com/appbaseio/abc/importer/adaptor/elasticsearch/clients/v7"
)
