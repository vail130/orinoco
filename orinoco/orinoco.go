package main

import (
	"fmt"
	"os"

	"../sieve"
	"../tap"
)

func help() {
	fmt.Println("usage: orinoco <command> [<args>]")
	fmt.Println("")
	fmt.Println("These are valid commands:")
	fmt.Println("  sieve   Data stream stats and pub-sub server")
	fmt.Println("  tap     Data stream client")
	fmt.Println("")
}

func main() {
	if len(os.Args) == 1 || os.Args[1] == "help" {
		help()
		os.Exit(0)
	}

	// TODO Get port from command line flag
	port := "9966"
	if port == "" {
		port = "9966"
	}

	if os.Args[1] == "sieve" {
		sieve.Sieve(port)
	} else if os.Args[1] == "tap" {
		tap.Tap(port)
	} else {
		help()
	}
}
