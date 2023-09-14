package state

import (
	"sync"

	"github.com/veertuinc/anka-prometheus-exporter/src/types"
)

var lock = &sync.Mutex{}

type State struct {
	TemplatesMap map[string]types.Template
}

var once sync.Once

var state *State

func GetState() *State {
	once.Do(func() {
		state = &State{
			TemplatesMap: make(map[string]types.Template),
		}
	})
	return state
}

func (state *State) GetTemplatesMap() map[string]types.Template {
	lock.Lock()
	defer lock.Unlock()
	return state.TemplatesMap
}

func (state *State) SetTemplatesMap(templates []types.Template) {
	lock.Lock()
	defer lock.Unlock()
	for _, templateV := range templates {
		state.TemplatesMap[templateV.UUID] = templateV
	}
}
