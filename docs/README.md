# OPA DynamoDB

Scalable policy store with real-time policy updates for use by small and enterprise scale teams wanting to use Open Policy Agent.

OPA DynamoDB adds custom functionality to rego policies to query data from DynamoDB.

OPA has several strategies for managing policies at scale and accepting internal data which you can [read about here](https://www.openpolicyagent.org/docs/latest/external-data/). This repository implements [Option 5](https://www.openpolicyagent.org/docs/latest/external-data/#option-5-pull-data-during-evaluation) using DynamoDB as the external data source. This implementation also removes the current limitations described by OPA.

 - Using this runtime you can test your policies against external data
 - AWS credentials can be infered by the credentials chain in Goland AWS SDK
 - Retry logic and caching are implemented by the AWS SDK and this implementation

# Examples

Read the [Getting Started](quickstart.md) for examples
