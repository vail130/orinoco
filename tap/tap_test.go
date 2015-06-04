package tap_test

import (
	"bufio"
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

func TestTap(t *testing.T) { check.TestingT(t) }

type TapTestSuite struct{}

var _ = check.Suite(&TapTestSuite{})

func (s *TapTestSuite) SetUpTest(c *check.C) {
	httputils.Delete("http://localhost:9966/streams")

	logPath, _ := filepath.Abs("../artifacts/tap.log")
	os.Remove(logPath)
	os.Create(logPath)
}

func (s *TapTestSuite) TestTapOutputsDataStreamLogFile(c *check.C) {
	testData := make([][]byte, 0)
	for i := 0; i < 10; i++ {
		jsonBytes := []byte(stringutils.Concat(`{"a":1}`, "\n"))
		httputils.PostDataToUrl("http://localhost:9966/streams/test", "application/json", jsonBytes)
		testData = append(testData, jsonBytes)
	}

	time.Sleep(1 * time.Second)

	logPath, _ := filepath.Abs("../artifacts/tap.log")
	file, err := os.OpenFile(logPath, os.O_RDONLY, 0666)
	defer file.Close()
	c.Assert(err, check.IsNil)

	scanner := bufio.NewScanner(file)
	i := 0
	for scanner.Scan() {
		messageData := scanner.Bytes()
		if i < len(testData) {
			var stream sieve.Stream
			json.Unmarshal(messageData, &stream)
			c.Assert(stream.Data, check.Equals, string(testData[i]))
		}
		i++
	}
}
