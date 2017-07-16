# Elasticsearch

The [elasticsearch](https://www.elastic.co/) adaptor sends data to an ES cluster.

Right now, the only supported version is ES 2 (the same that runs on Appbase) but we may have support for newer ES versions soon.
ES can be both a source and a sink.


### Example

```ini
src.type=elasticsearch
src.uri=https://user:pass@es_cluster/index

dest.type=elasticsearch
dest.uri=abctests
# uri field can have appbase appname as well as full URI
```

#### About IDs

If your table has a column named `_id`, then it will be automatically used as elasticsearch ID. 

If this is not the case, use a transform function to set `_id` in the document.

If no field named `_id` goes to ElasticSearch sink, an auto-generated ID is used.
