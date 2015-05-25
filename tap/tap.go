package tap

import (
	"fmt"
	"log"
	"os"
	"path"

	"golang.org/x/net/websocket"
	
	"../stringutils"
)

func Tap(host string, origin string, logDir string) {
	url := stringutils.Concat(host, "/subscribe")
	
	ws, err := websocket.Dial(url, "", origin)
	if err != nil {
		log.Fatal(err)
	}
	
	fmt.Println("Connecting to Sieve server at", host, "from origin", origin)
	
	var logDestination string
	var loggingPermissions os.FileMode = 0666
	
	if len(logDir) > 0 {
		logDestination = logDir
		os.MkdirAll(logDir, loggingPermissions)
	} else {
		logDestination = "standard out"
	}
	fmt.Println("Logging to", logDestination)

	for {
		var msg = make([]byte, 1024*2)
		var n int
		n, err = ws.Read(msg)
		if err == nil {
			eventDataString := stringutils.Concat(string(msg[:n]), "\n")
			fmt.Print(eventDataString)
			
			if len(logDir) > 0 {
				logFilePath := path.Join(logDir, "orinoco.log")
				if f, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_APPEND|os.O_RDWR, loggingPermissions); err == nil {
					f.WriteString(eventDataString)
					f.Close()
				}
			}
		}
	}
}
