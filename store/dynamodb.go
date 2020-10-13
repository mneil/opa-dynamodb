// Copyright 2020 Michael Neil

// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at

// 	http://www.apache.org/licenses/LICENSE-2.0

// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package store

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	log "github.com/sirupsen/logrus"
)

// IService is a service interface for the DynamoStore struct. This allows mocking the service
type IService interface {
	Query(input *dynamodb.QueryInput) (*dynamodb.QueryOutput, error)
}

// DynamoStore is a backend for policies from dynamodb
type DynamoStore struct {
	svc          IService
	TableName    string
	PartitionKey string
	SortKey      string
}

// NewDynamoStore creates a new DynamoStore
func NewDynamoStore(session *session.Session, table string) *DynamoStore {
	svc := dynamodb.New(session)
	return &DynamoStore{
		svc:       svc,
		TableName: table,
	}
}

// Get returns policy data from dynamo
func (dynamo *DynamoStore) Get(namespace string, principal string) (interface{}, error) {
	input := &dynamodb.QueryInput{
		ExpressionAttributeNames: map[string]*string{
			"#PK": aws.String(dynamo.PartitionKey),
			"#SK": aws.String(dynamo.SortKey),
		},
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":pk": {
				S: aws.String(namespace),
			},
			":sk": {
				S: aws.String(principal),
			},
		},
		KeyConditionExpression: aws.String("#PK = :pk AND #SK = :sk"),
		TableName:              aws.String(dynamo.TableName),
	}
	result, err := dynamo.svc.Query(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			log.Error(aerr.Code(), aerr.Error())
		} else {
			log.Error(err.Error())
		}
		return "", err
	}
	itemLength := len(result.Items)
	if itemLength == 0 {
		return "", nil
	}
	items := make([]map[string]interface{}, len(result.Items))
	for index, item := range result.Items {
		var tmpItem map[string]interface{}
		dynamodbattribute.UnmarshalMap(item, &tmpItem)
		delete(tmpItem, dynamo.PartitionKey)
		delete(tmpItem, dynamo.SortKey)
		items[index] = tmpItem
	}
	return items, nil
}
