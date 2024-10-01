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
	var clientCaFilePath string
	var clientCertPath string
	var clientCertKeyPath string
	var clientSkipTLSVerification bool
	var useClientTLS bool
	var uakId string
	var uakPath string
	var uakString string

	var webListenAddresses string
	flag.StringVar(&webListenAddresses, "web.listen-address", "", "Address on which to expose metrics and web interface. Examples: `:2112` or `[::1]:2112` for http, `vsock://:2112` for vsock")
	envflag.StringVar(&webListenAddresses, "WEB_LISTEN_ADDRESS", "", "Address on which to expose metrics and web interface. Examples: `:2112` or `[::1]:2112` for http, `vsock://:2112` for vsock")

	var webConfigFile string
	flag.StringVar(&webConfigFile, "web.config.file", "", "Path to configuration file that can enable server TLS or authentication. See: https://github.com/prometheus/exporter-toolkit/blob/master/docs/web-configuration.md")
	envflag.StringVar(&webConfigFile, "WEB_CONFIG_FILE", "", "Path to configuration file that can enable server TLS or authentication. See: https://github.com/prometheus/exporter-toolkit/blob/master/docs/web-configuration.md")

	flag.StringVar(&controllerAddress, "controller-address", "", "Controller address to monitor (url as arg) (required)")
	flag.StringVar(&controllerUsername, "controller-username", "", "Controller basic auth username (username as arg)")
	flag.StringVar(&controllerPassword, "controller-password", "", "Controller basic auth password (password as arg)")
	flag.IntVar(&intervalSeconds, "interval", DEFAULT_INTERVAL_SECONDS, "Seconds to wait between data requests to controller (int as arg)")
	// flag.IntVar(&port, "port", 2112, "Port to server /metrics endpoint (int as arg)")
	flag.BoolVar(&disableOptimizeInterval, "disable-interval-optimizer", false, "Optimize interval according to /metric api requests received (no args)")
	flag.BoolVar(&useClientTLS, "client-tls", false, "Enable client TLS (no args)")
	flag.BoolVar(&clientSkipTLSVerification, "client-skip-tls-verification", false, "Skip client TLS verification (no args)")
	flag.StringVar(&clientCaFilePath, "client-ca-cert", "", "Path to client CA PEM/x509 file (cert file path as arg)")
	flag.StringVar(&clientCertPath, "client-cert", "", "Path to client cert PEM/x509 file (cert file path as arg)")
	flag.StringVar(&clientCertKeyPath, "client-cert-key", "", "Path to client key PEM/x509 file (cert file path as arg)")
	flag.StringVar(&uakId, "uak-id", "", "UAK ID you wish to use for Controller requests (string as arg)")
	flag.StringVar(&uakPath, "uak-path", "", "Path to the UAK file used for Controller requests (path as arg) (supersedes -uak-string)")
	flag.StringVar(&uakString, "uak-string", "", "String form (cat myUAK.pem | sed '1,1d' | sed '$d' | tr -d '\\n') of the key file contents for Controller requests (string as arg)")

	envPrefix := "ANKA_PROMETHEUS_EXPORTER_"
	envflag.StringVar(&controllerAddress, "CONTROLLER_ADDRESS", "", "Controller address to monitor (url as arg) (required)")
	envflag.StringVar(&controllerUsername, "CONTROLLER_USERNAME", "", "Controller basic auth username (username as arg)")
	envflag.StringVar(&controllerPassword, "CONTROLLER_PASSWORD", "", "Controller basic auth password (password as arg)")
	envflag.IntVar(&intervalSeconds, "INTERVAL", DEFAULT_INTERVAL_SECONDS, "Seconds to wait between data requests to controller (int as arg)")
	// envflag.IntVar(&port, "PORT", 2112, "Port to server /metrics endpoint (int as arg)")
	envflag.BoolVar(&disableOptimizeInterval, "DISABLE_INTERVAL_OPTIMIZER", false, "Optimize interval according to /metric api requests received (no args)")
	envflag.BoolVar(&useClientTLS, "CLIENT_TLS", false, "Enable client TLS (no args)")
	envflag.BoolVar(&clientSkipTLSVerification, "CLIENT_SKIP_TLS_VERIFICATION", false, "Skip client TLS verification (no args)")
	envflag.StringVar(&clientCaFilePath, "CLIENT_CA_CERT", "", "Path to client CA PEM/x509 file (cert file path as arg)")
	envflag.StringVar(&clientCertPath, "CLIENT_CERT", "", "Path to client cert PEM/x509 file (cert file path as arg)")
	envflag.StringVar(&clientCertKeyPath, "CLIENT_CERT_KEY", "", "Path to client key PEM/x509 file (cert file path as arg)")
	envflag.StringVar(&uakId, "UAK_ID", "", "UAK ID you wish to use for Controller requests (string as arg)")
	envflag.StringVar(&uakPath, "UAK_PATH", "", "Path to the UAK file used for Controller requests (path as arg) (supersedes -uak-string)")
	envflag.StringVar(&uakString, "UAK_STRING", "", "String form (cat myUAK.pem | sed '1,1d' | sed '$d' | tr -d '\\n') of the key file contents for Controller requests (string as arg)")

	flag.Parse()
	envflag.ParsePrefix(envPrefix)

	if webListenAddresses == "" {
		webListenAddresses = ":2112"
	}

	if controllerAddress == "" {
		log.Fatal(fmt.Sprintf("controller address not supplied (%sCONTROLLER_ADDRESS=\"http://{address}:{port}\" or --controller-address http://{address}:{port})", envPrefix))
	}

	if len(flag.Args()) > 0 {
		log.Fatal(fmt.Sprintf("one of your flags included a value when one wasn't needed: %s", flag.Args()[0]))
	}

	log.Info(fmt.Sprintf("Starting Prometheus Exporter for Anka (%s)", version))

	clientTLSCerts := client.ClientTLSCerts{
		UseTLS:              useClientTLS,
		Cert:                clientCertPath,
		CertKey:             clientCertKeyPath,
		CACert:              clientCaFilePath,
		SkipTLSVerification: clientSkipTLSVerification,
	}

	clientUAK := client.UAK{
		ID:        uakId,
		KeyPath:   uakPath,
		KeyString: uakString,
	}

	client, err := client.NewClient(controllerAddress, controllerUsername, controllerPassword, intervalSeconds, clientTLSCerts, clientUAK)
	if err != nil {
		log.Fatal(fmt.Sprintf("Error creating client: %s", err.Error()))
	}
	client.Init()

	prometheusRegistry := prometheus.NewRegistry()

	// Create each metric that we later populate
	for _, m := range metrics.MetricsHolder {
		prometheusRegistry.Register(m.GetPrometheusMetric())
		client.Register(m.GetEvent(), m.GetEventHandler())
	}

	srv := server.NewServer(
		prometheusRegistry,
		webListenAddresses,
		version,
		webConfigFile,
	)
	if !disableOptimizeInterval {
		srv.SetIntervalUpdateFunc(client.UpdateInterval)
	}
	srv.Init()
}
