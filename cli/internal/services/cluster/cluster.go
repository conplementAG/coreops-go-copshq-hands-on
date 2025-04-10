package cluster

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

func (s *Service) Connect() error {
	return errors.New("not implemented")
}

func (s *Service) CreateNamespace() error {
	return errors.New("not implemented")
}
