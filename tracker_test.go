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

const elapsed = 5 * time.Minute

func TestTracker(t *testing.T) {
	tests := []struct {
		name       string
		recordTest []CheckResult
		expected   []EventKind
	}{
		{"opens on 3rd failure", []CheckResult{result("api", false, at.Add(elapsed)), result("api", false, at.Add(elapsed)), result("api", false, at.Add(elapsed*2))}, []EventKind{EventNone, EventNone, EventOpened}},
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
				if got.Kind == EventOpened && got.Incident.openedAt != check.At {
					t.Errorf("Incident opened at.Add(elapsed) %v, want %v", got.Incident.openedAt, check.At)
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
		{"opens 3 incidents in order", []CheckResult{result("api", false, at.Add(elapsed)), result("api", false, at.Add(elapsed*2)), result("api", false, at.Add(elapsed*3)), result("portal", false, at.Add(elapsed*4)), result("portal", false, at.Add(elapsed*5)), result("portal", false, at.Add(elapsed*6)), result("mailserver", false, at.Add(elapsed*7)), result("mailserver", false, at.Add(elapsed*8)), result("mailserver", false, at.Add(elapsed*9))}, []wantIncident{{1, "api"}, {2, "portal"}, {3, "mailserver"}}},
		{"opens 3 incidents staggered", []CheckResult{result("api", false, at.Add(elapsed)), result("portal", false, at.Add(elapsed*2)), result("mailserver", false, at.Add(elapsed*3)), result("api", false, at.Add(elapsed*4)), result("portal", false, at.Add(elapsed*5)), result("mailserver", false, at.Add(elapsed*6)), result("api", false, at.Add(elapsed*7)), result("portal", false, at.Add(elapsed*8)), result("mailserver", false, at.Add(elapsed*9))}, []wantIncident{{1, "api"}, {2, "portal"}, {3, "mailserver"}}},
		{"opens 3 incidents unordered", []CheckResult{result("mailserver", false, at.Add(elapsed)), result("portal", false, at.Add(elapsed*2)), result("api", false, at.Add(elapsed*3)), result("portal", false, at.Add(elapsed*4)), result("portal", false, at.Add(elapsed*5)), result("api", false, at.Add(elapsed*6)), result("mailserver", false, at.Add(elapsed*7)), result("mailserver", false, at.Add(elapsed*8)), result("api", false, at.Add(elapsed*9))}, []wantIncident{{1, "portal"}, {2, "mailserver"}, {3, "api"}}},
		{"opens 2 incidents unordered", []CheckResult{result("mailserver", false, at.Add(elapsed)), result("portal", false, at.Add(elapsed*2)), result("api", false, at.Add(elapsed*3)), result("portal", false, at.Add(elapsed*4)), result("portal", false, at.Add(elapsed*5)), result("api", false, at.Add(elapsed*6)), result("mailserver", false, at.Add(elapsed*7)), result("mailserver", false, at.Add(elapsed*8))}, []wantIncident{{1, "portal"}, {2, "mailserver"}}},
		{"opens 1 incident, 4 records", []CheckResult{result("mailserver", false, at.Add(elapsed)), result("mailserver", false, at.Add(elapsed*2)), result("mailserver", false, at.Add(elapsed*3)), result("mailserver", false, at.Add(elapsed*4))}, []wantIncident{{1, "mailserver"}}},
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

func TestIncidentResolveCount(t *testing.T) {
	tests := []struct {
		name           string
		recordTest     []CheckResult
		expectedOpened int
		expectedClosed int
	}{
		{"1 incident creat.Add(elapsed)ed and 0 resolved", []CheckResult{result("api", false, at.Add(elapsed)), result("api", false, at.Add(elapsed*2)), result("api", false, at.Add(elapsed*3))}, 1, 0},
		{"1 incident creat.Add(elapsed)ed and 1 resolved", []CheckResult{result("api", false, at.Add(elapsed)), result("api", false, at.Add(elapsed*2)), result("api", false, at.Add(elapsed*3)), result("api", true, at.Add(elapsed*4))}, 1, 1},
		{"2 incident creat.Add(elapsed)ed and 1 resolved", []CheckResult{result("api", false, at.Add(elapsed)), result("api", false, at.Add(elapsed*2)), result("api", false, at.Add(elapsed*3)), result("api", true, at.Add(elapsed*4)), result("api", false, at.Add(elapsed*5)), result("api", false, at.Add(elapsed*6)), result("api", false, at.Add(elapsed*7))}, 2, 1},
		{"2 incident creat.Add(elapsed)ed and 1 resolved, 2 services", []CheckResult{result("api", false, at.Add(elapsed)), result("api", false, at.Add(elapsed*2)), result("api", false, at.Add(elapsed*3)), result("db", false, at.Add(elapsed*4)), result("db", false, at.Add(elapsed*5)), result("db", false, at.Add(elapsed*6)), result("api", true, at.Add(elapsed*7))}, 2, 1},
		{"2 incident creat.Add(elapsed)ed and 2 resolved", []CheckResult{result("api", false, at.Add(elapsed)), result("api", false, at.Add(elapsed*2)), result("api", false, at.Add(elapsed*3)), result("api", true, at.Add(elapsed*4)), result("api", false, at.Add(elapsed*5)), result("api", false, at.Add(elapsed*6)), result("api", false, at.Add(elapsed*7)), result("api", true, at.Add(elapsed*8))}, 2, 2},
		{"no indicents opened", []CheckResult{result("api", true, at.Add(elapsed))}, 0, 0},
		{"2 fail 1 ok, no incidents", []CheckResult{result("api", false, at.Add(elapsed)), result("api", false, at.Add(elapsed*2)), result("api", true, at.Add(elapsed*3))}, 0, 0},
		{"2 fail 1 ok 2 fail, no incidents", []CheckResult{result("api", false, at.Add(elapsed)), result("api", false, at.Add(elapsed*2)), result("api", true, at.Add(elapsed*3)), result("api", false, at.Add(elapsed*4)), result("api", false, at.Add(elapsed*5))}, 0, 0},
		{"1 incident with 2 succeeds", []CheckResult{result("api", false, at.Add(elapsed)), result("api", false, at.Add(elapsed*2)), result("api", false, at.Add(elapsed*3)), result("api", true, at.Add(elapsed*4)), result("api", true, at.Add(elapsed*5))}, 1, 1},
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
					if !got.Incident.openedAt.Equal(check.At) {
						t.Errorf("incident opened at %v, want %v", got.Incident.openedAt, check.At)
					}
					opened++
				case EventResolved:
					if !got.Incident.resolvedAt.Equal(check.At) {
						t.Errorf("incident resolved at %v, want %v", got.Incident.resolvedAt, check.At)
					}
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
	}{
		{"1 incident 1 service", []CheckResult{result("api", false, at.Add(elapsed)), result("api", false, at.Add(elapsed*2)), result("api", false, at.Add(elapsed*3)), result("api", true, at.Add(elapsed*4))}},
		{"1 incident 2 service", []CheckResult{result("api", false, at.Add(elapsed)), result("db", false, at.Add(elapsed*2)), result("api", false, at.Add(elapsed*3)), result("api", false, at.Add(elapsed*4)), result("api", true, at.Add(elapsed*5))}},
		{"2 incident 1 service", []CheckResult{result("api", false, at.Add(elapsed)), result("api", false, at.Add(elapsed*2)), result("api", false, at.Add(elapsed*3)), result("api", true, at.Add(elapsed*4)), result("api", false, at.Add(elapsed*5)), result("api", false, at.Add(elapsed*6)), result("api", false, at.Add(elapsed*7)), result("api", true, at.Add(elapsed*8))}},
		{"3 incident 2 service", []CheckResult{result("api", false, at.Add(elapsed)), result("api", false, at.Add(elapsed*2)), result("api", false, at.Add(elapsed*3)), result("api", true, at.Add(elapsed*4)), result("db", false, at.Add(elapsed*5)), result("db", false, at.Add(elapsed*6)), result("db", false, at.Add(elapsed*7)), result("db", true, at.Add(elapsed*8)), result("api", false, at.Add(elapsed*9)), result("api", false, at.Add(elapsed*10)), result("api", false, at.Add(elapsed*11)), result("api", true, at.Add(elapsed*12))}},
		{"0 incident 1 service", []CheckResult{result("api", false, at.Add(elapsed)), result("api", false, at.Add(elapsed*2)), result("api", true, at.Add(elapsed*3))}},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tracker := NewTracker(3, idGenerator())
			for _, check := range tc.recordTest {
				got, err := tracker.Record(check)
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				switch got.Kind {
				case EventOpened:
					if !got.Incident.openedAt.Equal(check.At) {
						t.Errorf("opened at = %v, want %v", got.Incident.openedAt, check.At)
					}
				case EventResolved:
					if !got.Incident.resolvedAt.Equal(check.At) {
						t.Errorf("resolved at = %v, want %v", got.Incident.resolvedAt, check.At)
					}
				}
			}
		})
	}
}

func TestEmptyServiceError(t *testing.T) {
	tests := []struct {
		name       string
		recordTest []CheckResult
		expected   error
	}{
		{"No service name 1 fail", []CheckResult{result("", false, at.Add(elapsed))}, ErrEmptyService},
		{"No service name 1 succeed", []CheckResult{result("", true, at.Add(elapsed))}, ErrEmptyService},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tracker := NewTracker(3, idGenerator())
			for _, check := range tc.recordTest {
				_, err := tracker.Record(check)
				if !errors.Is(err, ErrEmptyService) {
					t.Errorf("got %v want %v", err, ErrEmptyService)
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
