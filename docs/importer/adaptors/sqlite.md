# SQLITE

SQLITE adaptor works for a SQLITE database.

A basic config.env looks like the following.
Notice how the database information is being passed as the source.

### SQLITE has no username or password. it is standalone database file.

```ini
src_type=sqlite
src_uri=./data.db?_busy_timeout=5000

dest_type=elasticsearch
dest_uri=https://USERID:PASS@scalr.api.appbase.io/APPNAME
```

For the destination URI, instead of using your user-id and password, you could also use your admin API key.

```
https://admin-API-key@scalr.api.appbase.io/APPNAME
```

You can find your admin API key inside your app page at appbase.io under Security -> API Credentials.

For other types of source URIs that are supported, visit [mattn/go-sqlite3](https://github.com/mattn/go-sqlite3/tree/master/_example)'s GitHub page. 

