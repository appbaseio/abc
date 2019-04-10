# MSSQL

MSSQL adaptor works for Microsoft SQL Server.

A basic config.env looks like the following.
Notice how the database information is being passed as the source.

```ini
src_type=mssql
src_uri=sqlserver://USER:PASSWORD@SERVER:PORT?database=DBNAME

dest_type=elasticsearch
dest_uri=https://USERID:PASS@scalr.api.appbase.io/APPNAME
```

For the destination URI, instead of using your user-id and password, you could also use your admin API key.

```
https://admin-API-key@scalr.api.appbase.io/APPNAME
```

You can find your admin API key inside your app page at appbase.io under Security -> API Credentials.

For other types of source URIs that are supported, visit [go-mssqldb](https://github.com/denisenkom/go-mssqldb#connection-parameters-and-dsn)'s GitHub page. 
