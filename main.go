package main

import (
	"os"

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
	log.Info("Running OPA")
	if err := cmd.RootCommand.Execute(); err != nil {
		log.Error(err)
		os.Exit(1)
	}
}
