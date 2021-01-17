package store

import (
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockedDynamo struct {
	mock.Mock
}

func (m *MockedDynamo) Query(input *dynamodb.QueryInput) (*dynamodb.QueryOutput, error) {
	args := m.Called(input)
	out := args.Get(0).(*dynamodb.QueryOutput)
	return out, args.Error(1)
}

type AwsError struct {
	code string
}

func (e *AwsError) Code() string {
	return e.code
}
func (e *AwsError) Error() string {
	return e.code
}
func NewAwsError(code string) *AwsError {
	return &AwsError{
		code: code,
	}
}

func TestNewDynamoStore(t *testing.T) {
	fakeCreds := credentials.NewStaticCredentials("a", "b", "c")
	config := &aws.Config{
		Credentials: fakeCreds,
	}
	session := session.Must(session.NewSession(config))
	store := NewDynamoStore(session, "FooBar")
	assert.Equal(t, "", store.PartitionKey)
	assert.Equal(t, "", store.SortKey)
	assert.Equal(t, "FooBar", store.TableName)
}

func TestGet(t *testing.T) {
	cases := []struct {
		name      string
		output    *dynamodb.QueryOutput
		err       error
		namespace string
		principal string
		expect    interface{}
	}{
		{
			name: "Query returns empty output",
			output: &dynamodb.QueryOutput{
				Items: []map[string]*dynamodb.AttributeValue{},
			},
			err:       nil,
			namespace: "foo",
			principal: "bar",
			expect:    "",
		},
		{
			name: "Throughput AWS error returned",
			output: &dynamodb.QueryOutput{
				Items: []map[string]*dynamodb.AttributeValue{},
			},
			err:       NewAwsError(dynamodb.ErrCodeProvisionedThroughputExceededException),
			namespace: "",
			principal: "",
			expect:    "",
		},
		{
			name: "return a random non aws error",
			output: &dynamodb.QueryOutput{
				Items: []map[string]*dynamodb.AttributeValue{},
			},
			err:       errors.New("random error"),
			namespace: "",
			principal: "",
			expect:    "",
		},
		{
			name: "Query returns simple key,value",
			output: &dynamodb.QueryOutput{
				Items: []map[string]*dynamodb.AttributeValue{
					{
						"foo": &dynamodb.AttributeValue{S: aws.String("bar")},
					},
				},
			},
			err:       nil,
			namespace: "foo",
			principal: "bar",
			expect: map[string]interface{}{
				"foo": "bar",
			},
		},
	}
	// our actual tests is here in the loop
	for _, c := range cases {
		mocked := &MockedDynamo{}
		store := &DynamoStore{
			svc:       mocked,
			TableName: "Foo",
		}
		mocked.On("Query", mock.Anything).Return(c.output, c.err)
		res, err := store.Get(c.namespace, c.principal)
		assert.Equal(t, c.expect, res, c.name)
		assert.Equal(t, c.err, err, c.name)
	}

}
