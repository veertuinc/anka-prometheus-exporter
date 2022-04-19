package events

type Event int

const (
	EventNodeUpdated              = 1
	EventRegistryDiskDataUpdated  = 2
	EventVmDataUpdated            = 3
	EventRegistryTemplatesUpdated = 4
)
