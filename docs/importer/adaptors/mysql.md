# MYSQL

MYSQL adaptor works for a MySQL Server.

A basic config.env looks like the following.
Notice how the database information is being passed as the source.

```ini
src_type=mysql
src_uri=USER:PASSWORD@tcp(HOST:PORT)/DBNAME

dest_type=elasticsearch
dest_uri=https://USER:PASSWORD@SERVER/INDEX
```

For other types of source URIs that are supported, visit [go-sql-driver/mysql](https://github.com/go-sql-driver/mysql#examples)'s GitHub page. 
