package client

import (
	"fmt"
	"reflect"
	"runtime"
	"sync"
	"sync/atomic"
	"time"

	"github.com/veertuinc/anka-prometheus-exporter/src/events"
	"github.com/veertuinc/anka-prometheus-exporter/src/log"
)

const (
	MaxIntervalSeconds = 60
)

type Client struct {
	events              map[events.Event][]func(interface{}) error
	communicator        *Communicator
	timeoutSeconds      int64
	errorTimeoutSeconds int
	eventsMutex         sync.Mutex
}

func NewClient(addr string, interval int, certs TLSCerts) (*Client, error) {
	var logger = log.GetLogger()

	communicator, err := NewCommunicator(addr, certs)
	if err != nil {
		return nil, err
	}

	c := &Client{
		events: map[events.Event][]func(interface{}) error{
			events.EventNodeUpdated:              make([]func(interface{}) error, 0),
			events.EventRegistryDiskDataUpdated:  make([]func(interface{}) error, 0),
			events.EventVmDataUpdated:            make([]func(interface{}) error, 0),
			events.EventRegistryTemplatesUpdated: make([]func(interface{}) error, 0),
		},
		communicator:        communicator,
		timeoutSeconds:      int64(interval),
		errorTimeoutSeconds: 10,
	}
	if err := c.communicator.TestConnection(); err != nil {
		logger.Fatal(err)
		return nil, err
	}
	return c, nil
}

func (c *Client) Init() {
	// We must first populate the data from the Controller API that is going to be stored in state before we attempt to create metrics from it
	// Order matters here since GetVmsData for example relies on RegistryTemplatesData
	c.communicator.GetRegistryTemplatesData()
	go c.initDataLoop(c.communicator.GetNodesData, events.EventNodeUpdated)
	go c.initDataLoop(c.communicator.GetVmsData, events.EventVmDataUpdated)
	go c.initDataLoop(c.communicator.GetRegistryDiskData, events.EventRegistryDiskDataUpdated)
	go c.initDataLoop(c.communicator.GetRegistryTemplatesData, events.EventRegistryTemplatesUpdated)
}

func (c *Client) Register(ev events.Event, eventHandler func(interface{}) error) error {
	c.eventsMutex.Lock()
	defer c.eventsMutex.Unlock()
	val, ok := c.events[ev]
	if ok {
		c.events[ev] = append(val, eventHandler)
	} else {
		return fmt.Errorf("no such event id: %v", ev)
	}
	return nil
}

func (c *Client) UpdateInterval(i int64) {
	if i > 1 {
		if i > MaxIntervalSeconds {
			atomic.StoreInt64(&c.timeoutSeconds, MaxIntervalSeconds)
		} else {
			atomic.StoreInt64(&c.timeoutSeconds, i-1)
		}
	}
}

// Loops over each eventHandler inside the metrics/metric_*.go files and populates the values for each metric.
func (c *Client) initDataLoop(f func() (interface{}, error), ev events.Event) {
	var log = log.GetLogger()
	for {
		if log.GetLevel().String() == "debug" {
			log.Debugln("Requesting data for: " + runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name())
		}
		data, err := f()
		if err != nil {
			log.Errorf("could not get data: %+v", err)
			time.Sleep(time.Duration(c.errorTimeoutSeconds) * time.Second)
			continue
		}
		c.eventsMutex.Lock()
		events := c.events[ev]
		c.eventsMutex.Unlock()
		for _, eventHandler := range events {
			if err := eventHandler(data); err != nil {
				log.Errorf("ignoring event handler failure for event id %+v - Error: %+v", ev, err)
			}
		}
		time.Sleep(time.Duration(atomic.LoadInt64(&c.timeoutSeconds)) * time.Second)
	}
}
