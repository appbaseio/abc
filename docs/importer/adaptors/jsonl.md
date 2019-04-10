JSONL adaptor
=============

The jsonl adaptor reads the data from jsonl file.

Here is how the configuration looks like:

```ini
src_type=jsonl
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

The `file.json` should contain JSON lines text format, also called newline-delimited JSON.
For more information on jsonl format checkout this [link](http://jsonlines.org/).

Example:

    {"name": "Gilbert", "wins": [["straight", "7♣"], ["one pair", "10♥"]]}
    {"name": "Alexa", "wins": [["two pair", "4♠"], ["two pair", "9♠"]]}
    {"name": "May", "wins": []}
    {"name": "Deloise", "wins": [["three of a kind", "5♣"]]}