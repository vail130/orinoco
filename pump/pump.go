package pump

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"path"

	"golang.org/x/net/websocket"
	
	"github.com/vail130/orinoco/stringutils"
)

func writeLogsToSocket(ws *websocket.Conn, boundary string, logPath string) {
	for {
		for {
			// TODO get message from logs
		}
		
		// TODO write message to socket
	}
}

func Pump(host string, port string, origin string, boundary string, logPath string) {
	url := stringutils.Concat("ws://", host, ":", port, "/publish")
	ws, err := websocket.Dial(url, "", origin)
	if err != nil {
		log.Fatal(err)
	}
	
	writeLogsToSocket(ws, boundary, logPath)
}

