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
./abc import --log.level=info --type=csv --typename=csvTypeName --source="file.csv" "https://USER:PASS@scalr.api.appbase.io/APPNAME"
```

```sh
./abc import --log.level=info --type=postgres -t --replication_slot="standby_replication_slot" --source="postgresql://USER:PASS@HOST:PORT/DBNAME" "https://USER:PASS@scalr.api.appbase.io/APPNAME"
```
