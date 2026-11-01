package main

import "fmt"

type ConsoleNotifier struct{}

func (cn ConsoleNotifier) Notify(i Incident) {
	fmt.Printf("Notification for incident %d:%s\n", i.id, i.title)
	fmt.Printf("%s", formattedSummaryLong(i))
	fmt.Println("End Notification")
}

type ShortNotifier struct{}

func (sn ShortNotifier) Notify(i Incident) {
	fmt.Printf("Notification: %s", formattedSummaryShort(i))
}

type Notifier interface {
	Notify(i Incident)
}

// Creates a formatted string that summarized the incident
func formattedSummaryLong(i Incident) string {
	var s string
	if i.status == "Resolved" {
		s = fmt.Sprintf("Incident #%d is resolved.\n", i.id)
	} else {
		s = fmt.Sprintf("Incident #%d\nTitle: %s\nDescription: %s\nService: %s\nSeverity: %d\nStatus: %s\n", i.id, i.title, i.description, i.service, i.severity, i.status)
	}
	return s
}
func formattedSummaryShort(i Incident) string {
	var s string
	if i.status == "Resolved" {
		s = fmt.Sprintf("Incident #%d - Resolved.\n", i.id)
	} else {
		s = fmt.Sprintf("Incident #%d - Title: %s\n", i.id, i.title)
	}
	return s
}

func notifyWith(n Notifier) func(Incident) {
	return func(i Incident) {
		n.Notify(i)
	}
}
