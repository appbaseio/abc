# MYSQL

MYSQL adaptor works for a MySQL Server.

A basic pipeline.js looks like the following.
Notice how the database information is being passed as the source.

```js
var source = mysql({
  "uri": "USER:PASSWORD@tcp(HOST:PORT)/DBNAME"
})

var sink = elasticsearch({
  "uri": "https://USER:PASSWORD@SERVER/INDEX"
})

t.Source("source", source, "/.*/").Save("sink", sink, "/.*/")
```

For other types of source URIs that are supported, visit [go-sql-driver/mysql](https://github.com/go-sql-driver/mysql#examples)'s GitHub page. 
