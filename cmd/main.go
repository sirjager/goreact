package main

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/rs/zerolog"

	"github.com/sirjager/goreact/config"

	"github.com/sirjager/goreact/cmd/gateway"
	"github.com/sirjager/goreact/cmd/grpc"

	"github.com/sirjager/goreact/internal/service"
)

var logr zerolog.Logger
var startTime time.Time
var serviceName string

func init() {
	serviceName = "goreact"
	startTime = time.Now()
	logr = zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr, NoColor: false})
	logr = logr.With().Timestamp().Logger()
	logr = logr.With().Str("service", strings.ToLower(serviceName)).Logger()
}

func main() {
	config, err := config.LoadConfigs(".", "example")
	if err != nil {
		logr.Fatal().Err(err).Msg("failed to load configurations")
	}
	config.StartTime = startTime
	config.ServiceName = serviceName

	errs := make(chan error)
	go handleSignals(errs)

	srvic, err := service.NewGoReactService(logr, config)
	if err != nil {
		logr.Fatal().Err(err).Msg("failed to create service")
	}

	if config.GatewayPort != "" {
		go gateway.RunServer(srvic)
	}

	if config.GrpcPort != "" {
		go grpc.RunServer(srvic)
	}

	logr.Error().Err(<-errs).Msg("stopped server")
}

func handleSignals(errs chan error) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	errs <- fmt.Errorf("%s", <-c)
}
