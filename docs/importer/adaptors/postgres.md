# PostgreSQL

The [PostgreSQL](https://www.postgresql.org/) adaptor is capable of reading and tailing databases.


### Example

```ini
src_type=postgres
src_uri=postgres://127.0.0.1:5432/test
tail=false

dest_type=elasticsearch
dest_uri=https://USERID:PASS@scalr.api.appbase.io/APPNAME
```

For the destination URI, instead of using your user-id and password, you could also use your admin API key.

```
https://admin-API-key@scalr.api.appbase.io/APPNAME
```

You can find your admin API key inside your app page at appbase.io under Security -> API Credentials.


### Notes

1. When using postgres as source, append `?sslmode=disable` to the uri to disable ssl otherwise you will have to put a certificate. 

#### Tailing

1. When using `tail: true` with postgres, you might get the error `must be superuser or replication role to use replication slots`. You need a [REPLICATION role](https://www.postgresql.org/docs/9.1/static/sql-createrole.html) for this to work. (See [ALTER role](https://www.postgresql.org/docs/9.0/static/sql-alterrole.html))
2. When tailing, you might have to [create replication slots](https://medium.com/@tk512/replication-slots-in-postgresql-b4b03d277c75). Also set `wal_level`. 
```ini
wal_level=logical
max_replication_slots=1
``` 
3. Create a logical replication slot for the source database. ([Logical](https://www.postgresql.org/docs/9.5/static/logicaldecoding-example.html))
```sql
select * from pg_create_logical_replication_slot('standby_replication_slot', 'test_decoding');
SELECT * FROM pg_replication_slots ;
```
4. Make sure you see the `database name` in replication slot row. Now update `replication_slot` parameter of pipeline.js to 'standby_replication_slot'

* [Delete Replication slots](https://stackoverflow.com/questions/30854961/)
