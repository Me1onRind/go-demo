package main

import (
	"context"
	"testing"
	"time"

	"github.com/Me1onRind/go-demo/protobuf/pb"
	"go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/naming/resolver"
	"google.golang.org/grpc"
)

func Test_Greet(t *testing.T) {

	cli, err := clientv3.NewFromURL("http://localhost:2379")
	if err != nil {
		t.Fatal(err)
	}
	builder, err := resolver.NewBuilder(cli)
	if err != nil {
		t.Fatal(err)
	}
	conn, err := grpc.Dial("etcd:///service/go-demo",
		grpc.WithResolvers(builder),
		grpc.WithBalancerName("round_robin"),
		grpc.WithInsecure(), grpc.WithTimeout(time.Second))
	if err != nil {
		t.Fatal(err)
	}

	fooClient := pb.NewFooClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()
	resp, err := fooClient.Greet(ctx, &pb.GreetReq{
		MyName: "Bar",
		Msg:    "Hello, World",
	})
	if err != nil {
		t.Fatal(err)
	}

	t.Log(resp.Msg)
}
