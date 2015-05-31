package pump_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strconv"
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
	httputils.Delete("http://localhost:9966/streams")
}

func (s *PumpTestSuite) TestPumpConsumesLogFile(c *check.C) {
	logPath, _ := filepath.Abs("../pump.log")
	if f, err := os.Create(logPath); err != nil {
		f.Close()
	}

	file, err := os.OpenFile(logPath, os.O_APPEND|os.O_RDWR, 0666)
	if err == nil {
		for i := 0; i < 10; i++ {
			jsonBytes := []byte(stringutils.Concat(`{"stream":"test","timestamp":"2015-05-29T21:59:3`, strconv.Itoa(i), `Z","data":{"a":1}}`, "\n"))
			_, err = file.WriteAt(jsonBytes, int64(i*len(jsonBytes)))
		}
	}
	file.Close()

	time.Sleep(1 * time.Second)

	data, _ := httputils.GetDataFromUrl("http://localhost:9966/streams/test?timestamp=2015-05-29T21:59:31Z")

	var streamSummary sieve.StreamSummary
	json.Unmarshal(data, &streamSummary)

	c.Assert(streamSummary.Stream, check.Equals, "test")
	c.Assert(streamSummary.Timestamp, check.Equals, "2015-05-29T21:59:31Z")
	c.Assert(streamSummary.MinuteToDate, check.Equals, 10)
}
