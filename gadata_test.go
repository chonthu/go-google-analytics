package gadata

import (
	"log"
	"strconv"
	"testing"
)

// TestGetData verifies Google Analytics API response to
// a basic request
func TestGetData(t *testing.T) {

	gaTemp := new(GAData)

	// initialise GAData object
	gaTemp.Init()

	testRequest := GaRequest{"ga:43047246",
		"2014-01-01",
		"2014-01-02",
		"ga:visits",
		"ga:day",
		"",
		"",
		"",
		100}

	result := gaTemp.GetData(1, &testRequest)
	log.Println(result)
}

// TestBatchGet checks the batch processing functionality
func TestBatchGet(t *testing.T) {
	var testRequests []*GaRequest
	gaTemp := new(GAData)
	gaTemp.Init()
	for i := 0; i < 5; i++ {
		testRequests = append(testRequests, &GaRequest{"ga:43047246",
			"2014-01-0" + strconv.Itoa(i+1),
			"2014-01-0" + strconv.Itoa(i+2),
			"ga:visits",
			"ga:day",
			"",
			"",
			"",
			100})
	}
	if results, err := gaTemp.BatchGet(testRequests); err == nil {
		for a, b := range results {
			log.Printf("results: %d: %s", a, b)
		}
	}

}
