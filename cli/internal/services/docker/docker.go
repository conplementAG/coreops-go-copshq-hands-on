package docker

import (
	"fmt"
	"time"
	"workshop/internal/common"
	"workshop/internal/constants"

	"github.com/conplementag/cops-hq/v2/pkg/commands"
	"github.com/spf13/viper"
)

type Service struct {
	executor   commands.Executor
	newVersion string
}

func NewService(executor commands.Executor) *Service {
	return &Service{executor: executor, newVersion: fmt.Sprintf("local_%s", time.Now().Format("20060102_150405"))}
}

func (s *Service) BuildAndPush() (string, error) {
	err := s.build()
	if err == nil {
		err = s.push()
	}

	return s.GetImageTag(), err
}

func (s *Service) build() error {
	cluster := viper.GetString(constants.Cluster)
	appName := viper.GetString(constants.AppName)
	environment := viper.GetString(constants.EnvironmentTag)
	fullImageName := fmt.Sprintf("%sneucopsacr.azurecr.io/%s-%s", cluster, appName, environment)

	imageTag := s.GetImageTag()

	cmd := fmt.Sprintf("docker build %s -f %s/Dockerfile -t %s:%s ", common.GetDockerFilePath(), common.GetDockerFilePath(), fullImageName, imageTag)

	err := s.executor.ExecuteTTY(cmd)

	return err
}

func (s *Service) push() error {
	cluster := viper.GetString(constants.Cluster)
	appName := viper.GetString(constants.AppName)
	environment := viper.GetString(constants.EnvironmentTag)
	fullImageName := fmt.Sprintf("%sneucopsacr.azurecr.io/%s-%s", cluster, appName, environment)

	imageTag := s.GetImageTag()

	s.loginToRegistry()

	cmd := fmt.Sprintf("docker push %s:%s ", fullImageName, imageTag)

	err := s.executor.ExecuteTTY(cmd)

	return err
}

func (s *Service) loginToRegistry() error {
	cluster := viper.GetString(constants.Cluster)
	registryName := fmt.Sprintf("%sneucopsacr", cluster)

	cmd := fmt.Sprintf("az acr login --name %s", registryName)
	_, err := s.executor.Execute(cmd)

	return err
}

func (s *Service) GetImageTag() string {
	version := viper.GetString(constants.ImageTag)
	if version == "" {
		version = s.newVersion
	}
	return version
}
