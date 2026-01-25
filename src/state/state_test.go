package state

import (
	"sync"
	"testing"

	"github.com/veertuinc/anka-prometheus-exporter/src/types"
)

func TestGetState_Singleton(t *testing.T) {
	// Get state multiple times and verify it's the same instance
	state1 := GetState()
	state2 := GetState()

	if state1 != state2 {
		t.Error("GetState() should return the same instance (singleton pattern)")
	}
}

func TestGetState_NotNil(t *testing.T) {
	state := GetState()
	if state == nil {
		t.Error("GetState() should not return nil")
	}
}

func TestState_GetTemplatesMap_InitiallyEmpty(t *testing.T) {
	state := GetState()
	// Clear any existing data for clean test
	state.TemplatesMap = make(map[string]types.Template)

	templates := state.GetTemplatesMap()
	if templates == nil {
		t.Error("GetTemplatesMap() should not return nil")
	}
	if len(templates) != 0 {
		t.Errorf("Expected empty templates map, got %d entries", len(templates))
	}
}

func TestState_SetTemplatesMap(t *testing.T) {
	state := GetState()
	// Clear any existing data
	state.TemplatesMap = make(map[string]types.Template)

	templates := []types.Template{
		{
			UUID: "uuid-1",
			Name: "macos-ventura",
			Size: 50000000,
			Tags: []types.TemplateTag{
				{Name: "v1.0", Size: 10000000},
			},
		},
		{
			UUID: "uuid-2",
			Name: "macos-sonoma",
			Size: 60000000,
			Tags: []types.TemplateTag{
				{Name: "latest", Size: 12000000},
			},
		},
	}

	state.SetTemplatesMap(templates)

	result := state.GetTemplatesMap()
	if len(result) != 2 {
		t.Fatalf("Expected 2 templates, got %d", len(result))
	}

	if result["uuid-1"].Name != "macos-ventura" {
		t.Errorf("Expected template name macos-ventura, got %s", result["uuid-1"].Name)
	}
	if result["uuid-2"].Name != "macos-sonoma" {
		t.Errorf("Expected template name macos-sonoma, got %s", result["uuid-2"].Name)
	}
}

func TestState_SetTemplatesMap_UpdateExisting(t *testing.T) {
	state := GetState()
	// Clear any existing data
	state.TemplatesMap = make(map[string]types.Template)

	// Initial set
	templates := []types.Template{
		{UUID: "uuid-1", Name: "macos-ventura", Size: 50000000},
	}
	state.SetTemplatesMap(templates)

	// Update with new data
	updatedTemplates := []types.Template{
		{UUID: "uuid-1", Name: "macos-ventura-updated", Size: 55000000},
	}
	state.SetTemplatesMap(updatedTemplates)

	result := state.GetTemplatesMap()
	if result["uuid-1"].Name != "macos-ventura-updated" {
		t.Errorf("Expected updated name, got %s", result["uuid-1"].Name)
	}
	if result["uuid-1"].Size != 55000000 {
		t.Errorf("Expected updated size, got %d", result["uuid-1"].Size)
	}
}

func TestState_ConcurrentAccess(t *testing.T) {
	state := GetState()
	// Clear any existing data
	state.TemplatesMap = make(map[string]types.Template)

	var wg sync.WaitGroup
	numGoroutines := 100

	// Concurrent writes
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			templates := []types.Template{
				{UUID: "uuid-concurrent", Name: "test-template", Size: uint(id)},
			}
			state.SetTemplatesMap(templates)
		}(i)
	}

	// Concurrent reads
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_ = state.GetTemplatesMap()
		}()
	}

	wg.Wait()

	// Verify state is still valid after concurrent access
	result := state.GetTemplatesMap()
	if result == nil {
		t.Error("TemplatesMap should not be nil after concurrent access")
	}
}

func TestState_SetTemplatesMap_EmptySlice(t *testing.T) {
	state := GetState()
	// Set some initial data
	state.TemplatesMap = make(map[string]types.Template)
	state.TemplatesMap["existing"] = types.Template{UUID: "existing", Name: "existing-template"}

	// Set empty slice - should not remove existing data
	state.SetTemplatesMap([]types.Template{})

	result := state.GetTemplatesMap()
	if len(result) != 1 {
		t.Errorf("Expected 1 template (existing should remain), got %d", len(result))
	}
}

func TestState_SetTemplatesMap_PreservesTags(t *testing.T) {
	state := GetState()
	state.TemplatesMap = make(map[string]types.Template)

	templates := []types.Template{
		{
			UUID: "uuid-with-tags",
			Name: "template-with-tags",
			Size: 100000,
			Tags: []types.TemplateTag{
				{Name: "v1.0", Size: 10000},
				{Name: "v2.0", Size: 20000},
				{Name: "latest", Size: 25000},
			},
		},
	}

	state.SetTemplatesMap(templates)

	result := state.GetTemplatesMap()
	if len(result["uuid-with-tags"].Tags) != 3 {
		t.Errorf("Expected 3 tags, got %d", len(result["uuid-with-tags"].Tags))
	}
	if result["uuid-with-tags"].Tags[0].Name != "v1.0" {
		t.Errorf("Expected first tag v1.0, got %s", result["uuid-with-tags"].Tags[0].Name)
	}
}
