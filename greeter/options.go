package greeter

import "github.com/prometheus/client_golang/prometheus"

type ClientParameters struct {
	metricsRegistry *prometheus.Registry

	timeHistogram bool
	timeSummary   bool
}

type ClientOption func(*ClientParameters)

func EnableMetrics(reg *prometheus.Registry, timeHistogram, timeSummary bool) ClientOption {
	return func(prms *ClientParameters) {
		prms.metricsRegistry = reg
		prms.timeHistogram = timeHistogram
		prms.timeSummary = timeSummary
	}
}
