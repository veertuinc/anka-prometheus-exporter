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
	DEFAULT_INTERVAL_SECONDS = 15
)

var (
	version string
)

func main() {

	var controllerAddress string
	var controllerUsername string
	var controllerPassword string
	var intervalSeconds int
	var disableOptimizeInterval bool
	var port int
	var caFilePath string
	var clientCertPath string
	var clientCertKeyPath string
	var skipTLSVerification bool
	var useTLS bool
	var uakId string
	var uakPath string
	var uakString string

	flag.StringVar(&controllerAddress, "controller-address", "", "Controller address to monitor (url as arg) (required)")
	flag.StringVar(&controllerUsername, "controller-username", "", "Controller basic auth username (username as arg)")
	flag.StringVar(&controllerPassword, "controller-password", "", "Controller basic auth password (password as arg)")
	flag.IntVar(&intervalSeconds, "interval", DEFAULT_INTERVAL_SECONDS, "Seconds to wait between data requests to controller (int as arg)")
	flag.IntVar(&port, "port", 2112, "Port to server /metrics endpoint (int as arg)")
	flag.BoolVar(&disableOptimizeInterval, "disable-interval-optimizer", false, "Optimize interval according to /metric api requests received (no args)")
	flag.BoolVar(&useTLS, "tls", false, "Enable TLS (no args)")
	flag.BoolVar(&skipTLSVerification, "skip-tls-verification", false, "Skip TLS verification (no args)")
	flag.StringVar(&caFilePath, "ca-cert", "", "Path to ca PEM/x509 file (cert file path as arg)")
	flag.StringVar(&clientCertPath, "client-cert", "", "Path to client cert PEM/x509 file (cert file path as arg)")
	flag.StringVar(&clientCertKeyPath, "client-cert-key", "", "Path to client key PEM/x509 file (cert file path as arg)")
	flag.StringVar(&uakId, "uak-id", "", "ID for the UAK you wish to use to make Controller requests (string as arg)")
	flag.StringVar(&uakPath, "uak-path", "", "Path to the UAK you wish to use for Controller requests (path as arg) (takes priority over -uak-string if both are specified)")
	flag.StringVar(&uakString, "uak-string", "", "String form (cat myUAK.pem | sed '1,1d' | sed '$d' | tr -d '\\n') of the key file contents you wish to use to make Controller requests")

	envPrefix := "ANKA_PROMETHEUS_EXPORTER_"
	envflag.StringVar(&controllerAddress, "CONTROLLER_ADDRESS", "", "Controller address to monitor (url as arg) (required)")
	envflag.StringVar(&controllerUsername, "CONTROLLER_USERNAME", "", "Controller basic auth username (username as arg)")
	envflag.StringVar(&controllerPassword, "CONTROLLER_PASSWORD", "", "Controller basic auth password (password as arg)")
	envflag.IntVar(&intervalSeconds, "INTERVAL", DEFAULT_INTERVAL_SECONDS, "Seconds to wait between data requests to controller (int as arg)")
	envflag.IntVar(&port, "PORT", 2112, "Port to server /metrics endpoint (int as arg)")
	envflag.BoolVar(&disableOptimizeInterval, "DISABLE_INTERVAL_OPTIMIZER", false, "Optimize interval according to /metric api requests received (no args)")
	envflag.BoolVar(&useTLS, "TLS", false, "Enable TLS (no args)")
	envflag.BoolVar(&skipTLSVerification, "SKIP_TLS_VERIFICATION", false, "Skip TLS verification (no args)")
	envflag.StringVar(&caFilePath, "CA_CERT", "", "Path to ca PEM/x509 file (cert file path as arg)")
	envflag.StringVar(&clientCertPath, "CLIENT_CERT", "", "Path to client cert PEM/x509 file (cert file path as arg)")
	envflag.StringVar(&clientCertKeyPath, "CLIENT_CERT_KEY", "", "Path to client key PEM/x509 file (cert file path as arg)")
	envflag.StringVar(&uakId, "UAK_ID", "", "ID for the UAK you wish to use to make Controller requests (string as arg)")
	envflag.StringVar(&uakPath, "UAK_PATH", "", "Path to the UAK you wish to use for Controller requests (takes priority over -uak-string if both are specified)")
	envflag.StringVar(&uakString, "UAK_STRING", "", "String form (cat myUAK.pem | sed '1,1d' | sed '$d' | tr -d '\\n') of the key file contents you wish to use to make Controller requests")

	flag.Parse()
	envflag.ParsePrefix(envPrefix)

	if controllerAddress == "" {
		log.Fatal(fmt.Errorf("controller address not supplied (%sCONTROLLER_ADDRESS=\"http://{address}:{port}\" or --controller-address http://{address}:{port})", envPrefix))
	}

	if len(flag.Args()) > 0 {
		log.Fatal(fmt.Errorf("one of your flags included a value when one wasn't needed: %s", flag.Args()[0]))
	}

	log.Info(fmt.Sprintf("Starting Prometheus Exporter for Anka (%s)", version))

	clientTLSCerts := client.TLSCerts{
		UseTLS:              useTLS,
		ClientCert:          clientCertPath,
		ClientCertKey:       clientCertKeyPath,
		CACert:              caFilePath,
		SkipTLSVerification: skipTLSVerification,
	}

	clientUAK := client.UAK{
		ID:        uakId,
		KeyPath:   uakPath,
		KeyString: uakString,
	}

	client, err := client.NewClient(controllerAddress, controllerUsername, controllerPassword, intervalSeconds, clientTLSCerts, clientUAK)
	if err != nil {
		log.Fatal(err)
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
