package client

import (
	"fmt"
	"time"
	"github.com/veertuinc/anka-prometheus/events"
)

const (
	MAX_INTERVAL_SECONDS = 60
)

type Client struct {
	events map[events.Event][]func(interface{}) error
	communicator  *Communicator
	timeoutSeconds	int64
	errorTimeoutSeconds		int
}

func NewClient(addr string, interval int) (*Client, error) {
	c := &Client{
		events: map[events.Event][]func(interface{}) error {
			events.EVENT_NODE_UPDATED: make( []func(interface{}) error, 0),
			events.EVENT_REGISTRY_DATA_UPDATED: make( []func(interface{}) error, 0),
			events.EVENT_VM_DATA_UPDATED: make( []func(interface{}) error, 0),
		},
		communicator: NewCommunicator(addr),
		timeoutSeconds: int64(interval),
		errorTimeoutSeconds: 5,
	}
	if err := c.communicator.TestConnection(); err != nil {
		return nil, err
	}
	return c, nil
}

func (this *Client) Init() {
	go this.initDataLoop(this.communicator.GetNodesData, events.EVENT_NODE_UPDATED)
	go this.initDataLoop(this.communicator.GetVmsData, events.EVENT_VM_DATA_UPDATED)
	go this.initDataLoop(this.communicator.GetRegistryData, events.EVENT_REGISTRY_DATA_UPDATED)
}

func (this *Client) Register(ev events.Event, eventHandler func(interface{}) error) error {
	if val, ok := this.events[ev]; ok {
		this.events[ev] = append(val, eventHandler)
	} else {
		return fmt.Errorf("no such event id: ", ev)
	}
	return nil
}

func (this *Client) UpdateInterval(i int64) {
	if i > 1 {
		if i > MAX_INTERVAL_SECONDS {
			this.timeoutSeconds = MAX_INTERVAL_SECONDS
		} else {
			this.timeoutSeconds = i - 1
		}
	}
}

func (this *Client) initDataLoop(f func() (interface{}, error), ev events.Event) {
	for {
		data, err := f()
		if err != nil {
			fmt.Println("Could not get data. Error: ", err)
			time.Sleep(time.Duration(this.errorTimeoutSeconds) * time.Second)
			continue
		}
		
		this.executeEventHandlers(ev, data)
		time.Sleep(time.Duration(this.timeoutSeconds) * time.Second)
	}
}

func (this *Client) executeEventHandlers(ev events.Event, data interface{}) {
	for _, eventHandler := range this.events[ev] {
		if err := eventHandler(data); err != nil {
			fmt.Println("ignoring event handler failure for event id ", ev, "Error: ", err)
		}
	}
}