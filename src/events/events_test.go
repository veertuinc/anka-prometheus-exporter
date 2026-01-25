package events

import (
	"testing"
)

func TestEventConstants(t *testing.T) {
	// Verify event constants have expected values
	tests := []struct {
		name     string
		event    Event
		expected int
	}{
		{"EVENT_NODE_UPDATED", EVENT_NODE_UPDATED, 1},
		{"EVENT_REGISTRY_DISK_DATA_UPDATED", EVENT_REGISTRY_DISK_DATA_UPDATED, 2},
		{"EVENT_VM_DATA_UPDATED", EVENT_VM_DATA_UPDATED, 3},
		{"EVENT_REGISTRY_TEMPLATES_UPDATED", EVENT_REGISTRY_TEMPLATES_UPDATED, 4},
		{"EVENT_STATUS_UPDATED", EVENT_STATUS_UPDATED, 5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if int(tt.event) != tt.expected {
				t.Errorf("Expected %s = %d, got %d", tt.name, tt.expected, int(tt.event))
			}
		})
	}
}

func TestEventType(t *testing.T) {
	// Verify Event type can be used as expected
	var e Event = EVENT_NODE_UPDATED

	// Test that Event can be compared
	if e != EVENT_NODE_UPDATED {
		t.Error("Event comparison failed")
	}

	// Test that different events are not equal
	if e == EVENT_VM_DATA_UPDATED {
		t.Error("Different events should not be equal")
	}
}

func TestEventUniqueness(t *testing.T) {
	// Verify all events have unique values
	events := []Event{
		EVENT_NODE_UPDATED,
		EVENT_REGISTRY_DISK_DATA_UPDATED,
		EVENT_VM_DATA_UPDATED,
		EVENT_REGISTRY_TEMPLATES_UPDATED,
		EVENT_STATUS_UPDATED,
	}

	seen := make(map[Event]bool)
	for _, e := range events {
		if seen[e] {
			t.Errorf("Duplicate event value found: %d", e)
		}
		seen[e] = true
	}
}

func TestEventAsMapKey(t *testing.T) {
	// Test that Event can be used as a map key (common use case)
	eventMap := make(map[Event]string)

	eventMap[EVENT_NODE_UPDATED] = "node_updated"
	eventMap[EVENT_VM_DATA_UPDATED] = "vm_data_updated"

	if eventMap[EVENT_NODE_UPDATED] != "node_updated" {
		t.Error("Failed to use Event as map key")
	}

	if len(eventMap) != 2 {
		t.Errorf("Expected map length 2, got %d", len(eventMap))
	}
}
