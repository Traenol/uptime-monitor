package main

import "fmt"

func main() {
	create := IncidentFactory()
	inc := create("First Incident", "Baby's first incident", "A Basic service", 0)
	inc2 := create("Second Incident", "More Robust", "Postgres", 4)

	sendConsole := notifyWith(ConsoleNotifier{})
	sendShort := notifyWith(ShortNotifier{})

	//inc = inc.Resolve()
	sendConsole(inc)
	sendConsole(inc2)
	sendShort(inc)
	sendShort(inc2)

	fmt.Printf("%v\n%+v\n%T\n", inc, inc, inc)
}
