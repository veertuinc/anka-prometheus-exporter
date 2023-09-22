package server

import (
	"fmt"
	"math"
	"net/http"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/veertuinc/anka-prometheus-exporter/src/log"
)

type Server struct {
	lastInterval       int64
	lastRequestTime    int64
	registry           *prometheus.Registry
	intervalChangeFunc func(i int64)
	lock               *sync.Mutex
	port               int
}

func NewServer(promReg *prometheus.Registry, port int) *Server {
	return &Server{
		lastInterval:       0,
		lastRequestTime:    time.Now().Unix(),
		registry:           promReg,
		intervalChangeFunc: nil,
		lock:               &sync.Mutex{},
		port:               port,
	}
}

func (server *Server) Init() {
	var log = log.GetLogger()
	log.Info(fmt.Sprintf("Serving metrics at /metrics and :%d", server.port))
	http.HandleFunc("/metrics", server.handleRequest())
	http.ListenAndServe(fmt.Sprintf(":%d", server.port), nil)
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
