package main

import (
	"errors"
	"fmt"
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

func (t *Tracker) Record(r CheckResult) (Event, error) {
	event := Event{
		Kind: EventNone,
	}

	if !r.OK {
		if r.Service == "" {
			return Event{}, errors.New("no service name received.")
		}
		t.failures[r.Service]++
		if t.failures[r.Service] == t.threshold {
			event.Kind = EventOpened
			inc := Incident{
				id:          t.nextID(),
				title:       fmt.Sprintf("%s is down.", r.Service),
				description: fmt.Sprintf("%s has failed healthchecks %d times. This meets or exceeds its defined threshold: %d", r.Service, t.failures[r.Service], t.threshold),
				service:     r.Service,
				severity:    0, // Placeholder
				status:      "Triggered",
				openedAt:    r.At,
			}
			t.open[r.Service] = &inc
			event.Incident = inc
		}
	} else {
		fmt.Printf("Resotred, count was at %d for %s\n", t.failures[r.Service], r.Service)
		t.failures[r.Service] = 0
		event.Incident = t.open[r.Service].Resolve(r.At)
		t.open[r.Service] = &event.Incident
	}
	return event, nil
}
