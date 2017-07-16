# JSON adaptor

The json adaptor reads data from a json file.

Here is how a configuration file looks like-

```ini
src.type=json
src.uri=/full/path/to/file.json
typename=typename

dest.type=elasticsearch
dest.uri=appname
```

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
