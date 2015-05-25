package tap

import (
	"fmt"
	"log"

	"golang.org/x/net/websocket"
	
	"../stringutils"
)

func Tap(port string) {
	origin := "http://localhost/"	
	url := stringutils.Concat("ws://localhost:", port, "/subscribe")

	ws, err := websocket.Dial(url, "", origin)
	if err != nil {
		log.Fatal(err)
	}

	for {
		var msg = make([]byte, 1024*2)
		var n int
		n, err = ws.Read(msg)
		if err == nil {
			fmt.Println(string(msg[:n]))
		}
	}
}
