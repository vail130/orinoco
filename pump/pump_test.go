package pump_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"gopkg.in/check.v1"

	"github.com/vail130/orinoco/httputils"
	"github.com/vail130/orinoco/sieve"
)

func TestPump(t *testing.T) { check.TestingT(t) }

type PumpTestSuite struct{}

var _ = check.Suite(&PumpTestSuite{})

func (s *PumpTestSuite) SetUpTest(c *check.C) {
	httputils.Delete("http://localhost:9966/events")
}

func (s *PumpTestSuite) TestPumpConsumesLogFile(c *check.C) {
	logPath, _ := filepath.Abs("../pump.log")
	
	if file, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY, 0666); err != nil {
		_, err := file.Write([]byte(`{"event":"test","timestamp":"","data":{"a":1}}\n`))
		file.Close()
		if err != nil {
			c.Assert(err, check.IsNil)
		}
	}

	time.Sleep(time.Second)

	data, _ := httputils.GetDataFromUrl("http://localhost:9966/events/test")

	var eventSummary sieve.EventSummary
	json.Unmarshal(data, &eventSummary)

	c.Assert(eventSummary.Event, check.Equals, "test")
}
