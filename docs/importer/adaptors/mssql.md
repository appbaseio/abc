# MSSQL

MSSQL adaptor works for Microsoft SQL Server.

A basic config.env looks like the following.
Notice how the database information is being passed as the source.

```ini
src_type=mssql
src_uri=sqlserver://USER:PASSWORD@SERVER:PORT?database=DBNAME

dest_type=elasticsearch
dest_uri=https://USER:PASSWORD@SERVER/INDEX
```

For other types of source URIs that are supported, visit [go-mssqldb](https://github.com/denisenkom/go-mssqldb#connection-parameters-and-dsn)'s GitHub page. 
