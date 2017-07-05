# apps

`apps` command gives detailed information on the user apps. 
It includes the following additional stats.

1. Number of API calls in the month
2. Total number of records
3. Storage size

```sh
> abc apps

+------+----------------+-----------+---------+---------+
|  ID  |      NAME      | API CALLS | RECORDS | STORAGE |
+------+----------------+-----------+---------+---------+
| 2442 | appname        |        42 |      42 |    4242 |
| 2484 | another        |        84 |      84 |    8484 |
| .... | .........      |       ... |     ... |     ... |
+------+----------------+-----------+---------+---------+
```
