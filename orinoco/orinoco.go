package orinoco

import (
	"io/ioutil"
	"log"
	"sync"

	"gopkg.in/yaml.v2"

	"github.com/vail130/orinoco/litmus"
	"github.com/vail130/orinoco/pump"
	"github.com/vail130/orinoco/sieve"
	"github.com/vail130/orinoco/stringutils"
	"github.com/vail130/orinoco/tap"
	"github.com/vail130/orinoco/timeutils"
)

type OrinocoConfig struct {
	SaveFiles string                         `yaml:"save_files"`
	Streams   map[string](map[string]string) `yaml:"streams"`
	Triggers  map[string]litmus.Trigger      `yaml:"triggers"`
}

func Orinoco(configPath string) {
	configData, err := ioutil.ReadFile(configPath)
	if err != nil {
		log.Fatal(err)
	}

	var config OrinocoConfig
	yaml.Unmarshal(configData, &config)

	saveConsumedLogFiles := stringutils.StringToBool(config.SaveFiles)

	for {
		var wg sync.WaitGroup
		for streamName, streamMap := range config.Streams {
			for source, destination := range streamMap {
				wg.Add(1)
				stream := Stream{
					streamName,
					source,
					destination,
				}
				go stream.Process(&wg, saveConsumedLogFiles)
			}
		}
		wg.Wait()

		now := timeutils.UtcNow()
		streamMap := sieve.GetStreamMapForTime(now)
		streamSummaries := sieve.GetAllStreamSummaries(now, streamMap)

		litmus.EvaluateStreamSummaries(config.Triggers, streamSummaries)
	}
}

type Stream struct {
	Name        string
	Source      string
	Destination string
}

func (stream *Stream) Process(wg *sync.WaitGroup, saveConsumedLogFiles bool) {
	defer wg.Done()
	streamHandler := func(streamName string, data []byte) {
		t := sieve.GetTimestampForRequest(nil, data)

		broadcastMessage := func(data []byte) {
			tap.LogMessage(stream.Destination, data)
		}
		sieve.ProcessStream(streamName, t, data, broadcastMessage)
	}
	
	logStreamer := pump.LogStreamer{
		stream.Name,
		saveConsumedLogFiles,
		stream.Source,
		"",
	}
	logStreamer.ConsumeLogs(streamHandler)
}
