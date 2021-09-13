# Redis adaptor

- The [redis](http://redis.io/) adaptor is capable of reading different data structures stored in a Redis database and transfer them to Elasticsearch.


## Future Enhancements
- Currently bitmaps and hyperloglogs data structures are not supported. Further investigation will be done to figure out how to add support for these two data structures.
- Implement a writer for Redis.

## Usage
`abc import --src_type=redis --src_uri="<redis db url>" "<destination url>"`

## Configure
The redis client accepts the following options to configure your redis server.
- uri (default: "redis://127.0.0.1:6379/0")
- db (default: 0)
- readTimeout (default: 3s)
- writeTimeout (default: 3s)
- dialTimeout (default: 5s)
- tlsConfig

## Example
`abc import --src_type=redis --src_uri="redis://localhost:6379/0" appbase-redis-demo`