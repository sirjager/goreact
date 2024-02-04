package grpc

import (
	"net"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/sirjager/goreact/internal/service"

	rpc "github.com/sirjager/goreact/stubs/go"
)

func RunServer(srvic *service.Service) {
	listener, err := net.Listen("tcp", ":"+srvic.Config.GrpcPort)
	if err != nil {
		srvic.Logr.Fatal().Err(err).Msg("unable to listen grpc tcp server")
	}

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(
			grpc_middleware.ChainUnaryServer(
				Logger(srvic.Logr),
			),
		),
		grpc.StreamInterceptor(
			grpc_middleware.ChainStreamServer(
				StreamLogger(srvic.Logr),
			),
		),

		grpc.MaxRecvMsgSize(1024*1024), // bytes * Kilobytes * Megabytes
	)

	rpc.RegisterGoReactServer(grpcServer, srvic)

	reflection.Register(grpcServer)

	srvic.Logr.Info().Msgf("started grpc server at %s", listener.Addr().String())

	err = grpcServer.Serve(listener)
	if err != nil {
		srvic.Logr.Fatal().Err(err).Msg("unable to serve gRPC server")
	}
}
