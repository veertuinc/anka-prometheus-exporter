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

	var log = log.GetLogger()

	var controllerAddress string
	var intervalSeconds int
	var disableOptimizeInterval bool
	var port int
	var caFilePath string
	var clientCertPath string
	var clientCertKeyPath string
	var skipTLSVerification bool
	var useTLS bool

	flag.StringVar(&controllerAddress, "controller_address", "", "Controller address to monitor (url as arg)")
	flag.IntVar(&intervalSeconds, "interval", DEFAULT_INTERVAL_SECONDS, "Seconds to wait between data requests to controller (int as arg)")
	flag.IntVar(&port, "port", 2112, "Port to server /metrics endpoint (int as arg)")
	flag.BoolVar(&disableOptimizeInterval, "disable_interval_optimizer", false, "Optimize interval according to /metric api requests receieved (no args)")
	flag.BoolVar(&useTLS, "tls", false, "Enable TLS (no args)")
	flag.BoolVar(&skipTLSVerification, "skip_tls_verification", false, "Skip TLS verification (no args)")
	flag.StringVar(&caFilePath, "ca_cert", "", "Path to ca PEM/x509 file (cert file path as arg)")
	flag.StringVar(&clientCertPath, "client_cert", "", "Path to client cert PEM/x509 file (cert file path as arg)")
	flag.StringVar(&clientCertKeyPath, "client_cert_key", "", "Path to client key PEM/x509 file (cert file path as arg)")
	flag.Parse()

	if len(flag.Args()) > 0 {
		log.Fatalln("One of your flags included a value when one wasn't needed. The value we found:", flag.Args()[0])
	}

	if controllerAddress == "" {
		fmt.Println("Controller address not supplied")
		return
	}

	// fmt.Printf("%+v\n", flag.Lookup("controller_address"))
	// fmt.Printf("%+v\n", flag.Lookup("interval"))
	// fmt.Printf("%+v\n", flag.Lookup("port"))
	// fmt.Printf("%+v\n", flag.Lookup("disable_interval_optimizer"))
	// fmt.Printf("%+v\n", flag.Lookup("tls"))
	// fmt.Printf("%+v\n", flag.Lookup("skip_tls_verification"))
	// fmt.Printf("%+v\n", flag.Lookup("ca_cert"))
	// fmt.Printf("%+v\n", flag.Lookup("client_cert"))
	// fmt.Printf("%+v\n", flag.Lookup("client_cert_key"))

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
