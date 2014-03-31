##GAData

Lightweight Golang library for pulling Google Analytics API data.
Built for use with Core Reporting API (v3):

https://developers.google.com/analytics/devguides/reporting/core/v3/reference

### Authentication
In order to authenticate this library for use with your Google Analytics account, an oauth2 token needs to be generated. For a new project login to [Google Developers Console](https://console.developers.google.com) and Create Project. Add Analytics API to list of APIs,  create a new Client ID and download it in JSON format.

### Usage
Example query flow:

```
go
import (
    "fmt"
    "github.com/vly/go-gadata"
)

func main() {
	// initialise GAData
    gaTest := new(GAData)
    
    // create channel to receive results
    ch := make(chan string)
    
    // initialise instance incl. authentication
    gaTest.Init(ch)
    
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
    
    // launch data pull in the background
	go gaTemp.GetData(&testRequest)
	
    // wait for server response to come back
    select {    
	case msg := <-ch:
		fmt.Printf("received: %s", msg)
		
	// time out after 10 seconds, shouldn't take longer
	// than that unless you are running a complex request.
	case <-time.After(10 * 1e9):
		fmt.Println("I give up waiting...")
	}
}
```

### Testing
Unit tests are included with this library, use `go test ./...` to run through the set provided. 

### Changelog
#### 0.1.0:
- Initial release