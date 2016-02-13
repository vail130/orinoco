package litmus_test

import (
	"encoding/json"
	"testing"
	"time"

	"gopkg.in/check.v1"

	"github.com/vail130/orinoco/httputils"
	"github.com/vail130/orinoco/sieve"
	"github.com/vail130/orinoco/stringutils"
)

func TestLitmus(t *testing.T) { check.TestingT(t) }

type LitmusTestSuite struct{}

var _ = check.Suite(&LitmusTestSuite{})

func (s *LitmusTestSuite) SetUpTest(c *check.C) {
	httputils.Delete("http://localhost:9966/streams/")
}

func (s *LitmusTestSuite) TestLitmusTriggersCustomStream(c *check.C) {
	timestamp := "2015-01-01T01:00:42Z"
	jsonBytes := []byte(stringutils.Concat(`{"a":1,"timestamp":"`, timestamp, `"}`))

	// 11 times
	httputils.PostDataToUrl("http://localhost:9966/streams/test2/", "application/json", jsonBytes)
	httputils.PostDataToUrl("http://localhost:9966/streams/test2/", "application/json", jsonBytes)
	httputils.PostDataToUrl("http://localhost:9966/streams/test2/", "application/json", jsonBytes)
	httputils.PostDataToUrl("http://localhost:9966/streams/test2/", "application/json", jsonBytes)
	httputils.PostDataToUrl("http://localhost:9966/streams/test2/", "application/json", jsonBytes)
	httputils.PostDataToUrl("http://localhost:9966/streams/test2/", "application/json", jsonBytes)
	httputils.PostDataToUrl("http://localhost:9966/streams/test2/", "application/json", jsonBytes)
	httputils.PostDataToUrl("http://localhost:9966/streams/test2/", "application/json", jsonBytes)
	httputils.PostDataToUrl("http://localhost:9966/streams/test2/", "application/json", jsonBytes)
	httputils.PostDataToUrl("http://localhost:9966/streams/test2/", "application/json", jsonBytes)
	httputils.PostDataToUrl("http://localhost:9966/streams/test2/", "application/json", jsonBytes)

	httputils.PutDataToUrl(stringutils.Concat("http://localhost:9966/litmus/triggers/evaluate/?timestamp=", timestamp), "application/json", jsonBytes)

	time.Sleep(time.Second + time.Millisecond * 250)

	data, err := httputils.GetDataFromUrl("http://localhost:9966/streams/test2_stream_more_than_10_per_minute/")
	c.Assert(err, check.IsNil)

	var streamSummary sieve.StreamSummary
	json.Unmarshal(data, &streamSummary)

	c.Assert(streamSummary.Stream, check.Equals, "test2_stream_more_than_10_per_minute")
	c.Assert(streamSummary.MinuteToDate > 0, check.Equals, true)
}
