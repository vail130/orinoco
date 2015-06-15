package orinoco

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"

	"github.com/vail130/orinoco/litmus"
	"github.com/vail130/orinoco/pump"
	"github.com/vail130/orinoco/sieve"
)

type Config struct {
	Port          string                `yaml:"port"`
	MinBatchSize  int                   `yaml:"min_batch_size"`
	MaxBatchDelay int                   `yaml:"max_batch_delay"`
	Streams       [](map[string]string) `yaml:"streams"`
	Triggers      []litmus.Trigger      `yaml:"triggers"`
}

func Orinoco(configPath string) {
	configData, err := ioutil.ReadFile(configPath)
	if err != nil {
		log.Fatal(err)
	}

	var config Config
	yaml.Unmarshal(configData, &config)

	go litmus.Litmus(config.Triggers)

	eventChannel := make(chan *sieve.Event)
	go pump.Pump(config.MinBatchSize, config.MaxBatchDelay, config.Streams, eventChannel)
	sieve.Sieve(config.Port, eventChannel)
}
