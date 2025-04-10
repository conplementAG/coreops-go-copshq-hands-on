package main

import (
	"errors"

	"github.com/conplementag/cops-hq/v2/pkg/hq"
)

func main() {
	hq := hq.NewQuiet("workshop-cli", "0.0.1", "workshop-cli.log")

	createCommands(hq)

	hq.Run()
}

func createCommands(hq hq.HQ) {
	// BUILD AND PUSH THE IMAGE FIRST!
	deployCommand := hq.GetCli().AddBaseCommand("deploy", "deploy command", "deploys a helm chart", func() {
		// 1. login to azure

		// 2. establish cluster connection

		// 3. create / update namespace

		// 4. deploy helm chart
	})

	deployCommand.AddPersistentParameterString("environment-tag", "", true, "e", "environment")
	// deployCommand.AddPersistentParameterString("cluster", "", true, "", "cluster name")
	// deployCommand.AddPersistentParameterString("cluster-connection-string", "", true, "c", "cluster connection string")
	// deployCommand.AddPersistentParameterString("users", "", true, "u", "cluster namespace users")
	// deployCommand.AddPersistentParameterString("appName", "cp-notes", false, "a", "application name")
	// deployCommand.AddPersistentParameterString("image-tag", "latest", false, "t", "app image tag")
}

func connectToCluster(hq hq.HQ, cluster string, connectionString string) error {
	// copsctl connect -e <cluster> -c <connection-string> -a
	return errors.New("not implemented")
}

func createCopsNamespace(hq hq.HQ, environment string, users string) error {
	// namespace := fmt.Sprintf("ws-%s", environment)
	// copsctl namespace create -n <namespace> -c cp-workshop -p cp-workshop -u <user>
	return errors.New("not implemented")
}

func deployHelmChart(hq hq.HQ, environment string, appName string) error {
	// namespace := fmt.Sprintf("ws-%s", environment)
	// helmPath := common.GetHelmPath()
	// helm upgrade --wait --timeout 15m --namespace <namespace> --install <appName> <chart-path>
	return errors.New("not implemented")
}
