package infrastructure

import (
	"errors"
	"fmt"
	"workshop/internal/common"
	"workshop/internal/constants"

	"github.com/conplementag/cops-hq/v2/pkg/commands"
	"github.com/conplementag/cops-hq/v2/pkg/naming/resources"
	"github.com/conplementag/cops-hq/v2/pkg/recipes/terraform"
	"github.com/spf13/viper"
)

type Service struct {
	executor commands.Executor
}

func NewService(executor commands.Executor) *Service {
	return &Service{executor: executor}
}

func (s *Service) Deploy() error {

	tf, err := s.prepareTerraformProject()
	if err != nil {
		return err
	}

	planOnly := viper.GetBool(constants.PlanOnly)
	useExistingPlan := viper.GetBool(constants.UseExistingPlan)
	autoApprove := viper.GetBool(constants.AutoApprove)

	err = tf.DeployFlow(planOnly, useExistingPlan, autoApprove)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) Destroy(name string) error {
	return errors.New("not implemented")
}

func (s *Service) prepareTerraformProject() (terraform.Terraform, error) {
	environment := viper.GetString(constants.EnvironmentTag)
	appName := viper.GetString(constants.AppName)
	subscription := viper.GetString(constants.SubscriptionId)
	tenant := viper.GetString(constants.TenantId)
	regionLong := viper.GetString(constants.RegionLong)

	// 1. init terraform project
	storageSettings := terraform.DefaultBackendStorageSettings
	storageSettings.BlobContainerKey = fmt.Sprintf("%s.terraform.tfstate", environment)

	namingService, err := common.GetNamingService()
	if err != nil {
		return nil, err
	}

	iacResourceGroupName, err := namingService.GenerateResourceName(resources.ResourceGroup, "iac")
	if err != nil {
		return nil, err
	}

	iacStorageAccountName, err := namingService.GenerateResourceName(resources.StorageAccount, "iac")
	if err != nil {
		return nil, err
	}

	terraformPath := common.GetTerraformPath()

	tf := terraform.New(s.executor,
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

	err = tf.Init()
	if err != nil {
		return nil, err
	}

	appResourceGroupName, err := namingService.GenerateResourceName(resources.ResourceGroup, "app")
	if err != nil {
		return nil, err
	}

	appStorageAccountName, err := namingService.GenerateResourceName(resources.StorageAccount, "app")
	if err != nil {
		return nil, err
	}

	appKeyVaultName, err := namingService.GenerateResourceName(resources.KeyVault, "app")
	if err != nil {
		return nil, err
	}

	clientId := viper.GetString("client_id")
	clientSecret := viper.GetString("client_secret")

	// 1.2 set variables
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
	})

	return tf, err
}
