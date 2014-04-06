##GAData

Lightweight Golang library for pulling Google Analytics API data.
Built for use with Core Reporting API (v3):

https://developers.google.com/analytics/devguides/reporting/core/v3/reference

### Authentication
In order to authenticate this library for use with your Google Analytics account, an oauth2 token needs to be generated. For a new project login to [Google Developers Console](https://console.developers.google.com) and Create Project. Add Analytics API to list of APIs,  create a new Client ID and download it in JSON format.

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
		result := gaTemp.GetData(&testRequest)
		fmt.Printf("results: %s\n", result)
	}
}
```

Example multiple requests flow:

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
    ffor i := 0; i < 5; i++ {
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