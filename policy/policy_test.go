package policy

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/open-policy-agent/opa/runtime"
	"github.com/stretchr/testify/assert"
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

func TestPolicyDataIntegration(t *testing.T) {
	// requires dynamo db connection and runs local opa server
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}
	addr := "http://127.0.0.1:8080"
	// register our dynamo function
	RegisterDynamodbPolicy()
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
		assert.Equal(t, resp.StatusCode, http.StatusOK, c.name)
		body, err = ioutil.ReadAll(resp.Body)
		var v struct {
			result struct {
				allow bool
			}
		}
		json.Unmarshal(body, &v)
		assert.Equal(t, v.result.allow, c.allow, c.name)
	}

}
