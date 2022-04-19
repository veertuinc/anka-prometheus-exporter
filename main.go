package main

import (
	"flag"
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/veertuinc/anka-prometheus-exporter/envflag"
	"github.com/veertuinc/anka-prometheus-exporter/src/client"
	"github.com/veertuinc/anka-prometheus-exporter/src/log"
	"github.com/veertuinc/anka-prometheus-exporter/src/metrics"
	"github.com/veertuinc/anka-prometheus-exporter/src/server"
)

const (
	DefaultIntervalSeconds = 15
)

var (
	version string
)

func main() {
	var logger = log.GetLogger()

	var controllerAddress string
	var intervalSeconds int
	var disableOptimizeInterval bool
	var port int
	var caFilePath string
	var clientCertPath string
	var clientCertKeyPath string
	var skipTLSVerification bool
	var useTLS bool

	flag.StringVar(&controllerAddress, "controller-address", "", "Controller address to monitor (url as arg) (required)")
	flag.IntVar(&intervalSeconds, "interval", DefaultIntervalSeconds, "Seconds to wait between data requests to controller (int as arg)")
	flag.IntVar(&port, "port", 2112, "Port to server /metrics endpoint (int as arg)")
	flag.BoolVar(&disableOptimizeInterval, "disable-interval-optimizer", false, "Optimize interval according to /metric api requests receieved (no args)")
	flag.BoolVar(&useTLS, "tls", false, "Enable TLS (no args)")
	flag.BoolVar(&skipTLSVerification, "skip-tls-verification", false, "Skip TLS verification (no args)")
	flag.StringVar(&caFilePath, "ca-cert", "", "Path to ca PEM/x509 file (cert file path as arg)")
	flag.StringVar(&clientCertPath, "client-cert", "", "Path to client cert PEM/x509 file (cert file path as arg)")
	flag.StringVar(&clientCertKeyPath, "client-cert-key", "", "Path to client key PEM/x509 file (cert file path as arg)")

	envPrefix := "ANKA_PROMETHEUS_EXPORTER_"
	envflag.StringVar(&controllerAddress, "CONTROLLER_ADDRESS", "", "Controller address to monitor (url as arg) (required)")
	envflag.IntVar(&intervalSeconds, "INTERVAL", DefaultIntervalSeconds, "Seconds to wait between data requests to controller (int as arg)")
	envflag.IntVar(&port, "PORT", 2112, "Port to server /metrics endpoint (int as arg)")
	envflag.BoolVar(&disableOptimizeInterval, "DISABLE_INTERVAL_OPTIMIZER", false, "Optimize interval according to /metric api requests receieved (no args)")
	envflag.BoolVar(&useTLS, "TLS", false, "Enable TLS (no args)")
	envflag.BoolVar(&skipTLSVerification, "SKIP_TLS_VERIFICATION", false, "Skip TLS verification (no args)")
	envflag.StringVar(&caFilePath, "CA_CERT", "", "Path to ca PEM/x509 file (cert file path as arg)")
	envflag.StringVar(&clientCertPath, "CLIENT_CERT", "", "Path to client cert PEM/x509 file (cert file path as arg)")
	envflag.StringVar(&clientCertKeyPath, "CLIENT_CERT_KEY", "", "Path to client key PEM/x509 file (cert file path as arg)")
	flag.Parse()
	envflag.ParsePrefix(envPrefix)

	if controllerAddress == "" {
		logger.Fatalf(fmt.Errorf("controller address not supplied (%sCONTROLLER_ADDRESS=\"http://{address}:{port}\" or --controller-address http://{address}:{port})", envPrefix).Error())
	}

	if len(flag.Args()) > 0 {
		logger.Fatalf("one of your flags included a value when one wasn't needed: %s", flag.Args()[0])
	}

	logger.Infof("Starting Prometheus Exporter for Anka (%s)", version)

	clientTLSCerts := client.TLSCerts{
		UseTLS:              useTLS,
		ClientCert:          clientCertPath,
		ClientCertKey:       clientCertKeyPath,
		CACert:              caFilePath,
		SkipTLSVerification: skipTLSVerification,
	}

	client, err := client.NewClient(controllerAddress, intervalSeconds, clientTLSCerts)
	if err != nil {
		logger.Fatalf(err.Error())
	}
	client.Init()

	prometheusRegistry := prometheus.NewRegistry()

	// Create each metric that we later populate
	for _, m := range metrics.MetricsHolder {
		prometheusRegistry.Register(m.GetPrometheusMetric()) // TODO: error should be handled here
		client.Register(m.GetEvent(), m.GetEventHandler())   // TODO: error should be handled here
	}

	srv := server.NewServer(prometheusRegistry, port)
	if !disableOptimizeInterval {
		srv.SetIntervalUpdateFunc(client.UpdateInterval)
	}
	srv.Init()
}
