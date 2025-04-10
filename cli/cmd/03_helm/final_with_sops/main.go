package main

import (
	"fmt"
	"workshop/internal/common"
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

		err = hq.LoadConfigFile(fmt.Sprintf("../../config/%s.yaml", environment))
		// or load per convention
		// err = hq.LoadEnvironmentConfigFile()
		if err != nil {
			panic(err)
		}

		cluster := viper.GetString(constants.Cluster)
		clusterConnectionString := viper.GetString(constants.ClusterConnectionString)
		users := viper.GetString(constants.Users)
		appName := viper.GetString(constants.AppName)

		// 2. establish cluster connection
		err = connectToCluster(hq, cluster, clusterConnectionString)
		if err != nil {
			panic(err)
		}

		// 3. create / update namespace
		err = createCopsNamespace(hq, environment, users)
		if err != nil {
			panic(err)
		}

		// 4. deploy helm chart
		err = deployHelmChart(hq, environment, appName)
		if err != nil {
			panic(err)
		}
	})

	deployCommand.AddPersistentParameterString(constants.EnvironmentTag, "", true, "e", "environment")
	deployCommand.AddPersistentParameterString(constants.ImageTag, "latest", false, "t", "app image tag")
}

func connectToCluster(hq hq.HQ, cluster string, connectionString string) error {
	cmd := fmt.Sprintf("copsctl connect -e %s -c \"%s\" -a", cluster, connectionString)
	result, err := hq.GetExecutor().Execute(cmd)
	logrus.Info(result)
	return err
}

func createCopsNamespace(hq hq.HQ, environment string, users string) error {
	namespace := fmt.Sprintf("ws-%s", environment)
	cmd := fmt.Sprintf("copsctl namespace create -n %s -u %s -c cp-workshop -p cp-workshop", namespace, users)
	result, err := hq.GetExecutor().Execute(cmd)
	logrus.Info(result)
	return err
}

func deployHelmChart(hq hq.HQ, environment string, appName string) error {
	namespace := fmt.Sprintf("ws-%s", environment)
	helmPath := common.GetHelmPath()
	cmd := fmt.Sprintf("helm upgrade --wait --timeout 15m --namespace %s --install %s %s", namespace, appName, helmPath)
	result, err := hq.GetExecutor().Execute(cmd)
	logrus.Info(result)
	return err
}
