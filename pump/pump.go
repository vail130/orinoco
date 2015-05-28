package pump

import (
	"bufio"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/websocket"
	
	"github.com/vail130/orinoco/stringutils"
)

func sendEventOverSocket(ws *websocket.Conn, data []byte) {
	var err error
	
	for i := 0; i < 3; i++ {
		err = ws.WriteMessage(websocket.TextMessage, data)
		if err != nil {
			log.Println(err)
		} else {
			break
		}
	}
}

func writeLogsToSocket(ws *websocket.Conn, boundary string, logPath string) {
	var err error
	
	now := time.Now()
	consumingPath := stringutils.Concat(logPath, ".", string(now.Unix()), ".consuming")
	consumedPath := stringutils.Concat(logPath, ".", string(now.Unix()), ".consumed")
	err = os.Rename(logPath, consumingPath)
	
	file, err := os.Open(logPath)
	defer file.Close()
	
	if err != nil {
		return
	}
	
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
	    data := scanner.Bytes()
		data = append(data, []byte(boundary)...)
		go sendEventOverSocket(ws, data)
	}
	
	if err = scanner.Err(); err != nil {
	    log.Fatal(err)
		return
	}
	
	err = os.Rename(consumingPath, consumedPath)
}

func Pump(host string, port string, origin string, boundary string, logPath string) {
	url := stringutils.Concat("ws://", host, ":", port, "/publish")
	headers := make(http.Header)
	headers["origin"] = []string{origin}
	ws, _, err := websocket.DefaultDialer.Dial(url, headers)
	if err != nil {
		log.Fatal(err)
	}
	
	for {
		writeLogsToSocket(ws, boundary, logPath)
		time.Sleep(time.Second)
	}
}

