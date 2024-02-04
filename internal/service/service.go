package service

import (
	"github.com/rs/zerolog"

	"github.com/sirjager/goreact/config"

	rpc "github.com/sirjager/goreact/stubs/go"
)

type Service struct {
	rpc.UnimplementedGoReactServer
	Config          config.Config
	Logr            zerolog.Logger
}

func NewGoReactService(
	Logr zerolog.Logger,
	config config.Config,
) (*Service, error) {

	return &Service{
		Logr:            Logr,
		Config:          config,
	}, nil
}
