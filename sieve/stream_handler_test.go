package sieve_test

import (
	"encoding/json"
	"testing"

	"gopkg.in/check.v1"

	"github.com/vail130/orinoco/httputils"
	"github.com/vail130/orinoco/sieve"
	"github.com/vail130/orinoco/stringutils"
)

func TestSieveStreams(t *testing.T) { check.TestingT(t) }

type SieveStreamsTestSuite struct{}

var _ = check.Suite(&SieveStreamsTestSuite{})

const streamUrl = "http://localhost:9966/streams/"
const testStreamUrl = "http://localhost:9966/streams/test/"

func (s *SieveStreamsTestSuite) SetUpTest(c *check.C) {
	httputils.Delete(streamUrl)
}

func (s *SieveStreamsTestSuite) TestSieveUsesTimestampProvidedForTestStream(c *check.C) {
	var url string

	timestampString := "2015-05-27T09:29:13Z"

	jsonData := []byte(stringutils.Concat(`{"timestamp":"`, timestampString, `"}`))
	httputils.PostDataToUrl(testStreamUrl, "application/json", jsonData)

	url = stringutils.Concat(testStreamUrl, "?timestamp=", timestampString)
	data, _ := httputils.GetDataFromUrl(url)

	var streamSummary sieve.StreamSummary
	json.Unmarshal(data, &streamSummary)

	c.Assert(streamSummary.Timestamp, check.Equals, timestampString)
}

func (s *SieveStreamsTestSuite) TestSieveCalculatesSecondToDateCorrectly(c *check.C) {
	timestampString := "2015-01-01T01:00:42Z"
	jsonData := []byte(stringutils.Concat(`{"timestamp":"`, timestampString, `"}`))

	httputils.PostDataToUrl(testStreamUrl, "application/json", jsonData)
	httputils.PostDataToUrl(testStreamUrl, "application/json", jsonData)
	httputils.PostDataToUrl(testStreamUrl, "application/json", jsonData)

	url := stringutils.Concat(testStreamUrl, "?timestamp=", timestampString)
	data, _ := httputils.GetDataFromUrl(url)

	var streamSummary sieve.StreamSummary
	json.Unmarshal(data, &streamSummary)
	c.Assert(streamSummary.SecondToDate, check.Equals, 3)
}

func (s *SieveStreamsTestSuite) TestSieveCalculatesMinuteToDateCorrectly(c *check.C) {
	httputils.PostDataToUrl(testStreamUrl, "application/json", []byte(`{"timestamp":"2015-01-01T01:00:01Z"}`))
	httputils.PostDataToUrl(testStreamUrl, "application/json", []byte(`{"timestamp":"2015-01-01T01:00:15Z"}`))
	httputils.PostDataToUrl(testStreamUrl, "application/json", []byte(`{"timestamp":"2015-01-01T01:00:22Z"}`))
	httputils.PostDataToUrl(testStreamUrl, "application/json", []byte(`{"timestamp":"2015-01-01T01:00:38Z"}`))
	httputils.PostDataToUrl(testStreamUrl, "application/json", []byte(`{"timestamp":"2015-01-01T01:00:42Z"}`))

	url := stringutils.Concat(testStreamUrl, "?timestamp=2015-01-01T01:00:50Z")
	data, _ := httputils.GetDataFromUrl(url)

	var streamSummary sieve.StreamSummary
	json.Unmarshal(data, &streamSummary)
	c.Assert(streamSummary.MinuteToDate, check.Equals, 5)
}

