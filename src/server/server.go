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

func (s *Server) Init() {
	var log = log.GetLogger()
	log.Info(fmt.Sprintf("Serving metrics at /metrics and :%d", s.port))
	http.HandleFunc("/metrics", s.handleRequest())
	http.ListenAndServe(fmt.Sprintf(":%d", s.port), nil)
}

func (s *Server) handleRequest() func(http.ResponseWriter, *http.Request) {
	handler := promhttp.HandlerFor(s.registry, promhttp.HandlerOpts{})

	return func(w http.ResponseWriter, r *http.Request) {
		if s.intervalChangeFunc != nil {
			go s.handleInterval()
		}
		handler.ServeHTTP(w, r)
	}
}

func (s *Server) handleInterval() {
	s.lock.Lock()
	defer s.lock.Unlock()

	timeStamp := time.Now().Unix()
	interval := timeStamp - s.lastRequestTime

	if math.Abs(float64(interval-s.lastInterval)) > 1 {
		s.intervalChangeFunc(interval)
		s.lastInterval = interval
	}
	s.lastRequestTime = timeStamp
}

func (s *Server) SetIntervalUpdateFunc(f func(i int64)) {
	if f != nil {
		s.intervalChangeFunc = f
	}
}
