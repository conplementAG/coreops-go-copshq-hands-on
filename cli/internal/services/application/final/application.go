package application

import (
	"errors"
	"fmt"
	"workshop/internal/common"
	"workshop/internal/constants"

	"github.com/conplementag/cops-hq/v2/pkg/commands"
	"github.com/conplementag/cops-hq/v2/pkg/naming/resources"
	"github.com/conplementag/cops-hq/v2/pkg/recipes/helm"
	"github.com/spf13/viper"
)

type Service struct {
	executor commands.Executor
}

func NewService(executor commands.Executor) *Service {
	return &Service{executor: executor}
}

func (s *Service) Deploy() error {
	deploymentSettings := helm.DefaultDeploymentSettings
	deploymentSettings.Wait = true

	helmDir := common.GetHelmPath()

	appName := viper.GetString(constants.AppName)
	environmentTag := viper.GetString(constants.EnvironmentTag)
	namespace := fmt.Sprintf("ws-%s", environmentTag)

	helmRecipe := helm.NewWithSettings(s.executor, namespace, appName, helmDir, deploymentSettings)
	vars, err := s.getHelmVariables()
	if err != nil {
		return err
	}
	helmRecipe.SetVariables(vars)
	err = helmRecipe.Deploy()

	return err
}

func (s *Service) Destroy(name string) error {
	return errors.New("not implemented")
}

func (s *Service) getHelmVariables() (map[string]any, error) {

	cluster := viper.GetString(constants.Cluster)
	environment := viper.GetString(constants.EnvironmentTag)
	appName := viper.GetString(constants.AppName)
	imageTag := viper.GetString(constants.ImageTag)

	namingService, err := common.GetNamingService()
	if err != nil {
		return nil, err
	}

	appKeyVaultName, err := namingService.GenerateResourceName(resources.KeyVault, "app")
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

	helmVariables := map[string]any{
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
	}

	return helmVariables, nil
}
