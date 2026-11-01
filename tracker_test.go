package main

import "testing"

func TestTracker(t *testing.T) {
	tests := []struct {
		name       string
		recordTest []CheckResult
		expected   []EventKind
	}{
		{"opens on 3rd failure", []CheckResult{result("api", false), result("api", false), result("api", false)}, []EventKind{EventNone, EventNone, EventOpened}},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tracker := NewTracker(3, idGenerator())
			for i, check := range tc.recordTest {
				if got, _ := tracker.Record(check); got.Kind != tc.expected[i] {
					t.Errorf("%v = %v, want %v\n", check, got.Kind.String(), tc.expected[i].String())
				}
			}
		})
	}
}