func (s *SieveStreamsTestSuite) TestSieveCalculatesHourToDateCorrectly(c *check.C) {
	httputils.PostDataToUrl(testStreamUrl, "application/json", []byte(`{"timestamp":"2015-01-01T01:04:00Z"}`))
	httputils.PostDataToUrl(testStreamUrl, "application/json", []byte(`{"timestamp":"2015-01-01T01:15:00Z"}`))
	httputils.PostDataToUrl(testStreamUrl, "application/json", []byte(`{"timestamp":"2015-01-01T01:22:00Z"}`))
	httputils.PostDataToUrl(testStreamUrl, "application/json", []byte(`{"timestamp":"2015-01-01T01:38:00Z"}`))
	httputils.PostDataToUrl(testStreamUrl, "application/json", []byte(`{"timestamp":"2015-01-01T01:42:00Z"}`))

	url := stringutils.Concat(testStreamUrl, "?timestamp=2015-01-01T01:50:00Z")
	data, _ := httputils.GetDataFromUrl(url)

	var streamSummary sieve.StreamSummary
	json.Unmarshal(data, &streamSummary)
	c.Assert(streamSummary.HourToDate, check.Equals, 5)
}

func (s *SieveStreamsTestSuite) TestSieveCalculatesProjectedThisMinuteCorrectly(c *check.C) {
	timestampString := "2015-01-01T01:00:29Z"
	jsonData := []byte(stringutils.Concat(`{"timestamp":"`, timestampString, `"}`))

	httputils.PostDataToUrl(testStreamUrl, "application/json", jsonData)
	httputils.PostDataToUrl(testStreamUrl, "application/json", jsonData)
	httputils.PostDataToUrl(testStreamUrl, "application/json", jsonData)

	url := stringutils.Concat(testStreamUrl, "?timestamp=", timestampString)
	data, _ := httputils.GetDataFromUrl(url)

	var streamSummary sieve.StreamSummary
	json.Unmarshal(data, &streamSummary)
	c.Assert(streamSummary.ProjectedThisMinute, check.Equals, float32(6))
}

func (s *SieveStreamsTestSuite) TestSieveCalculatesProjectedThisHourCorrectly(c *check.C) {
	timestampString := "2015-01-01T01:29:00Z"
	jsonData := []byte(stringutils.Concat(`{"timestamp":"`, timestampString, `"}`))

	httputils.PostDataToUrl(testStreamUrl, "application/json", jsonData)
	httputils.PostDataToUrl(testStreamUrl, "application/json", jsonData)
	httputils.PostDataToUrl(testStreamUrl, "application/json", jsonData)

	url := stringutils.Concat(testStreamUrl, "?timestamp=", timestampString)
	data, _ := httputils.GetDataFromUrl(url)

	var streamSummary sieve.StreamSummary
	json.Unmarshal(data, &streamSummary)
	c.Assert(streamSummary.ProjectedThisHour, check.Equals, float32(6))
}

func (s *SieveStreamsTestSuite) TestSieveCalculatesTrailingAveragePerSecondCorrectly(c *check.C) {
	httputils.PostDataToUrl(testStreamUrl, "application/json", []byte(`{"timestamp":"2015-01-01T01:00:00Z"}`))
	httputils.PostDataToUrl(testStreamUrl, "application/json", []byte(`{"timestamp":"2015-01-01T01:00:00Z"}`))
	httputils.PostDataToUrl(testStreamUrl, "application/json", []byte(`{"timestamp":"2015-01-01T01:00:01Z"}`))
	httputils.PostDataToUrl(testStreamUrl, "application/json", []byte(`{"timestamp":"2015-01-01T01:00:01Z"}`))
	httputils.PostDataToUrl(testStreamUrl, "application/json", []byte(`{"timestamp":"2015-01-01T01:00:01Z"}`))
	httputils.PostDataToUrl(testStreamUrl, "application/json", []byte(`{"timestamp":"2015-01-01T01:00:01Z"}`))
	httputils.PostDataToUrl(testStreamUrl, "application/json", []byte(`{"timestamp":"2015-01-01T01:00:02Z"}`))
	httputils.PostDataToUrl(testStreamUrl, "application/json", []byte(`{"timestamp":"2015-01-01T01:00:02Z"}`))
	httputils.PostDataToUrl(testStreamUrl, "application/json", []byte(`{"timestamp":"2015-01-01T01:00:02Z"}`))

	url := stringutils.Concat(testStreamUrl, "?timestamp=2015-01-01T01:00:03Z")
	data, _ := httputils.GetDataFromUrl(url)

	var streamSummary sieve.StreamSummary
	json.Unmarshal(data, &streamSummary)
	c.Assert(streamSummary.TrailingAveragePerSecond, check.Equals, float32(9)/float32(3))
}

