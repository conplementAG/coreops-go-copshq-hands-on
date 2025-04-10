package main

import (
	"workshop/internal/constants"
	"workshop/internal/services/application"
	"workshop/internal/services/cluster"
	"workshop/internal/services/infrastructure"

	"github.com/conplementag/cops-hq/v2/pkg/hq"
	"github.com/conplementag/cops-hq/v2/pkg/recipes/azure_login"
	"github.com/spf13/viper"
)

func main() {
	hq := hq.NewQuiet("workshop-cli", "0.0.1", "workshop-cli.log")

	createCommands(hq)

	hq.Run()
}

func createCommands(hq hq.HQ) {
	deployCommand := hq.GetCli().AddBaseCommand("deploy", "deploy command", "deploys a helm chart", func() {
		var err error

		login := azure_login.New(hq.GetExecutor())
		err = login.Login()
		if err != nil {
			panic(err)
		}

		err = hq.LoadEnvironmentConfigFile()
		if err != nil {
			panic(err)
		}

		clusterService := cluster.NewService(hq.GetExecutor())
		err = clusterService.Connect()
		if err != nil {
			panic(err)
		}

		err = clusterService.CreateNamespace()
		if err != nil {
			panic(err)
		}

		login.SetSubscription(viper.GetString(constants.SubscriptionId))

		terraformService := infrastructure.NewService(hq.GetExecutor())
		err = terraformService.Deploy()
		if err != nil {
			panic(err)
		}

		// if viper.GetBool(constants.BuildImage) {
		// 	dockerService := docker.NewService(hq.GetExecutor())
		// 	imageTag, err := dockerService.BuildAndPush()
		// 	if err != nil {
		// 		panic(err)
		// 	}
		// 	viper.Set(constants.ImageTag, imageTag)
		// }

		applicationService := application.NewService(hq.GetExecutor())
		err = applicationService.Deploy()
		if err != nil {
			panic(err)
		}
	})

	deployCommand.AddPersistentParameterString(constants.EnvironmentTag, "", true, "e", "environment")
	deployCommand.AddPersistentParameterString(constants.ImageTag, "", false, "t", "app image tag")

	// deployCommand.AddParameterBool(constants.BuildImage, false, false, "", "build application image")
	// deployCommand.AddParameterBool(constants.DeployInfrastructure, false, false, "", "deploy infrastructure")
	// deployCommand.AddParameterBool(constants.DeployApplication, false, false, "", "deploy application")

	deployCommand.AddParameterBool(constants.AutoApprove, false, false, "", "auto approve terraform changes")
	deployCommand.AddParameterBool(constants.PlanOnly, false, false, "", "create terraform plan only")
	deployCommand.AddParameterBool(constants.UseExistingPlan, false, false, "", "execute existing terraform plan")

}

func getDeploymentConfiguration() (bool, bool) {
	deployInfrastructure := viper.GetBool(constants.DeployInfrastructure)
	deployApplication := viper.GetBool(constants.DeployApplication)

	// If both flags are not set, default to true for both
	if !deployInfrastructure && !deployApplication {
		deployInfrastructure = true
		deployApplication = true
	}

	return deployInfrastructure, deployApplication
}
