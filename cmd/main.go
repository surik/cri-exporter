package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/spf13/cobra"
	criexporter "github.com/surik/cri-exporter"
	"github.com/surik/cri-exporter/internal/app"
)

const (
	criEndpointFlag   = "container-runtime-endpoint"
	bindAddrFlag      = "bind-addr"
	metricsNamePrefix = "metrics-name-prefix"
)

var rootCmd = &cobra.Command{
	Use:     "cri-exporter",
	Short:   "A CRI exporter to export image info to Prometheus",
	Version: criexporter.Version,
	RunE:    run,
}

func run(cmd *cobra.Command, args []string) error {
	criEndpoint, err := cmd.Flags().GetString(criEndpointFlag)
	if err != nil {
		return err
	}

	bindAddr, err := cmd.Flags().GetString(bindAddrFlag)
	if err != nil {
		return err
	}

	metricsNamePrefix, err := cmd.Flags().GetString(metricsNamePrefix)
	if err != nil {
		return err
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	app, err := app.NewApp(ctx, criEndpoint, metricsNamePrefix, bindAddr)
	if err != nil {
		return err
	}

	app.Run()

	<-ctx.Done()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return app.Shutdown(ctx)
}

func main() {
	rootCmd.PersistentFlags().String(criEndpointFlag, "unix:///var/run/cri-dockerd.sock", "The endpoint of container runtime service")
	rootCmd.PersistentFlags().String(bindAddrFlag, ":9000", "The address to bind the HTTP server to expose metrics")
	rootCmd.PersistentFlags().String(metricsNamePrefix, "cri", "The prefix for the metrics name")

	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("Whoops. There was an error while executing your CLI '%s'", err)
	}
}
