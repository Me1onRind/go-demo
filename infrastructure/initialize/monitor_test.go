package initialize

import (
	"context"
	"fmt"
	"time"

	//"net/http"
	"testing"

	"github.com/Me1onRind/go-demo/config"
	"github.com/Me1onRind/go-demo/infrastructure/middleware"

	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
)

func TestPrometheusServer(t *testing.T) {
	InitPrometheusServer()
	time.Sleep(time.Millisecond * 10)
	f := middleware.GrpcMonitor()
	f(context.Background(), nil, &grpc.UnaryServerInfo{
		FullMethod: "test/method",
	}, func(ctx context.Context, req interface{}) (interface{}, error) { return nil, nil })
	client := resty.New()
	resp, err := client.R().Get(fmt.Sprintf("http://%s/metrics", config.RemoteConfig.Prometheus.Addr()))
	if assert.Empty(t, err) {
		fmt.Printf("%s\n", resp.Body())
	}
}
