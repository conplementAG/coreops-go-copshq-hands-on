package main

import (
	"errors"
	"fmt"
	"workshop/internal/constants"

	"github.com/conplementag/cops-hq/v2/pkg/hq"
	"github.com/conplementag/cops-hq/v2/pkg/recipes/azure_login"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func main() {
	hq := hq.NewQuiet("workshop-cli", "0.0.1", "workshop-cli.log")

	createCommands(hq)

	hq.Run()
}

func createCommands(hq hq.HQ) {
	// BUILD AND PUSH THE IMAGE FIRST!
	deployCommand := hq.GetCli().AddBaseCommand("deploy", "deploy command", "deploys a helm chart", func() {
		var err error
		environment := viper.GetString("environment-tag")

		// 1. login to azure
		login := azure_login.New(hq.GetExecutor())
		err = login.Login()
		if err != nil {
			panic(err)
		}

		// 2. load config
		err = hq.LoadEnvironmentConfigFile()
		if err != nil {
			panic(err)
		}

		cluster := viper.GetString(constants.Cluster)
		clusterConnectionString := viper.GetString(constants.ClusterConnectionString)
		users := viper.GetString(constants.Users)
		appName := viper.GetString(constants.AppName)
		imageTag := viper.GetString(constants.ImageTag)

		// 3. establish cluster connection
		err = connectToCluster(hq, cluster, clusterConnectionString)
		if err != nil {
			panic(err)
		}

		// 4. create / update namespace
		err = createCopsNamespace(hq, environment, users)
		if err != nil {
			panic(err)
		}

		// 5. deploy helm chart
		err = deployHelmChart(hq, cluster, environment, appName, imageTag)
		if err != nil {
			panic(err)
		}
	})

	deployCommand.AddPersistentParameterString(constants.EnvironmentTag, "", true, "e", "environment")
	deployCommand.AddPersistentParameterString(constants.ImageTag, "", true, "t", "app image tag")
}

func connectToCluster(hq hq.HQ, cluster string, connectionString string) error {
	// copsctl.New
	// copsctl.Connect
	return errors.New("not implemented")
}

func deployHelmChart(hq hq.HQ, cluster string, environment string, appName string, imageTag string) error {
	// namespace := fmt.Sprintf("ws-%s", environment)

	// helm.NewWithSettings

	// helm.SetVariables

	// helm.Deploy
	return errors.New("not implemented")
}

func createCopsNamespace(hq hq.HQ, environment string, users string) error {
	namespace := fmt.Sprintf("ws-%s", environment)
	cmd := fmt.Sprintf("copsctl namespace create -n %s -u %s -c cp-workshop -p cp-workshop", namespace, users)
	result, err := hq.GetExecutor().Execute(cmd)
	logrus.Info(result)
	return err
}
