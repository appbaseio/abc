package all

import (
	// ensures init functions get called
	_ "github.com/appbaseio/abc/importer/adaptor/elasticsearch/clients/v7"
	_ "github.com/appbaseio/abc/importer/adaptor/elasticsearch/clients/v8"
)
