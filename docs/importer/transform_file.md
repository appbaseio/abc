# Transform file

A transform file can be specified with the `import` command which implements transforms when data is moved from source to sink.

The most basic form of transform file is the following. It does nothing but move everything from source to sink.

```js
t.Source("source", source, "/.*/").Save("sink", sink, "/.*/")
```

But we can add [transforms](transforms/) in it to manipulate data that is going to the sink.

```js
t.Source("source", source, "/.*/")
	.Transform(pretty({"spaces":0}))
	// more transforms
	.Save("sink", sink, "/.*/")
```

It can also be used to specify mappings to use in ElasticSearch.
To specify mapping, you use the `Mapping` method. It takes only a single argument which is an object containing mapping data.

```js
t.Source("source", source, "/.*/")
	.Mapping({
		"TypeName": {
			"properties": {
				"name": { "type": "string" },
				"age": { "type": "integer" },
				// more properties
			}
		},
		"AnotherType": {
			"properties": {
				// ....
			}
		}
	})
	.Transform(pretty({"spaces":0}))
	// transforms
	.Save("sink", sink, "/.*/")
```

Note that mapping are set on a type level so the mapping object should contain type and the properties to apply to that type (like we have `TypeName` and `AnotherType` here).
Also the type name used is for the sink, so the type name should be consistent with the namespace that is generated after going through 
all the [transforms](transforms/) i.e. if you have a transform that 
changes namespace in any way, the type names used in mapping should take care of that.


#### src_filter and transform file

src_filter option will not work when you are using a transform file. 
In that case, put the filter in the namespace selector in `.Source` method.

Example

```js
t.Source("source", source, "/.*log/").Save("sink", sink, "/.*/")
```
