package main

import (
	"os"

	"gopkg.in/alecthomas/kingpin.v2"

	"github.com/vail130/orinoco/litmus"
	"github.com/vail130/orinoco/pump"
	"github.com/vail130/orinoco/sieve"
	"github.com/vail130/orinoco/tap"
)

const orinocoMessageBoundary = "____OrInOcO____"

var (
	app = kingpin.New("orinoco", "A data stream monitoring services.")

	sieveApp  = app.Command("sieve", "Run a data stream stats and pub-sub server.")
	sievePort = sieveApp.Flag("port", "Port for sieve to listen on.").Short('p').Default("9966").String()

	tapApp    = app.Command("tap", "Run a data stream client to subscribe to sieve.")
	tapConfig = tapApp.Flag("config", "Path to configuration file.").Short('c').String()

	pumpApp    = app.Command("pump", "Run a data stream client to pump data to sieve.")
	pumpConfig = pumpApp.Flag("config", "Path to configuration file. This overrides other flags").Short('c').String()

	litmusApp    = app.Command("litmus", "Run a data stream monitoring daemon.")
	litmusConfig = litmusApp.Flag("config", "Path to configuration file.").Short('c').String()
)

func main() {
	app.Version("0.0.4")

	switch kingpin.MustParse(app.Parse(os.Args[1:])) {

	case sieveApp.FullCommand():
		sieve.Sieve(*sievePort, orinocoMessageBoundary)

	case tapApp.FullCommand():
		tap.Tap(*tapConfig, orinocoMessageBoundary)

	case pumpApp.FullCommand():
		pump.Pump(*pumpConfig)

	case litmusApp.FullCommand():
		litmus.Litmus(*litmusConfig)
	}
}
