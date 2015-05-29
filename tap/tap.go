package tap

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"os"
	"path"

	"github.com/gorilla/websocket"

	"github.com/vail130/orinoco/stringutils"
)

var loggingPermissions os.FileMode = 0666

func logMessage(message []byte, logPath string) {
	if len(logPath) > 0 {
		if f, err := os.OpenFile(logPath, os.O_CREATE|os.O_APPEND|os.O_EXCL, loggingPermissions); err == nil {
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
			_, partialMessage, err := ws.ReadMessage()
			if err != nil {
				log.Fatal(err)
			}

			fullMessage = append(fullMessage, partialMessage...)

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
	headers := make(http.Header)
	headers["origin"] = []string{origin}
	ws, _, err := websocket.DefaultDialer.Dial(url, headers)
	if err != nil {
		log.Fatal(err)
	}

	if len(logPath) > 0 {
		os.MkdirAll(path.Dir(logPath), loggingPermissions)
	}

	readFromSocket(ws, boundary, logPath)
}
