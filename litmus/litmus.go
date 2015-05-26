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

func triggerEvent(event string, trigger Trigger, metricValue float64) {
	fmt.Println("Trigger", event, trigger.Endpoint, metricValue)
}

func getDataFromUrl(url string) []byte {
	response, err := http.Get(url)
	defer response.Body.Close()
	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		data = make([]byte, 0)
	}
	return data
}

func evaluateTriggerForEventSummary(event string, trigger Trigger, eventSummary sieve.EventSummary, conditionRegexp *regexp.Regexp) {
	if trigger.Event == "*" || eventSummary.Event == trigger.Event {
		fieldName := stringutils.UnderscoreToTitle(trigger.Metric)
		reflectedValue := reflect.ValueOf(eventSummary)
		metricValue := float64(reflectedValue.FieldByName(fieldName).Float())

		matches := conditionRegexp.FindAllStringSubmatch(trigger.Condition, -1)
		num, err := strconv.ParseFloat(matches[0][2], 64)
		if err != nil {
			return
		}

		condition := matches[0][1]

		if (condition == "==" && metricValue == num) ||
			(condition == ">" && metricValue > num) ||
			(condition == ">=" && metricValue >= num) ||
			(condition == "<" && metricValue < num) ||
			(condition == "<=" && metricValue <= num) {
			triggerEvent(event, trigger, metricValue)
		}
	}
}

func monitorSieve(url string, triggerMap map[string]Trigger) {
	conditionRegexp, _ := regexp.Compile(`([=<>]+)([0-9.]+)`)

	for {
		data := getDataFromUrl(url)

		var eventSummaries []sieve.EventSummary
		json.Unmarshal(data, &eventSummaries)

		for event, trigger := range triggerMap {
			for j := 0; j < len(eventSummaries); j++ {
				eventSummary := eventSummaries[j]
				evaluateTriggerForEventSummary(event, trigger, eventSummary, conditionRegexp)
			}
		}

		time.Sleep(time.Second)
	}
}

func Litmus(host string, port string, configPath string) {
	configData, err := ioutil.ReadFile(configPath)
	if err != nil {
		log.Fatal(err)
	}

	var config Config
	yaml.Unmarshal(configData, &config)

	if host == "" {
		host = config.Host
	}

	if port == "" {
		port = config.Port
	}

	url := stringutils.Concat("http://", host, ":", port, "/events")

	monitorSieve(url, config.Triggers)
}
