/*
	GAData is a library for retrieving Google Analytics API (v3) data
*/

package gadata

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// Base constants
const (
	DataEndpoint string = "https://www.googleapis.com/analytics/v3/data/ga"
	Limit        int    = 5 // max requests / sec guard
)

// GaRequest is the Google Analytics request structure
type GaRequest struct {
	Id         string `json:"ids"`
	StartDate  string `json:"start-date"`
	EndDate    string `json:"end-date"`
	Metrics    string `json:"metrics"`
	Dimensions string `json:"dimensions"`
	Filters    string `json:"filters"`
	Segments   string `json:"segment"`
	Sort       string `json:"sort"`
	MaxResults int    `json:"max-results"`
	Attempts   int
}

// clipEmpty removes empty keys from struct
// ...so bad
func clipEmpty(out *url.Values, key string) {
	if len(out.Get(key)) == 0 {
		out.Del(key)
	}
}

// Gen random number
func randomOffset(min, max int) int {
	rand.Seed(time.Now().Unix())
	out := rand.Intn(max-min) + min
	return out
}

// ToURLValues converts struct to url.Values struct
func (a *GaRequest) ToURLValues() (out url.Values) {
	out = url.Values{"ids": {a.Id},
		"start-date":  {a.StartDate},
		"end-date":    {a.EndDate},
		"metrics":     {a.Metrics},
		"dimensions":  {a.Dimensions},
		"filters":     {a.Filters},
		"segment":     {a.Segments},
		"sort":        {a.Sort},
		"max-results": {strconv.Itoa(a.MaxResults)},
	}
	clipEmpty(&out, "filters")
	clipEmpty(&out, "segment")
	if len(a.Sort) == 0 {
		clipEmpty(&out, "sort")
	}

	return
}

// Initial returned response
type GaResponse struct {
	Data string
	Pos  int
}

// Processed GA response
type CleanResponse struct {
	Columns []struct {
		Name  string `json:"name"`
		CType string `json:"columnType"`
		DType string `json:"dataType"`
	} `json:"columnHeaders"`
	Total map[string]string `json:"totalsForAllResults"`
	Rows  [][]string        `json:"rows"`
}

func (rawResponse GaResponse) Process() (data CleanResponse, ok bool) {
	if err := json.Unmarshal([]byte(rawResponse.Data), &data); err == nil {
		ok = true
	} else {
		fmt.Printf(err.Error())
	}
	return
}

// GAData is the primary Google Analytics API pull structure
type GAData struct {
	Auth     *OauthData
	Request  *GaRequest
	Response *GaResponse
}

// Initialise the GAData connection, ready to make a new request
func (g *GAData) Init() {
	g.Auth = new(OauthData)
	g.Auth.InitConnection()
}

// GetData queries GA API endpoint, returns response
func (g *GAData) GetData(key int, request *GaRequest) *GaResponse {
	out := request.ToURLValues().Encode()
	client := &http.Client{}
	req, err := http.NewRequest("GET", fmt.Sprintf("%s?%s", DataEndpoint, out), nil)
	checkError(err)
	req.Header.Add("Authorization", "Bearer "+g.Auth.tokens.AccessToken)
	resp, err := client.Do(req)
	// resp, err := http.Get(fmt.Sprintf("%s?%s", DataEndpoint, out))
	checkError(err)
	defer resp.Body.Close()
	contents, err := ioutil.ReadAll(resp.Body)
	checkError(err)

	// pass the response
	response := new(GaResponse)
	response.Data = string(contents)
	response.Pos = key
	if strings.Contains(response.Data, "Invalid Credentials") {
		log.Printf(response.Data)
		if err = g.Auth.refreshToken(); err == nil {
			return g.GetData(key, request)
		}
	} else if strings.Contains(response.Data, "userRateLimitExceeded") { // Retry if hitting limit max 5 times
		if request.Attempts < 5 {
			time.Sleep(time.Duration(randomOffset(1, 10)) * time.Second)
			request.Attempts += 1
			g.GetData(key, request)
		}

	}
	if strings.Contains(response.Data, "\"error\"") {
		log.Println(response.Data)
	}
	return response
}

// BatchGet runs all queries in parellel and returns the results (or times out)
func (g *GAData) BatchGet(requests []*GaRequest) (responses []*CleanResponse, err error) {
	responses = make([]*CleanResponse, len(requests))
	ch := make(chan *GaResponse)
	// Start with a single request to ensure connectivity / token validity
	_, ok := g.GetData(0, requests[0]).Process()
	if ok {
		for a, b := range requests {
			// if we hit max requests limit / sec, sleep for 1 sec.
			if a%Limit == 0 {
				time.Sleep(time.Duration(randomOffset(1, 5)) * time.Second)
			}
			go func(x int, z *GaRequest) { ch <- g.GetData(x, z) }(a, b)
		}

		for i := 0; i < len(requests); i++ {
			select {
			case result := <-ch:
				if out, ok := result.Process(); ok {
					responses[result.Pos] = &out
				}
			// 60 sec timeout
			case <-time.After(60 * time.Second):
				return
			}
		}
	}

	return
}
