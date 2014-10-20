package gadata

import (
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"reflect"
	"strings"
)

// Suggests resolution to common errors
func errorSolution(err string) (string, bool) {
	switch {
	default:
		return "", false
	case strings.Contains(err, "no such host"):
		return "Can't contact Google Analytics server, no internet connect?", true
	case strings.Contains(err, "Error unmarshalling config file"):
		return "Bad client_secret JSON file, try generating a new one", true
	}

}

// Check for normal errors
func checkError(err error) {
	if err != nil {
		log.Fatalln(err.Error())
	}
}

// BrowserOpen opens a URL is the OS' default web browser
func BrowserOpen(url string) error {
	return exec.Command("open", url).Run()
}

// WebCallback listens on a predefined port for a oauth response
// sends back via channel once it receives a response and shuts down.
func WebCallback(ch chan string) {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.RequestURI, "/?") {
			block := strings.SplitAfter(r.RequestURI, "/?")[1]
			if !strings.Contains(block, "code=") {
				ch <- block
			} else {
				ch <- strings.SplitAfter(block, "code=")[1]
				fmt.Fprintln(w, "Authentication completed, you can close this window.")
				close(ch)
				return
			}
		}
		fmt.Fprintln(w, "Error encountered during authentication.")
		return
	})

	log.Fatalf("Server exited: %v", http.ListenAndServe(ReturnURI, nil))
}

// Diffmap data structure
type DiffMap [][]int

// Results Row data structure
type SeriesData struct {
	Data    *[]interface{}
	Headers *[]string
	Size    int
}

// Initiate Results row
func Series(data *map[string]interface{}) *SeriesData {
	output := new(SeriesData)

	i := 0
	keys := make([]string, len(*data))
	values := make([]interface{}, len(*data))
	for k, v := range *data {
		keys[i] = k
		values[i] = v
		i++
	}
	output.Data = &values
	output.Headers = &keys
	output.Size = len(keys)
	return output
}

func (s *SeriesData) Get(key string) (loc int, val interface{}, ok bool) {
	for i, a := range *s.Headers {
		if a == key {
			val = (*s.Data)[i]
			ok = true
			loc = i
			return
		}
	}
	return
}

func (s *SeriesData) Add(key string, val interface{}) (ok bool) {
	if _, _, has := s.Get(key); !has {
		*s.Headers = append(*s.Headers, key)
		*s.Data = append(*s.Data, val)
		ok = true
		s.Size += 1
	}
	return
}

func (s *SeriesData) Delete(key string) (ok bool) {
	if loc, _, has := s.Get(key); has {
		*s.Headers = append((*s.Headers)[:loc], (*s.Headers)[loc+1:]...)
		*s.Data = append((*s.Data)[:loc], (*s.Data)[loc+1:]...)
		ok = true
		s.Size -= 1
	}
	return
}

// If row is made up of numeric numbers, get total
func (s *SeriesData) Sum() (total float64, ok bool) {
	total = 0.0
	for k, v := range *s.Data {
		if reflect.TypeOf(v) == reflect.TypeOf(1) || reflect.TypeOf(v) == reflect.TypeOf(1.1) {
			total += float64(reflect.ValueOf(v).Float())
		} else {
			log.Println(k, " failed")
			return
		}
	}
	ok = true
	return
}

// Set or replace column index
func (s *SeriesData) SetIndex(newIndex *[]string) (ok bool) {
	if len(*newIndex) == len(*s.Data) {
		s.Headers = newIndex
		ok = true
	}
	return
}

// Get column index
func (s *SeriesData) Index() *[]string {
	return s.Headers
}

// Swap items
func (s *SeriesData) Swap(i int, j int) {
	(*s.Data)[i], (*s.Data)[j] = (*s.Data)[j], (*s.Data)[i]
	(*s.Headers)[i], (*s.Headers)[j] = (*s.Headers)[j], (*s.Headers)[i]
}

// Find difference between Series
func (s *SeriesData) Diff(s2 *SeriesData) (out *DiffMap, ok bool) {
	if len((*s.Data)) != len((*s.Data)) {
		return
	}

	out = new(DiffMap)
	for i, v := range *s.Data {
		for j, z := range *s2.Data {
			if v == z {
				(*out)[i] = []int{i, j}
			}
		}
	}
	ok = true
	return
}

// Apply DiffMap
func (s *SeriesData) ApplyDiff(diff *DiffMap) {
	for _, b := range *diff {
		s.Swap(b[0], b[1])
	}
}

// Set of results
type DataFrameData struct {
	Index *[]string
	Data  *[]*SeriesData
	Size  int
}

// Create new DataFrame
func DataFrame(data *map[string]*SeriesData) *DataFrameData {
	output := new(DataFrameData)
	i := 0
	keys := make([]string, len(*data))
	values := make([]*SeriesData, len(*data))
	for k, v := range *data {
		keys[i] = k
		values[i] = v
		i++
	}
	output.Index = &keys
	output.Data = &values
	output.Size = len(*data)

	return output
}

// Daterange object
type Daterange struct {
	Start     string
	End       string
	Precision string
	Segments  []*Daterange
	IsSplit   bool
}

// Get Time objects
func (d Daterange) GetDates() {

}

// Split daterange into months
func (d Daterange) SplitMonths() bool {

	return false
}

// Split daterange into weeks
func (d Daterange) SplitWeeks() bool {

	return false
}

// Split daterange into days
func (d Daterange) SplitDays() bool {
	return false

}
