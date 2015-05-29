package pump

import (
	"bufio"
	"fmt"
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
	var err error

	now := time.Now()
	consumingPath := stringutils.Concat(logPath, ".", strconv.Itoa(int(now.Unix())), ".consuming")
	consumedPath := stringutils.Concat(logPath, ".", strconv.Itoa(int(now.Unix())), ".consumed")
	err = os.Rename(logPath, consumingPath)

	file, err := os.Open(logPath)
	defer file.Close()

	if err != nil {
		return
	}

	fmt.Println("PUMP CONNECTED")

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {

		data := scanner.Bytes()
		fmt.Println(stringutils.Concat("PUMPING DATA:", string(data)))
		data = append(data, []byte(boundary)...)
		sendEventOverSocket(ws, data)
	}

	if err = scanner.Err(); err != nil {
		log.Fatal(err)
		return
	}

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
	}
}
