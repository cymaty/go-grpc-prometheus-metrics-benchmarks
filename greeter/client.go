package greeter

import (
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"google.golang.org/grpc"
)

type Client struct {
	grpcClientConn *grpc.ClientConn

	GreeterClient
}

func NewClient(addr string, opts ...ClientOption) (*Client, error) {
	clientParams := new(ClientParameters)
	for _, opt := range opts {
		opt(clientParams)
	}

	dialOpts := []grpc.DialOption{grpc.WithInsecure()}

	if clientParams.metricsRegistry != nil {
		var counterOpts []grpc_prometheus.CounterOption
		clientMetrics := grpc_prometheus.NewClientMetrics(counterOpts...)

		if clientParams.timeHistogram {
			var histogramOpts []grpc_prometheus.HistogramOption
			clientMetrics.EnableClientHandlingTimeHistogram(histogramOpts...)
		}

		withUnaryInterceptor := grpc.WithUnaryInterceptor(clientMetrics.UnaryClientInterceptor())

		clientParams.metricsRegistry.MustRegister(clientMetrics)

		if clientParams.timeSummary {
			foo := NewClientMetrics()
			clientParams.metricsRegistry.MustRegister(foo)

			withUnaryInterceptor = grpc.WithUnaryInterceptor(grpc_middleware.ChainUnaryClient(
				clientMetrics.UnaryClientInterceptor(),
				foo.UnaryClientInterceptor(),
			))
		}

		dialOpts = append(
			dialOpts,
			grpc.WithStreamInterceptor(clientMetrics.StreamClientInterceptor()),
			withUnaryInterceptor,
		)
	}

	conn, err := grpc.Dial(addr, dialOpts...)
	if err != nil {
		return nil, err
	}

	return &Client{
		grpcClientConn: conn,

		GreeterClient: NewGreeterClient(conn),
	}, nil
}

func (cl *Client) Close() error {
	return cl.grpcClientConn.Close()
}
