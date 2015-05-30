package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	ga "github.com/chonthu/go-google-analytics"
	"io"
	"log"
	"os"
	"strings"
)

type Data struct {
	Headers []map[string]string `json:"columnHeaders"`
	Rows    [][]string          `json:"rows"`
	Total   map[string]string   `json:"totalsForAllResults"`
	Count   int                 `json:"totalResults"`
	Sampled bool                `json:"containsSampledData"`
}

func (d *Data) GetSeries() {
	for _, b := range d.Rows {
		fmt.Printf("%s, ", b[1])
	}
}

func UrlFilter(url string) string {
	return strings.SplitAfter(url, "//")[1]
}

func FlushBatch(gaTest *ga.GAData, requests []*ga.GaRequest) {
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

func ProcessURL(line *[]string) (req *ga.GaRequest) {
	if (*line)[1] == "E" {
		req = &ga.GaRequest{"ga:43047246",
			"2014-01-01",
			"2014-09-28",
			"ga:uniqueEvents",
			"ga:date",
			"ga:eventLabel==" + UrlFilter((*line)[3]),
			"",
			"",
			500}
	} else {
		req = &ga.GaRequest{"ga:43047246",
			"2014-01-01",
			"2014-09-28",
			"ga:visits",
			"ga:date",
			"ga:pagePath==" + UrlFilter((*line)[3]),
			"",
			"",
			500}
	}
	return
}

func main() {

	// initialise GAData
	gaTest := new(ga.GAData)
	// var requests []*gadata.GaRequest

	// initialise instance incl. authentication
	gaTest.Init()

	// i := 0
	if file, err := os.Open("sample.csv"); err == nil {
		defer file.Close()
		r := csv.NewReader(file)
		for {
			d, err := r.Read()
			if err == io.EOF {
				break
			} else if err != nil {
				log.Fatalln(err.Error())
			}
			test := new(Data)
			result := gaTest.GetData(1, ProcessURL(&d))
			if ok := json.Unmarshal([]byte((*result).Data), test); ok == nil {
				fmt.Printf("%s, %s, %s, ", d[0], d[4], d[3])
				test.GetSeries()
				fmt.Printf("\n")
			}
			// i += 1
			// if i%10 == 0 {
			// 	FlushBatch(gaTest, requests)
			// 	requests = make([]*gadata.GaRequest, 0)
			// }
		}
	} else {
		log.Fatalln("Failed to read csv file.")
	}

}
