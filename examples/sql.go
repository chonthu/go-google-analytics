package main

import (
	"encoding/json"
	"fmt"
	ga "github.com/chonthu/go-google-analytics"
)

func main() {
	// initialise Client
	analytics := new(ga.Client)

	// initialise instance incl. authentication
	analytics.Init()

	// build a basic GA query, replace your ga ID
	testSQL := ga.SQL("SELECT visits,day FROM 56659181 WHERE start_date > 201500607 AND end_date > today ORDER BY visits LIMIT 100")

	// launch data
	result := analytics.Get(1, &testSQL)
	reponse := new(ga.ResponseData)
	json.Unmarshal(result.Data, &reponse)
	fmt.Println(reponse)
}
