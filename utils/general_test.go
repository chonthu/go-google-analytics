package utils

import (
	"log"
	// "reflect"
	"testing"
)

func TestCreateSeries(t *testing.T) {
	testData := map[string]interface{}{
		"a": 1,
		"b": "b",
		"c": 1.3,
	}

	series := Series(&testData)

	keys := series.Headers
	for _, a := range *keys {
		if _, s, ok := series.Get(a); ok {
			if s != testData[a] {
				t.Error("Something went wrong creating Series")
			}
			log.Println(s)
		}
	}
}

func TestSeriesSetIndex(t *testing.T) {
	testData := map[string]interface{}{
		"a": 1,
		"b": "b",
		"c": 1.3,
	}
	newIndex := []string{"d", "e", "f"}
	series := Series(&testData)
	series.SetIndex(&newIndex)

	if series.Index() != &newIndex {
		t.Error("Failed to update Series index")
	}
}

func TestSeriesAdd(t *testing.T) {
	testData := map[string]interface{}{
		"a": 1,
		"b": "b",
		"c": 1.3,
	}
	series := Series(&testData)

	series.Add("d", 1.3)

	if _, _, ok := series.Get("d"); !ok {
		t.Error("Failed to add item")
	}
}

func TestSeriesDelete(t *testing.T) {
	testData := map[string]interface{}{
		"a": 1,
		"b": "b",
		"c": 1.3,
	}
	series := Series(&testData)

	series.Delete("a")

	if _, _, ok := series.Get("a"); ok {
		t.Error("Failed to delete item")
	}
}
