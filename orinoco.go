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

	pumpApp       = app.Command("pump", "Run a data stream client to pump data to sieve.")
	pumpLogPath   = pumpApp.Flag("logpath", "Log file to consume to pump to sieve.").Short('l').String()
	pumpUrl       = pumpApp.Flag("url", "Sieve endpoint to post events from log file.").Short('u').String()
	pumpConfig    = pumpApp.Flag("config", "Path to configuration file. This overrides other flags").Short('c').String()
	pumpSaveFiles = pumpApp.Flag("save-files", "Disable removal of consumed log files").Default("0").String()

	sieveApp      = app.Command("sieve", "Run a data stream stats and pub-sub server.")
	sievePort     = sieveApp.Flag("port", "Port for sieve to listen on.").Short('p').Default("9966").String()
	sieveBoundary = sieveApp.Flag("boundary", "Designated boundary between messages.").Short('b').Default(orinocoMessageBoundary).String()

	tapApp      = app.Command("tap", "Run a data stream client to subscribe to sieve.")
	tapHost     = tapApp.Flag("host", "Sieve host to connect to.").Short('h').Default("localhost").String()
	tapPort     = tapApp.Flag("port", "Port to use to connect to sieve.").Short('p').Default("9966").String()
	tapOrigin   = tapApp.Flag("origin", "Origin from which to connect to sieve.").Short('o').Default("http://localhost/").String()
	tapBoundary = tapApp.Flag("boundary", "Designated boundary between messages.").Short('b').Default(orinocoMessageBoundary).String()
	tapLogPath  = tapApp.Flag("logpath", "File to log data stream to. Omitting this flag will log to standard out.").Short('l').String()

	litmusApp    = app.Command("litmus", "Run a data stream monitoring daemon.")
	litmusUrl    = litmusApp.Flag("url", "Sieve host url.").Short('u').Default("http://localhost:9966").String()
	litmusConfig = litmusApp.Flag("config", "Path to configuration file.").Short('c').String()
)

func main() {
	app.Version("0.0.3")

	switch kingpin.MustParse(app.Parse(os.Args[1:])) {

	case pumpApp.FullCommand():
		pump.Pump(*pumpLogPath, *pumpUrl, *pumpConfig, *pumpSaveFiles)

	case sieveApp.FullCommand():
		sieve.Sieve(*sievePort, *sieveBoundary)

	case tapApp.FullCommand():
		tap.Tap(*tapHost, *tapPort, *tapOrigin, *tapBoundary, *tapLogPath)

	case litmusApp.FullCommand():
		litmus.Litmus(*litmusUrl, *litmusConfig)
	}
}
