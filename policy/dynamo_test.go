package policy

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetenv(t *testing.T) {
	bar := getenv("FOO", "BAR")
	assert.Equal(t, "BAR", bar)
	os.Setenv("BAZ", "QUX")
	baz := getenv("BAZ", "QUUX")
	assert.Equal(t, "QUX", baz)
}

func TestNewDynamoConfigFromEnv(t *testing.T) {
	config := NewDynamoConfigFromEnv()
	assert.Equal(t, "http://dynamodb:8000/", config.Endpoint)
	assert.Equal(t, "OpaDynamoDB", config.TableName)
	assert.Equal(t, "PK", config.PartitionKey)
	assert.Equal(t, "SK", config.SortKey)
}
