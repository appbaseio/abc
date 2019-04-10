# Redis adaptor

- The [redis](http://redis.io/) adaptor is capable of reading different data structures stored in a Redis database and transfer them to Elasticsearch.

**Note:** Two redis data structures (bitmaps and hyperloglogs) are not supported by the redis adaptor yet. They are indexed as string types however but the appbaseio browser will show non-human readable values for these data structures and hence it is recommended not to index these data structures from your redis DB.

## Usage
`abc import --src_type=redis --src_uri="<redis db url>" "<destination url>"`


## Example
`abc import --src_type=redis --src_uri="redis://localhost:6379/0" https://USERID:PASS@scalr.api.appbase.io/APPNAME`

For the destination URI, instead of using your user-id and password, you could also use your admin API key.

```
https://admin-API-key@scalr.api.appbase.io/APPNAME
```

You can find your admin API key inside your app page at appbase.io under Security -> API Credentials.