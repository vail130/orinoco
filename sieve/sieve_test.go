package sieve_test

import (
	"encoding/json"
	"testing"
	
	"github.com/vail130/orinoco/httputils"
	"github.com/vail130/orinoco/stringutils"
	"github.com/vail130/orinoco/sieve"
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
	
	url = "http://localhost:9966/events/test"
	data, _ := httputils.GetDataFromUrl(url)
	
	var eventSummary sieve.EventSummary
	json.Unmarshal(data, &eventSummary)
	
	if eventSummary.Event != "test" {
		t.Error(stringutils.Concat(string(eventSummary.Event), " != test"))
	}
}

func TestSieveUsesTimestampProvidedForTestEvent(t *testing.T) {
	var url string
	
	timestampString := "2015-05-27T09:29:13-04:00"
	
	url = "http://localhost:9966/events/test"
    jsonData := []byte(stringutils.Concat(`{"timestamp":"`, timestampString, `"}`))
	httputils.PostDataToUrl(url, "application/json", jsonData)
	
	url = "http://localhost:9966/events/test"
	data, _ := httputils.GetDataFromUrl(url)
	
	var eventSummary sieve.EventSummary
	json.Unmarshal(data, &eventSummary)
	
	if eventSummary.Timestamp != timestampString {
		t.Error(stringutils.Concat(string(eventSummary.Timestamp), " != ", timestampString))
	}
}

func TestSieveReturnsNullForAllEventsAfterResetting(t *testing.T) {
    jsonData := []byte(`{"a":1}`)
	httputils.PostDataToUrl("http://localhost:9966/events/test", "application/json", jsonData)
	
	httputils.Delete("http://localhost:9966/events")
	
	data, _ := httputils.GetDataFromUrl("http://localhost:9966/events")
	if string(data) != "null" {
		t.Error(stringutils.Concat(string(data), " != null"))
	}
}

func TestSieveReturnsNullForIndividualEventAfterResetting(t *testing.T) {
    jsonData := []byte(`{"a":1}`)
	httputils.PostDataToUrl("http://localhost:9966/events/test", "application/json", jsonData)
	
	httputils.Delete("http://localhost:9966/events")
	
	data, _ := httputils.GetDataFromUrl("http://localhost:9966/events/test")
	if string(data) != "null" {
		t.Error(stringutils.Concat(string(data), " != null"))
	}
}
