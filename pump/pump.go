package pump

import (
	"io/ioutil"
	"log"
	"sync"
	"time"

	"gopkg.in/yaml.v2"

	"github.com/vail130/orinoco/funcutils"
	"github.com/vail130/orinoco/httputils"
	"github.com/vail130/orinoco/stringutils"
)

type Config struct {
	Host      string            `yaml:"host"`
	Port      string            `yaml:"port"`
	SaveFiles string            `yaml:"save_files"`
	Streams   map[string]string `yaml:"streams"`
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
			logStreamer := LogStreamer{
				streamName,
				saveConsumedLogFiles,
				logPath,
				"",
			}

			go funcutils.ExecuteWithWaitGroup(&wg, func() { logStreamer.ConsumeLogs(streamHandler) })
		}
		time.Sleep(time.Second)
		wg.Wait()
	}
}
