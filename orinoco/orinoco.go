package main

import "fmt"
import "os"

import "../sieve"
import "../tap"

func help() {
	fmt.Println("usage: orinoco <command> [<args>]")
	fmt.Println("")
	fmt.Println("These are valid commands:")
	fmt.Println("  sieve   Data stream analysis service with HTTP API")
	fmt.Println("  tap     Client to connect to sieve")
	fmt.Println("")
}

func main() {
	if len(os.Args) == 1 || os.Args[1] == "help" {
		help()
		os.Exit(0)
	}

	if os.Args[1] == "sieve" {
		sieve.Sieve()
	} else if os.Args[1] == "tap" {
		tap.Tap()
	} else {
		help()
	}
}
