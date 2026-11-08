package main

import (
	"fmt"
	"time"
)

func main() {
	var at = time.Date(2026, 11, 1, 12, 0, 0, 0, time.UTC)

	const elapsed = 5 * time.Minute
	const threshold = 3
	tracker, err := NewTracker(threshold, idGenerator())
	sendConsole := notifyWith(ConsoleNotifier{})
	sendShort := notifyWith(ShortNotifier{})

	if err != nil {
		fmt.Printf("error received: %v", err)
		return
	}

	tests := []CheckResult{
		{"api", false, at}, {"api", false, at.Add(elapsed)}, {"api", false, at.Add(elapsed * 2)},
		{"api", true, at.Add(elapsed * 3)}, {"api", false, at.Add(elapsed * 4)}, {"db", false, at.Add(elapsed * 5)},
		{"db", true, at.Add(elapsed * 6)}, {"api", true, at.Add(elapsed * 7)}, {"db", false, at.Add(elapsed * 8)},
		{"db", false, at.Add(elapsed * 9)}, {"api", false, at.Add(elapsed * 10)}, {"db", false, at.Add(elapsed * 11)},
	}
	for _, test := range tests {
		event, err := tracker.Record(test)
		if err != nil {
			fmt.Printf("error received recording event: %v, event: %v", err, test)
			return
		}
		switch event.Kind {
		case EventOpened:
			sendConsole(event.Incident)
		case EventResolved:
			sendShort(event.Incident)
		}
	}
}
