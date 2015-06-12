package orinoco

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"

	"github.com/vail130/orinoco/litmus"
	"github.com/vail130/orinoco/sieve"
)

type Config struct {
	Port     string           `yaml:"port"`
	Streams  [](map[string]string)         `yaml:"streams"`
	Triggers []litmus.Trigger `yaml:"triggers"`
}

func Orinoco(configPath string) {
	configData, err := ioutil.ReadFile(configPath)
	if err != nil {
		log.Fatal(err)
	}

	var config Config
	yaml.Unmarshal(configData, &config)

	go litmus.Litmus(config.Triggers)

	sieve.Sieve(config.Port, config.Streams)
}
