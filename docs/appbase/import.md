# Import command

Import command can be used to import data from any supported database source into appbase.io/ES cluster. 
It goes like - 

```
abc import --src.uri {URI} --src.type {DBType} --tail [URI|Appname]
```

To view the complete list of input parameters supported, use -

```
abc import --help
```

At the time of writing, the list of parameters supported looks like -

```
--config=                                    Path to external config file, if specified, only that is used
--log.level="info"                           Only log messages with the given severity or above. Valid levels: [debug, info, error]
--replication-slot=standby_replication_slot  [postgres] replication slot to use
--src.filter=.*                              Namespace filter for source
--src.type=postgres                          type of source database
--src.uri=http://user:pass@host:port/db      url of source database
--tail=false                                 allow tail feature
--test=false                                 if set to true, only pipeline is created and sync is not started. Useful for checking your configuration
--transform-file=                            transform file to use
--typename=mytype                            [csv] typeName to use
```

Note that you only need to set the parameters that are required for the source database type. For example, you don't set `replication_slot` when taking CSV as the source. 

**Note** - Help for [transform-file](../importer/transform_file.md) is available here.


## Examples


### CSV

```sh
abc import --src.type=csv --typename=csvTypeName --src.uri="file.csv" "https://USER:PASS@scalr.api.appbase.io/APPNAME"
```


### ElasticSearch

```sh
abc import --src.type=elasticsearch --src.uri="http://USER:PASS@HOST:PORT/INDEX" "https://USER:PASS@scalr.api.appbase.io/APPNAME"
```

We can also use an Appbase app as source.

```sh
abc import --src.type=elasticsearch --src.uri="https://USER:PASS@scalr.api.appbase.io/APPNAME2" "https://USER:PASS@scalr.api.appbase.io/APPNAME"
```


### MongoDB

```sh
abc import --src.type=mongodb -t --src.uri="mongodb://USER:PASS@HOST:PORT/DB" "https://USER:PASS@scalr.api.appbase.io/APPNAME"
```


### MSSQL

```sh
abc import --src.type=mssql --src.uri="sqlserver://USER:PASSWORD@SERVER:PORT?database=DBNAME" "https://USER:PASS@scalr.api.appbase.io/APPNAME"
```

For more source URL patterns, see [go-mssqldb](https://github.com/denisenkom/go-mssqldb#connection-parameters-and-dsn)'s GitHub page. 


### MySQL

```sh
abc import --src.type=mysql --src.uri="USER:PASS@tcp(HOST:PORT)/DBNAME" "https://USER:PASS@scalr.api.appbase.io/APPNAME"
```

For more source URL patterns, see [go-sql-driver/mysql](https://github.com/go-sql-driver/mysql#examples)'s GitHub page. 


### Postgres

```sh
abc import --src.type=postgres -t --replication-slot="standby_replication_slot" --src.uri="postgresql://USER:PASS@HOST:PORT/DBNAME" "https://USER:PASS@scalr.api.appbase.io/APPNAME"
```

### Using a config file

```sh
abc import --config=test.env
```

File extension doesn't matter. 
The file `test.env` should be an INI/ENV like file with key value pair containing the values of attributes required for importing.
Example of a test.env file is --

```ini
src.type=csv
src.uri=/full/path/to/file.csv
typename=csvTypeName

dest.type=elasticsearch
dest.uri=https://USER:PASS@scalr.api.appbase.io/APPNAME
```

Note that the key names are same as what we have in `import` parameters. 
Only exception is that all hyphens in a key name are to be replaced by underscores. 
(e.g. `replication-slot` becomes `replication_slot` in the config file)

