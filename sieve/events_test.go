package sieve_test

import (
	"encoding/json"
	"testing"

	"gopkg.in/check.v1"

	"github.com/vail130/orinoco/httputils"
	"github.com/vail130/orinoco/sieve"
	"github.com/vail130/orinoco/stringutils"
)

func TestSieveEvents(t *testing.T) { check.TestingT(t) }

type SieveEventsTestSuite struct{}

var _ = check.Suite(&SieveEventsTestSuite{})

func (s *SieveEventsTestSuite) SetUpTest(c *check.C) {
	httputils.Delete("http://localhost:9966/events")
}

func (s *SieveEventsTestSuite) TestSieveUsesTimestampProvidedForTestEvent(c *check.C) {
	var url string

	timestampString := "2015-05-27T09:29:13Z"

	url = "http://localhost:9966/events/test"
	jsonData := []byte(stringutils.Concat(`{"timestamp":"`, timestampString, `"}`))
	httputils.PostDataToUrl(url, "application/json", jsonData)

	url = stringutils.Concat("http://localhost:9966/events/test?timestamp=", timestampString)
	data, _ := httputils.GetDataFromUrl(url)

	var eventSummary sieve.EventSummary
	json.Unmarshal(data, &eventSummary)

	c.Assert(eventSummary.Timestamp, check.Equals, timestampString)
}

func (s *SieveEventsTestSuite) TestSieveCalculatesSecondToDateCorrectly(c *check.C) {
	var url string

	timestampString := "2015-01-01T01:00:42Z"
	url = "http://localhost:9966/events/test"
	jsonData := []byte(stringutils.Concat(`{"timestamp":"`, timestampString, `"}`))

	httputils.PostDataToUrl(url, "application/json", jsonData)
	httputils.PostDataToUrl(url, "application/json", jsonData)
	httputils.PostDataToUrl(url, "application/json", jsonData)

	url = stringutils.Concat("http://localhost:9966/events/test?timestamp=", timestampString)
	data, _ := httputils.GetDataFromUrl(url)

	var eventSummary sieve.EventSummary
	json.Unmarshal(data, &eventSummary)
	c.Assert(eventSummary.SecondToDate, check.Equals, 3)
}

func (s *SieveEventsTestSuite) TestSieveCalculatesMinuteToDateCorrectly(c *check.C) {
	var url string

	url = "http://localhost:9966/events/test"
	httputils.PostDataToUrl(url, "application/json", []byte(`{"timestamp":"2015-01-01T01:00:01Z"}`))
	httputils.PostDataToUrl(url, "application/json", []byte(`{"timestamp":"2015-01-01T01:00:15Z"}`))
	httputils.PostDataToUrl(url, "application/json", []byte(`{"timestamp":"2015-01-01T01:00:22Z"}`))
	httputils.PostDataToUrl(url, "application/json", []byte(`{"timestamp":"2015-01-01T01:00:38Z"}`))
	httputils.PostDataToUrl(url, "application/json", []byte(`{"timestamp":"2015-01-01T01:00:42Z"}`))

	url = "http://localhost:9966/events/test?timestamp=2015-01-01T01:00:50Z"
	data, _ := httputils.GetDataFromUrl(url)

	var eventSummary sieve.EventSummary
	json.Unmarshal(data, &eventSummary)
	c.Assert(eventSummary.MinuteToDate, check.Equals, 5)
}

func (s *SieveEventsTestSuite) TestSieveCalculatesHourToDateCorrectly(c *check.C) {
	var url string

	url = "http://localhost:9966/events/test"
	httputils.PostDataToUrl(url, "application/json", []byte(`{"timestamp":"2015-01-01T01:04:00Z"}`))
	httputils.PostDataToUrl(url, "application/json", []byte(`{"timestamp":"2015-01-01T01:15:00Z"}`))
	httputils.PostDataToUrl(url, "application/json", []byte(`{"timestamp":"2015-01-01T01:22:00Z"}`))
	httputils.PostDataToUrl(url, "application/json", []byte(`{"timestamp":"2015-01-01T01:38:00Z"}`))
	httputils.PostDataToUrl(url, "application/json", []byte(`{"timestamp":"2015-01-01T01:42:00Z"}`))

	url = "http://localhost:9966/events/test?timestamp=2015-01-01T01:50:00Z"
	data, _ := httputils.GetDataFromUrl(url)

	var eventSummary sieve.EventSummary
	json.Unmarshal(data, &eventSummary)
	c.Assert(eventSummary.HourToDate, check.Equals, 5)
}

func (s *SieveEventsTestSuite) TestSieveCalculatesProjectedThisMinuteCorrectly(c *check.C) {
	var url string

	timestampString := "2015-01-01T01:00:29Z"
	url = "http://localhost:9966/events/test"
	jsonData := []byte(stringutils.Concat(`{"timestamp":"`, timestampString, `"}`))

	httputils.PostDataToUrl(url, "application/json", jsonData)
	httputils.PostDataToUrl(url, "application/json", jsonData)
	httputils.PostDataToUrl(url, "application/json", jsonData)

	url = stringutils.Concat("http://localhost:9966/events/test?timestamp=", timestampString)
	data, _ := httputils.GetDataFromUrl(url)

	var eventSummary sieve.EventSummary
	json.Unmarshal(data, &eventSummary)
	c.Assert(eventSummary.ProjectedThisMinute, check.Equals, float32(6))
}

