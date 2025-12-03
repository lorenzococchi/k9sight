package k8s

import (
	"testing"
)

func TestLabelsMatch(t *testing.T) {
	tests := []struct {
		name     string
		selector map[string]string
		labels   map[string]string
		expected bool
	}{
		{
			name:     "exact match",
			selector: map[string]string{"app": "nginx"},
			labels:   map[string]string{"app": "nginx"},
			expected: true,
		},
		{
			name:     "selector subset of labels",
			selector: map[string]string{"app": "nginx"},
			labels:   map[string]string{"app": "nginx", "env": "prod", "version": "v1"},
			expected: true,
		},
		{
			name:     "selector not in labels",
			selector: map[string]string{"app": "nginx"},
			labels:   map[string]string{"app": "redis"},
			expected: false,
		},
		{
			name:     "selector key missing from labels",
			selector: map[string]string{"app": "nginx", "env": "prod"},
			labels:   map[string]string{"app": "nginx"},
			expected: false,
		},
		{
			name:     "empty selector matches everything",
			selector: map[string]string{},
			labels:   map[string]string{"app": "nginx"},
			expected: true,
		},
		{
			name:     "empty labels with non-empty selector",
			selector: map[string]string{"app": "nginx"},
			labels:   map[string]string{},
			expected: false,
		},
		{
			name:     "both empty",
			selector: map[string]string{},
			labels:   map[string]string{},
			expected: true,
		},
		{
			name:     "multiple selector labels all match",
			selector: map[string]string{"app": "nginx", "env": "prod"},
			labels:   map[string]string{"app": "nginx", "env": "prod", "version": "v1"},
			expected: true,
		},
		{
			name:     "multiple selector labels partial match",
			selector: map[string]string{"app": "nginx", "env": "prod"},
			labels:   map[string]string{"app": "nginx", "env": "staging"},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := labelsMatch(tt.selector, tt.labels)
			if result != tt.expected {
				t.Errorf("labelsMatch(%v, %v) = %v, want %v", tt.selector, tt.labels, result, tt.expected)
			}
		})
	}
}

func TestJoinStrings(t *testing.T) {
	tests := []struct {
		name     string
		input    []string
		sep      string
		expected string
	}{
		{
			name:     "empty slice",
			input:    []string{},
			sep:      ", ",
			expected: "",
		},
		{
			name:     "single element",
			input:    []string{"hello"},
			sep:      ", ",
			expected: "hello",
		},
		{
			name:     "two elements",
			input:    []string{"hello", "world"},
			sep:      ", ",
			expected: "hello, world",
		},
		{
			name:     "multiple elements",
			input:    []string{"a", "b", "c", "d"},
			sep:      "-",
			expected: "a-b-c-d",
		},
		{
			name:     "empty separator",
			input:    []string{"a", "b", "c"},
			sep:      "",
			expected: "abc",
		},
		{
			name:     "nil slice",
			input:    nil,
			sep:      ", ",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := joinStrings(tt.input, tt.sep)
			if result != tt.expected {
				t.Errorf("joinStrings(%v, %q) = %q, want %q", tt.input, tt.sep, result, tt.expected)
			}
		})
	}
}

func TestAllResourceTypes(t *testing.T) {
	// Verify AllResourceTypes contains expected types
	expectedTypes := map[ResourceType]bool{
		ResourceDeployments:  true,
		ResourceStatefulSets: true,
		ResourceDaemonSets:   true,
		ResourceJobs:         true,
		ResourceCronJobs:     true,
		ResourcePods:         true,
	}

	if len(AllResourceTypes) != len(expectedTypes) {
		t.Errorf("AllResourceTypes has %d types, expected %d", len(AllResourceTypes), len(expectedTypes))
	}

	for _, rt := range AllResourceTypes {
		if !expectedTypes[rt] {
			t.Errorf("Unexpected resource type in AllResourceTypes: %s", rt)
		}
	}
}
