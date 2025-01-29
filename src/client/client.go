package client

import (
	"fmt"
	"io"
	"reflect"
	"runtime"
	"sync"
	"sync/atomic"
	"time"

	"github.com/veertuinc/anka-prometheus-exporter/src/events"
	"github.com/veertuinc/anka-prometheus-exporter/src/log"
)

const (
	MAX_INTERVAL_SECONDS = 60
)

type Client struct {
	events              map[events.Event][]func(interface{}) error
	communicator        *Communicator
	timeoutSeconds      int64
	errorTimeoutSeconds int
	eventsMutex         sync.Mutex
}

func NewClient(addr, username, password string, interval int, certs ClientTLSCerts, uak UAK) (*Client, error) {
	communicator, err := NewCommunicator(addr, username, password, certs, uak)
	if err != nil || communicator == nil {
		log.Error("Failed to create communicator")
		return nil, fmt.Errorf("failed to create communicator: %v", err)
	}
	c := &Client{
		events: map[events.Event][]func(interface{}) error{
			events.EVENT_NODE_UPDATED:               make([]func(interface{}) error, 0),
			events.EVENT_REGISTRY_DISK_DATA_UPDATED: make([]func(interface{}) error, 0),
			events.EVENT_VM_DATA_UPDATED:            make([]func(interface{}) error, 0),
			events.EVENT_REGISTRY_TEMPLATES_UPDATED: make([]func(interface{}) error, 0),
			events.EVENT_STATUS_UPDATED:             make([]func(interface{}) error, 0),
		},
		communicator:        communicator,
		timeoutSeconds:      int64(interval),
		errorTimeoutSeconds: 10,
	}
	if err := c.communicator.TestConnection(); err != nil {
		response, err := c.communicator.getResponse("/api/v1/status", "", "")
		if err != nil {
			log.Error(fmt.Sprintf("Error getting response: %s", err.Error()))
		} else {
			bodyBytes := make([]byte, 1024)
			for {
				n, err := response.Body.Read(bodyBytes)
				if n > 0 {
					log.Error(fmt.Sprintf("call to %s returned %d code and body of '%s'", response.Request.URL, response.StatusCode, string(bodyBytes[:n])))
					return nil, fmt.Errorf("failed to test connection")
				}
				if err != nil {
					if err != io.EOF {
						log.Error(fmt.Sprintf("Error reading response body: %s", err.Error()))
						return nil, fmt.Errorf("failed to test connection")
					}
					break
				}
			}
		}
		return nil, err
	}
	return c, nil
}

func (client *Client) Init() {
	// We must first populate the data from the Controller API that is going to be stored in state before we attempt to create metrics from it
	// Order matters here since GetVmsData for example relies on RegistryTemplatesData
	_, err := client.communicator.GetRegistryTemplatesData()
	if err != nil {
		log.Error(fmt.Sprintf("Error getting registry templates data: %v", err))
	}
	go client.initDataLoop(client.communicator.GetNodesData, events.EVENT_NODE_UPDATED)
	go client.initDataLoop(client.communicator.GetVmsData, events.EVENT_VM_DATA_UPDATED)
	go client.initDataLoop(client.communicator.GetRegistryDiskData, events.EVENT_REGISTRY_DISK_DATA_UPDATED)
	go client.initDataLoop(client.communicator.GetRegistryTemplatesData, events.EVENT_REGISTRY_TEMPLATES_UPDATED)
	go client.initDataLoop(client.communicator.GetStatus, events.EVENT_STATUS_UPDATED)
}

func (client *Client) Register(ev events.Event, eventHandler func(interface{}) error) error {
	client.eventsMutex.Lock()
	defer client.eventsMutex.Unlock()
	val, ok := client.events[ev]
	if ok {
		client.events[ev] = append(val, eventHandler)
	} else {
		return fmt.Errorf("no such event id: %v", ev)
	}
	return nil
}

func (client *Client) UpdateInterval(i int64) {
	if i > 1 {
		if i > MAX_INTERVAL_SECONDS {
			atomic.StoreInt64(&client.timeoutSeconds, MAX_INTERVAL_SECONDS)
		} else {
			atomic.StoreInt64(&client.timeoutSeconds, i-1)
		}
	}
}

// Loops over each eventHandler inside of the metrics/metric_*.go files and populates the values for each metric
func (client *Client) initDataLoop(f func() (interface{}, error), ev events.Event) {
	for {
		log.Debug("Requesting data for: " + runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name())
		data, err := f()
		if err != nil {
			log.Error(fmt.Sprintf("could not get data: %+v", err))
			time.Sleep(time.Duration(client.errorTimeoutSeconds) * time.Second)
			continue
		}
		client.eventsMutex.Lock()
		events := client.events[ev]
		client.eventsMutex.Unlock()
		for _, eventHandler := range events {
			if err := eventHandler(data); err != nil {
				log.Error(fmt.Sprintf("ignoring event handler failure for event id %+v - Error: %+v", ev, err))
			}
		}
		log.Debug("Finished requesting data for: " + runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name())
		time.Sleep(time.Duration(atomic.LoadInt64(&client.timeoutSeconds)) * time.Second)
	}
}
