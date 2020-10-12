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

// General setup for Integration tests

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/open-policy-agent/opa/runtime"
)

func startLocalServer(ctx context.Context, addr string) error {
	parsedURL, err := url.Parse(addr)
	splitURL := strings.Split(parsedURL.Host, ":")
	port, err := strconv.Atoi(splitURL[1])
	Runtime, err := runtime.NewRuntime(ctx, runtime.Params{
		Addrs: &[]string{
			fmt.Sprintf("%s:%d", splitURL[0], port+1),
		},
		InsecureAddr: parsedURL.Host,
		Paths: []string{
			filepath.Join("..", "testdata", "attestors"),
		},
	})
	if err != nil {
		return err
	}
	go Runtime.StartServer(ctx)
	delay := time.Duration(10) * time.Millisecond
	retries := 300 // wait 3 seconds for server to start
	for i := 0; i < retries; i++ {
		if _, err := http.Get(
			addr,
		); err == nil {
			return nil
		}
		time.Sleep(delay)
	}
	return fmt.Errorf("Failed to start OPA server")
}

func createDynamoDBTable() {
	config := &aws.Config{}
	dynamoConfig := NewDynamoConfigFromEnv()
	config.Endpoint = &dynamoConfig.Endpoint
	session := session.Must(session.NewSession(config))
	svc := dynamodb.New(session)
	svc.CreateTable(&dynamodb.CreateTableInput{
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String(dynamoConfig.PartitionKey),
				AttributeType: aws.String("S"),
			},
			{
				AttributeName: aws.String(dynamoConfig.SortKey),
				AttributeType: aws.String("S"),
			},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String(dynamoConfig.PartitionKey),
				KeyType:       aws.String("HASH"),
			},
			{
				AttributeName: aws.String(dynamoConfig.SortKey),
				KeyType:       aws.String("RANGE"),
			},
		},
		BillingMode: aws.String("PAY_PER_REQUEST"),
		TableName:   aws.String(dynamoConfig.TableName),
	})
}

func inflateDynamoDB() {
	dataPath := filepath.Join("..", "testdata", "dynamodb.json")
	content, err := ioutil.ReadFile(dataPath)
	if err != nil {
		panic(err)
	}
	var data []map[string]interface{}
	json.Unmarshal(content, &data)
	config := &aws.Config{}
	dynamoConfig := NewDynamoConfigFromEnv()
	config.Endpoint = &dynamoConfig.Endpoint
	session := session.Must(session.NewSession(config))
	svc := dynamodb.New(session)
	for _, item := range data {
		dynamoItem, _ := dynamodbattribute.MarshalMap(item)
		svc.PutItem(&dynamodb.PutItemInput{
			Item:      dynamoItem,
			TableName: aws.String(dynamoConfig.TableName),
		})
	}

}
func deleteDynamoDBTable() {
	config := &aws.Config{}
	dynamoConfig := NewDynamoConfigFromEnv()
	config.Endpoint = &dynamoConfig.Endpoint
	session := session.Must(session.NewSession(config))
	svc := dynamodb.New(session)
	svc.DeleteTable(&dynamodb.DeleteTableInput{
		TableName: aws.String(dynamoConfig.TableName),
	})

}

func TestMain(m *testing.M) {
	flag.Parse()
	if !testing.Short() {
		createDynamoDBTable()
		inflateDynamoDB()
	}
	code := m.Run()
	if !testing.Short() {
		// deflate db
		deleteDynamoDBTable()
	}
	os.Exit(code)
}
