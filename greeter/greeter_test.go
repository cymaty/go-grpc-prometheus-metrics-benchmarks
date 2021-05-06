package greeter_test

import (
	"bufio"
	"context"
	"net"
	"os"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/expfmt"
	"google.golang.org/grpc"

	"go-lab/grpc/greeter"
)

// https://prometheus.io/docs/concepts/metric_types/
// https://github.com/grpc-ecosystem/go-grpc-prometheus
// https://dev.to/aleksk1ng/go-grpc-clean-architecture-microservice-with-prometheus-grafana-monitoring-and-jaeger-opentracing-51om
// https://grpc.io/docs/languages/go/basics/
// https://github.com/grpc/grpc-go/blob/master/examples/helloworld/greeter_server/main.go

// https://pkg.go.dev/google.golang.org/protobuf
// https://golang.org/pkg/testing/
// https://pkg.go.dev/github.com/prometheus/client_golang@v1.10.0/prometheus#hdr-A_Basic_Example
// https://github.com/prometheus/client_golang/blob/master/prometheus/examples_test.go

func TestGRPC(t *testing.T) {
	l, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		t.Error(err)
	}

	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	greeter.RegisterGreeterServer(grpcServer, &greeter.Server{})
	go func() {
		grpcServer.Serve(l)
	}()

	greeterClient, err := greeter.NewClient(l.Addr().String())
	if err != nil {
		t.Fatal(err)
	}

	req := &greeter.HelloRequest{
		Name: "Foo",
	}
	resp, err := greeterClient.SayHello(context.Background(), req)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(resp.Message)
}

func BenchmarkMetrics(b *testing.B) {
	addr := runServer(b)

	tCases := map[string]struct {
		ClientFn func(*testing.B, *prometheus.Registry) *greeter.Client
	}{
		"NoMetrics": {
			ClientFn: func(b *testing.B, _ *prometheus.Registry) *greeter.Client {
				greeterClient, err := greeter.NewClient(addr)
				if err != nil {
					b.Fatal(err)
				}
				return greeterClient
			},
		},

		"Metrics": {
			ClientFn: func(b *testing.B, metricsRegistry *prometheus.Registry) *greeter.Client {
				greeterClient, err := greeter.NewClient(addr, greeter.EnableMetrics(metricsRegistry, false, false))
				if err != nil {
					b.Fatal(err)
				}
				return greeterClient
			},
		},

		"TimeHistogram": {
			ClientFn: func(b *testing.B, metricsRegistry *prometheus.Registry) *greeter.Client {
				greeterClient, err := greeter.NewClient(addr, greeter.EnableMetrics(metricsRegistry, true, false))
				if err != nil {
					b.Fatal(err)
				}
				return greeterClient
			},
		},

		"TimeSummary": {
			ClientFn: func(b *testing.B, metricsRegistry *prometheus.Registry) *greeter.Client {
				greeterClient, err := greeter.NewClient(addr, greeter.EnableMetrics(metricsRegistry, false, true))
				if err != nil {
					b.Fatal(err)
				}
				return greeterClient
			},
		},
	}

	for tn, tc := range tCases {
		b.Run(tn, func(b *testing.B) {
			metricsRegistry := prometheus.NewRegistry()
			b.Cleanup(func() {
				// outputMetrics(b, metricsRegistry)
			})

			greeterClient := tc.ClientFn(b, metricsRegistry)

			req := &greeter.HelloRequest{
				Name: "Foo",
			}

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_, err := greeterClient.SayHello(context.Background(), req)
				if err != nil {
					b.Fatal(err)
				}
			}
			b.StopTimer()
		})
	}
}

func runServer(b *testing.B) string {
	b.Helper()

	l, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		b.Error(err)
	}

	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)

	greeter.RegisterGreeterServer(grpcServer, &greeter.Server{})

	go func() {
		if err := grpcServer.Serve(l); err != nil {
			b.Fatal(err)
		}
	}()

	return l.Addr().String()
}

func outputMetrics(b *testing.B, reg *prometheus.Registry) {
	b.Helper()

	mfs, err := reg.Gather()
	if err != nil {
		b.Fatal(err)
	}

	w := bufio.NewWriter(os.Stdout)
	enc := expfmt.NewEncoder(w, expfmt.FmtOpenMetrics)

	for _, mf := range mfs {
		if err := enc.Encode(mf); err != nil {
			b.Fatal(b)
		}
	}

	if closer, ok := enc.(expfmt.Closer); ok {
		// Check err
		closer.Close()

		// if handleError(closer.Close()) {
		// 	return
		// }
	}

	w.Flush()
}
