package main

import (
	"errors"
	"testing"
	"time"
)

type wantIncident struct {
	id      int
	service string
}

var at = time.Date(2026, 11, 1, 12, 0, 0, 0, time.UTC)

func TestTracker(t *testing.T) {
	tests := []struct {
		name       string
		recordTest []CheckResult
		expected   []EventKind
	}{
		{"opens on 3rd failure", []CheckResult{result("api", false, at), result("api", false, at), result("api", false, at)}, []EventKind{EventNone, EventNone, EventOpened}},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tracker := NewTracker(3, idGenerator())
			for i, check := range tc.recordTest {
				got, err := tracker.Record(check)
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if got.Kind != tc.expected[i] {
					t.Errorf("%v = %v, want %v.\n", check, got.Kind.String(), tc.expected[i].String())
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
		{"opens 3 incidents in order", []CheckResult{result("api", false, at), result("api", false, at), result("api", false, at), result("portal", false, at), result("portal", false, at), result("portal", false, at), result("mailserver", false, at), result("mailserver", false, at), result("mailserver", false, at)}, []wantIncident{{1, "api"}, {2, "portal"}, {3, "mailserver"}}},
		{"opens 3 incidents staggered", []CheckResult{result("api", false, at), result("portal", false, at), result("mailserver", false, at), result("api", false, at), result("portal", false, at), result("mailserver", false, at), result("api", false, at), result("portal", false, at), result("mailserver", false, at)}, []wantIncident{{1, "api"}, {2, "portal"}, {3, "mailserver"}}},
		{"opens 3 incidents unordered", []CheckResult{result("mailserver", false, at), result("portal", false, at), result("api", false, at), result("portal", false, at), result("portal", false, at), result("api", false, at), result("mailserver", false, at), result("mailserver", false, at), result("api", false, at)}, []wantIncident{{1, "portal"}, {2, "mailserver"}, {3, "api"}}},
		{"opens 2 incidents unordered", []CheckResult{result("mailserver", false, at), result("portal", false, at), result("api", false, at), result("portal", false, at), result("portal", false, at), result("api", false, at), result("mailserver", false, at), result("mailserver", false, at)}, []wantIncident{{1, "portal"}, {2, "mailserver"}}},
		{"opens 1 incident, 4 records", []CheckResult{result("mailserver", false, at), result("mailserver", false, at), result("mailserver", false, at), result("mailserver", false, at)}, []wantIncident{{1, "mailserver"}}},
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
		{"opens 3 incidents, closes 2", []CheckResult{result("api", false, at), result("api", false, at), result("api", false, at), result("portal", false, at), result("portal", false, at), result("portal", false, at), result("mailserver", false, at), result("mailserver", false, at), result("mailserver", false, at)}, []wantIncident{{1, "api"}, {2, "portal"}, {3, "mailserver"}}},
		{"opens 3 incidents staggered", []CheckResult{result("api", false, at), result("portal", false, at), result("mailserver", false, at), result("api", false, at), result("portal", false, at), result("mailserver", false, at), result("api", false, at), result("portal", false, at), result("mailserver", false, at)}, []wantIncident{{1, "api"}, {2, "portal"}, {3, "mailserver"}}},
		{"opens 3 incidents unordered", []CheckResult{result("mailserver", false, at), result("portal", false, at), result("api", false, at), result("portal", false, at), result("portal", false, at), result("api", false, at), result("mailserver", false, at), result("mailserver", false, at), result("api", false, at)}, []wantIncident{{1, "portal"}, {2, "mailserver"}, {3, "api"}}},
		{"opens 2 incidents unordered", []CheckResult{result("mailserver", false, at), result("portal", false, at), result("api", false, at), result("portal", false, at), result("portal", false, at), result("api", false, at), result("mailserver", false, at), result("mailserver", false, at)}, []wantIncident{{1, "portal"}, {2, "mailserver"}}},
		{"opens 1 incident, 4 records", []CheckResult{result("mailserver", false, at), result("mailserver", false, at), result("mailserver", false, at), result("mailserver", false, at)}, []wantIncident{{1, "mailserver"}}},
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
		{"2 failures, 1 success, 2 failures", []CheckResult{result("api", false, at), result("api", false, at), result("api", true, at), result("api", false, at), result("api", false, at)}, []wantIncident{}},
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

func TestIncidentResolveCount(t *testing.T) {
	tests := []struct {
		name           string
		recordTest     []CheckResult
		expectedOpened int
		expectedClosed int
	}{
		{"1 incident created and 0 resolved", []CheckResult{result("api", false, at), result("api", false, at), result("api", false, at)}, 1, 0},
		{"1 incident created and 1 resolved", []CheckResult{result("api", false, at), result("api", false, at), result("api", false, at), result("api", true, at)}, 1, 1},
		{"2 incident created and 1 resolved", []CheckResult{result("api", false, at), result("api", false, at), result("api", false, at), result("api", true, at), result("api", false, at), result("api", false, at), result("api", false, at)}, 2, 1},
		{"2 incident created and 1 resolved, 2 services", []CheckResult{result("api", false, at), result("api", false, at), result("api", false, at), result("db", false, at), result("db", false, at), result("db", false, at), result("api", true, at)}, 2, 1},
		{"2 incident created and 2 resolved", []CheckResult{result("api", false, at), result("api", false, at), result("api", false, at), result("api", true, at), result("api", false, at), result("api", false, at), result("api", false, at), result("api", true, at)}, 2, 2},
		{"no indicents opened", []CheckResult{result("api", true, at)}, 0, 0},
		{"2 fail 1 ok, no incidents", []CheckResult{result("api", false, at), result("api", false, at), result("api", true, at)}, 0, 0},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			opened, closed := 0, 0
			tracker := NewTracker(3, idGenerator())
			for _, check := range tc.recordTest {
				got, err := tracker.Record(check)
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				switch got.Kind {
				case EventOpened:
					opened++
				case EventResolved:
					closed++
				}
			}
			if opened != tc.expectedOpened {
				t.Errorf("opened = %d want %d", opened, tc.expectedOpened)
			}
			if closed != tc.expectedClosed {
				t.Errorf("closed = %d want %d", closed, tc.expectedClosed)
			}
		})
	}

}

func TestIncidentOpenAndResolveTimes(t *testing.T) {
	tests := []struct {
		name       string
		recordTest []CheckResult
		openedAt   []time.Time
		resolvedAt []time.Time
	}{
		{"1 incident 1 service", []CheckResult{result("api", false, at), result("api", false, at), result("api", false, at), result("api", true, at)}, []time.Time{time.Date(2026, 11, 1, 12, 0, 0, 0, time.UTC)}, []time.Time{time.Date(2026, 11, 1, 12, 0, 0, 0, time.UTC)}},
		{"1 incident 2 service", []CheckResult{result("api", false, at), result("db", false, at), result("api", false, at), result("api", false, at), result("api", true, at)}, []time.Time{time.Date(2026, 11, 1, 12, 0, 0, 0, time.UTC)}, []time.Time{time.Date(2026, 11, 1, 12, 0, 0, 0, time.UTC)}},
		{"2 incident 1 service", []CheckResult{result("api", false, at), result("api", false, at), result("api", false, at), result("api", true, at), result("api", false, at), result("api", false, at), result("api", false, at), result("api", true, at)}, []time.Time{time.Date(2026, 11, 1, 12, 0, 0, 0, time.UTC), time.Date(2026, 11, 1, 12, 0, 0, 0, time.UTC)}, []time.Time{time.Date(2026, 11, 1, 12, 0, 0, 0, time.UTC), time.Date(2026, 11, 1, 12, 0, 0, 0, time.UTC)}},
		{"3 incident 2 service", []CheckResult{result("api", false, at), result("api", false, at), result("api", false, at), result("api", true, at), result("db", false, at), result("db", false, at), result("db", false, at), result("db", true, at), result("api", false, at), result("api", false, at), result("api", false, at), result("api", true, at)}, []time.Time{time.Date(2026, 11, 1, 12, 0, 0, 0, time.UTC), time.Date(2026, 11, 1, 12, 0, 0, 0, time.UTC), time.Date(2026, 11, 1, 12, 0, 0, 0, time.UTC)}, []time.Time{time.Date(2026, 11, 1, 12, 0, 0, 0, time.UTC), time.Date(2026, 11, 1, 12, 0, 0, 0, time.UTC), time.Date(2026, 11, 1, 12, 0, 0, 0, time.UTC)}},
		{"0 incident 1 service", []CheckResult{result("api", false, at), result("api", false, at), result("api", true, at)}, []time.Time{time.Date(2026, 11, 1, 12, 0, 0, 0, time.UTC)}, []time.Time{time.Date(2026, 11, 1, 12, 0, 0, 0, time.UTC)}},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Skip("TODO")
		})
	}
}

func TestEmptyServiceError(t *testing.T) {
	var emptyServiceName = errors.New("no service name received")
	tests := []struct {
		name       string
		recordTest []CheckResult
		expected   string
	}{
		{"No service name", []CheckResult{result("", false, at)}, emptyServiceName.Error()},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tracker := NewTracker(3, idGenerator())
			for _, check := range tc.recordTest {
				_, err := tracker.Record(check)
				if errors.Is(err, emptyServiceName) {
					t.Errorf("got %v want %v", err, emptyServiceName)
				}
			}
		})
	}

}

func result(service string, ok bool, atTime time.Time) CheckResult {
	return CheckResult{
		Service: service,
		OK:      ok,
		At:      atTime,
	}
}
