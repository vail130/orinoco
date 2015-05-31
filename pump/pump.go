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
	SieveUrl  string            `yaml:"url"`
	SaveFiles string            `yaml:"save_files"`
	Streams   map[string]string `yaml:"streams"`
}

var saveConsumedLogFiles bool

func sendEventOverHttp(url string, data []byte) {
	_, err := httputils.PostDataToUrl(url, "application/json", data)
	if err != nil {
		log.Fatalln(err)
	}
}

func consumeLogs(logPath string, url string) {
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
			sendEventOverHttp(url, messageData)
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

	saveConsumedLogFiles = stringutils.StringToBool(config.SaveFiles)

	for {
		var wg sync.WaitGroup
		for logPath, stream := range config.Streams {
			streamUrl := stringutils.Concat(config.SieveUrl, "/streams/", stream)
			wg.Add(1)
			go func(logPath string, streamUrl string) {
				defer wg.Done()
				consumeLogs(logPath, streamUrl)
			}(logPath, streamUrl)
		}
		time.Sleep(time.Second)
		wg.Wait()
	}
}
