package gadata

import (
	"fmt"
	"testing"
	"time"
)

// TestGetData verifies Google Analytics API response to
// a basic request
func TestGetData(t *testing.T) {
	gaTemp := new(GAData)
	gaTemp.Auth = new(OauthData)
	ch := make(chan string)
	gaTemp.Ch = ch
	gaTemp.Auth.InitConnection()

	testRequest := GaRequest{"ga:43047246",
		"2014-01-01",
		"2014-01-02",
		"ga:visits",
		"ga:day",
		"",
		"",
		"ga:day",
		100}
	// launch data pull in the background
	go gaTemp.GetData(&testRequest)

	var msg string
	select {
	case msg = <-ch:
		fmt.Printf("received %s", msg)
	// time out after 10 seconds
	case <-time.After(10 * 1e9):
		fmt.Println("I give up waiting...")
	}
}
