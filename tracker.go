package main

import (
	"time"
)

type CheckResult struct {
	Service string
	OK      bool
	At      time.Time
}

type EventKind int

const (
	EventNone EventKind = iota
	EventOpened
	EventResolved
)

func (kind EventKind) String() string {
	switch kind {
	case EventNone:
		return "EventNone"
	case EventOpened:
		return "EventOpened"
	case EventResolved:
		return "EventResolved"
	default:
		return ""
	}
}

type Event struct {
	Kind     EventKind
	Incident Incident
}

type Tracker struct {
	threshold int
	nextID    func() int
	failures  map[string]int
	open      map[string]*Incident
}

func NewTracker(threshold int, nextID func() int) *Tracker {
	return &Tracker{
		threshold: threshold,
		nextID:    nextID,
		failures:  make(map[string]int),
		open:      make(map[string]*Incident),
	}
}

func result(service string, ok bool) CheckResult {
	return CheckResult{
		Service: service,
		OK:      ok,
		At:      time.Date(2026, 11, 1, 12, 0, 0, 0, time.UTC),
	}
}
func (t *Tracker) Record(r CheckResult) (Event, error) {
	event := Event{
		Kind: EventNone,
	}

	if !r.OK {
		t.failures[r.Service]++
	} else {
		t.failures[r.Service] = 0
	}

	if t.failures[r.Service] == t.threshold {
		event.Kind = EventOpened
		// create an incident
	}
	return event, nil
}
