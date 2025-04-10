package orchestrator

import (
	"workshop/internal/common"
	"workshop/internal/services/application"
	"workshop/internal/services/cluster"
	"workshop/internal/services/infrastructure"

	"github.com/conplementag/cops-hq/v2/pkg/hq"
	"github.com/conplementag/cops-hq/v2/pkg/naming"
	"github.com/conplementag/cops-hq/v2/pkg/recipes/azure_login"
)

type Orchestrator struct {
	hq                    hq.HQ
	azureLoginService     *azure_login.Login
	applicationService    *application.Service
	infrastructureService *infrastructure.Service
	clusterService        *cluster.Service
	namingService         *naming.Service
}

func New(hq hq.HQ) *Orchestrator {
	namingService, _ := common.GetNamingService()
	executor := hq.GetExecutor()

	return &Orchestrator{
		hq:                    hq,
		azureLoginService:     azure_login.New(executor),
		applicationService:    application.NewService(executor),
		infrastructureService: infrastructure.NewService(executor),
		clusterService:        cluster.NewService(executor),
		namingService:         namingService,
	}
}

func (o *Orchestrator) Deploy() error {
	var err error

	// 1. login to azure
	err = o.azureLoginService.Login()
	if err != nil {
		panic(err)
	}

	// 2. load config
	err = o.hq.LoadEnvironmentConfigFile()
	if err != nil {
		panic(err)
	}

	err = o.clusterService.Connect()
	if err != nil {
		panic(err)
	}

	// 4. create / update namespace
	err = o.clusterService.CreateNamespace()
	if err != nil {
		panic(err)
	}

	// 5. deploy terraform project
	err = o.infrastructureService.Deploy()
	if err != nil {
		panic(err)
	}

	// 6. deploy helm chart
	err = o.applicationService.Deploy()
	if err != nil {
		panic(err)
	}

	return nil
}
