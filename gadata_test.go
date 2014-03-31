package gadata

import (
	"log"
	"testing"
	"time"
)

// TestGetData verifies Google Analytics API response to
// a basic request
func TestGetData(t *testing.T) {
	gaTemp := new(GAData)

	// create the comms channel and initialise GAData object
	ch := make(chan string)
	gaTemp.Init(ch)

	testRequest := GaRequest{"ga:43047246",
		"2014-01-01",
		"2014-01-02",
		"ga:visits",
		"ga:day",
		"",
		"",
		"",
		100}
	// launch data pull in the background
	go gaTemp.GetData(&testRequest)

	var msg string
	select {
	case msg = <-ch:
		log.Printf("received %s", msg)
	// time out after 10 seconds
	case <-time.After(10 * 1e9):
		log.Println("I give up waiting...")
	}
}

func TestWasteTime(t *testing.T) {
	// waste some time
	time.Sleep(3000 * time.Millisecond)
}
