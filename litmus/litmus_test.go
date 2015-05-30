package litmus_test

import (
	"encoding/json"
	"os"
	"path/filepath"
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
	httputils.Delete("http://localhost:9966/events")

	logPath, _ := filepath.Abs("../tap.log")
	os.Remove(logPath)
	os.Create(logPath)
}

func (s *LitmusTestSuite) TestLitmusTriggersCustomEvent(c *check.C) {
	testData := make([][]byte, 0)
	for i := 0; i < 11; i++ {
		jsonBytes := []byte(stringutils.Concat(`{"a":1}`, "\n"))
		httputils.PostDataToUrl("http://localhost:9966/events/test_litmus", "application/json", jsonBytes)
		testData = append(testData, jsonBytes)
	}

	time.Sleep(1 * time.Second)

	data, err := httputils.GetDataFromUrl("http://localhost:9966/events/test_litmus_event_more_than_10_per_minute")
	c.Assert(err, check.IsNil)

	var eventSummary sieve.EventSummary
	json.Unmarshal(data, &eventSummary)

	c.Assert(eventSummary.Event, check.Equals, "test_litmus_event_more_than_10_per_minute")
	c.Assert(eventSummary.SecondToDate, check.Equals, 1)
}
