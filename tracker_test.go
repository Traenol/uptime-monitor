package main

import (
	"testing"
	"time"
)

type wantIncident struct {
	id      int
	service string
}

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
				if got, err := tracker.Record(check); got.Kind != tc.expected[i] {
					if err != nil {
						t.Errorf("%v = %v, want %v.\n", check, got.Kind.String(), tc.expected[i].String())
					}
				}
			}
		})
	}
}

func TestIncidentGetsUniqueID(t *testing.T) {
	tests := []struct {
		name       string
		recordTest []CheckResult
		expected   []wantIncident
	}{
		{"opens 3 incidents in order", []CheckResult{result("api", false), result("api", false), result("api", false), result("portal", false), result("portal", false), result("portal", false), result("mailserver", false), result("mailserver", false), result("mailserver", false)}, []wantIncident{{1, "api"}, {2, "portal"}, {3, "mailserver"}}},
		{"opens 3 incidents staggered", []CheckResult{result("api", false), result("portal", false), result("mailserver", false), result("api", false), result("portal", false), result("mailserver", false), result("api", false), result("portal", false), result("mailserver", false)}, []wantIncident{{1, "api"}, {2, "portal"}, {3, "mailserver"}}},
		{"opens 3 incidents unordered", []CheckResult{result("mailserver", false), result("portal", false), result("api", false), result("portal", false), result("portal", false), result("api", false), result("mailserver", false), result("mailserver", false), result("api", false)}, []wantIncident{{1, "portal"}, {2, "mailserver"}, {3, "api"}}},
		{"opens 2 incidents unordered", []CheckResult{result("mailserver", false), result("portal", false), result("api", false), result("portal", false), result("portal", false), result("api", false), result("mailserver", false), result("mailserver", false)}, []wantIncident{{1, "portal"}, {2, "mailserver"}}},
		{"opens 1 incident, 4 records", []CheckResult{result("mailserver", false), result("mailserver", false), result("mailserver", false), result("mailserver", false)}, []wantIncident{{1, "mailserver"}}},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tracker := NewTracker(3, idGenerator())
			opened := 0
			for _, check := range tc.recordTest {
				got, err := tracker.Record(check)
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if got.Kind != EventOpened {
					continue
				}
				if opened >= len(tc.expected) {
					t.Fatalf("more incidents opened than expected (%d)", len(tc.expected))
				}
				want := tc.expected[opened]
				if got.Incident.id != want.id {
					t.Errorf("incident id = %d want %d", got.Incident.id, want.id)
				}
				if got.Incident.service != want.service {
					t.Errorf("service name = %q want %q", got.Incident.service, want.service)
				}
				opened++
			}
			if opened != len(tc.expected) {
				t.Errorf("opened %d incidents, want %d", opened, len(tc.expected))
			}
		})
	}
}

func TestIncidentResolveAfterSuccess(t *testing.T) {
	tests := []struct {
		name       string
		recordTest []CheckResult
		expected   []wantIncident
	}{
		{"opens 3 incidents, closes 2", []CheckResult{result("api", false), result("api", false), result("api", false), result("portal", false), result("portal", false), result("portal", false), result("mailserver", false), result("mailserver", false), result("mailserver", false)}, []wantIncident{{1, "api"}, {2, "portal"}, {3, "mailserver"}}},
		{"opens 3 incidents staggered", []CheckResult{result("api", false), result("portal", false), result("mailserver", false), result("api", false), result("portal", false), result("mailserver", false), result("api", false), result("portal", false), result("mailserver", false)}, []wantIncident{{1, "api"}, {2, "portal"}, {3, "mailserver"}}},
		{"opens 3 incidents unordered", []CheckResult{result("mailserver", false), result("portal", false), result("api", false), result("portal", false), result("portal", false), result("api", false), result("mailserver", false), result("mailserver", false), result("api", false)}, []wantIncident{{1, "portal"}, {2, "mailserver"}, {3, "api"}}},
		{"opens 2 incidents unordered", []CheckResult{result("mailserver", false), result("portal", false), result("api", false), result("portal", false), result("portal", false), result("api", false), result("mailserver", false), result("mailserver", false)}, []wantIncident{{1, "portal"}, {2, "mailserver"}}},
		{"opens 1 incident, 4 records", []CheckResult{result("mailserver", false), result("mailserver", false), result("mailserver", false), result("mailserver", false)}, []wantIncident{{1, "mailserver"}}},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {

		})
	}
}

func TestIncidentReset(t *testing.T) {
	tests := []struct {
		name       string
		recordTest []CheckResult
		expected   []wantIncident
	}{
		{"2 failures, 1 success, 2 failures", []CheckResult{result("api", false), result("api", false), result("api", true), result("api", false), result("api", false)}, []wantIncident{}},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tracker := NewTracker(3, idGenerator())
			for _, check := range tc.recordTest {
				got, err := tracker.Record(check)
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if got.Kind == EventOpened {
					t.Errorf("incident opened below threshold.")
				}
			}
		})
	}
}

func TestSecondIncidentAfterResolve(t *testing.T) {

}

func TestIncidentOpenResolveTimes(t *testing.T) {

}

func TestEmptyServiceError(t *testing.T) {

}

func result(service string, ok bool) CheckResult {
	return CheckResult{
		Service: service,
		OK:      ok,
		At:      time.Date(2026, 11, 1, 12, 0, 0, 0, time.UTC),
	}
}
