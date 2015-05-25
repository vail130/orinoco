package main

import (
	"os"
	
	"gopkg.in/alecthomas/kingpin.v2"

	"../sieve"
	"../tap"
)

var (
	app = kingpin.New("orinoco", "A data stream monitoring services.")
	
	sieveApp = app.Command("sieve", "Run a data stream stats and pub-sub server.")
	port = sieveApp.Flag("port", "Port for sieve to listen on.").Short('p').Default("9966").String()
	
	tapApp = app.Command("tap", "Run a data stream client.")
	host = tapApp.Flag("host", "Sieve host to connect to.").Short('h').Default("ws://localhost:9966").String()
	origin = tapApp.Flag("origin", "Origin from which to connect to sieve.").Short('o').Default("http://localhost/").String()
	logPath = tapApp.Flag("logpath", "File to log data stream to. Omitting this flag will log to standard out.").Short('l').String()
)

func main() {
	app.Version("0.0.3")
	
	switch kingpin.MustParse(app.Parse(os.Args[1:])) {
		
	case sieveApp.FullCommand():
		sieve.Sieve(*port)
		
	case tapApp.FullCommand():
		tap.Tap(*host, *origin, *logPath)
	}
}
