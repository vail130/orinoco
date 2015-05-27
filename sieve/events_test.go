package sieve_test

import (
	"encoding/json"
	"testing"
	
	"gopkg.in/check.v1"
	
	"github.com/vail130/orinoco/httputils"
	"github.com/vail130/orinoco/stringutils"
	"github.com/vail130/orinoco/sieve"
)

func TestSieveEvents(t *testing.T) { check.TestingT(t) }

type SieveEventsTestSuite struct{}
var _ = check.Suite(&SieveEventsTestSuite{})

func (s *SieveEventsTestSuite) SetUpTest(c *check.C) {
	httputils.Delete("http://localhost:9966/events")
}

func (s *SieveEventsTestSuite) TestSieveUsesTimestampProvidedForTestEvent(c *check.C) {
	var url string
	
	timestampString := "2015-05-27T09:29:13-04:00"
	
	url = "http://localhost:9966/events/test"
    jsonData := []byte(stringutils.Concat(`{"timestamp":"`, timestampString, `"}`))
	httputils.PostDataToUrl(url, "application/json", jsonData)
	
	url = stringutils.Concat("http://localhost:9966/events/test?timestamp=", timestampString)
	data, _ := httputils.GetDataFromUrl(url)
	
	var eventSummary sieve.EventSummary
	json.Unmarshal(data, &eventSummary)
	
	c.Assert(eventSummary.Timestamp, check.Equals, timestampString)
}
