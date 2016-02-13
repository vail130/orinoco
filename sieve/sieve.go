package sieve

import (
	"os"

	"github.com/vail130/orinoco/stringutils"
)

type SieveConfig struct {
	IsTestEnv bool
	EventChannel chan *Event
}

var sieveConfig SieveConfig

func Sieve(eventChannel chan *Event) {
	sieveConfig = SieveConfig{
		stringutils.StringToBool(os.Getenv("TEST")),
		eventChannel,
	}
}
