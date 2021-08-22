package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/Me1onRind/go-demo/internal/controller/foo_controller"
	"github.com/Me1onRind/go-demo/internal/core/initialize"
	"github.com/Me1onRind/go-demo/internal/core/middleware"
	"github.com/Me1onRind/go-demo/internal/core/register"
	"github.com/Me1onRind/go-demo/protobuf/pb"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"google.golang.org/grpc"
)

func Init() {
	initFuncs := []func() error{
		initialize.InitOpentracking("go-grpc-demo", "0.0.1"),
	}

	for _, v := range initFuncs {
		if err := v(); err != nil {
			panic(err)
		}
	}
}

func registerService(s *grpc.Server) {
	pb.RegisterFooServer(s, foo_controller.NewFooController())
}

func main() {
	Init()
	addr := "127.0.0.1:8080"
	ctx := context.Background()

	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer(grpc_middleware.WithUnaryServerChain(
		middleware.GrpcContext(),
		middleware.GrpcRecover(),
		middleware.GrpcTracer(),
		middleware.GrpcLogger(),
	))
	registerService(s)

	if err := register.Register(ctx, "go-demo", addr); err != nil {
		log.Fatalf("register %s failed:%v", "go-demo", err)
	}

	fmt.Printf("start grpc server:%s\n", addr)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
