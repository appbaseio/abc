package all

import (
	// ensures init functions get called
	_ "github.com/aviaryan/abc/adaptor/elasticsearch/clients/v1"
	_ "github.com/aviaryan/abc/adaptor/elasticsearch/clients/v2"
	_ "github.com/aviaryan/abc/adaptor/elasticsearch/clients/v5"
)
