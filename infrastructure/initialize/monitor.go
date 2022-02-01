package initialize

import (
	"log"
	"net/http"

	"github.com/Me1onRind/go-demo/config"
	"github.com/Me1onRind/go-demo/global/prometheus_vec"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func InitPromethuesServer() error {

	prometheus.MustRegister(prometheus_vec.ReqTotalCounterVec)

	log.Printf("Prometheus Server listen address: %s", config.RemoteConfig.Prometheus.Addr())
	http.Handle("/metrics", promhttp.Handler())
	return http.ListenAndServe(config.RemoteConfig.Prometheus.Addr(), nil)
}
