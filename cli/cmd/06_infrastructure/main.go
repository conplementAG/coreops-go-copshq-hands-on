package main

import (
	"errors"
	"fmt"
	"time"
	"workshop/internal/common"
	"workshop/internal/constants"

	"github.com/conplementag/cops-hq/v2/pkg/hq"
	"github.com/conplementag/cops-hq/v2/pkg/naming"
	"github.com/conplementag/cops-hq/v2/pkg/recipes/azure_login"
	"github.com/conplementag/cops-hq/v2/pkg/recipes/copsctl"
	"github.com/conplementag/cops-hq/v2/pkg/recipes/helm"
	"github.com/conplementag/cops-hq/v2/pkg/recipes/terraform"
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

		login.SetSubscription(viper.GetString(constants.SubscriptionId))

		// 5. deploy terraform project
		err = deployTerraformProject(hq, environment, appName)
		if err != nil {
			panic(err)
		}

		// 6. deploy helm chart
		err = deployHelmChart(hq, cluster, environment, appName, imageTag)
		if err != nil {
			panic(err)
		}
	})

	deployCommand.AddPersistentParameterString(constants.EnvironmentTag, "", true, "e", "environment")
	deployCommand.AddPersistentParameterString(constants.ImageTag, "latest", false, "t", "app image tag")
}

func getNamingService() (*naming.Service, error) {
	regionLong := viper.GetString(constants.RegionLong)
	env := viper.GetString(constants.EnvironmentTag)
	namingService, err := naming.New("ws", regionLong, env, "", "g")
	if err != nil {
		return nil, err
	}
	return namingService, nil
}

func deployTerraformProject(hq hq.HQ, environment string, appName string) error {
	// subscription := viper.GetString(constants.SubscriptionId)
	// tenant := viper.GetString(constants.TenantId)
	// regionLong := viper.GetString(constants.RegionLong)

	// 1. init terraform project
	storageSettings := terraform.DefaultBackendStorageSettings
	storageSettings.BlobContainerKey = fmt.Sprintf("%s.terraform.tfstate", environment)

	// 1.1 naming service

	// terraformPath := common.GetTerraformPath()

	// 1.2 terraform backend

	// 2. plan
	// DeployFlow(planOnly bool, useExistingPlan bool, autoApprove bool)

	// 3. apply

	return errors.New("not implemented")
}

func connectToCluster(hq hq.HQ, cluster string, connectionString string) error {
	copsctl := copsctl.New(hq.GetExecutor())
	err := copsctl.Connect(cluster, connectionString, false, false)
	return err
}

func createCopsNamespace(hq hq.HQ, environment string, users string) error {
	namespace := fmt.Sprintf("ws-%s", environment)
	cmd := fmt.Sprintf("copsctl namespace create -n %s -u %s -c cp-workshop -p cp-workshop", namespace, users)
	result, err := hq.GetExecutor().Execute(cmd)
	logrus.Info(result)
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
