# Kafka adapter

- Under development
- Tested
- SSL support not available
- Supports multiple Topic consumption
- Consumes from all the topics present on the Kafka cluster if no Topic name is given

# Usage
`abc.exe  import --src_type=kafka --src_uri="<host:port>/<topic1,topic2,topic3>" "<destination URI>"`

## Example:
`abc.exe  import --src_type=kafka --src_uri="kafka://localhost:9092/user_log,partner_log" "https://xxxx-xxxx@scalr.api.appbase.io/appName"`
