package pump_test

import (
	"testing"
	
	"gopkg.in/check.v1"
	
	"github.com/vail130/orinoco/httputils"
)

func TestPump(t *testing.T) { check.TestingT(t) }

type PumpTestSuite struct{}
var _ = check.Suite(&PumpTestSuite{})

func (s *PumpTestSuite) SetUpTest(c *check.C) {
	httputils.Delete("http://localhost:9966/events")
}

func (s *PumpTestSuite) TestPumpConsumesLogFile(c *check.C) {
    
}
