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

	return
}

// We don't know what it will be (sort of)
type GaResponse interface{}

// GAData is the primary Google Analytics API pull structure
type GAData struct {
	Auth     *OauthData
	Request  *GaRequest
	Response *GaResponse
	Ch       chan string
}

// GetData queries GA API endpoint, passing the response via a channel
func (g *GAData) GetData(request *GaRequest) (err error) {
	g.Request = request
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
	// pass the response to the channel
	g.Ch <- string(contents)
	return
}
