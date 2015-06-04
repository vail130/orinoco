package orinoco_test

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

func TestOrinoco(t *testing.T) { check.TestingT(t) }

type OrinocoTestSuite struct{}

var _ = check.Suite(&OrinocoTestSuite{})

func (s *OrinocoTestSuite) TestOrinocoConsumesLogFile(c *check.C) {
	pumpTestLitmuslogPath, _ := filepath.Abs("../artifacts/pump-test_litmus.log")
	if f, err := os.Create(pumpTestLitmuslogPath); err != nil {
		f.Close()
	}

	file, err := os.OpenFile(pumpTestLitmuslogPath, os.O_APPEND|os.O_RDWR, 0666)
	if err == nil {
		for i := 0; i < 20; i++ {
			jsonBytes := []byte(stringutils.Concat(`{"a":1}`, "\n"))
			_, err = file.WriteAt(jsonBytes, int64(i*len(jsonBytes)))
		}
	}
	file.Close()

	time.Sleep(time.Second)

	data, _ := httputils.GetDataFromUrl("http://localhost:9966/streams/test_litmus_stream_more_than_10_per_minute")

	var streamSummary sieve.StreamSummary
	json.Unmarshal(data, &streamSummary)

	c.Assert(streamSummary.Stream, check.Equals, "test_litmus_stream_more_than_10_per_minute")
	c.Assert(streamSummary.ProjectedThisHour > 0, check.Equals, true)
}
