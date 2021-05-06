package greeter

import (
	context "context"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc"
)

type ClientMetrics struct {
	clientHandledSummary *prometheus.SummaryVec
}

func NewClientMetrics() *ClientMetrics {
	return &ClientMetrics{
		clientHandledSummary: prometheus.NewSummaryVec(
			prometheus.SummaryOpts{
				Name:       "grpc_client_handling2_seconds",
				Help:       "A summary metric used to record Sprocket API client requests latency.",
				Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
				// ConstLabels: prometheus.Labels{appLabel: appName},
			},
			[]string{},
			// []string{operationLabel},
		),

		// clientHandledSummary: prometheus.HistogramOpts{
		// 	Name: "grpc_client_handling_seconds",
		// 	Help: "Histogram of response latency (seconds) of the gRPC until it is finished by the application.",
		// 	// Buckets: prom.DefBuckets,
		// },
	}
}

// Describe sends the super-set of all possible descriptors of metrics
// collected by this Collector to the provided channel and returns once
// the last descriptor has been sent.
func (m *ClientMetrics) Describe(ch chan<- *prometheus.Desc) {
	m.clientHandledSummary.Describe(ch)
}

// Collect is called by the Prometheus registry when collecting
// metrics. The implementation sends each collected metric via the
// provided channel and returns once the last metric has been sent.
func (m *ClientMetrics) Collect(ch chan<- prometheus.Metric) {
	m.clientHandledSummary.Collect(ch)
}

func (m *ClientMetrics) UnaryClientInterceptor() func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		start := time.Now()
		err := invoker(ctx, method, req, reply, cc, opts...)
		durationSec := time.Since(start).Seconds()

		m.clientHandledSummary.WithLabelValues().Observe(durationSec)

		return err
	}
}
