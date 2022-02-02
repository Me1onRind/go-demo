package initialize

import (
	"context"
	"log"
	"net/http"

	"github.com/Me1onRind/go-demo/config"
	"github.com/Me1onRind/go-demo/global/prometheus_vec"
	"github.com/Me1onRind/go-demo/infrastructure/logger"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
)

func InitPrometheusServer(ctx context.Context) error {

	prometheus.MustRegister(prometheus_vec.ReqTotalCounterVec)

	log.Printf("Prometheus Server listen address: %s", config.RemoteConfig.Prometheus.Addr())
	http.Handle("/metrics", promhttp.Handler())
	go func() {
		if err := http.ListenAndServe(config.RemoteConfig.Prometheus.Addr(), nil); err != nil {
			logger.CtxError(ctx, "Prometheus Server ListenAndServe fail", zap.Error(err))
		}
	}()
	return nil
}
