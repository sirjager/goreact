package gateway

import (
	"context"
	"net"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/rakyll/statik/fs"
	"google.golang.org/protobuf/encoding/protojson"

	"github.com/sirjager/goreact/internal/service"
	"github.com/sirjager/goreact/statik/docs"
	"github.com/sirjager/goreact/statik/web"

	rpc "github.com/sirjager/goreact/stubs/go"
)

func RunServer(srvic *service.Service) {
	opts := []runtime.ServeMuxOption{}
	mux := http.NewServeMux()

	opts = append(opts, runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
		MarshalOptions:   protojson.MarshalOptions{UseProtoNames: true},
		UnmarshalOptions: protojson.UnmarshalOptions{DiscardUnknown: false},
	}))

	allowedHeaders := []string{
		//
	}

	opts = append(opts, runtime.WithIncomingHeaderMatcher(AllowedHeaders(allowedHeaders)))

	grpcMux := runtime.NewServeMux(opts...)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := rpc.RegisterGoReactHandlerServer(ctx, grpcMux, srvic); err != nil {
		srvic.Logr.Fatal().Err(err).Msg("can not register handler server")
	}
	mux.Handle("/api/", grpcMux)

	// File server for swagger documentations
	docsFS, err := fs.NewWithNamespace(docs.Docs)
	if err != nil {
		srvic.Logr.Fatal().Err(err).Msg("can not serve statik docs server")
	}
	docsHandler := http.StripPrefix("/api/docs/", http.FileServer(docsFS))
	mux.Handle("/api/docs/", docsHandler)

	// File server for web app
	webFS, err := fs.NewWithNamespace(web.Web)
	if err != nil {
		srvic.Logr.Fatal().Err(err).Msg("can not serve statik web server")
	}
	webHandler := http.StripPrefix("/", http.FileServer(webFS))
	mux.Handle("/", webHandler)

	listener, err := net.Listen("tcp", ":"+srvic.Config.GatewayPort)
	if err != nil {
		srvic.Logr.Fatal().Err(err).Msg("unable to start rest gateway server")
	}

	srvic.Logr.Info().Msgf("started rest server at %s", listener.Addr().String())

	handler := Logger(srvic.Logr, mux)

	if err = http.Serve(listener, handler);err != nil {
		srvic.Logr.Fatal().Err(err).Msg("unable to serve http server")
  }
}
