# delete

`delete` command allows you to delete an app. It's syntax goes as-

```sh
abc delete [AppID|AppName]
```

#### Example

```sh
# delete an app called LatestApp
abc delete LatestApp
```

```sh
# delete an app with ID 1303
abc delete 1303
```

## cluster flag

To delete a cluster instead of an app we just need to pass a cluster flag along with the cluser ID. For ex `abc delete --cluster clusterID`
