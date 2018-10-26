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
abc app [-c|--creds] [-m|--metrics] [-a| --analytics] [ID|Appname]
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


### Accessing app analytics (only for paid users)

Paid users can access analytics data for a particular app. Currently the following analytics endpoints are supported:
- overview
- noresultsearches
- popularresults
- popularsearches
- popularfilters
- geoip
- latency

#### Example

By default the command will ping the overview endpoint

```sh
> abc app -a MyCoolApp

Analytics(Overview) Results:
No Result Searches
+-------+------+
| COUNT | KEY  |
+-------+------+
|   1   | blah |
|   1   | gurr |
+-------+------+
No Result Searches
+-------+--------------+
| COUNT |     KEY      |
+-------+--------------+
|   1   |     gru      |
|  ...  |     ...      |
|   1   | wonder woman |
|   1   |  wonderland  |
+-------+--------------+
Search Volume Results
+-------+---------------+---------------------+
| COUNT |      KEY      |     DATE-AS-STR     |
+-------+---------------+---------------------+
|   7   | 1540512000000 | 2018/10/26 00:00:00 |
+-------+---------------+---------------------+
```

To ping other analytics endpoints the `endpoint` flag can be used. For example

```sh
> abc app -a --endpoint=latency MyCoolApp

Analytics(Latency) Results:
+-------+-----+
| COUNT | KEY |
+-------+-----+
|   2   |  0  |
|  ..   | ..  |
|   0   | 10  |
+-------+-----+
```