package policy

import (
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/open-policy-agent/opa/rego"
	"github.com/open-policy-agent/opa/types"

	"github.com/mneil/opa-dynamodb/store"
	log "github.com/sirupsen/logrus"
)

// DynamoConfig holds information about the dynamo table
type DynamoConfig struct {
	Endpoint     string
	TableName    string
	PartitionKey string
	SortKey      string
}

func getenv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

// NewDynamoConfigFromEnv creates a new dynamo config from environment variables with defaults
func NewDynamoConfigFromEnv() *DynamoConfig {
	return &DynamoConfig{
		TableName:    getenv("DYNAMO_TABLE", "OpaDynamoDB"),
		PartitionKey: getenv("DYNAMO_PK", "PK"),
		SortKey:      getenv("DYNAMO_SK", "SK"),
		Endpoint:     getenv("ENDPOINT_URL", ""),
	}
}

// RegisterDynamodbPolicy registers a new function dynamodb.policy with Rego
func RegisterDynamodbPolicy(dynamoConfig *DynamoConfig) {
	config := &aws.Config{}
	if dynamoConfig.Endpoint != "" {
		log.Warnf("Using custom endpoint %s", dynamoConfig.Endpoint)
		config.Endpoint = &dynamoConfig.Endpoint
	}
	Session := session.Must(session.NewSession(config))
	Store := store.NewDynamoStore(Session, dynamoConfig.TableName)
	Store.PartitionKey = dynamoConfig.PartitionKey
	Store.SortKey = dynamoConfig.SortKey
	policy := NewPolicy("dynamo", Store)
	log.Info("Registering dynamodb.polcy")
	rego.RegisterBuiltin2(
		&rego.Function{
			Name:    "dynamodb.policy",
			Decl:    types.NewFunction(types.Args(types.S, types.S), types.A),
			Memoize: true,
		},
		policy.Get,
	)
}