func (s *SieveEventsTestSuite) TestSieveCalculatesProjectedThisHourCorrectly(c *check.C) {
	var url string

	timestampString := "2015-01-01T01:29:00Z"
	url = "http://localhost:9966/events/test"
	jsonData := []byte(stringutils.Concat(`{"timestamp":"`, timestampString, `"}`))

	httputils.PostDataToUrl(url, "application/json", jsonData)
	httputils.PostDataToUrl(url, "application/json", jsonData)
	httputils.PostDataToUrl(url, "application/json", jsonData)

	url = stringutils.Concat("http://localhost:9966/events/test?timestamp=", timestampString)
	data, _ := httputils.GetDataFromUrl(url)

	var eventSummary sieve.EventSummary
	json.Unmarshal(data, &eventSummary)
	c.Assert(eventSummary.ProjectedThisHour, check.Equals, float32(6))
}

func (s *SieveEventsTestSuite) TestSieveCalculatesTrailingAveragePerSecondCorrectly(c *check.C) {
	var url string

	url = "http://localhost:9966/events/test"
	httputils.PostDataToUrl(url, "application/json", []byte(`{"timestamp":"2015-01-01T01:00:00Z"}`))
	httputils.PostDataToUrl(url, "application/json", []byte(`{"timestamp":"2015-01-01T01:00:00Z"}`))
	httputils.PostDataToUrl(url, "application/json", []byte(`{"timestamp":"2015-01-01T01:00:01Z"}`))
	httputils.PostDataToUrl(url, "application/json", []byte(`{"timestamp":"2015-01-01T01:00:01Z"}`))
	httputils.PostDataToUrl(url, "application/json", []byte(`{"timestamp":"2015-01-01T01:00:01Z"}`))
	httputils.PostDataToUrl(url, "application/json", []byte(`{"timestamp":"2015-01-01T01:00:01Z"}`))
	httputils.PostDataToUrl(url, "application/json", []byte(`{"timestamp":"2015-01-01T01:00:02Z"}`))
	httputils.PostDataToUrl(url, "application/json", []byte(`{"timestamp":"2015-01-01T01:00:02Z"}`))
	httputils.PostDataToUrl(url, "application/json", []byte(`{"timestamp":"2015-01-01T01:00:02Z"}`))

	url = "http://localhost:9966/events/test?timestamp=2015-01-01T01:00:03Z"
	data, _ := httputils.GetDataFromUrl(url)

	var eventSummary sieve.EventSummary
	json.Unmarshal(data, &eventSummary)
	c.Assert(eventSummary.TrailingAveragePerSecond, check.Equals, float32(9)/float32(3))
}

func (s *SieveEventsTestSuite) TestSieveCalculatesTrailingAveragePerMinuteCorrectly(c *check.C) {
	var url string

	url = "http://localhost:9966/events/test"
	httputils.PostDataToUrl(url, "application/json", []byte(`{"timestamp":"2015-01-01T01:00:00Z"}`))
	httputils.PostDataToUrl(url, "application/json", []byte(`{"timestamp":"2015-01-01T01:00:00Z"}`))
	httputils.PostDataToUrl(url, "application/json", []byte(`{"timestamp":"2015-01-01T01:01:00Z"}`))
	httputils.PostDataToUrl(url, "application/json", []byte(`{"timestamp":"2015-01-01T01:01:00Z"}`))
	httputils.PostDataToUrl(url, "application/json", []byte(`{"timestamp":"2015-01-01T01:01:00Z"}`))
	httputils.PostDataToUrl(url, "application/json", []byte(`{"timestamp":"2015-01-01T01:01:00Z"}`))
	httputils.PostDataToUrl(url, "application/json", []byte(`{"timestamp":"2015-01-01T01:02:00Z"}`))
	httputils.PostDataToUrl(url, "application/json", []byte(`{"timestamp":"2015-01-01T01:02:00Z"}`))
	httputils.PostDataToUrl(url, "application/json", []byte(`{"timestamp":"2015-01-01T01:02:00Z"}`))

	url = "http://localhost:9966/events/test?timestamp=2015-01-01T01:03:00Z"
	data, _ := httputils.GetDataFromUrl(url)

	var eventSummary sieve.EventSummary
	json.Unmarshal(data, &eventSummary)
	c.Assert(eventSummary.TrailingAveragePerMinute, check.Equals, float32(9)/float32(3))
}

