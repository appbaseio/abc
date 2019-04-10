# MYSQL

MYSQL adaptor works for a MySQL Server.

A basic config.env looks like the following.
Notice how the database information is being passed as the source.

```ini
src_type=mysql
src_uri=USER:PASSWORD@tcp(HOST:PORT)/DBNAME

dest_type=elasticsearch
dest_uri=https://USERID:PASS@scalr.api.appbase.io/APPNAME
```

For the destination URI, instead of using your user-id and password, you could also use your admin API key.

```
https://admin-API-key@scalr.api.appbase.io/APPNAME
```

You can find your admin API key inside your app page at appbase.io under Security -> API Credentials.

Syncing only a particular table is possible with the `--src_filter` switch.

For other types of source URIs that are supported, visit [go-sql-driver/mysql](https://github.com/go-sql-driver/mysql#examples)'s GitHub page. 
