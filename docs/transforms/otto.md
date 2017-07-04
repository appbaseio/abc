# otto function

`otto()` creates a JavaScript VM that receives and sends data through the defined javascript function for processing. The parameter passed to the function has been converted from a go map[string]interface{} to a JS object of the following form:

```JSON
{
    "ns":"message.namespace",
    "ts":12345, // time represented in milliseconds since epoch
    "op":"insert",
    "data": {
        "id": "abcdef",
        "name": "hello world"
    }
}
```

***NOTE*** when working with data from MongoDB, the _id field will be represented in the following fashion:

```JSON
{
    "ns":"message.namespace",
    "ts":12345, // time represented in milliseconds since epoch
    "op":"insert",
    "data": {
        "_id": {
            "$oid": "54a4420502a14b9641000001"
        },
        "name": "hello world"
    }
}
```

### configuration

```javascript
otto({"filename": "/path/to/transform.js"})
// transform() is also available for backwards compatibility reasons but may be removed in future versions
// transform({"filename": "/path/to/transform.js"})
```

### example

message in
```JSON
{
    "_id": 0,
    "name": "abc",
    "type": "function"
}
```

config
```javascript
otto({"filename":"transform.js"})
```

transform function (i.e. `transform.js`)
```javascript
module.exports=function(doc) {
    doc["data"]["name_type"] = doc["data"]["name"] + " " + doc["data"]["type"];
    return doc
}
```

message out
```JSON
{
    "_id": 0,
    "name": "abc",
    "type": "function",
    "name_type": "abc function"
}
```