package main

import (
	"encoding/json"
	"github.com/vly/go-gadata"
	"log"
	"strconv"
)

type Data struct {
	Headers []map[string]string `json:"columnHeaders"`
	Rows    [][]string          `json:"rows"`
	Total   map[string]string   `json:"totalsForAllResults"`
	Count   int                 `json:"totalResults"`
	Sampled bool                `json:"containsSampledData"`
}

func main() {

	// initialise GAData
	gaTest := new(gadata.GAData)
	var requests []*gadata.GaRequest

	// initialise instance incl. authentication
	gaTest.Init()

	i := 0
	for 
	requests = append(requests, &gadata.GaRequest{"ga:43047246",
		"2014-01-0" + strconv.Itoa(i+1),
		"2014-01-0" + strconv.Itoa(i+2),
		"ga:visits",
		"ga:day",
		"",
		"",
		"",
		100})

	if results, err := gaTest.BatchGet(requests); err == nil {
		for a, b := range results {
			test := new(Data)
			if ok := json.Unmarshal([]byte(b), test); ok == nil {
				log.Printf("result: %d: %v and %v", a, test, test.Headers)
			}
			log.Printf("results: %d: %s", a, b)
		}
	}
}
