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

package main

import (
	"os"

	"github.com/mneil/opa-dynamodb/policy"
	"github.com/open-policy-agent/opa/cmd"
	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)
}
func main() {
	log.Info("Entering application")
	policy.RegisterDynamodbPolicy(policy.NewDynamoConfigFromEnv())
	log.Info("Running OPA")
	if err := cmd.RootCommand.Execute(); err != nil {
		log.Error(err)
		os.Exit(1)
	}
}
