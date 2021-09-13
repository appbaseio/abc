# CSV adaptor

CSV adaptor works for csv files.

A basic pipeline.js looks like the following.
We have an additional parameter `typeName` in csv adaptor because csv files only have data and no concept of tables / types. 
So we need to define it manually.

```js
var source = csv({
  "uri": "/full/local/path/to/file.csv",
  "typeName": "type_name_to_use"
})

var sink = elasticsearch({
  "uri": "https://USER:PASSWORD@SERVER/INDEX"
})

t.Source("source", source, "/.*/").Save("sink", sink, "/.*/")
```
