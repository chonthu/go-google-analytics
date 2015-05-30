package main

import (
	"encoding/json"
	"fmt"
	ga "github.com/chonthu/go-google-analytics"
)

func main() {
	// initialise Client
	analytics := new(ga.Client)

	// initialise instance incl. authentication
	analytics.Init()

	// build a basic GA query, replace your ga ID
	testRequest := ga.Request{
		"ga:56659181", // GA id
		"2014-01-01",  // start date
		"2014-01-02",  // end date
		"ga:visits",   // metrics
		"ga:day",      // dimensions
		"",            // filters
		"",            // segments
		"",            // sort
		100,           // results no.
		5,             // number of attempts
	}

	// launch data
	result := analytics.Get(1, &testRequest)
	reponse := new(ga.ResponseData)
	json.Unmarshal(result.Data, &reponse)
	fmt.Println(reponse)
}
