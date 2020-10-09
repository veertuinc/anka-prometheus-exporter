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

func (this *Server) Init() {
	var log = log.Init()
	log.Info(fmt.Sprintf("Serving metrics at /metrics and :%d", this.port))
	http.HandleFunc("/metrics", this.handleRequest())
	http.ListenAndServe(fmt.Sprintf(":%d", this.port), nil)
}

func (this *Server) handleRequest() func(http.ResponseWriter, *http.Request) {
	handler := promhttp.HandlerFor(this.registry, promhttp.HandlerOpts{})

	return func(w http.ResponseWriter, r *http.Request) {
		if this.intervalChangeFunc != nil {
			go this.handleInterval()
		}
		handler.ServeHTTP(w, r)
	}
}

func (this *Server) handleInterval() {
	this.lock.Lock()
	defer this.lock.Unlock()

	timeStamp := time.Now().Unix()
	interval := timeStamp - this.lastRequestTime

	if math.Abs(float64(interval-this.lastInterval)) > 1 {
		this.intervalChangeFunc(interval)
		this.lastInterval = interval
	}
	this.lastRequestTime = timeStamp
}

func (this *Server) SetIntervalUpdateFunc(f func(i int64)) {
	if f != nil {
		this.intervalChangeFunc = f
	}
}
