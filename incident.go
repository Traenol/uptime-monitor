package main

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
}

func (i Incident) Resolve() Incident {
	i.status = "Resolved"
	return i
}

func idGenerator() func() int {
	next := 0
	return func() int {
		next++
		return next
	}
}
