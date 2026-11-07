package main

import "time"

func IncidentFactory() func(title, desc, service string, severity int) Incident {
	nextID := idGenerator()
	return func(title, desc, service string, severity int) Incident {
		return Incident{
			id:          nextID(),
			title:       title,
			description: desc,
			service:     service,
			severity:    severity,
			status:      "Triggered",
			openedAt:    time.Now(),
		}
	}
}

type Incident struct {
	id          int
	title       string
	description string
	service     string
	severity    int
	status      string
	openedAt    time.Time
	resolvedAt  time.Time
}

func (i Incident) Resolve(t time.Time) Incident {
	i.status = "Resolved"
	i.resolvedAt = t
	return i
}

func idGenerator() func() int {
	next := 0
	return func() int {
		next++
		return next
	}
}
