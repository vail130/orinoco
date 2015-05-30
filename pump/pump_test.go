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
	"github.com/vail130/orinoco/stringutils"
)

func TestPump(t *testing.T) { check.TestingT(t) }

type PumpTestSuite struct{}

var _ = check.Suite(&PumpTestSuite{})

func (s *PumpTestSuite) SetUpTest(c *check.C) {
	httputils.Delete("http://localhost:9966/events")
}

func (s *PumpTestSuite) TestPumpConsumesLogFile(c *check.C) {
	logPath, _ := filepath.Abs("../pump.log")
	now := time.Now()
	jsonString := stringutils.Concat(`{"event":"test","timestamp":"`, now.Format(time.RFC3339), `","data":{"a":1}}`, "\n")
	
	file, err := os.Create(logPath)
	if err != nil {
		c.Assert(err, check.IsNil)
	}
	defer file.Close()
	_, err = file.WriteString(jsonString)

	time.Sleep(time.Second)

	data, _ := httputils.GetDataFromUrl("http://localhost:9966/events/test")

	var eventSummary sieve.EventSummary
	json.Unmarshal(data, &eventSummary)

	c.Assert(eventSummary.Event, check.Equals, "test")
}
