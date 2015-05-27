package stringutils_test

import (
	"testing"
	
	"github.com/vail130/orinoco/stringutils"
)

func TestConcatReturnsCombinedStrings(t *testing.T) {
	result := stringutils.Concat("a", "b", "c")
	if result != "abc" {
		t.FailNow()
	}
}

func TestUnderscoreToTitleReturnsTransformedString(t *testing.T) {
	result := stringutils.UnderscoreToTitle("a_b_c")
	if result != "ABC" {
		t.FailNow()
	}
}
