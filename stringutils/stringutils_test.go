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

func TestStringToBoolWorks(t *testing.T) {
	var result bool

	result = stringutils.StringToBool("")
	if result {
		t.FailNow()
	}

	result = stringutils.StringToBool("f")
	if result {
		t.FailNow()
	}

	result = stringutils.StringToBool("false")
	if result {
		t.FailNow()
	}

	result = stringutils.StringToBool("0")
	if result {
		t.FailNow()
	}

	result = stringutils.StringToBool("t")
	if !result {
		t.FailNow()
	}

	result = stringutils.StringToBool("true")
	if !result {
		t.FailNow()
	}

	result = stringutils.StringToBool("1")
	if !result {
		t.FailNow()
	}

	result = stringutils.StringToBool("asdf")
	if !result {
		t.FailNow()
	}
}
