# Redis adaptor

- The [redis](http://redis.io/) adaptor is capable of reading different data structures stored in a Redis database and transfer them to Elasticsearch.

**Note:** Two redis data structures (bitmaps and hyperloglogs) are not supported by the redis adaptor yet. They are indexed as string types however but the appbaseio browser will show non-human readable values for these data structures and hence it is recommended not to index these data structures from your redis DB.

## Usage
`abc import --src_type=redis --src_uri="<redis db url>" "<destination url>"`


## Example
`abc import --src_type=redis --src_uri="redis://localhost:6379/0" appbase-redis-demo`