package server

import (
	"fmt"
	"math"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/exporter-toolkit/web"
	"github.com/veertuinc/anka-prometheus-exporter/src/log"
)

type Server struct {
	lastInterval       int64
	lastRequestTime    int64
	registry           *prometheus.Registry
	intervalChangeFunc func(i int64)
	lock               *sync.Mutex
	webListenAddress   string
	version            string
	configFile         string
}

func NewServer(
	promReg *prometheus.Registry,
	webListenAddress string,
	version string,
	configFile string,
) *Server {
	return &Server{
		lastInterval:       0,
		lastRequestTime:    time.Now().Unix(),
		registry:           promReg,
		intervalChangeFunc: nil,
		lock:               &sync.Mutex{},
		webListenAddress:   webListenAddress,
		version:            version,
		configFile:         configFile,
	}
}

func (server *Server) Init() {
	log.Info(fmt.Sprintf("Serving metrics at %s/metrics", server.webListenAddress))

	http.HandleFunc("/metrics", server.handleRequest())

	landingConfig := web.LandingConfig{
		HeaderColor: "#7e57c2",
		Name:        "Anka Prometheus Exporter",
		Description: "Prometheus Exporter for Anka Build Cloud Controller",
		Version:     server.version,
		Links: []web.LandingLinks{
			{
				Address: "/metrics",
				Text:    "Metrics",
			},
		},
	}
	landingPage, err := web.NewLandingPage(landingConfig)
	if err != nil {
		log.Error(fmt.Sprintf("Error creating landing page: %s", err.Error()))
		os.Exit(1)
	}
	http.Handle("/", landingPage)

	if err := web.ListenAndServe(&http.Server{}, &web.FlagConfig{
		WebListenAddresses: &[]string{server.webListenAddress},
		WebConfigFile:      &server.configFile,
	}, log.Logger); err != nil {
		log.Error(fmt.Sprintf("Error starting web server: %s", err.Error()))
		os.Exit(1)
	}
}

func (server *Server) handleRequest() func(http.ResponseWriter, *http.Request) {
	handler := promhttp.HandlerFor(server.registry, promhttp.HandlerOpts{})

	return func(w http.ResponseWriter, r *http.Request) {
		if server.intervalChangeFunc != nil {
			go server.handleInterval()
		}
		handler.ServeHTTP(w, r)
	}
}

func (server *Server) handleInterval() {
	server.lock.Lock()
	defer server.lock.Unlock()

	timeStamp := time.Now().Unix()
	interval := timeStamp - server.lastRequestTime

	if math.Abs(float64(interval-server.lastInterval)) > 1 {
		server.intervalChangeFunc(interval)
		server.lastInterval = interval
	}
	server.lastRequestTime = timeStamp
}

func (server *Server) SetIntervalUpdateFunc(f func(i int64)) {
	if f != nil {
		server.intervalChangeFunc = f
	}
}
