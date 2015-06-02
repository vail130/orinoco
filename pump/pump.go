package pump

import (
	"bufio"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"sync"
	"time"

	"gopkg.in/yaml.v2"

	"github.com/vail130/orinoco/httputils"
	"github.com/vail130/orinoco/stringutils"
)

type Config struct {
	Host      string            `yaml:"host"`
	Port      string            `yaml:"port"`
	SaveFiles string            `yaml:"save_files"`
	Streams   map[string]string `yaml:"streams"`
}

type StreamHandler func(string, []byte)

func ConsumeLogs(logPath string, streamName string, streamHandler StreamHandler, saveConsumedLogFiles bool) {
	now := time.Now()
	unixTimeStamp := strconv.FormatInt(now.Unix(), 10)
	base64UUID, err := stringutils.GetBase64UUID()
	if err != nil {
		log.Fatalln(err)
	}

	uniquePath := stringutils.Concat(logPath, ".", unixTimeStamp, ".", base64UUID)
	consumingPath := stringutils.Concat(uniquePath, ".consuming")
	err = os.Rename(logPath, consumingPath)
	if err != nil {
		if !os.IsNotExist(err) {
			log.Fatalln(err)
		}
		return
	}

	file, err := os.OpenFile(consumingPath, os.O_RDONLY, 0666)
	if err != nil {
		log.Fatalln(err)
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		messageData := scanner.Bytes()
		if string(messageData) != "null" {
			streamHandler(streamName, messageData)
		}
	}

	if err = scanner.Err(); err != nil {
		log.Fatalln(err)
	}

	file.Close()

	if !saveConsumedLogFiles {
		os.Remove(consumingPath)
	}
}

func Pump(configPath string) {
	configData, err := ioutil.ReadFile(configPath)
	if err != nil {
		log.Fatal(err)
	}

	var config Config
	yaml.Unmarshal(configData, &config)

	sieveUrl := stringutils.Concat("http://", config.Host, ":", config.Port)
	saveConsumedLogFiles := stringutils.StringToBool(config.SaveFiles)

	streamHandler := func(streamName string, data []byte) {
		streamUrl := stringutils.Concat(sieveUrl, "/streams/", streamName)
		_, err := httputils.PostDataToUrl(streamUrl, "application/json", data)
		if err != nil {
			log.Fatalln(err)
		}
	}

	for {
		var wg sync.WaitGroup
		for logPath, streamName := range config.Streams {
			wg.Add(1)
			go func(logPath string, streamName string, streamHandler StreamHandler, saveConsumedLogFiles bool) {
				defer wg.Done()
				ConsumeLogs(logPath, streamName, streamHandler, saveConsumedLogFiles)
			}(logPath, streamName, streamHandler, saveConsumedLogFiles)
		}
		time.Sleep(time.Second)
		wg.Wait()
	}
}
