package main

import (
	"context"
	"flag"
	"net/http"

	"github.com/golang/glog"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"

	gw_v1 "github.com/pepese/grpc-gateway-sample/proto/dest/helloworld/v1"
	gw_v2 "github.com/pepese/grpc-gateway-sample/proto/dest/helloworld/v2"
)

var (
	// command-line options:
	// gRPC server endpoint
	grpcServerEndpoint = flag.String("grpc-server-endpoint", "localhost:8080", "gRPC server endpoint")
)

func run() error {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Register gRPC server endpoint
	// Note: Make sure the gRPC server is running properly and accessible
	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}

	for _, f := range []func(context.Context, *runtime.ServeMux, string, []grpc.DialOption) error{
		gw_v1.RegisterHelloServiceHandlerFromEndpoint,
		gw_v2.RegisterHelloServiceHandlerFromEndpoint,
	} {
		if err := f(ctx, mux, *grpcServerEndpoint, opts); err != nil {
			return err
		}
	}

	// Start HTTP server (and proxy calls to gRPC server endpoint)
	return http.ListenAndServe(":8081", mux)
}

func main() {
	flag.Parse()
	defer glog.Flush()

	if err := run(); err != nil {
		glog.Fatal(err)
	}
}
