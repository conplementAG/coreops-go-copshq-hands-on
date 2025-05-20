package cluster

import (
	"errors"
	"fmt"
	"workshop/internal/constants"

	"github.com/conplementag/cops-hq/v2/pkg/commands"
	"github.com/conplementag/cops-hq/v2/pkg/recipes/copsctl"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Service struct {
	executor commands.Executor
}

func NewService(executor commands.Executor) *Service {
	return &Service{executor: executor}
}

func (s *Service) Connect() error {
	copsctlService := copsctl.New(s.executor)

	clusterName := viper.GetString(constants.Cluster)
	clusterConnectionString := viper.GetString(constants.ClusterConnectionString)
	isTechnicalAccountConnect := false
	connectToSecondaryCluster := false

	return copsctlService.Connect(clusterName, clusterConnectionString, isTechnicalAccountConnect, connectToSecondaryCluster)
}

func (s *Service) CreateNamespace() error {
	// copsctlService := copsctl.New(s.executor)
	// info, err := copsctlService.GetEnvironmentInfo()
	// if err != nil {
	// 	return err
	// }

	// todo: use GetStringSlice to enable multiple users
	namespaceAdminUsers := viper.GetString(constants.Users)
	if namespaceAdminUsers == "" {
		return errors.New("namespace admin users have to be set, please check your config")
	}

	// coreOpsMandatoryTechnicalAccount := info.TechnicalAccountName + "." + info.TechnicalAccountNamespace

	// appName := viper.GetString(constants.AppName)
	environmentTag := viper.GetString(constants.EnvironmentTag)
	name := fmt.Sprintf("ws-%s", environmentTag)

	logrus.Info("Creating namespace " + name)
	_, err := s.executor.Execute(fmt.Sprintf("copsctl namespace create -n %s -u %s -c cp-workshop -p cp-workshop",
		name, namespaceAdminUsers))

	return err
}
