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
	// BUILD AND PUSH THE IMAGE FIRST!
	deployCommand := hq.GetCli().AddBaseCommand("deploy", "deploy command", "deploys a helm chart", func() {
		var err error

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

		// 3. connect to cluster
		clusterService := cluster.NewService(hq.GetExecutor())
		err = clusterService.Connect()
		if err != nil {
			panic(err)
		}

		// 4. create / update namespace
		err = clusterService.CreateNamespace()
		if err != nil {
			panic(err)
		}

		login.SetSubscription(viper.GetString(constants.SubscriptionId))

		// 5. deploy terraform project
		terraformService := infrastructure.NewService(hq.GetExecutor())
		err = terraformService.Deploy()
		if err != nil {
			panic(err)
		}

		// 6. deploy helm chart
		applicationService := application.NewService(hq.GetExecutor())
		err = applicationService.Deploy()
		if err != nil {
			panic(err)
		}
	})

	deployCommand.AddPersistentParameterString(constants.EnvironmentTag, "", true, "e", "environment")
	deployCommand.AddPersistentParameterString(constants.ImageTag, "latest", false, "t", "app image tag")

	deployCommand.AddParameterBool(constants.AutoApprove, false, false, "", "auto approve terraform changes")
	deployCommand.AddParameterBool(constants.PlanOnly, false, false, "", "create terraform plan only")
	deployCommand.AddParameterBool(constants.UseExistingPlan, false, false, "", "execute existing terraform plan")
}
