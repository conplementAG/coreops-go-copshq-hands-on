package infrastructure

import (
	"errors"

	"github.com/conplementag/cops-hq/v2/pkg/commands"
	"github.com/conplementag/cops-hq/v2/pkg/recipes/terraform"
)

type Service struct {
	executor commands.Executor
}

func NewService(executor commands.Executor) *Service {
	return &Service{executor: executor}
}

func (s *Service) Deploy() error {
	return errors.New("not implemented")
}

func (s *Service) Destroy(name string) error {
	return errors.New("not implemented")
}

func (s *Service) prepareTerraformProject() (terraform.Terraform, error) {

	return nil, errors.New("not implemented")
}
