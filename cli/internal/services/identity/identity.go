package identity

import (
	"fmt"
	"strings"
	"workshop/internal/services/cluster"

	"github.com/conplementag/cops-hq/v2/pkg/commands"
	"github.com/conplementag/cops-hq/v2/pkg/recipes/copsctl"
)

const audiences = "api://AzureADTokenExchange"

type Service struct {
	executor       commands.Executor
	clusterService *cluster.Service
}

func NewService(hqExecutor commands.Executor) *Service {
	return &Service{
		executor:       hqExecutor,
		clusterService: cluster.NewService(hqExecutor),
	}
}

func (s *Service) GetClientId(resourceGroup string, managedIdentityName string) (string, error) {
	command := fmt.Sprintf("az identity show -n %s -g %s --query clientId -otsv", managedIdentityName, resourceGroup)
	result, err := s.executor.ExecuteSilent(command)
	if err != nil {
		return "", err
	}
	return strings.TrimRight(result, "\r"), nil
}

func (s *Service) EnsureFederatedCredential(identityName string, federatedCredentialName string, resourceGroupName string, namespaceName string) (string, error) {
	clusterInfo, err := copsctl.New(s.executor).GetClusterInfo()
	if err != nil {
		return "", err
	}

	// the name of the federated-credential will be suffixed with a cluster specific description, because in blue/green cluster migration scenarios the Oidc issuer url changes
	federatedCredentialFullName := federatedCredentialName + "-" + clusterInfo.Description
	// the subject must match the full name of serviceaccount name in the cluster
	subject := "system:serviceaccount:" + namespaceName + ":" + federatedCredentialName
	command := fmt.Sprintf("az identity federated-credential create --identity-name \"%s\" --name \"%s\" --resource-group \"%s\" --audiences \"%s\" --issuer \"%s\" --subject \"%s\"",
		identityName,
		federatedCredentialFullName,
		resourceGroupName,
		audiences,
		clusterInfo.OidcIssuerUrl,
		subject)

	_, err = s.executor.ExecuteSilent(command)
	if err != nil {
		return "", err
	}

	return s.GetClientId(resourceGroupName, identityName)
}

func (s *Service) DeleteFederatedCredential(identityName string, federatedCredentialName string, resourceGroupName string) error {
	clusterInfo, err := copsctl.New(s.executor).GetClusterInfo()
	if err != nil {
		return err
	}

	// the name of the federated-credential will be suffixed with a cluster specific description, because in blue/green cluster migration scenarios the Oidc issuer url changes
	federatedCredentialFullName := federatedCredentialName + "-" + clusterInfo.Description
	command := fmt.Sprintf("az identity federated-credential delete --identity-name \"%s\" --name \"%s\" --resource-group \"%s\" --yes",
		identityName,
		federatedCredentialFullName,
		resourceGroupName)

	_, err = s.executor.ExecuteSilent(command)

	return err
}
