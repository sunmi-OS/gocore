package http_request

import "github.com/prometheus/client_golang/prometheus"

const (
	clientNamespace = "http_client"
)

var (
	clientReqDur = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: clientNamespace,
		Subsystem: "requests",
		Name:      "duration_ms",
		Help:      "http client requests duration(ms).",
		Buckets:   []float64{5, 10, 25, 50, 100, 250, 500, 1000, 2500, 5000, 10000},
	}, []string{"path"})
	clientReqCodeTotal = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: clientNamespace,
		Subsystem: "requests",
		Name:      "code_total",
		Help:      "http client requests error count.",
	}, []string{"path", "code"})
	clientSendBytes = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: clientNamespace,
		Subsystem: "bandwith",
		Name:      "send",
	}, []string{"path"})
	clientRecvBytes = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: clientNamespace,
		Subsystem: "bandwith",
		Name:      "recv",
	}, []string{"path"})
)

func init() {
	prometheus.MustRegister(clientReqDur)
	prometheus.MustRegister(clientReqCodeTotal)
	prometheus.MustRegister(clientSendBytes)
	prometheus.MustRegister(clientRecvBytes)
}
