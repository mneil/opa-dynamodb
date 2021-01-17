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
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/open-policy-agent/opa/ast"
	"github.com/open-policy-agent/opa/rego"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockedStore struct {
	mock.Mock
}

func (m *MockedStore) Get(namespace string, principal string) (interface{}, error) {
	args := m.Called(namespace, principal)
	return args.Get(0), args.Error(1)
}

func TestPolicyGet(t *testing.T) {
	cases := []struct {
		name      string
		namespace string
		principal string
		output    interface{}
		err       error
		expect    *ast.Term
	}{
		{
			name:      "ok",
			namespace: "foo/bar",
			principal: "baz",
			output:    []map[string]string{},
			err:       nil,
			expect:    &ast.Term{Value: ast.Value(nil)},
		},
	}
	for _, c := range cases {
		store := &MockedStore{}
		policy := NewPolicy("foo", store)
		store.On("Get", c.namespace, c.principal).Return(c.output, c.err)
		bctx := rego.BuiltinContext{}

		res, err := policy.Get(bctx, ast.StringTerm(c.namespace), ast.StringTerm(c.principal))
		assert.Equal(t, c.expect, res, c.name)
		if c.err != nil {
			assert.Nil(t, err, c.name)
		}
	}

}

func TestPolicyDataIntegration(t *testing.T) {
	// requires dynamo db connection and runs local opa server
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}
	addr := "http://127.0.0.1:8080"
	// register our dynamo function
	RegisterDynamodbPolicy(NewDynamoConfigFromEnv())
	ctx := context.Background()
	defer ctx.Done()
	// start our opa server
	err := startLocalServer(ctx, addr)
	assert.Nil(t, err)
	// the test cases to check against
	cases := []struct {
		name      string
		policy    string
		principal string
		namespace string
		action    string
		object    string
		allow     bool
	}{
		{
			name:      "rbac not allow bob",
			policy:    "rbac/authz",
			principal: "bob",
			namespace: "foo/bar",
			action:    "read",
			object:    "server123",
			allow:     false,
		},
		{
			name:      "rbac allow alice",
			policy:    "rbac/authz",
			principal: "alice",
			namespace: "foo/bar",
			action:    "read",
			object:    "server123",
			allow:     true,
		},
	}
	// our actual tests is here in the loop
	for _, c := range cases {
		body, err := json.Marshal(struct {
			Input interface{} `json:"input"`
		}{
			Input: struct {
				Principal string `json:"principal"`
				Namespace string `json:"namespace"`
				Action    string `json:"action"`
				Object    string `json:"object"`
			}{
				Principal: c.principal,
				Namespace: c.namespace,
				Action:    c.action,
				Object:    c.object,
			},
		})
		resp, err := http.Post(
			fmt.Sprintf("%s/v1/data/%s", addr, c.policy),
			"application/json",
			bytes.NewBuffer(body),
		)
		assert.Nil(t, err, c.name)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusOK, resp.StatusCode, c.name)
		body, err = ioutil.ReadAll(resp.Body)
		v := Response{}
		json.Unmarshal(body, &v)
		fmt.Print("THE RESULT")
		fmt.Print(string(body))
		assert.Equal(t, c.allow, v.Result.Allow, c.name)
	}
}

type Result struct {
	Allow bool `json:"allow"`
}

type Response struct {
	Result Result `json:"result"`
}
