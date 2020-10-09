package main

import (
	"flag"
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/veertuinc/anka-prometheus-exporter/src/client"
	"github.com/veertuinc/anka-prometheus-exporter/src/log"
	"github.com/veertuinc/anka-prometheus-exporter/src/metrics"
	"github.com/veertuinc/anka-prometheus-exporter/src/server"
)

const (
	DEFAULT_INTERVAL_SECONDS = 10
)

var (
	version string
)

func main() {
	var controllerAddress string
	var intervalSeconds int
	var disableOptimizeInterval bool
	var port int
	var caFilePath string
	var clientCertPath string
	var clientCertKeyPath string
	var skipTLSVerification bool
	var useTLS bool

	flag.StringVar(&controllerAddress, "controller_address", "", "Controller address to monitor")
	flag.IntVar(&intervalSeconds, "interval", DEFAULT_INTERVAL_SECONDS, "Seconds to wait between data requests to controller")
	flag.IntVar(&port, "port", 2112, "Port to server /metrics endpoint")
	flag.BoolVar(&disableOptimizeInterval, "disable_interval_optimizer", false, "Optimize interval according to /metric api requests receieved")
	flag.BoolVar(&useTLS, "tls", false, "Use TLS")
	flag.BoolVar(&skipTLSVerification, "skip_tls_verification", false, "Skip TLS verification")
	flag.StringVar(&caFilePath, "ca_cert", "", "Path to ca PEM/x509 file")
	flag.StringVar(&clientCertPath, "client_cert", "", "Path to client cert PEM/x509 file")
	flag.StringVar(&clientCertKeyPath, "client_cert_key", "", "Path to client key PEM/x509 file")
	flag.Parse()

	if controllerAddress == "" {
		fmt.Println("Controller address not supplied")
		return
	}

	var log = log.Init()

	log.Infof("Starting Prometheus Exporter for Anka (%s)", version)

	clientTLSCerts := client.TLSCerts{
		UseTLS:              useTLS,
		ClientCert:          clientCertPath,
		ClientCertKey:       clientCertKeyPath,
		CACert:              caFilePath,
		SkipTLSVerification: skipTLSVerification,
	}

	client, err := client.NewClient(controllerAddress, intervalSeconds, clientTLSCerts)
	if err != nil {
		fmt.Println(err)
		return
	}
	client.Init()

	prometheusRegistry := prometheus.NewRegistry()

	// Create each metric that we later populate
	for _, m := range metrics.MetricsHolder {
		prometheusRegistry.Register(m.GetPrometheusMetric())
		client.Register(m.GetEvent(), m.GetEventHandler())
	}

	srv := server.NewServer(prometheusRegistry, port)
	if !disableOptimizeInterval {
		srv.SetIntervalUpdateFunc(client.UpdateInterval)
	}
	srv.Init()
}
