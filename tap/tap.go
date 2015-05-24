package tap

import (
	"bytes"
	"fmt"
	"log"

	"golang.org/x/net/websocket"
)

func Tap(port string) {
	origin := "http://localhost/"

	var buffer bytes.Buffer
	buffer.WriteString("ws://localhost:")
	buffer.WriteString(port)
	buffer.WriteString("/subscribe")
	url := buffer.String()

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
