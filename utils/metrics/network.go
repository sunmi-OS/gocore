package metrics

import "github.com/prometheus/client_golang/prometheus"

const (
	namespace = "network"
)

var (
	NetworkBytes = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "bytes",
			Help:      "statistic network bytes ",
		}, []string{"type"})
)

func init() {
	prometheus.MustRegister(NetworkBytes)
}
