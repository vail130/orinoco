package tap

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"path"

	"golang.org/x/net/websocket"
	
	"github.com/vail130/orinoco/stringutils"
)

var loggingPermissions os.FileMode = 0666

func logMessage(message []byte, logPath string) {
	if len(logPath) > 0 {
		if f, err := os.OpenFile(logPath, os.O_CREATE|os.O_APPEND|os.O_RDWR, loggingPermissions); err == nil {
			f.Write(append(message, []byte("\n")...))
			f.Close()
		}
	} else {
		fmt.Println(string(message))
	}
}

func readFromSocket(ws *websocket.Conn, boundary string, logPath string) {
	boundaryBytes := []byte(boundary)
	var leftoverMessage []byte
	
	for {
		fullMessage := leftoverMessage
		leftoverMessage = make([]byte, 0)
		
		for {
			var partialMessage = make([]byte, 2048)
			var n int
			n, err := ws.Read(partialMessage)
			if err == nil {
				fullMessage = append(fullMessage, partialMessage[:n]...)
			}
			
			if bytes.Index(fullMessage, boundaryBytes) > -1 {
				messagePieces := bytes.SplitN(fullMessage, boundaryBytes, 1)
				fullMessage = messagePieces[0][:len(messagePieces[0])-len(boundaryBytes)]
				if len(messagePieces) > 1 {
					leftoverMessage = messagePieces[1]
				}
				break
			}
		}
		
		logMessage(fullMessage, logPath)
	}
}

func Tap(host string, port string, origin string, boundary string, logPath string) {
	url := stringutils.Concat("ws://", host, ":", port, "/subscribe")
	ws, err := websocket.Dial(url, "", origin)
	if err != nil {
		log.Fatal(err)
	}
	
	if len(logPath) > 0 {
		os.MkdirAll(path.Dir(logPath), loggingPermissions)
	}
	
	readFromSocket(ws, boundary, logPath)
}
