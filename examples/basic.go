package main

import (
	"fmt"
	ga "github.com/chonthu/go-google-analytics"
)

func main() {
	// initialise GAData
	analtyics := new(ga.GAData)

	// initialise instance incl. authentication
	analtyics.Init()

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
	result := analtyics.Get(1, &testRequest)
	fmt.Printf("results: %s\n", result)
}
