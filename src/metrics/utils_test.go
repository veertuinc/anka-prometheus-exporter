package metrics

import (
	"testing"

	"github.com/veertuinc/anka-prometheus-exporter/src/types"
)

func TestIntMapFromTwoStringSlices(t *testing.T) {
	tests := []struct {
		name        string
		outer       []string
		inner       []string
		expectedLen int
	}{
		{
			name:        "architectures and states",
			outer:       []string{"amd64", "arm64"},
			inner:       []string{"Active", "Offline"},
			expectedLen: 2,
		},
		{
			name:        "empty slices",
			outer:       []string{},
			inner:       []string{},
			expectedLen: 0,
		},
		{
			name:        "single outer, multiple inner",
			outer:       []string{"arch1"},
			inner:       []string{"state1", "state2", "state3"},
			expectedLen: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := intMapFromTwoStringSlices(tt.outer, tt.inner)
			if len(result) != tt.expectedLen {
				t.Errorf("Expected outer map length %d, got %d", tt.expectedLen, len(result))
			}

			// Verify inner maps are initialized correctly
			for _, outerKey := range tt.outer {
				innerMap, ok := result[outerKey]
				if !ok {
					t.Errorf("Expected key %s in outer map", outerKey)
					continue
				}
				if len(innerMap) != len(tt.inner) {
					t.Errorf("Expected inner map length %d for key %s, got %d", len(tt.inner), outerKey, len(innerMap))
				}
				for _, innerKey := range tt.inner {
					val, ok := innerMap[innerKey]
					if !ok {
						t.Errorf("Expected key %s in inner map for %s", innerKey, outerKey)
					}
					if val != 0 {
						t.Errorf("Expected initial value 0, got %d", val)
					}
				}
			}
		})
	}
}

func TestUniqueThisStringArray(t *testing.T) {
	tests := []struct {
		name     string
		input    []string
		expected []string
	}{
		{
			name:     "no duplicates",
			input:    []string{"a", "b", "c"},
			expected: []string{"a", "b", "c"},
		},
		{
			name:     "with duplicates",
			input:    []string{"a", "b", "a", "c", "b"},
			expected: []string{"a", "b", "c"},
		},
		{
			name:     "all duplicates",
			input:    []string{"a", "a", "a"},
			expected: []string{"a"},
		},
		{
			name:     "empty array",
			input:    []string{},
			expected: []string{},
		},
		{
			name:     "single element",
			input:    []string{"only"},
			expected: []string{"only"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := uniqueThisStringArray(tt.input)
			if len(result) != len(tt.expected) {
				t.Errorf("Expected length %d, got %d", len(tt.expected), len(result))
			}
			// Check that all expected values are present
			for i, v := range tt.expected {
				if result[i] != v {
					t.Errorf("Expected %s at index %d, got %s", v, i, result[i])
				}
			}
		})
	}
}

func TestUniqueNodeGroupsArray(t *testing.T) {
	tests := []struct {
		name     string
		input    []types.NodeGroup
		expected int
	}{
		{
			name: "no duplicates",
			input: []types.NodeGroup{
				{Id: "group-1", Name: "Group 1"},
				{Id: "group-2", Name: "Group 2"},
			},
			expected: 2,
		},
		{
			name: "with duplicates by ID",
			input: []types.NodeGroup{
				{Id: "group-1", Name: "Group 1"},
				{Id: "group-1", Name: "Group 1 Duplicate"},
				{Id: "group-2", Name: "Group 2"},
			},
			expected: 2,
		},
		{
			name:     "empty array",
			input:    []types.NodeGroup{},
			expected: 0,
		},
		{
			name: "all duplicates",
			input: []types.NodeGroup{
				{Id: "same-id", Name: "Name 1"},
				{Id: "same-id", Name: "Name 2"},
				{Id: "same-id", Name: "Name 3"},
			},
			expected: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := uniqueNodeGroupsArray(tt.input)
			if len(result) != tt.expected {
				t.Errorf("Expected %d unique groups, got %d", tt.expected, len(result))
			}
		})
	}
}

func TestCreateGaugeMetric(t *testing.T) {
	name := "test_metric"
	help := "Test metric description"

	gauge := CreateGaugeMetric(name, help)

	if gauge == nil {
		t.Fatal("CreateGaugeMetric returned nil")
	}

	// Verify we can set a value without panic
	gauge.Set(42.0)
}

func TestCreateGaugeMetricVec(t *testing.T) {
	name := "test_metric_vec"
	help := "Test metric vector description"
	labels := []string{"label1", "label2"}

	gaugeVec := CreateGaugeMetricVec(name, help, labels)

	if gaugeVec == nil {
		t.Fatal("CreateGaugeMetricVec returned nil")
	}

	// Verify we can set values with labels without panic
	gaugeVec.WithLabelValues("value1", "value2").Set(100.0)
}

