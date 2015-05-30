package ga

import (
	"log"
	"strconv"
	"testing"
)

// TestGet verifies Google Analytics API response to
// a basic request
func TestGet(t *testing.T) {

	gaTemp := new(GAData)

	// initialise GAData object
	gaTemp.Init()

	testRequest := Request{"ga:23949588",
		"2014-01-01",
		"2014-01-02",
		"ga:visits",
		"ga:day",
		"",
		"",
		"",
		100,
		5}

	result := gaTemp.Get(1, &testRequest)
	log.Println(result)
}

// TestBatchGet checks the batch processing functionality
func TestBatchGet(t *testing.T) {
	var testRequests []*Request
	gaTemp := new(GAData)
	gaTemp.Init()
	for i := 0; i < 10; i++ {
		testRequests = append(testRequests, &Request{"ga:23949588",
			"2014-01-0" + strconv.Itoa(i+1),
			"2014-01-0" + strconv.Itoa(i+2),
			"ga:visits",
			"ga:day",
			"",
			"",
			"",
			100,
			5})
	}
	if results, err := gaTemp.BatchGet(testRequests); err == nil {
		for a, b := range results {
			log.Printf("results: %d: %s", a, b)
		}
	}

}