func (s *SieveEventsTestSuite) TestSieveCalculatesTrailingAveragePerHourCorrectly(c *check.C) {
	var url string

	url = "http://localhost:9966/events/test"
	httputils.PostDataToUrl(url, "application/json", []byte(`{"timestamp":"2015-01-01T01:00:00Z"}`))
	httputils.PostDataToUrl(url, "application/json", []byte(`{"timestamp":"2015-01-01T01:00:00Z"}`))
	httputils.PostDataToUrl(url, "application/json", []byte(`{"timestamp":"2015-01-01T01:00:00Z"}`))
	httputils.PostDataToUrl(url, "application/json", []byte(`{"timestamp":"2015-01-01T01:00:00Z"}`))
	httputils.PostDataToUrl(url, "application/json", []byte(`{"timestamp":"2015-01-01T02:00:00Z"}`))
	httputils.PostDataToUrl(url, "application/json", []byte(`{"timestamp":"2015-01-01T02:00:00Z"}`))
	httputils.PostDataToUrl(url, "application/json", []byte(`{"timestamp":"2015-01-01T02:00:00Z"}`))
	httputils.PostDataToUrl(url, "application/json", []byte(`{"timestamp":"2015-01-01T03:00:00Z"}`))
	httputils.PostDataToUrl(url, "application/json", []byte(`{"timestamp":"2015-01-01T03:00:00Z"}`))

	url = "http://localhost:9966/events/test?timestamp=2015-01-01T04:00:00Z"
	data, _ := httputils.GetDataFromUrl(url)

	var eventSummary sieve.EventSummary
	json.Unmarshal(data, &eventSummary)
	c.Assert(eventSummary.TrailingAveragePerHour, check.Equals, float32(9)/float32(3))
}

func (s *SieveEventsTestSuite) TestSieveCalculatesTrailingChangePerSecondCorrectly(c *check.C) {
	var url string

	url = "http://localhost:9966/events/test"
	httputils.PostDataToUrl(url, "application/json", []byte(`{"timestamp":"2015-01-01T01:00:00Z"}`))
	httputils.PostDataToUrl(url, "application/json", []byte(`{"timestamp":"2015-01-01T01:00:00Z"}`))
	httputils.PostDataToUrl(url, "application/json", []byte(`{"timestamp":"2015-01-01T01:00:01Z"}`))
	httputils.PostDataToUrl(url, "application/json", []byte(`{"timestamp":"2015-01-01T01:00:01Z"}`))
	httputils.PostDataToUrl(url, "application/json", []byte(`{"timestamp":"2015-01-01T01:00:01Z"}`))
	httputils.PostDataToUrl(url, "application/json", []byte(`{"timestamp":"2015-01-01T01:00:01Z"}`))

	url = "http://localhost:9966/events/test?timestamp=2015-01-01T01:00:02Z"
	data, _ := httputils.GetDataFromUrl(url)

	var eventSummary sieve.EventSummary
	json.Unmarshal(data, &eventSummary)
	c.Assert(eventSummary.ChangePerSecond, check.Equals, 2)
}

func (s *SieveEventsTestSuite) TestSieveCalculatesTrailingChangePerMinuteCorrectly(c *check.C) {
	var url string

	url = "http://localhost:9966/events/test"
	httputils.PostDataToUrl(url, "application/json", []byte(`{"timestamp":"2015-01-01T01:01:00Z"}`))
	httputils.PostDataToUrl(url, "application/json", []byte(`{"timestamp":"2015-01-01T01:01:00Z"}`))
	httputils.PostDataToUrl(url, "application/json", []byte(`{"timestamp":"2015-01-01T01:01:00Z"}`))
	httputils.PostDataToUrl(url, "application/json", []byte(`{"timestamp":"2015-01-01T01:01:00Z"}`))
	httputils.PostDataToUrl(url, "application/json", []byte(`{"timestamp":"2015-01-01T01:02:00Z"}`))
	httputils.PostDataToUrl(url, "application/json", []byte(`{"timestamp":"2015-01-01T01:02:00Z"}`))

	url = "http://localhost:9966/events/test?timestamp=2015-01-01T01:03:00Z"
	data, _ := httputils.GetDataFromUrl(url)

	var eventSummary sieve.EventSummary
	json.Unmarshal(data, &eventSummary)
	c.Assert(eventSummary.ChangePerMinute, check.Equals, -2)
}

func (s *SieveEventsTestSuite) TestSieveCalculatesTrailingChangePerHourCorrectly(c *check.C) {
	var url string

	url = "http://localhost:9966/events/test"
	httputils.PostDataToUrl(url, "application/json", []byte(`{"timestamp":"2015-01-01T02:00:00Z"}`))
	httputils.PostDataToUrl(url, "application/json", []byte(`{"timestamp":"2015-01-01T02:00:00Z"}`))
	httputils.PostDataToUrl(url, "application/json", []byte(`{"timestamp":"2015-01-01T02:00:00Z"}`))
	httputils.PostDataToUrl(url, "application/json", []byte(`{"timestamp":"2015-01-01T03:00:00Z"}`))
	httputils.PostDataToUrl(url, "application/json", []byte(`{"timestamp":"2015-01-01T03:00:00Z"}`))

	url = "http://localhost:9966/events/test?timestamp=2015-01-01T04:00:00Z"
	data, _ := httputils.GetDataFromUrl(url)

	var eventSummary sieve.EventSummary
	json.Unmarshal(data, &eventSummary)
	c.Assert(eventSummary.ChangePerHour, check.Equals, -1)
}
