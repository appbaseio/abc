# app

`app` command display detailed information about a single app.

The most basic syntax to use it is -

```sh
abc app AppName
```

It can also take `AppID` in place of `AppName`.

```sh
abc app AppID
```

The above commands will only give basic details of the app such as its name, ID, ElasticSearch version etc.

If you want more details such as credentials and metrics, you can pass in defined switches. The full version of `app` command looks like-

```sh
abc app [-c|--creds] [-m|--metrics] [ID|Appname]
```

#### Example

```sh
> abc app -cm MyCoolApp

ID:         2000       
Name:       MyCoolApp
ES Version: 2.2.0

Admin API key:      USER:PASS
Read-only API key:  USER:PASS

Storage:    4242
Records:    42
+-------+-----------+---------+
| DATE  | API CALLS | RECORDS |
+-------+-----------+---------+
| 13-05 |     8     |    2    |
| 14-05 |     10    |    2    |
| ..... |     ..    |   ...   |
+-------+-----------+---------+
| TOTAL |     18    |    4    |
+-------+-----------+---------+
```

⭐️ It shows API calls and number of records created on a per day basis (just like the graph you see on dashboard.appbase.io). 
