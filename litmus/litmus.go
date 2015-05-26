package litmus

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"

	"gopkg.in/yaml.v2"

	"../sieve"
	"../stringutils"
)

type Trigger struct {
	Event     string `yaml:"event"`
	Metric    string `yaml:"metric"`
	Condition string `yaml:"condition"`
	Endpoint  string `yaml:"endpoint"`
}

type Config struct {
	Host     string             `yaml:"host"`
	Port     string             `yaml:"port"`
	Triggers map[string]Trigger `yaml:"triggers"`
}

func triggerEvent(event string, trigger Trigger) {
	fmt.Println("Trigger", event, trigger.Endpoint)
}

func Litmus(host string, port string, configPath string) {
	configData, err := ioutil.ReadFile(configPath)
	if err != nil {
		log.Fatal(err)
	}

	var config Config
	yaml.Unmarshal(configData, &config)

	url := stringutils.Concat("http://", host, ":", port, "/events")
	fmt.Println(url)

	conditionRegexp, err := regexp.Compile(`([=<>]+)([0-9.]+)`)

	for {
		response, err := http.Get(url)

		defer response.Body.Close()
		data, err := ioutil.ReadAll(response.Body)
		if err != nil {
			data = make([]byte, 0)
		}

		var eventSummaries []sieve.EventSummary
		json.Unmarshal(data, &eventSummaries)

		for event, trigger := range config.Triggers {
			for j := 0; j < len(eventSummaries); j++ {
				eventSummary := eventSummaries[j]
				if trigger.Event == "*" || eventSummary.Event == trigger.Event {
					reflectedValue := reflect.ValueOf(eventSummary)
					
					fieldNameParts := strings.Split(trigger.Metric, "_")
					for k := 0; k < len(fieldNameParts); k++ {
						fieldNameParts[k] = strings.Title(fieldNameParts[k])
					}
					fieldName := strings.Join(fieldNameParts, "")
					
					metricValue := float64(reflectedValue.FieldByName(fieldName).Float())

					matches := conditionRegexp.FindAllStringSubmatch(trigger.Condition, -1)
					num, err := strconv.ParseFloat(matches[0][2], 64)
					if err != nil {
						continue
					}

					condition := matches[0][1]

					if (condition == "==" && metricValue == num) ||
						(condition == ">" && metricValue > num) ||
						(condition == ">=" && metricValue >= num) ||
						(condition == "<" && metricValue < num) ||
						(condition == "<=" && metricValue <= num) {
						go triggerEvent(event, trigger)
					}
				}
			}
		}

		time.Sleep(time.Second)
	}
}
