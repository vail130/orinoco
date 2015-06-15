package sieve_test

import (
	"encoding/json"
	"testing"

	"gopkg.in/check.v1"

	"github.com/vail130/orinoco/httputils"
	"github.com/vail130/orinoco/sieve"
	"github.com/vail130/orinoco/stringutils"
)

func TestSieveApi(t *testing.T) { check.TestingT(t) }

type SieveApiTestSuite struct{}

var _ = check.Suite(&SieveApiTestSuite{})

func (s *SieveApiTestSuite) SetUpTest(c *check.C) {
	httputils.Delete(streamUrl)
}

func (s *SieveApiTestSuite) TestSieveReturnsNullWhenNoStreamsHaveBeenPosted(c *check.C) {
	data, _ := httputils.GetDataFromUrl(streamUrl)
	c.Assert(string(data), check.Equals, "null")
}

func (s *SieveApiTestSuite) TestSieveReturnsStreamSummaryAfterReceivingData(c *check.C) {
	jsonData := []byte(`{"a":1}`)
	httputils.PostDataToUrl(testStreamUrl, "application/json", jsonData)

	data, _ := httputils.GetDataFromUrl(testStreamUrl)

	var streamSummary sieve.StreamSummary
	json.Unmarshal(data, &streamSummary)

	c.Assert(streamSummary.Stream, check.Equals, "test")
}

func (s *SieveApiTestSuite) TestSieveReturnsNullForAllStreamsAfterResetting(c *check.C) {
	jsonData := []byte(`{"a":1}`)
	httputils.PostDataToUrl(testStreamUrl, "application/json", jsonData)

	httputils.Delete(streamUrl)

	data, _ := httputils.GetDataFromUrl(streamUrl)
	c.Assert(string(data), check.Equals, "null")
}

func (s *SieveApiTestSuite) TestSieveReturnsNullForIndividualStreamAfterResetting(c *check.C) {
	jsonData := []byte(`{"a":1}`)
	httputils.PostDataToUrl(testStreamUrl, "application/json", jsonData)

	httputils.Delete(streamUrl)

	data, _ := httputils.GetDataFromUrl(stringutils.Concat(streamUrl, "/test"))
	c.Assert(string(data), check.Equals, "null")
}
