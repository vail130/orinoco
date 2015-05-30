package pump

import (
	"bufio"
	"log"
	"net/http"
	"os"
	"strconv"
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
	now := time.Now()
	unixTimeStamp := strconv.FormatInt(now.Unix(), 10)
	
	consumingPath := stringutils.Concat(logPath, ".", unixTimeStamp, ".consuming")
	err := os.Rename(logPath, consumingPath)
	if err != nil {
		return
	}
	
	if file, err := os.OpenFile(consumingPath, os.O_RDONLY, 0666); err == nil {
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			data := scanner.Bytes()
			data = append(data, []byte(boundary)...)
			go sendEventOverSocket(ws, data)
		}
	
		if err = scanner.Err(); err != nil {
			log.Println(err)
		}
		
		file.Close()
	}
	
	consumedPath := stringutils.Concat(logPath, ".", unixTimeStamp, ".consumed")
	os.Rename(consumingPath, consumedPath)
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
