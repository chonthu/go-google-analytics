package ga

import (
	"encoding/json"
	"fmt"
	"github.com/chonthu/go-google-analytics/utils"
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
	StdEndpoint string = "https://www.googleapis.com/analytics/v3/data/ga" // standard endpoint
	Limit       int    = 5                                                 // max requests / sec guard
)

// Request is the Google Analytics request structure
type Request struct {
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
func (a *Request) ToURLValues() (out url.Values) {
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

type ResponseData struct {
	Id      string
	Query   Request
	Kind    string `json:kind`
	Columns []struct {
		Name  string `json:"name"`
		CType string `json:"columnType"`
		DType string `json:"dataType"`
	} `json:"columnHeaders"`
	Total map[string]string `json:"totalsForAllResults"`
	Rows  [][]string        `json:"rows"`
}

// Initial returned response
type Response struct {
	Data []byte
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

func (rawResponse Response) Process() (data CleanResponse, ok bool) {
	if err := json.Unmarshal([]byte(rawResponse.Data), &data); err == nil {
		ok = true
	} else {
		fmt.Printf(err.Error())
	}
	return
}

// Client is the primary Google Analytics API pull structure
type Client struct {
	Auth     *utils.OauthData
	Request  *Request
	Response *Response
}

// Initialise the Client connection, ready to make a new request
func (g *Client) Init() {
	g.Auth = new(utils.OauthData)
	g.Auth.InitConnection()
}

// Get queries GA API endpoint, returns response
func (g *Client) Get(key int, request *Request) *Response {
	out := request.ToURLValues().Encode()
	client := &http.Client{}
	req, err := http.NewRequest("GET", fmt.Sprintf("%s?%s", StdEndpoint, out), nil)
	utils.CheckError(err)
	req.Header.Add("Authorization", "Bearer "+g.Auth.Tokens.AccessToken)
	resp, err := client.Do(req)
	// resp, err := http.Get(fmt.Sprintf("%s?%s", StdEndpoint, out))
	utils.CheckError(err)
	defer resp.Body.Close()
	contents, err := ioutil.ReadAll(resp.Body)
	utils.CheckError(err)

	// pass the response
	response := new(Response)
	raw := string(contents)

	if strings.Contains(raw, "Invalid Credentials") {
		log.Printf(raw)
		if err = g.Auth.RefreshToken(); err == nil {
			return g.Get(key, request)
		}
	} else if strings.Contains(raw, "userRateLimitExceeded") { // Retry if hitting limit max 5 times
		if request.Attempts < 5 {
			time.Sleep(time.Duration(randomOffset(1, 10)) * time.Second)
			request.Attempts += 1
			g.Get(key, request)
		}
	} else if strings.Contains(raw, "\"error\"") {
		log.Println(raw)
	}

	response.Pos = key
	response.Data = contents

	return response
}

// BatchGet runs all queries in parellel and returns the results (or times out)
func (g *Client) BatchGet(requests []*Request) (responses []*CleanResponse, err error) {
	responses = make([]*CleanResponse, len(requests))
	ch := make(chan *Response)
	// Start with a single request to ensure connectivity / token validity
	_, ok := g.Get(0, requests[0]).Process()
	if ok {
		for a, b := range requests {
			// if we hit max requests limit / sec, sleep for 1 sec.
			if a%Limit == 0 {
				time.Sleep(time.Duration(randomOffset(1, 5)) * time.Second)
			}
			go func(x int, z *Request) { ch <- g.Get(x, z) }(a, b)
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
