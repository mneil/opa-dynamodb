package policy

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

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
		data      []map[string]string
		allow     bool
	}{
		{
			name:      "rbac not allow bob",
			policy:    "rbac/authz",
			principal: "baz",
			namespace: "foo/bar",
			data: []map[string]string{
				{
					"user":   "bob",
					"action": "read",
					"object": "server123",
				},
			},
			allow: false,
		},
		{
			name:      "rbac allow alice",
			policy:    "rbac/authz",
			principal: "baz",
			namespace: "foo/bar",
			data: []map[string]string{
				{
					"user":   "alice",
					"action": "read",
					"object": "server123",
				},
			},
			allow: false,
		},
	}
	// our actual tests is here in the loop
	for _, c := range cases {
		body, err := json.Marshal(struct {
			Input interface{}
		}{
			Input: struct {
				principal string
				namespace string
				data      []map[string]string
			}{
				principal: c.principal,
				namespace: c.namespace,
				data:      c.data,
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
		var v struct {
			result struct {
				allow bool
			}
		}
		json.Unmarshal(body, &v)
		assert.Equal(t, c.allow, v.result.allow, c.name)
	}

}
