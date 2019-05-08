# Neo4J adaptor

The Neo4J adaptor reads data from a Neo4J database.

Here is how a configuration file looks like-

```ini
src_type=neo4j
src_uri=bolt://localhost:7687
src_username=username
src_password=password
src_realm=realm

dest_type=elasticsearch
dest_uri=https://USERID:PASS@scalr.api.appbase.io/APPNAME
```

For the destination URI, instead of using your user-id and password, you could also use your admin API key.

```
https://admin-API-key@scalr.api.appbase.io/APPNAME
```

You can find your admin API key inside your app page at appbase.io under Security -> API Credentials.

The adaptor would fetch all the nodes and relationships from the graph database.

Example usage:

```
abc import \
--src_type=neo4j \
--src_uri="bolt://localhost:7687" \
--src_username=neo4j \
--src_password=test \
"http://localhost:9206/some_index"
```
