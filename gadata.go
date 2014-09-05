/*
	GAData is a library for retrieving Google Analytics API (v3) data
*/

package gadata

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// Base constants
const (
	DataEndpoint string = "https://www.googleapis.com/analytics/v3/data/ga"
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
}

// clipEmpty removes empty keys from struct
// ...so bad
func clipEmpty(out *url.Values, key string) {
	if len(out.Get(key)) == 0 {
		out.Del(key)
	}
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

// We don't know what it will be (sort of)
type GaResponse struct {
	Data string
	Pos  int
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
	resp, _ := client.Do(req)
	// resp, err := http.Get(fmt.Sprintf("%s?%s", DataEndpoint, out))
	checkError(err)
	defer resp.Body.Close()
	contents, err := ioutil.ReadAll(resp.Body)
	checkError(err)

	// pass the response
	response := new(GaResponse)
	response.Data = string(contents)
	response.Pos = key
	if strings.Contains(response, "Invalid Credentials") {
		if err = g.Auth.refreshToken(); err == nil {
			g.GetData(key, request)
		}
	}
	return response
}

// BatchGet runs all queries in parellel and returns the results (or times out)
func (g *GAData) BatchGet(requests []*GaRequest) (responses []string, err error) {
	ch := make(chan *GaResponse)
	for a, b := range requests {
		go func(x int, z *GaRequest) { ch <- g.GetData(x, z) }(a, b)
	}
	responses = make([]string, len(requests))
	for i := 0; i < len(requests); i++ {
		select {
		case result := <-ch:
			responses[result.Pos] = result.Data
			//responses = append(responses, result)
		// 60 sec timeout
		case <-time.After(60 * time.Second):
			return
		}
	}

	return
}
