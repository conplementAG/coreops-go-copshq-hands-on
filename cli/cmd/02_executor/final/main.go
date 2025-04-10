package main

import (
	"encoding/json"
	"workshop/cmd/02_executor/version"

	"github.com/conplementag/cops-hq/v2/pkg/hq"
	"github.com/sirupsen/logrus"
)

func main() {
	hq := hq.NewQuiet("workshop-cli", "0.0.1", "workshop-cli.log")

	createCommands(hq)

	hq.Run()
}

func createCommands(hq hq.HQ) {
	hq.GetCli().AddBaseCommand("get-version", "get version command", "print version of az cli", func() {
		// 1. execute 'az version' command
		result, err := hq.GetExecutor().Execute("az version")
		if err != nil {
			logrus.Errorf("Error executing command: %v", err)
		}
		logrus.Info(result)

		// 2. unmarshal json result to AzureCliVersion
		var version version.AzureCliVersion
		err = json.Unmarshal([]byte(result), &version)
		if err != nil {
			logrus.Errorf("Error unmarshalling command result: %v", err)
		}
		logrus.Infof("%s", version.AzureCli)
	})
}
