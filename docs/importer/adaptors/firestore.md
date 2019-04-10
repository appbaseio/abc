# Cloud Firestore

- The [firestore](https://firebase.google.com/docs/firestore) adaptor is capable of reading collections from the firestore database.

- Checkout firestore [service accounts](https://console.cloud.google.com/iam-admin/serviceaccounts) for additional 
details about initializing the firestore SDK on your own server.  

Here is how a configuration file looks like:

```ini
src_type=firestore
sac_path="/path/to/service_account_credentials_file.json"
src_filter="<collection name/regex>"

dest_type=elasticsearch
dest_uri=https://USERID:PASS@scalr.api.appbase.io/APPNAME
```

For the destination URI, instead of using your user-id and password, you could also use your admin API key.

```
https://admin-API-key@scalr.api.appbase.io/APPNAME
```

You can find your admin API key inside your app page at appbase.io under Security -> API Credentials.

## Usage
`abc import --src_type=firestore --sac_path="/path/to/service_account_credentials_file.json" "<destination url>"`

## Filter the collections
`abc import --src_type=firestore --sac_path="/path/to/service_account_credentials_file.json" --src_filter="<collection name/regex>" "<destination url>"`

## Example
`abc import --src_type=firestore --sac_path="/home/johnappleseed/ServiceAccountKey.json" appbase-firestore-demo`