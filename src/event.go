package main

import (
	"fmt"
	"log"
	"os"
	"time"
)

func initEventLog() {
	file, err := os.OpenFile("events.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	check(err)
	log.SetOutput(file)
}

func trackEvent(event Event) {
	colref := firestoreClient.Collection("events")

	_, result, err := colref.Add(ctx, event)
	check(err)
	fmt.Println(result, event.String())
}

func makeEvent(number string, match MatchInterface) Event {

	//pointsRed := []string{"foo", "bar", "hello"}
	//ointsWhite := []string{"foo", "bar", "world"}

	// type = Start, ScoreRed, ScoreWhite, Hansoku, Encho, Finished

	//Pool 8, Fight 2 on Shiaijo A. Furl vs. Somame. Furl could achieve two men (M),
	//Somame could achieve a tsuki (T). Furl has already has a hikiwake as well (h)
	return Event{
		Type:        "Start",
		Description: match.Description(),
		Payload:     "",
		Date:        time.Now().Format(dateLayout),
	}

}