func TestConvertToStatusData(t *testing.T) {
	tests := []struct {
		name        string
		input       interface{}
		expectError bool
	}{
		{
			name: "valid status",
			input: types.Status{
				Status:  "Running",
				Version: "1.0.0",
			},
			expectError: false,
		},
		{
			name:        "invalid type - string",
			input:       "not a status",
			expectError: true,
		},
		{
			name:        "invalid type - int",
			input:       123,
			expectError: true,
		},
		{
			name:        "invalid type - nil",
			input:       nil,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ConvertToStatusData(tt.input)
			if tt.expectError {
				if err == nil {
					t.Error("Expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if result == nil {
					t.Error("Expected non-nil result")
				}
			}
		})
	}
}

func TestConvertToNodeData(t *testing.T) {
	tests := []struct {
		name        string
		input       interface{}
		expectError bool
		expectLen   int
	}{
		{
			name: "valid nodes",
			input: []types.Node{
				{NodeID: "node-1", NodeName: "Node 1"},
				{NodeID: "node-2", NodeName: "Node 2"},
			},
			expectError: false,
			expectLen:   2,
		},
		{
			name:        "empty nodes",
			input:       []types.Node{},
			expectError: false,
			expectLen:   0,
		},
		{
			name:        "invalid type",
			input:       "not nodes",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ConvertToNodeData(tt.input)
			if tt.expectError {
				if err == nil {
					t.Error("Expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if len(result) != tt.expectLen {
					t.Errorf("Expected length %d, got %d", tt.expectLen, len(result))
				}
			}
		})
	}
}

func TestConvertToRegistryDiskData(t *testing.T) {
	tests := []struct {
		name        string
		input       interface{}
		expectError bool
	}{
		{
			name: "valid registry disk",
			input: types.RegistryDisk{
				Total: 1000000000,
				Free:  500000000,
			},
			expectError: false,
		},
		{
			name:        "invalid type",
			input:       "not registry disk",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ConvertToRegistryDiskData(tt.input)
			if tt.expectError {
				if err == nil {
					t.Error("Expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if result == nil {
					t.Error("Expected non-nil result")
				}
			}
		})
	}
}

func TestConvertToRegistryTemplatesData(t *testing.T) {
	tests := []struct {
		name        string
		input       interface{}
		expectError bool
		expectLen   int
	}{
		{
			name: "valid templates",
			input: []types.Template{
				{UUID: "uuid-1", Name: "Template 1"},
			},
			expectError: false,
			expectLen:   1,
		},
		{
			name:        "invalid type",
			input:       "not templates",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ConvertToRegistryTemplatesData(tt.input)
			if tt.expectError {
				if err == nil {
					t.Error("Expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if len(result) != tt.expectLen {
					t.Errorf("Expected length %d, got %d", tt.expectLen, len(result))
				}
			}
		})
	}
}

func TestConvertToInstancesData(t *testing.T) {
	tests := []struct {
		name        string
		input       interface{}
		expectError bool
		expectLen   int
	}{
		{
			name: "valid instances",
			input: []types.Instance{
				{InstanceID: "instance-1"},
				{InstanceID: "instance-2"},
			},
			expectError: false,
			expectLen:   2,
		},
		{
			name:        "invalid type",
			input:       "not instances",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ConvertToInstancesData(tt.input)
			if tt.expectError {
				if err == nil {
					t.Error("Expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if len(result) != tt.expectLen {
					t.Errorf("Expected length %d, got %d", tt.expectLen, len(result))
				}
			}
		})
	}
}

func TestConvertMetricToGauge(t *testing.T) {
	gauge := CreateGaugeMetric("test_gauge", "Test gauge")

	result, err := ConvertMetricToGauge(gauge)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if result == nil {
		t.Error("Expected non-nil gauge")
	}

	// Test with wrong type
	gaugeVec := CreateGaugeMetricVec("test_vec", "Test vec", []string{"label"})
	_, err = ConvertMetricToGauge(gaugeVec)
	if err == nil {
		t.Error("Expected error for wrong type")
	}
}

func TestConvertMetricToGaugeVec(t *testing.T) {
	gaugeVec := CreateGaugeMetricVec("test_vec", "Test vec", []string{"label"})

	result, err := ConvertMetricToGaugeVec(gaugeVec)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if result == nil {
		t.Error("Expected non-nil gauge vec")
	}

	// Test with wrong type
	gauge := CreateGaugeMetric("test_gauge", "Test gauge")
	_, err = ConvertMetricToGaugeVec(gauge)
	if err == nil {
		t.Error("Expected error for wrong type")
	}
}
