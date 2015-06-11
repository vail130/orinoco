package litmus_test

import (
	"encoding/json"
	"testing"
	"time"

	"gopkg.in/check.v1"

	"github.com/vail130/orinoco/httputils"
	"github.com/vail130/orinoco/sieve"
)

func TestLitmus(t *testing.T) { check.TestingT(t) }

type LitmusTestSuite struct{}

var _ = check.Suite(&LitmusTestSuite{})

func (s *LitmusTestSuite) SetUpTest(c *check.C) {
	httputils.Delete("http://localhost:9966/streams")
}

func (s *LitmusTestSuite) TestLitmusTriggersCustomStream(c *check.C) {
	testData := make([][]byte, 0)
	for i := 0; i < 11; i++ {
		jsonBytes := []byte(`{"a":1}`)
		httputils.PostDataToUrl("http://localhost:9966/streams/test2", "application/json", jsonBytes)
		testData = append(testData, jsonBytes)
	}

	time.Sleep(1 * time.Second)

	data, err := httputils.GetDataFromUrl("http://localhost:9966/streams/test2_stream_more_than_10_per_minute")
	c.Assert(err, check.IsNil)

	var streamSummary sieve.StreamSummary
	json.Unmarshal(data, &streamSummary)

	c.Assert(streamSummary.Stream, check.Equals, "test2_stream_more_than_10_per_minute")
	c.Assert(streamSummary.MinuteToDate > 0, check.Equals, true)
}
