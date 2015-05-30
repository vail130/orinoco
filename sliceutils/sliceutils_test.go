package sliceutils_test

import (
	"testing"

	"github.com/vail130/orinoco/sliceutils"

	"gopkg.in/check.v1"
)

func TestSliceUtils(t *testing.T) { check.TestingT(t) }

type SliceUtilsTestSuite struct{}

var _ = check.Suite(&SliceUtilsTestSuite{})

func (s *SliceUtilsTestSuite) TestConcatByteSlicesReturnsCombinedSlices(c *check.C) {
	result := sliceutils.ConcatByteSlices([]byte("Hello "), []byte("to "), []byte("you!"))
	c.Assert(string(result), check.Equals, "Hello to you!")
}
