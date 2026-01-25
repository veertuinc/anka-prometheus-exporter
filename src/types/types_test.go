package types

import (
	"encoding/json"
	"testing"
)

func TestDefaultResponse_GetStatus(t *testing.T) {
	tests := []struct {
		name     string
		response DefaultResponse
		expected string
	}{
		{
			name:     "OK status",
			response: DefaultResponse{Status: "OK", Message: ""},
			expected: "OK",
		},
		{
			name:     "error status",
			response: DefaultResponse{Status: "ERROR", Message: "something went wrong"},
			expected: "ERROR",
		},
		{
			name:     "empty status",
			response: DefaultResponse{Status: "", Message: ""},
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.response.GetStatus(); got != tt.expected {
				t.Errorf("GetStatus() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestDefaultResponse_GetMessage(t *testing.T) {
	tests := []struct {
		name     string
		response DefaultResponse
		expected string
	}{
		{
			name:     "with message",
			response: DefaultResponse{Status: "ERROR", Message: "Authentication Required"},
			expected: "Authentication Required",
		},
		{
			name:     "empty message",
			response: DefaultResponse{Status: "OK", Message: ""},
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.response.GetMessage(); got != tt.expected {
				t.Errorf("GetMessage() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestStatusResponse_GetBody(t *testing.T) {
	sr := &StatusResponse{
		DefaultResponse: DefaultResponse{Status: "OK", Message: ""},
		Body: Status{
			Status:          "Running",
			Version:         "1.0.0",
			RegistryAddress: "http://registry:8089",
			RegistryStatus:  "Running",
			License:         "valid",
		},
	}

	body := sr.GetBody()
	status, ok := body.(Status)
	if !ok {
		t.Fatal("GetBody() did not return Status type")
	}
	if status.Version != "1.0.0" {
		t.Errorf("Expected version 1.0.0, got %s", status.Version)
	}
	if status.Status != "Running" {
		t.Errorf("Expected status Running, got %s", status.Status)
	}
}

func TestNodesResponse_GetBody(t *testing.T) {
	nr := &NodesResponse{
		DefaultResponse: DefaultResponse{Status: "OK", Message: ""},
		Body: []Node{
			{
				NodeID:   "node-1",
				NodeName: "test-node",
				CPU:      8,
				RAM:      32,
				VMCount:  2,
				State:    "Active",
				Capacity: 4,
				HostArch: "arm64",
			},
		},
	}

	body := nr.GetBody()
	nodes, ok := body.([]Node)
	if !ok {
		t.Fatal("GetBody() did not return []Node type")
	}
	if len(nodes) != 1 {
		t.Fatalf("Expected 1 node, got %d", len(nodes))
	}
	if nodes[0].NodeID != "node-1" {
		t.Errorf("Expected node ID node-1, got %s", nodes[0].NodeID)
	}
	if nodes[0].HostArch != "arm64" {
		t.Errorf("Expected arch arm64, got %s", nodes[0].HostArch)
	}
}

func TestRegistryDiskResponse_GetBody(t *testing.T) {
	rdr := &RegistryDiskResponse{
		DefaultResponse: DefaultResponse{Status: "OK", Message: ""},
		Body: RegistryDisk{
			Total: 1000000000,
			Free:  500000000,
		},
	}

	body := rdr.GetBody()
	disk, ok := body.(RegistryDisk)
	if !ok {
		t.Fatal("GetBody() did not return RegistryDisk type")
	}
	if disk.Total != 1000000000 {
		t.Errorf("Expected total 1000000000, got %d", disk.Total)
	}
	if disk.Free != 500000000 {
		t.Errorf("Expected free 500000000, got %d", disk.Free)
	}
}

func TestRegistryTemplateResponse_GetBody(t *testing.T) {
	rtr := &RegistryTemplateResponse{
		DefaultResponse: DefaultResponse{Status: "OK", Message: ""},
		Body: []Template{
			{
				UUID: "template-uuid-1",
				Name: "macos-ventura",
				Size: 50000000,
			},
			{
				UUID: "template-uuid-2",
				Name: "macos-sonoma",
				Size: 60000000,
			},
		},
	}

	body := rtr.GetBody()
	templates, ok := body.([]Template)
	if !ok {
		t.Fatal("GetBody() did not return []Template type")
	}
	if len(templates) != 2 {
		t.Fatalf("Expected 2 templates, got %d", len(templates))
	}
	if templates[0].Name != "macos-ventura" {
		t.Errorf("Expected template name macos-ventura, got %s", templates[0].Name)
	}
}

func TestRegistryTemplateTagsResponse_GetBody(t *testing.T) {
	rttr := &RegistryTemplateTagsResponse{
		DefaultResponse: DefaultResponse{Status: "OK", Message: ""},
		Body: RegistryTemplateTags{
			Versions: []TemplateTag{
				{Name: "v1.0", Size: 10000000},
				{Name: "v2.0", Size: 12000000},
			},
		},
	}

	body := rttr.GetBody()
	tags, ok := body.(RegistryTemplateTags)
	if !ok {
		t.Fatal("GetBody() did not return RegistryTemplateTags type")
	}
	if len(tags.Versions) != 2 {
		t.Fatalf("Expected 2 tags, got %d", len(tags.Versions))
	}
	if tags.Versions[0].Name != "v1.0" {
		t.Errorf("Expected tag name v1.0, got %s", tags.Versions[0].Name)
	}
}

func TestInstancesResponse_GetBody(t *testing.T) {
	ir := &InstancesResponse{
		DefaultResponse: DefaultResponse{Status: "OK", Message: ""},
		Body: []Instance{
			{
				InstanceID: "instance-1",
				Vm: VmData{
					State:        "Started",
					TemplateUUID: "template-uuid-1",
					TemplateName: "macos-ventura",
					Arch:         "arm64",
				},
			},
		},
	}

	body := ir.GetBody()
	instances, ok := body.([]Instance)
	if !ok {
		t.Fatal("GetBody() did not return []Instance type")
	}
	if len(instances) != 1 {
		t.Fatalf("Expected 1 instance, got %d", len(instances))
	}
	if instances[0].Vm.State != "Started" {
		t.Errorf("Expected state Started, got %s", instances[0].Vm.State)
	}
}

func TestStatusResponse_JSONUnmarshal(t *testing.T) {
	jsonData := `{
		"status": "OK",
		"message": "",
		"body": {
			"status": "Running",
			"version": "1.25.0",
			"registry_address": "http://registry:8089",
			"registry_status": "Running",
			"license": "enterprise"
		}
	}`

	var sr StatusResponse
	if err := json.Unmarshal([]byte(jsonData), &sr); err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	if sr.GetStatus() != "OK" {
		t.Errorf("Expected status OK, got %s", sr.GetStatus())
	}
	if sr.Body.Version != "1.25.0" {
		t.Errorf("Expected version 1.25.0, got %s", sr.Body.Version)
	}
}

func TestNodesResponse_JSONUnmarshal(t *testing.T) {
	jsonData := `{
		"status": "OK",
		"message": "",
		"body": [
			{
				"node_id": "abc123",
				"node_name": "mac-mini-1",
				"cpu_count": 8,
				"ram": 16,
				"vm_count": 2,
				"vcpu_count": 4,
				"vram": 8,
				"cpu_util": 25.5,
				"ram_util": 50.0,
				"free_disk_space": 100000000,
				"anka_disk_usage": 50000000,
				"disk_size": 500000000,
				"state": "Active",
				"capacity": 4,
				"host_arch": "arm64"
			}
		]
	}`

	var nr NodesResponse
	if err := json.Unmarshal([]byte(jsonData), &nr); err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	if nr.GetStatus() != "OK" {
		t.Errorf("Expected status OK, got %s", nr.GetStatus())
	}
	nodes := nr.GetBody().([]Node)
	if len(nodes) != 1 {
		t.Fatalf("Expected 1 node, got %d", len(nodes))
	}
	if nodes[0].NodeID != "abc123" {
		t.Errorf("Expected node_id abc123, got %s", nodes[0].NodeID)
	}
	if nodes[0].HostArch != "arm64" {
		t.Errorf("Expected host_arch arm64, got %s", nodes[0].HostArch)
	}
}

func TestPredefinedStates(t *testing.T) {
	// Test that predefined states contain expected values
	expectedNodeStates := []string{"Offline", "Inactive (Invalid License)", "Active", "Updating", "Drain Mode"}
	if len(NodeStates) != len(expectedNodeStates) {
		t.Errorf("Expected %d node states, got %d", len(expectedNodeStates), len(NodeStates))
	}
	for i, state := range expectedNodeStates {
		if NodeStates[i] != state {
			t.Errorf("Expected node state %s at index %d, got %s", state, i, NodeStates[i])
		}
	}

	expectedArchitectures := []string{"amd64", "arm64"}
	if len(Architectures) != len(expectedArchitectures) {
		t.Errorf("Expected %d architectures, got %d", len(expectedArchitectures), len(Architectures))
	}

	expectedInstanceStates := []string{"Scheduling", "Pulling", "Started", "Stopping", "Stopped", "Terminating", "Terminated", "Error", "Pushing"}
	if len(InstanceStates) != len(expectedInstanceStates) {
		t.Errorf("Expected %d instance states, got %d", len(expectedInstanceStates), len(InstanceStates))
	}
}
