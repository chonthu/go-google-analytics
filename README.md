##Google analytics Data pull

Lightweight Golang library for pulling Google Analytics API data.
Built for use with Core Reporting API (v3):

https://developers.google.com/analytics/devguides/reporting/core/v3/reference

Is being used for BAU report generation and collation. Pull requests welcome!

[ ![Codeship Status for vly/go-gadata](https://www.codeship.io/projects/ee9cdc60-9af7-0131-e5cd-7e7415696371/status?branch=master)](https://www.codeship.io/projects/17520)

### Authentication
In order to authenticate this library for use with your Google Analytics account, an oauth2 token needs to be generated. For a new project login to [Google Developers Console](https://console.developers.google.com) and Create Project. Add Analytics API to list of APIs,  create a new Client ID and download it in JSON format.
Place the client_secret.json in the root of your application.

### Usage

Example single request flow:

```
go
import (
    "fmt"
    "github.com/vly/go-gadata"
)

func main() {
	// initialise GAData
    gaTest := new(GAData)
	
		// initialise instance incl. authentication
    gaTest.Init()
    
    // build a basic GA query, replace your ga ID
    testRequest := GaRequest{"ga:43047246", // GA id
		                         "2014-01-01",  // start date
		                         "2014-01-02",  // end date 
		                         "ga:visits",   // metrics 
		                         "ga:day",      // dimensions
		                         "",            // filters
		                         "",            // segments
		                         "",            // sort
		                         100}           // results no.
    
    // launch data
		result := gaTemp.GetData(1, &testRequest)
		fmt.Printf("results: %s\n", result)
	}
}
```

Example multiple requests flow. 
Returns a sorted slice (array) of results...

```
go
import (
    "fmt"
    "github.com/vly/go-gadata"
)

func main() {
	// initialise GAData
    gaTest := new(gadata.GAData)
	
		// initialise instance incl. authentication
    gaTest.Init()
    
    // build a basic GA query, replace your ga ID
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
```

### Testing
Unit tests are included with this library, use `go test ./...` to run through the set provided. 

### Changelog
#### 0.1.0:
- Initial release