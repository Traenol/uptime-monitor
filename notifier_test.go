package main

import (
	"fmt"
	"testing"
	"time"
)

type fakeNotifier struct{ got []Incident }

func (f *fakeNotifier) Notify(i Incident) { f.got = append(f.got, i) }

var _ Notifier = ConsoleNotifier{}
var _ Notifier = ShortNotifier{}

func TestFormattedSummaryLong(t *testing.T) {
	tests := []struct {
		name          string
		recordTest    Incident
		expected      string
		expectedError error
	}{{
		"normal incident", Incident{
			id:          1,
			title:       fmt.Sprintf("%s is down.", "api"),
			description: fmt.Sprintf("%s has failed healthchecks %d times. This meets or exceeds its defined threshold: %d", "api", 3, 3),
			service:     "api",
			severity:    0, // Placeholder
			status:      "Triggered",
			openedAt:    time.Date(2026, 11, 1, 12, 0, 0, 0, time.UTC)},
		"Incident #1\nTitle: api is down.\nDescription: api has failed healthchecks 3 times. This meets or exceeds its defined threshold: 3\nService: api\nSeverity: 0\nStatus: Triggered\n", nil,
	}, {
		"broken incident", Incident{}, "", ErrInvalidIncident,
	}, {
		"resolved incident", Incident{
			id:          85,
			title:       fmt.Sprintf("%s is down.", "api"),
			description: fmt.Sprintf("%s has failed healthchecks %d times. This meets or exceeds its defined threshold: %d", "api", 3, 3),
			service:     "api",
			severity:    0, // Placeholder
			status:      "Resolved",
			openedAt:    time.Date(2026, 11, 1, 12, 0, 0, 0, time.UTC),
			resolvedAt:  time.Date(2026, 11, 1, 12, 30, 0, 0, time.UTC)},
		"Incident #85 is resolved.\n", nil,
	},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := formattedSummaryLong(tc.recordTest)
			if err != tc.expectedError {
				t.Errorf("unexpected error occured: %v", err)
			}
			if got != tc.expected {
				t.Errorf("summary = %s, want %s", got, tc.expected)
			}
		})
	}
}

func TestFormattedSummaryShort(t *testing.T) {
	tests := []struct {
		name          string
		recordTest    Incident
		expected      string
		expectedError error
	}{{
		"normal incident", Incident{
			id:          1,
			title:       fmt.Sprintf("%s is down.", "api"),
			description: fmt.Sprintf("%s has failed healthchecks %d times. This meets or exceeds its defined threshold: %d", "api", 3, 3),
			service:     "api",
			severity:    0, // Placeholder
			status:      "Triggered",
			openedAt:    time.Date(2026, 11, 1, 12, 0, 0, 0, time.UTC)},
		"Incident #1 - Title: api is down.\n", nil,
	}, {
		"broken incident", Incident{}, "", ErrInvalidIncident,
	}, {
		"resolved incident", Incident{
			id:          10,
			title:       fmt.Sprintf("%s is down.", "api"),
			description: fmt.Sprintf("%s has failed healthchecks %d times. This meets or exceeds its defined threshold: %d", "api", 3, 3),
			service:     "api",
			severity:    0, // Placeholder
			status:      "Resolved",
			openedAt:    time.Date(2026, 11, 1, 12, 0, 0, 0, time.UTC),
			resolvedAt:  time.Date(2026, 11, 1, 12, 30, 0, 0, time.UTC)},
		"Incident #10 - Resolved.\n", nil,
	},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := formattedSummaryShort(tc.recordTest)
			if err != tc.expectedError {
				t.Errorf("unexpected error occured: %v", err)
			}
			if got != tc.expected {
				t.Errorf("summary = %s, want %s", got, tc.expected)
			}
		})
	}
}

func TestNotifyWith(t *testing.T) {
	fake := fakeNotifier{}
	send := notifyWith(&fake)
	send(Incident{
		id:          1,
		title:       fmt.Sprintf("%s is down.", "api"),
		description: fmt.Sprintf("%s has failed healthchecks %d times. This meets or exceeds its defined threshold: %d", "api", 3, 3),
		service:     "api",
		severity:    0, // Placeholder
		status:      "Triggered",
		openedAt:    time.Date(2026, 11, 1, 12, 0, 0, 0, time.UTC),
	})
	if len(fake.got) < 1 {
		t.Errorf("Notifier did not receive incident")
	}
}
