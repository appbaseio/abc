# CSV

CSV adaptor works for csv files.

A basic config.env looks like the following.
We have an additional parameter `typename` in csv adaptor because csv files only have data and no concept of tables / types. 
So we need to define it manually.

```ini
src_type=csv
src_uri=/full/local/path/to/file.csv
typename=type_name_to_use

dest_type=elasticsearch
dest_uri=https://USERID:PASS@scalr.api.appbase.io/APPNAME
```

For the destination URI, instead of using your user-id and password, you could also use your admin API key.

```
https://admin-API-key@scalr.api.appbase.io/APPNAME
```

You can find your admin API key inside your app page at appbase.io under Security -> API Credentials.
