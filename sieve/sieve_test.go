package sieve_test

import (
	"testing"
	
	"github.com/vail130/orinoco/httputils"
	"github.com/vail130/orinoco/stringutils"
)

func TestSieveReturnsNullWhenNoEventsHaveBeenPosted(t *testing.T) {
    url := "http://localhost:9966/events"
	data, _ := httputils.GetDataFromUrl(url)
	if string(data) != "null" {
		t.Error(stringutils.Concat(string(data), " != null"))
	}
}

func TestSieveReturnsEventSummaryAfterReceivingData(t *testing.T) {
	var url string
	
	url = "http://localhost:9966/events/test"
    jsonData := []byte(`{"a":1}`)
	httputils.PostDataToUrl(url, "application/json", jsonData)
	
	url = "http://localhost:9966/events"
	data, _ := httputils.GetDataFromUrl(url)
	if string(data) == "null" {
		t.Error(stringutils.Concat(string(data), " == null"))
	}
}

func TestSieveReturnsNullAfterResettingIt(t *testing.T) {
    jsonData := []byte(`{"a":1}`)
	httputils.PostDataToUrl("http://localhost:9966/events/test", "application/json", jsonData)
	
	httputils.Delete("http://localhost:9966/events")
	
	data, _ := httputils.GetDataFromUrl("http://localhost:9966/events")
	if string(data) != "null" {
		t.Error(stringutils.Concat(string(data), " != null"))
	}
}
