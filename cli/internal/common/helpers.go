package common

import (
	"path/filepath"
	"workshop/internal/constants"

	"github.com/conplementag/cops-hq/v2/pkg/hq"
	"github.com/conplementag/cops-hq/v2/pkg/naming"
	"github.com/spf13/viper"
)

func GetHelmPath() string {
	return filepath.Join(hq.ProjectBasePath, "helm")
}

func GetTerraformPath() string {
	return filepath.Join(hq.ProjectBasePath, "terraform")
}

func GetConfigPath() string {
	return filepath.Join(hq.ProjectBasePath, "config")
}

func GetDockerFilePath() string {
	return filepath.Join(hq.ProjectBasePath, "..", "app")
}

// g is short for green. If we ever need to create a parallel infra for an environment (like in a
// disaster recovery scenario), change this to b (as in blue)
const color = "g"

func GetNamingService() (*naming.Service, error) {
	environmentTag := viper.GetString(constants.EnvironmentTag)
	region := viper.GetString(constants.RegionLong)
	namingService, err := naming.New("ws", region, environmentTag, "", color)
	return namingService, err
}
