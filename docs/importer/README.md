# Generic Importer Commands

### init

```
abc init [source adaptor name] [sink adaptor name]
```

Generates a basic `pipeline.js` file in the current directory.

_Example_
```
$ abc init mongodb elasticsearch
$ cat pipeline.js
var source = mongodb({
  "uri": "${MONGODB_URI}"
  // "timeout": "30s",
  // "tail": false,
  // "ssl": false,
  // "cacerts": ["/path/to/cert.pem"],
  // "wc": 1,
  // "fsync": false,
  // "bulk": false,
  // "collection_filters": "{}"
})

var sink = elasticsearch({
  "uri": "${ELASTICSEARCH_URI}"
  // "timeout": "10s", // defaults to 30s
  // "aws_access_key": "ABCDEF", // used for signing requests to AWS Elasticsearch service
  // "aws_access_secret": "ABCDEF" // used for signing requests to AWS Elasticsearch service
})

t.Source(source).Save(sink)
// t.Source("source", source).Save("sink", sink)
// t.Source("source", source, "namespace").Save("sink", sink, "namespace")
$
```

Edit the `pipeline.js` file to configure the source and sink nodes and also to set the namespace.

### about

`abc about`

Lists all the adaptors currently available.

_Example_

```
elasticsearch - an elasticsearch sink adaptor
file - an adaptor that reads / writes files
mongodb - a mongodb adaptor that functions as both a source and a sink
postgres - a postgres adaptor that functions as both a source and a sink
rabbitmq - an adaptor that handles publish/subscribe messaging with RabbitMQ 
rethinkdb - a rethinkdb adaptor that functions as both a source and a sink
```

Giving the name of an adaptor produces more detail, such as the sample configuration.

_Example_

```
abc about postgres
postgres - a postgres adaptor that functions as both a source and a sink

 Sample configuration:
{
  "uri": "${POSTGRESQL_URI}"
  // "debug": false,
  // "tail": false,
  // "replication_slot": "slot"
}
```

### run

```
abc run [-log.level "info"] <application.js>
```

Runs the pipeline script file which has its name given as the final parameter.

### test

```
abc test [-log.level "info"] <application.js>
```

Evaluates and connects the pipeline, sources and sinks. Establishes connections but does not run.
Prints out the state of connections at the end. Useful for debugging new configurations.

### xlog

The `xlog` command is useful for inspecting the current state of the commit log.
It contains 3 subcommands, `current`, `oldest`, and `offset`, as well as 
a required flag `-log_dir` which should be the path to where the commit log is stored.

***NOTE*** the command should only be run against the commit log when abc
is not actively running.

```
abc xlog -log_dir=/path/to/dir current
12345
```

Returns the most recent offset appended to the commit log.

```
abc xlog -log_dir=/path/to/dir oldest
0
```

Returns the oldest offset in the commit log.

```
abc xlog -log_dir=/path/to/dir show 0
offset    : 0
timestamp : 2017-05-16 11:00:20 -0400 EDT
mode      : COPY
op        : INSERT
key       : MyCollection
value     : {"_id":{"$oid":"58efd14b60d271d7457b4f24"},"i":0}
```

Prints out the entry stored at the provided offset.

### offset

The `offset` command provides access to current state of each consumer (i.e. sink)
offset. It contains 4 subcommands, `list`, `show`, `mark`, and `delete`, as well as 
a required flag `-log_dir` which should be the path to where the commit log is stored.

```
abc offset -log_dir=/path/to/dir list
+------+---------+
| SINK | OFFSET  |
+------+---------+
| sink | 1103003 |
+------+---------+
```

Lists all consumers and their associated offset in `log_dir`.

```
abc offset -log_dir=/path/to/dir show sink
+-------------------+---------+
|     NAMESPACE     | OFFSET  |
+-------------------+---------+
| newCollection     | 1102756 |
| testC             | 1103003 |
| MyCollection      |  999429 |
| anotherCollection | 1002997 |
+-------------------+---------+
```

Prints out each namespace and its associated offset.

```
abc offset -log_dir=/path/to/dir mark sink 1
OK
```

Rewrites the namespace offset map based on the provided offset.

```
abc offset -log_dir=/path/to/dir delete sink
OK
```

Removes the consumer (i.e. sink) log directory.