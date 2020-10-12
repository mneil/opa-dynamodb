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
