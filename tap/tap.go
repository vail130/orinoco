package tap

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"

	"github.com/gorilla/websocket"
	"gopkg.in/yaml.v2"

	"github.com/vail130/orinoco/sliceutils"
	"github.com/vail130/orinoco/stringutils"
)

type Config struct {
	Host    string `yaml:"host"`
	Port    string `yaml:"port"`
	Origin  string `yaml:"origin"`
	LogPath string `yaml:"log_path"`
}

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

func readFromSocket(ws *websocket.Conn, logPath string, boundary string) {
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

func Tap(configPath string, boundary string) {
	configData, err := ioutil.ReadFile(configPath)
	if err != nil {
		log.Fatal(err)
	}

	var config Config
	yaml.Unmarshal(configData, &config)

	url := stringutils.Concat("ws://", config.Host, ":", config.Port, "/subscribe")
	headers := make(http.Header)
	headers["origin"] = []string{config.Origin}
	ws, _, err := websocket.DefaultDialer.Dial(url, headers)
	if err != nil {
		log.Fatal(err)
	}

	if len(config.LogPath) > 0 {
		os.MkdirAll(path.Dir(config.LogPath), loggingPermissions)
		os.Create(config.LogPath)
	}

	readFromSocket(ws, config.LogPath, boundary)
}
