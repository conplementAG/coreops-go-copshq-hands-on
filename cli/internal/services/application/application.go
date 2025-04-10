package application

import (
	"errors"

	"github.com/conplementag/cops-hq/v2/pkg/commands"
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

func (s *Service) getHelmVariables() map[string]any {
	helmVariables := map[string]any{}

	return helmVariables
}
