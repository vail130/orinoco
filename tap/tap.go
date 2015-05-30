package tap

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path"

	"github.com/gorilla/websocket"

	"github.com/vail130/orinoco/sliceutils"
	"github.com/vail130/orinoco/stringutils"
)

var loggingPermissions os.FileMode = 0666

func logMessage(logPath string, message []byte) {
	if len(logPath) > 0 {
		file, err := os.OpenFile(logPath, os.O_WRONLY|os.O_APPEND, loggingPermissions)
		if err != nil {
			log.Fatalln(err)
			return
		}

		message = sliceutils.ConcatByteSlices(message, []byte("\n"))
		file.Write(message)
		file.Close()
	} else {
		fmt.Println(string(message))
	}
}

func readFromSocket(ws *websocket.Conn, boundary string, logPath string) {
	boundaryBytes := []byte(boundary)

	for {
		_, message, err := ws.ReadMessage()
		if err != nil {
			log.Fatalln(err)
			break
		}

		message = message[:len(message)-len(boundaryBytes)]
		logMessage(logPath, message)
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
		os.Create(logPath)
	}

	readFromSocket(ws, boundary, logPath)
}
