# JSON adaptor

The json adaptor reads data from a json file.

Here is how a configuration file looks like-

```ini
src.type=json
src.uri=/full/path/to/file.json

dest.type=elasticsearch
dest.uri=appname
```

The `file.json` should contain a json object with table namespaces as its keys.
The keys in return should contain an array containing the data that belongs to that namespace.

Example - 

```js
{
	"pokemons": [
		{
			"_id": 1,
			"name": "Raichu"
		},
		{
			"_id": 1,
			"name": "Bulbasaur"
		}
	],
	"bitbeasts": [
		{
			"_id": 1,
			"name": "Dragoon",
			"owner": "Tyson"
		},
		{
			"_id": 2,
			"name": "Dranzer",
			"owner": "Kai"
		}
	]
}
```
