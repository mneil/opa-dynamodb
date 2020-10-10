package store

import "github.com/aws/aws-sdk-go/aws/session"

// DynamoStore is a backend for policies from dynamodb
type DynamoStore struct {
}

// NewDynamoStore creates a new DynamoStore
func NewDynamoStore(session *session.Session) *DynamoStore {
	return &DynamoStore{}
}

// Get returns policy data from dynamo
func (dynamo *DynamoStore) Get(namespace string, principal string) (interface{}, error) {
	return "ok", nil
}
