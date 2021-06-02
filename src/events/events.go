package events

type Event int

const (
	EVENT_NODE_UPDATED               = 1
	EVENT_REGISTRY_DISK_DATA_UPDATED = 2
	EVENT_VM_DATA_UPDATED            = 3
	EVENT_REGISTRY_TEMPLATES_UPDATED = 4
)
