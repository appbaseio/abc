# create

`create` command allows you to create a new app.
It supports switches to allow which version of ElasticSearch (ES2 or ES6) you want to run your app on. 
You can also give an optional category to the app.

```sh
abc create [--es2|--es6] [--category=category] AppName
```

`--es2` is set to true by default so if you pass no `es*` switch, an ES 2 app is created. 

⭐️ Use `abc create --help` to view this setting.


#### Examples

```sh
# quickly create an app named 'NewApp' on ES2
abc create NewApp
```

```sh
# create an app named 'LatestApp' on ES6 servers with a custom category
abc create --es6 --category="bleeding-edge" LatestApp
```
