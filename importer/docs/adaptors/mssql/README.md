# MSSQL

MSSQL adaptor works for Microsoft SQL Server.

A basic pipeline.js looks like the following.
Notice how the database information is being passed as the source.

```js
var source = mssql({
  "uri": "sqlserver://USER:PASSWORD@SERVER:PORT?database=DBNAME"
})

var sink = elasticsearch({
  "uri": "https://USER:PASSWORD@SERVER/INDEX"
})

t.Source("source", source, "/.*/").Save("sink", sink, "/.*/")
```

For other types of source URIs that are supported, visit [go-mssqldb](https://github.com/denisenkom/go-mssqldb#connection-parameters-and-dsn)'s GitHub page. 
