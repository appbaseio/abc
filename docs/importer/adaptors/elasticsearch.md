# Elasticsearch

The [elasticsearch](https://www.elastic.co/) adaptor sends data to an ES cluster.

Right now, the only supported version is ES 2 (the same that runs on Appbase) but we may have support for newer ES versions soon.
ES can be both a source and a sink.


### Example

```ini
src_type=elasticsearch
src_uri=https://user:pass@es_cluster/index

dest_type=elasticsearch
dest_uri=https://USERID:PASS@scalr.api.appbase.io/APPNAME
```

For the destination URI, instead of using your user-id and password, you could also use your admin API key.

```
https://admin-API-key@scalr.api.appbase.io/APPNAME
```

You can find your admin API key inside your app page at appbase.io under Security -> API Credentials.

#### About IDs

If your table has a column named `_id`, then it will be automatically used as elasticsearch ID. 

If this is not the case, use a transform function to set `_id` in the document.

If no field named `_id` goes to ElasticSearch sink, an auto-generated ID is used.
