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
