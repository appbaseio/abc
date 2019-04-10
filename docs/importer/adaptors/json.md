# JSON adaptor

The json adaptor reads data from a json file.

Here is how a configuration file looks like-

```ini
src_type=json
src_uri=/full/path/to/file.json
typename=typename

dest_type=elasticsearch
dest_uri=https://USERID:PASS@scalr.api.appbase.io/APPNAME
```

For the destination URI, instead of using your user-id and password, you could also use your admin API key.

```
https://admin-API-key@scalr.api.appbase.io/APPNAME
```

You can find your admin API key inside your app page at appbase.io under Security -> API Credentials.

The `file.json` should contain a json array with individual rows as its contents.

Example - 

```js
[
	{
		"_id": 1,
		"name": "Raichu"
	},
	{
		"_id": 2,
		"name": "Bulbasaur"
	}
]
```
