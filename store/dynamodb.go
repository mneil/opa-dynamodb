package store

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	log "github.com/sirupsen/logrus"
)

// DynamoStore is a backend for policies from dynamodb
type DynamoStore struct {
	svc          *dynamodb.DynamoDB
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
			switch aerr.Code() {
			case dynamodb.ErrCodeProvisionedThroughputExceededException:
				log.Error(dynamodb.ErrCodeProvisionedThroughputExceededException, aerr.Error())
			case dynamodb.ErrCodeResourceNotFoundException:
				log.Error(dynamodb.ErrCodeResourceNotFoundException, aerr.Error())
			case dynamodb.ErrCodeRequestLimitExceeded:
				log.Error(dynamodb.ErrCodeRequestLimitExceeded, aerr.Error())
			case dynamodb.ErrCodeInternalServerError:
				log.Error(dynamodb.ErrCodeInternalServerError, aerr.Error())
			default:
				log.Error(aerr.Error())
			}
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
