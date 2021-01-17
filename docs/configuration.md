# Configuration

OPA DynamoDB can be configured using environment variable. This table describes the default settings and the variable that you can set to change them.

| Variable     | Default     | Required | Description                                          |
|--------------|-------------|----------|------------------------------------------------------|
| DYNAMO_TABLE | OpaDynamoDB | No       | The name of the table to get data from               |
| DYNAMO_PK    | PK          | No       | The hash (partition) key of the primary partition    |
| DYNAMO_SK    | SK          | No       | The sort (range) key of the primary partition        |
| ENDPOINT_URL | ""          | No       | DynamoDB API url. Useful for testing w/ dynamo local |
