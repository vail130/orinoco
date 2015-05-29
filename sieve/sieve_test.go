package sieve_test

import (
	"encoding/json"
	"testing"

	"gopkg.in/check.v1"

	"github.com/vail130/orinoco/httputils"
	"github.com/vail130/orinoco/sieve"
)

func TestSieveApi(t *testing.T) { check.TestingT(t) }

type SieveApiTestSuite struct{}

var _ = check.Suite(&SieveApiTestSuite{})

func (s *SieveApiTestSuite) SetUpTest(c *check.C) {
	httputils.Delete("http://localhost:9966/events")
}

func (s *SieveApiTestSuite) TestSieveReturnsNullWhenNoEventsHaveBeenPosted(c *check.C) {
	url := "http://localhost:9966/events"
	data, _ := httputils.GetDataFromUrl(url)
	c.Assert(string(data), check.Equals, "null")
}

func (s *SieveApiTestSuite) TestSieveReturnsEventSummaryAfterReceivingData(c *check.C) {
	var url string

	url = "http://localhost:9966/events/test"
	jsonData := []byte(`{"a":1}`)
	httputils.PostDataToUrl(url, "application/json", jsonData)

	url = "http://localhost:9966/events/test"
	data, _ := httputils.GetDataFromUrl(url)

	var eventSummary sieve.EventSummary
	json.Unmarshal(data, &eventSummary)

	c.Assert(eventSummary.Event, check.Equals, "test")
}

func (s *SieveApiTestSuite) TestSieveReturnsNullForAllEventsAfterResetting(c *check.C) {
	jsonData := []byte(`{"a":1}`)
	httputils.PostDataToUrl("http://localhost:9966/events/test", "application/json", jsonData)

	httputils.Delete("http://localhost:9966/events")

	data, _ := httputils.GetDataFromUrl("http://localhost:9966/events")
	c.Assert(string(data), check.Equals, "null")
}

func (s *SieveApiTestSuite) TestSieveReturnsNullForIndividualEventAfterResetting(c *check.C) {
	jsonData := []byte(`{"a":1}`)
	httputils.PostDataToUrl("http://localhost:9966/events/test", "application/json", jsonData)

	httputils.Delete("http://localhost:9966/events")

	data, _ := httputils.GetDataFromUrl("http://localhost:9966/events/test")
	c.Assert(string(data), check.Equals, "null")
}
