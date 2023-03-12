package app

import (
	"net/http"

	"github.com/Me1onRind/go-demo/internal/infrastructure/logger"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Metrics struct {
}

func NewMetrics() *Metrics {
	return &Metrics{}
}

func (m *Metrics) Run() error {
	server := http.NewServeMux()
	server.Handle("/metrics", promhttp.Handler())
	addr := "0.0.0.0:9002"
	logger.Infof("prometheus listen: %s", addr)
	if err := http.ListenAndServe(addr, server); err != nil {
		logger.Errorf("Prometheus listen err:[%s]", err)
		return err
	}
	return nil
}
