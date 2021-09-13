# Firestore adaptor

- The [firestore](https://firebase.google.com/docs/firestore) adaptor is capable of reading collections from the firestore database.

- Checkout firestore [service accounts](https://console.cloud.google.com/iam-admin/serviceaccounts) for additional 
details about initializing the firestore SDK on your own server.  

## Future Enhancements
- Make use of batch transactions if the number of documents in the database exceeds a particular threshold.
- Display improved and relevant logs.

## Usage
`abc import --src_type=firestore --sac_path="/path/to/service_account_credentials_file.json" "<destination url>"`

## Filter the collections
`abc import --src_type=firestore --sac_path="/path/to/service_account_credentials_file.json" --src_filter="<collection name/regex>" "<destination url>"`

## Example
`abc import --src_type=firestore --sac_path="/home/johnappleseed/ServiceAccountKey.json" appbase-firestore-demo`