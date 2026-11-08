package main

import (
	"time"
)

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
