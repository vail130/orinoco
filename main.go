package main

import (
	"gopkg.in/alecthomas/kingpin.v2"

	"github.com/vail130/orinoco/orinoco"
)

var (
	config = kingpin.Arg("config", "Path to configuration file.").Required().String()
)

func main() {
	kingpin.Version("0.0.6")
	kingpin.Parse()
	orinoco.Orinoco(*config)
}
