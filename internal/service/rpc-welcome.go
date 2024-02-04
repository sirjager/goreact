package service

import (
	"context"
	"fmt"

	rpc "github.com/sirjager/goreact/stubs/go"
)

func welcomeMessage(name string) string {
	return fmt.Sprintf("Welcome to %s Api", name)
}

func (s *Service) Welcome(ctx context.Context, req *rpc.WelcomeRequest) (*rpc.WelcomeResponse, error) {
	return &rpc.WelcomeResponse{Message: welcomeMessage(s.Config.ServiceName)}, nil
}
