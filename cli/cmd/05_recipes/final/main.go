package main

import (
	"fmt"
	"time"
	"workshop/internal/common"
	"workshop/internal/constants"

	"github.com/conplementag/cops-hq/v2/pkg/hq"
	"github.com/conplementag/cops-hq/v2/pkg/recipes/azure_login"
	"github.com/conplementag/cops-hq/v2/pkg/recipes/copsctl"
	"github.com/conplementag/cops-hq/v2/pkg/recipes/helm"
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
		environment := viper.GetString(constants.EnvironmentTag)

		login := azure_login.New(hq.GetExecutor())
		err = login.Login()
		if err != nil {
			panic(err)
		}

		err = hq.LoadEnvironmentConfigFile()
		if err != nil {
			panic(err)
		}

		cluster := viper.GetString(constants.Cluster)
		clusterConnectionString := viper.GetString(constants.ClusterConnectionString)
		users := viper.GetString(constants.Users)
		appName := viper.GetString(constants.AppName)
		imageTag := viper.GetString(constants.ImageTag)

		// 1. establish cluster connection
		err = connectToCluster(hq, cluster, clusterConnectionString)
		if err != nil {
			panic(err)
		}

		// 2. create / update namespace
		err = createCopsNamespace(hq, environment, users)
		if err != nil {
			panic(err)
		}
		// 3. deploy helm chart
		err = deployHelmChart(hq, cluster, environment, appName, imageTag)
		if err != nil {
			panic(err)
		}
	})

	deployCommand.AddPersistentParameterString(constants.EnvironmentTag, "", true, "e", "environment")
	deployCommand.AddPersistentParameterString(constants.ImageTag, "latest", false, "t", "app image tag")
}

func connectToCluster(hq hq.HQ, cluster string, connectionString string) error {
	copsctl := copsctl.New(hq.GetExecutor())
	err := copsctl.Connect(cluster, connectionString, false, false)
	return err
}

func deployHelmChart(hq hq.HQ, cluster string, environment string, appName string, imageTag string) error {
	namespace := fmt.Sprintf("ws-%s", environment)
	helmPath := common.GetHelmPath()

	helm := helm.NewWithSettings(hq.GetExecutor(), namespace, appName, helmPath, helm.DeploymentSettings{
		Wait:    true,
		Timeout: 15 * time.Minute,
	})

	helm.SetVariables(map[string]any{
		"EnvironmentTag":     environment,
		"AppImageTag":        imageTag,
		"AppImageRepository": fmt.Sprintf("%sneucopsacr.azurecr.io/%s-%s", cluster, appName, environment),
		"Host":               fmt.Sprintf("%s.%s.%s.conplement.cloud", appName, environment, cluster),
	})

	logrus.Info("Deploying helm chart...")
	err := helm.Deploy()
	return err
}

func createCopsNamespace(hq hq.HQ, environment string, users string) error {
	namespace := fmt.Sprintf("ws-%s", environment)
	cmd := fmt.Sprintf("copsctl namespace create -n %s -u %s -c cp-workshop -p cp-workshop", namespace, users)
	result, err := hq.GetExecutor().Execute(cmd)
	logrus.Info(result)
	return err
}
