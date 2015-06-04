package main

import (
	"os"

	"gopkg.in/alecthomas/kingpin.v2"

	"github.com/vail130/orinoco/litmus"
	"github.com/vail130/orinoco/orinoco"
	"github.com/vail130/orinoco/pump"
	"github.com/vail130/orinoco/sieve"
	"github.com/vail130/orinoco/tap"
)

const orinocoMessageBoundary = "____OrInOcO____"

var (
	app = kingpin.New("orinoco", "A data stream monitoring services.")

	sieveApp  = app.Command("sieve", "A stats and pub-sub server.")
	sievePort = sieveApp.Flag("port", "Port for sieve to listen on.").Short('p').Default("9966").String()

	tapApp    = app.Command("tap", "A client that consumes data streams from Sieve.")
	tapConfig = tapApp.Flag("config", "Path to configuration file.").Short('c').String()

	pumpApp    = app.Command("pump", "A client that feeds data streams to Sieve.")
	pumpConfig = pumpApp.Flag("config", "Path to configuration file").Short('c').String()

	litmusApp    = app.Command("litmus", "A client that monitors data stream statistics through Sieve.")
	litmusConfig = litmusApp.Flag("config", "Path to configuration file.").Short('c').String()

	orinocoApp    = app.Command("run", "Runs all services in-process.")
	orinocoConfig = orinocoApp.Flag("config", "Path to configuration file.").Short('c').String()
)

func main() {
	app.Version("0.0.5")

	switch kingpin.MustParse(app.Parse(os.Args[1:])) {

	case sieveApp.FullCommand():
		sieve.Sieve(*sievePort, orinocoMessageBoundary)

	case tapApp.FullCommand():
		tap.Tap(*tapConfig, orinocoMessageBoundary)

	case pumpApp.FullCommand():
		pump.Pump(*pumpConfig)

	case litmusApp.FullCommand():
		litmus.Litmus(*litmusConfig)

	case orinocoApp.FullCommand():
		orinoco.Orinoco(*orinocoConfig)
	}
}
