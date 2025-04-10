package main

import (
	"fmt"
	"time"
	"workshop/internal/common"
	"workshop/internal/constants"

	"github.com/conplementag/cops-hq/v2/pkg/hq"
	"github.com/conplementag/cops-hq/v2/pkg/naming"
	"github.com/conplementag/cops-hq/v2/pkg/naming/resources"
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

		// Bug in terraform recipes (init does not respect the given subscription)
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
	subscription := viper.GetString(constants.SubscriptionId)
	tenant := viper.GetString(constants.TenantId)
	regionLong := viper.GetString(constants.RegionLong)

	// 1. init terraform project
	storageSettings := terraform.DefaultBackendStorageSettings
	storageSettings.BlobContainerKey = fmt.Sprintf("%s.terraform.tfstate", environment)

	// 1.1 naming service
	namingService, err := getNamingService()
	if err != nil {
		return err
	}

	iacResourceGroupName, err := namingService.GenerateResourceName(resources.ResourceGroup, "iac")
	if err != nil {
		return err
	}

	iacStorageAccountName, err := namingService.GenerateResourceName(resources.StorageAccount, "iac")
	if err != nil {
		return err
	}

	terraformPath := common.GetTerraformPath()

	tf := terraform.New(hq.GetExecutor(),
		appName,
		subscription,
		tenant,
		regionLong,
		iacResourceGroupName,
		iacStorageAccountName,
		terraformPath,
		storageSettings,
		terraform.DefaultDeploymentSettings,
	)

	// 1.2 terraform backend
	err = tf.Init()
	if err != nil {
		return err
	}

	appResourceGroupName, err := namingService.GenerateResourceName(resources.ResourceGroup, "app")
	if err != nil {
		return err
	}

	appStorageAccountName, err := namingService.GenerateResourceName(resources.StorageAccount, "app")
	if err != nil {
		return err
	}

	appKeyVaultName, err := namingService.GenerateResourceName(resources.KeyVault, "app")
	if err != nil {
		return err
	}

	// workloadIdentityClientName, err := namingService.GenerateResourceName(resources.UserAssignedIdentity, "app")
	// if err != nil {
	// 	return err
	// }

	clientId := viper.GetString(constants.ClientId)
	clientSecret := viper.GetString(constants.ClientSecret)

	// 1.3 set variables
	err = tf.SetVariables(map[string]any{
		"app_name":        appName,
		"env_tag":         environment,
		"location":        regionLong,
		"tenant_id":       tenant,
		"subscription_id": subscription,
		"client_id":       clientId,
		"client_secret":   clientSecret,
		// resource names
		"resource_group_name":  appResourceGroupName,
		"storage_account_name": appStorageAccountName,
		"keyvault_name":        appKeyVaultName,
		// "identity_name":        workloadIdentityClientName,
	})

	if err != nil {
		return err
	}

	// 2. plan
	// DeployFlow(planOnly bool, useExistingPlan bool, autoApprove bool)
	err = tf.DeployFlow(true, false, false)
	if err != nil {
		return err
	}

	// 3. apply
	if hq.GetExecutor().AskUserToConfirmWithKeyword("Do you want to apply the changes?", "yes") {
		err = tf.DeployFlow(false, true, false)
		if err != nil {
			return err
		}
	}

	return nil
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

	namingService, err := getNamingService()
	if err != nil {
		return err
	}

	appKeyVaultName, err := namingService.GenerateResourceName(resources.KeyVault, "app")
	if err != nil {
		return err
	}

	appResourceGroupName, err := namingService.GenerateResourceName(resources.ResourceGroup, "app")
	if err != nil {
		return err
	}

	appStorageAccountName, err := namingService.GenerateResourceName(resources.StorageAccount, "app")
	if err != nil {
		return err
	}

	// identityService := identity.NewService(hq.GetExecutor())
	// workloadIdentityClientName, err := namingService.GenerateResourceName(resources.UserAssignedIdentity, "app")
	// if err != nil {
	// 	return err
	// }
	// workloadIdentityClientId, err := identityService.EnsureFederatedCredential(workloadIdentityClientName, workloadIdentityClientName, appResourceGroupName, namespace)
	// if err != nil {
	// 	return err
	// }

	helm.SetVariables(map[string]any{
		"EnvironmentTag":     environment,
		"AppImageTag":        imageTag,
		"AppImageRepository": fmt.Sprintf("%sneucopsacr.azurecr.io/%s-%s", cluster, appName, environment),
		"Host":               fmt.Sprintf("%s.%s.%s.conplement.cloud", appName, environment, cluster),
		"StorageAccountName": appStorageAccountName,
		"KeyVaultSync": map[string]any{
			"KeyVaultName":      appKeyVaultName,
			"ResourceGroupName": appResourceGroupName,
			"SubscriptionId":    viper.GetString(constants.SubscriptionId),
			"TenantId":          viper.GetString(constants.TenantId),
			"ClientId":          viper.GetString(constants.ClientId),
			"ClientSecret":      viper.GetString(constants.ClientSecret),
		},
		"WorkloadIdentity": map[string]any{
			// "ClientId":   workloadIdentityClientId,
			// "ClientName": workloadIdentityClientName,
		},
	})

	logrus.Info("Deploying helm chart...")
	err = helm.Deploy()
	return err
}
