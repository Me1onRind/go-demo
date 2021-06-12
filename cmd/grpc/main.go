package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/Me1onRind/go-demo/internal/controller"
	"github.com/Me1onRind/go-demo/internal/core/middleware"
	"github.com/Me1onRind/go-demo/internal/core/register"
	"github.com/Me1onRind/go-demo/protobuf/pb"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"google.golang.org/grpc"
)

func registerService(s *grpc.Server) {
	pb.RegisterFooServer(s, controller.NewFooController())
}

func main() {
	addr := "127.0.0.1:8080"
	ctx := context.Background()

	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer(grpc_middleware.WithUnaryServerChain(
		middleware.GrpcContext(),
		middleware.GrpcRecover(),
		middleware.GrpcLogger(),
	))
	registerService(s)

	if err := register.Register(ctx, "go-demo", addr); err != nil {
		log.Fatalf("register %s failed:%v", "go-demo", err)
	}

	fmt.Printf("start grpc server:%s", addr)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
