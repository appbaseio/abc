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

Note that you only need to set the parameters that are required for the source database type. For example, you don't set `replication_slot` when taking CSV as the source. 


### Examples

```sh
./abc import --log.level=info --src.type=csv --typename=csvTypeName --src.uri="file.csv" "https://USER:PASS@scalr.api.appbase.io/APPNAME"
```

```sh
./abc import --log.level=info --src.type=postgres -t --replication-slot="standby_replication_slot" --src.uri="postgresql://USER:PASS@HOST:PORT/DBNAME" "https://USER:PASS@scalr.api.appbase.io/APPNAME"
```
