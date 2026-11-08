package main

import (
	"errors"
	"fmt"
	"os"
)

type ConsoleNotifier struct{}

func (cn ConsoleNotifier) Notify(i Incident) {
	fmt.Printf("Notification for incident %d:%s\n", i.id, i.title)
	sum, err := formattedSummaryLong(i)
	if err != nil {
		fmt.Printf("error received: %v", err)
		os.Exit(1)
	}
	fmt.Printf("%s", sum)
	fmt.Println("End Notification")
}

type ShortNotifier struct{}

func (sn ShortNotifier) Notify(i Incident) {
	sum, err := formattedSummaryShort(i)
	if err != nil {
		fmt.Printf("error received: %v", err)
		os.Exit(1)
	}
	fmt.Printf("Notification: %s", sum)
}

type Notifier interface {
	Notify(i Incident)
}

var ErrInvalidIncident = errors.New("invalid incident")

// Creates a formatted string that summarized the incident
func formattedSummaryLong(i Incident) (string, error) {
	var inc Incident
	if i == inc {
		return "", ErrInvalidIncident
	}
	if i.status == "Resolved" {
		return fmt.Sprintf("Incident #%d is resolved.\n", i.id), nil
	} else {
		return fmt.Sprintf("Incident #%d\nTitle: %s\nDescription: %s\nService: %s\nSeverity: %d\nStatus: %s\n", i.id, i.title, i.description, i.service, i.severity, i.status), nil
	}
}
func formattedSummaryShort(i Incident) (string, error) {
	var inc Incident
	if i == inc {
		return "", ErrInvalidIncident
	}
	if i.status == "Resolved" {
		return fmt.Sprintf("Incident #%d - Resolved.\n", i.id), nil
	} else {
		return fmt.Sprintf("Incident #%d - Title: %s\n", i.id, i.title), nil
	}
}

func notifyWith(n Notifier) func(Incident) {
	return func(i Incident) {
		n.Notify(i)
	}
}
