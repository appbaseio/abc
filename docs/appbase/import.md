# Import command

Import command can be used to import data from any supported database source into appbase.io/ES cluster. 
It goes like - 

```
abc import --source {URI} --type {DBType} --tail [URI|Appname|Appid]
```

To view the complete list of input parameters supported, use -

```
abc import --help
```

At the time of writing, the list of parameters supported looks like -

```
--log.level="info"                           Only log messages with the given severity or above. Valid levels: [debug, info, error]
--replication-slot=standby_replication_slot  [postgres] replication slot to use
--src.filter=.*                              Namespace filter for source
--src.type=postgres                          type of source database
--src.uri=http://user:pass@host:port/db      url of source database
--tail=false                                 allow tail feature
--timeout=10s                                source timeout
--typename=mytype                            [csv] typeName to use
```

Note that you only need to set the parameters that are required for the source database type. For example, you don't set `replication_slot` when taking CSV as the source. 


## Examples


#### CSV

```sh
./abc import --src.type=csv --typename=csvTypeName --src.uri="file.csv" "https://USER:PASS@scalr.api.appbase.io/APPNAME"
```

#### Postgres

```sh
./abc import --src.type=postgres -t --replication-slot="standby_replication_slot" --src.uri="postgresql://USER:PASS@HOST:PORT/DBNAME" "https://USER:PASS@scalr.api.appbase.io/APPNAME"
```

#### MySQL

```sh
./abc import --src.type=mysql --src.uri="USER:PASS@tcp(HOST:PORT)/DBNAME" "https://USER:PASS@scalr.api.appbase.io/APPNAME"
```

For more source URL patterns, see [go-sql-driver/mysql](https://github.com/go-sql-driver/mysql#examples)'s GitHub page. 

