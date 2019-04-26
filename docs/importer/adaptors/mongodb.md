# MongoDB

The [MongoDB](https://www.mongodb.com/) adaptor is capable of reading/tailing collections and receiving data for inserts.

Here is how a configuration file looks like

```ini
src_type=mongodb
src_uri=mongodb://user:pass@SERVER:PORT/DBNAME
tail=false

dest_type=elasticsearch
dest_uri=APP_NAME
dest_uri=https://USERID:PASS@scalr.api.appbase.io/APPNAME
```

For the destination URI, instead of using your user-id and password, you could also use your admin API key.

```
https://admin-API-key@scalr.api.appbase.io/APPNAME
```

You can find your admin API key inside your app page at appbase.io under Security -> API Credentials.

### Notes

For **tailing** to work, the user accessing the database will need to have `oplog` access. 
An admin user should have oplog access by default though if you are using a mongo database from an third-party provider like 
[mlab](https://mlab.com), this might not be the case. It is best to contact provider's support team in this case.
Also, you can see the provider docs for how to get `oplog` access.

Read this excellent [blog post](http://www.sohamkamani.com/blog/2016/06/30/docker-mongo-replica-set/) for how to setup oplog access on your own cluster.
In case you are using mlab which is a very popular mongo provider, see the [oplog doc for mlab](http://docs.mlab.com/oplog/).

To enable resume tailing, use the `--log_dir` argument with `abc import`. Give it a path to store logs, which would help resume tailing.
