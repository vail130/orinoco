package sieve

import (
	"fmt"
	"net/http"
)

func EventHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Event")
}