func (s *SieveStreamsTestSuite) TestSieveCalculatesTrailingAveragePerMinuteCorrectly(c *check.C) {
	httputils.PostDataToUrl(testStreamUrl, "application/json", []byte(`{"timestamp":"2015-01-01T01:00:00Z"}`))
	httputils.PostDataToUrl(testStreamUrl, "application/json", []byte(`{"timestamp":"2015-01-01T01:00:00Z"}`))
	httputils.PostDataToUrl(testStreamUrl, "application/json", []byte(`{"timestamp":"2015-01-01T01:01:00Z"}`))
	httputils.PostDataToUrl(testStreamUrl, "application/json", []byte(`{"timestamp":"2015-01-01T01:01:00Z"}`))
	httputils.PostDataToUrl(testStreamUrl, "application/json", []byte(`{"timestamp":"2015-01-01T01:01:00Z"}`))
	httputils.PostDataToUrl(testStreamUrl, "application/json", []byte(`{"timestamp":"2015-01-01T01:01:00Z"}`))
	httputils.PostDataToUrl(testStreamUrl, "application/json", []byte(`{"timestamp":"2015-01-01T01:02:00Z"}`))
	httputils.PostDataToUrl(testStreamUrl, "application/json", []byte(`{"timestamp":"2015-01-01T01:02:00Z"}`))
	httputils.PostDataToUrl(testStreamUrl, "application/json", []byte(`{"timestamp":"2015-01-01T01:02:00Z"}`))

	url := stringutils.Concat(testStreamUrl, "?timestamp=2015-01-01T01:03:00Z")
	data, _ := httputils.GetDataFromUrl(url)

	var streamSummary sieve.StreamSummary
	json.Unmarshal(data, &streamSummary)
	c.Assert(streamSummary.TrailingAveragePerMinute, check.Equals, float32(9)/float32(3))
}

func (s *SieveStreamsTestSuite) TestSieveCalculatesTrailingAveragePerHourCorrectly(c *check.C) {
	httputils.PostDataToUrl(testStreamUrl, "application/json", []byte(`{"timestamp":"2015-01-01T01:00:00Z"}`))
	httputils.PostDataToUrl(testStreamUrl, "application/json", []byte(`{"timestamp":"2015-01-01T01:00:00Z"}`))
	httputils.PostDataToUrl(testStreamUrl, "application/json", []byte(`{"timestamp":"2015-01-01T01:00:00Z"}`))
	httputils.PostDataToUrl(testStreamUrl, "application/json", []byte(`{"timestamp":"2015-01-01T01:00:00Z"}`))
	httputils.PostDataToUrl(testStreamUrl, "application/json", []byte(`{"timestamp":"2015-01-01T02:00:00Z"}`))
	httputils.PostDataToUrl(testStreamUrl, "application/json", []byte(`{"timestamp":"2015-01-01T02:00:00Z"}`))
	httputils.PostDataToUrl(testStreamUrl, "application/json", []byte(`{"timestamp":"2015-01-01T02:00:00Z"}`))
	httputils.PostDataToUrl(testStreamUrl, "application/json", []byte(`{"timestamp":"2015-01-01T03:00:00Z"}`))
	httputils.PostDataToUrl(testStreamUrl, "application/json", []byte(`{"timestamp":"2015-01-01T03:00:00Z"}`))

	url := stringutils.Concat(testStreamUrl, "?timestamp=2015-01-01T04:00:00Z")
	data, _ := httputils.GetDataFromUrl(url)

	var streamSummary sieve.StreamSummary
	json.Unmarshal(data, &streamSummary)
	c.Assert(streamSummary.TrailingAveragePerHour, check.Equals, float32(9)/float32(3))
}

