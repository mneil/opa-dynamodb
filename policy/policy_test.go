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
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}
	addr := "http://127.0.0.1:8080"
	RegisterDynamodbPolicy()
	ctx := context.Background()
	defer ctx.Done()
	err := startLocalServer(ctx, addr)
	assert.Nil(t, err)
	body, err := json.Marshal(struct {
		Input interface{}
	}{
		Input: struct {
			Principal string
			Namespace string
			Data      []map[string]string
		}{
			Principal: "baz",
			Namespace: "foo/bar",
			Data: []map[string]string{
				{"server": "example"},
			},
		},
	})
	resp, err := http.Post(
		fmt.Sprintf("%s/v1/example", addr),
		"application/json",
		bytes.NewBuffer(body),
	)
	assert.Nil(t, err)
	defer resp.Body.Close()
	assert.Equal(t, resp.StatusCode, http.StatusOK, "request should be ok")
	body, err = ioutil.ReadAll(resp.Body)
	var v struct {
		result struct {
			allow bool
		}
	}
	json.Unmarshal(body, &v)
	assert.Equal(t, v.result.allow, false)
}
