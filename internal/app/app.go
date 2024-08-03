package app

import (
	"context"
	"log"
	"net"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type app struct {
	agent   Agent
	metrics *Metrics
	router  *http.ServeMux

	server   *http.Server
	listener net.Listener
}

// NewApp creates a new app instance.
func NewApp(ctx context.Context, criEndpoint, metricsNamePrefix, bindAddr string) (*app, error) {
	agent, err := NewAgent(ctx, criEndpoint)
	if err != nil {
		log.Fatal(err)
	}

	metrics := NewMetrics(agent, metricsNamePrefix)
	reg := prometheus.NewRegistry()

	reg.MustRegister(
		collectors.NewGoCollector(),
		collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}),
		collectors.NewBuildInfoCollector(),
		metrics,
	)

	router := http.NewServeMux()
	router.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{Registry: reg}))

	listener, err := net.Listen("tcp", bindAddr)
	if err != nil {
		return nil, err
	}

	server := &http.Server{
		Handler: router,
		Addr:    bindAddr,
	}

	return &app{
		agent:    agent,
		metrics:  metrics,
		router:   router,
		server:   server,
		listener: listener,
	}, nil
}

// Run starts the app.
func (a *app) Run() {
	go func() {
		_ = a.server.Serve(a.listener)
	}()
	log.Printf("Listening on %s", a.server.Addr)
}

// Shutdown stops the app.
func (a *app) Shutdown(ctx context.Context) error {
	return a.server.Shutdown(ctx)
}
