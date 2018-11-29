# create

`create` command allows you to create a new app.
It supports switches to allow which version of ElasticSearch (ES2 or ES6) you want to run your app on. 
You can also give an optional category to the app.

```sh
abc create [--es2|--es6] [--category=category] [--cluster|-c] [--interactive|-i] [--loc] [--vmsize] [--plan] [--ssh] [--provider] [--nodes] [--version] [--volume] AppName|ClusterName
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

## cluster flag

The clusters flag can be used to deploy a cluster instead of an app. There are two modes for creating a cluster:

1. **Interactive mode**

`abc create --cluster --interactive ClusterName`

This will fire up the interactive mode in which the user will have to answer questions regarding the cluster details and at the end the cluster will be deployed with that same configuration

- Some of the questions provide help text.
- Some questions have multiple options for the user to scroll through and choose from.
- In case of choosing the plugins for the cluster the user can select multiple options. For this particular case *spacebar* is used to select options instead of the *return* key

2. **Non-interactive mode**

In this case all the options that are mandatory for the cluster to be deployed are to be provided through command line flags

`abc create --cluster --loc"east-us-1" --vm-size="n1-standard-1" --plan="Growth" --ssh=<your-ssh-key> --provider="gke" --nodes=3 --volume=15 --version="5.6.9" ClusterName`