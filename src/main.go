package main

import (
	"fmt"
	"flag"
	"github.com/veertuinc/anka-prometheus/client"
	"github.com/veertuinc/anka-prometheus/server"
	"github.com/veertuinc/anka-prometheus/metrics"
	"github.com/prometheus/client_golang/prometheus"
)

const (
	DEFAULT_INTERVAL_SECONDS = 15
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

	clientTLSCerts := client.TLSCerts {
		UseTLS: useTLS,
		ClientCert: clientCertPath,
		ClientCertKey: clientCertKeyPath,
		CACert: caFilePath,
		SkipTLSVerification:skipTLSVerification,
	}

	client, err := client.NewClient(controllerAddress, intervalSeconds, clientTLSCerts)
	if err != nil {
		fmt.Println(err)
		return
	}
	client.Init()

	prometheusRegistry := prometheus.NewRegistry()

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