func (s *SieveStreamsTestSuite) TestSieveCalculatesTrailingChangePerSecondCorrectly(c *check.C) {
	httputils.PostDataToUrl(testStreamUrl, "application/json", []byte(`{"timestamp":"2015-01-01T01:00:00Z"}`))
	httputils.PostDataToUrl(testStreamUrl, "application/json", []byte(`{"timestamp":"2015-01-01T01:00:00Z"}`))
	httputils.PostDataToUrl(testStreamUrl, "application/json", []byte(`{"timestamp":"2015-01-01T01:00:01Z"}`))
	httputils.PostDataToUrl(testStreamUrl, "application/json", []byte(`{"timestamp":"2015-01-01T01:00:01Z"}`))
	httputils.PostDataToUrl(testStreamUrl, "application/json", []byte(`{"timestamp":"2015-01-01T01:00:01Z"}`))
	httputils.PostDataToUrl(testStreamUrl, "application/json", []byte(`{"timestamp":"2015-01-01T01:00:01Z"}`))

	url := stringutils.Concat(testStreamUrl, "?timestamp=2015-01-01T01:00:02Z")
	data, _ := httputils.GetDataFromUrl(url)

	var streamSummary sieve.StreamSummary
	json.Unmarshal(data, &streamSummary)
	c.Assert(streamSummary.ChangePerSecond, check.Equals, 2)
}

func (s *SieveStreamsTestSuite) TestSieveCalculatesTrailingChangePerMinuteCorrectly(c *check.C) {
	httputils.PostDataToUrl(testStreamUrl, "application/json", []byte(`{"timestamp":"2015-01-01T01:01:00Z"}`))
	httputils.PostDataToUrl(testStreamUrl, "application/json", []byte(`{"timestamp":"2015-01-01T01:01:00Z"}`))
	httputils.PostDataToUrl(testStreamUrl, "application/json", []byte(`{"timestamp":"2015-01-01T01:01:00Z"}`))
	httputils.PostDataToUrl(testStreamUrl, "application/json", []byte(`{"timestamp":"2015-01-01T01:01:00Z"}`))
	httputils.PostDataToUrl(testStreamUrl, "application/json", []byte(`{"timestamp":"2015-01-01T01:02:00Z"}`))
	httputils.PostDataToUrl(testStreamUrl, "application/json", []byte(`{"timestamp":"2015-01-01T01:02:00Z"}`))

	url := stringutils.Concat(testStreamUrl, "?timestamp=2015-01-01T01:03:00Z")
	data, _ := httputils.GetDataFromUrl(url)

	var streamSummary sieve.StreamSummary
	json.Unmarshal(data, &streamSummary)
	c.Assert(streamSummary.ChangePerMinute, check.Equals, -2)
}

func (s *SieveStreamsTestSuite) TestSieveCalculatesTrailingChangePerHourCorrectly(c *check.C) {
	httputils.PostDataToUrl(testStreamUrl, "application/json", []byte(`{"timestamp":"2015-01-01T02:00:00Z"}`))
	httputils.PostDataToUrl(testStreamUrl, "application/json", []byte(`{"timestamp":"2015-01-01T02:00:00Z"}`))
	httputils.PostDataToUrl(testStreamUrl, "application/json", []byte(`{"timestamp":"2015-01-01T02:00:00Z"}`))
	httputils.PostDataToUrl(testStreamUrl, "application/json", []byte(`{"timestamp":"2015-01-01T03:00:00Z"}`))
	httputils.PostDataToUrl(testStreamUrl, "application/json", []byte(`{"timestamp":"2015-01-01T03:00:00Z"}`))

	url := stringutils.Concat(testStreamUrl, "?timestamp=2015-01-01T04:00:00Z")
	data, _ := httputils.GetDataFromUrl(url)

	var streamSummary sieve.StreamSummary
	json.Unmarshal(data, &streamSummary)
	c.Assert(streamSummary.ChangePerHour, check.Equals, -1)
}
