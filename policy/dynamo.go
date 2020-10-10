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

// RegisterDynamodbPolicy registers a new function dynamodb.policy with Rego
func RegisterDynamodbPolicy() {
	config := &aws.Config{}
	if endpoint, ok := os.LookupEnv("ENDPOINT_URL"); ok {
		log.Warnf("Using custom endpoint %s", endpoint)
		config.Endpoint = aws.String(endpoint)
	}
	Session := session.Must(session.NewSession(config))
	policy := NewPolicy("dynamo", store.NewDynamoStore(Session))
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
