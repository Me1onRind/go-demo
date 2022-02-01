package prometheus_vec

import "github.com/prometheus/client_golang/prometheus"

var (
	ReqTotalCounterVec = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "request_total",
		},
		[]string{"method"},
	)
)
