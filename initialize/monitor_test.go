package initialize

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"sync"
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
	_ = InitPrometheusServer(context.Background())
	time.Sleep(time.Millisecond * 10)

	total := 10
	totalStr := fmt.Sprintf("request_total{method=\"test/method\"} %d", total)

	f := middleware.GrpcMonitor()
	wg := &sync.WaitGroup{}
	for i := 0; i < total; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, _ = f(context.Background(), nil, &grpc.UnaryServerInfo{
				FullMethod: "test/method",
			}, func(ctx context.Context, req interface{}) (interface{}, error) { return nil, nil })
		}()
	}
	wg.Wait()
	client := resty.New()
	resp, err := client.R().Get(fmt.Sprintf("http://%s/metrics", config.RemoteConfig.Prometheus.Addr()))
	if assert.Empty(t, err) {
		buf := bufio.NewReader(bytes.NewBuffer(resp.Body()))
		for {
			line, _, err := buf.ReadLine()
			if errors.Is(err, io.EOF) {
				break
			}
			if assert.Empty(t, err) {
				if string(line) == totalStr {
					t.Logf("%s", line)
					return
				}
			}
		}
		t.Fatalf("Prometheus Server not response %s", totalStr)
	}
}
