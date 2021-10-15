## ABC Import: Examples

Here, we will show some interesting examples for how to import to Elasticsearch or OpenSearch (appbase.io) clusters using the `abc import` command.


### Importing documents with an id

By default, documents imported via abc will not contain any id fields and Elasticsearch will auto-generate the ids at index time. However, if you're importing from the same source periodically (e.g. daily or weekly), you will likely duplicate the documents across imports.

`abc import` expects an `_id` field to be set in your data or via transform for it to set the id field at the time of imports. If your source data doesn't contain an `_id` field, it can be set using a transform file.

```bash
abc import --src_type=postgres --src_uri="postgresql://$user:$password@$host:5432/$db" --transform_file="transform.js" "https://$user:$password@$elasticsearch_url/$index"
```

where the contents of `transform.js` file look as follows:

```js
t.Source("source", source, "/.*/")
 .Transform(goja({"filename":"goja.js"}))
 .Save("sink", sink, "/.*/")
```

and the contents of `goja.js` file look as follows:

```js
function transform(doc) {
    doc["data"]["_id"] = doc["data"]["id"]; // We map the `id` field of the source data to the `_id` field using a transform function.
    return doc
}
```

### Importing specific tables/collections from a database

`abc import` supports filtering of specific tables (if you're importing from SQL) or collections (if you're importing from MongoDB) that can be used while importing.


```bash
abc import --src_type=postgres --src_uri="postgresql://$user:$password@$host:5432/$db" --transform_file="transform.js" "https://$user:$password@$elasticsearch_url/$index"
```

where the contents of `transform.js` file look as follows:

```js
t.Source("source", source, "/table1,table2,table3/")
 .Save("sink", sink, "/.*/")
```

Here, we are applying a source filtering regex expression to import data specifically from `table1`, `table2` or `table3` of `$db` database.


### Importing to a different destination index

`abc import` requires the destination to also provide a default index to import to. This is typically good, however there are cases where you may want to import the data into different indexes. We now support this functionality (starting `1.0.0-beta.3` release) using a goja transform.

```bash
abc import --src_type=postgres --src_uri="postgresql://$user:$password@$host:5432/$db" --transform_file="transform.js" "https://$user:$password@$elasticsearch_url/$index"
```

where the contents of `transform.js` file are:

```js
t.Source("source", source, "/table1,table2,table3/")
 .Transform(goja({"filename":"goja.js"}))
 .Save("sink", sink, "/.*/")
```

and the contents of `goja.js` file are:

```js
function transform(doc) {
    // doc["ns"] key contains the source namespace
    // doc["data"] key contains the source data in JSON form
    switch (doc["ns"]) {
        case "table1":
            doc["data"]["_index"] = "table1";
            break;
        case "table2":
            doc["data"]["_index"] = "table2";
            break;
        case "table3":
            doc["data"]["_index"] = "table3";
            break;
        default:
            doc["data"]["_index"] = "default-index";
            break;
    }
    return doc
}
```

Given the flexibility of being able to write a JS transform snippet, you can decide the destination index based on any number of factors, including the value of keys within doc["data"] itself.


### Importing from MongoDB Atlas

`abc import` supports the standard `mongodb://` protocol to import MongoDB data from. However, MongoDB Atlas now uses a `mongodb+srv://` connection protocol. Here's how you can get the standard `mongodb://` protocol connection string and import.

From your MongoDB Atlas cluster connection modal, choose the "Connect your application" method:

![](https://i.imgur.com/HFvogr3.png)

Next, select Node.js with driver version set to "2.2.12 or later".

![](https://i.imgur.com/27q0p7W.png)

This will show the connection string in the standard protocol format. Now, import with the following command:

```bash
abc import --src_type=mongodb --src_uri="mongodb://$user:$pass@$mongodb_uri/$db?ssl=true&authSource=admin" "https://$user:$password@$elasticsearch_url/$index"
```

abc only requires `ssl=true` (MongoDB Atlas connections only work over SSL) and `authSource=admin` query string parameters, others can be skipped.

