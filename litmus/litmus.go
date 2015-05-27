package litmus

import (
	"bytes"
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

	"github.com/vail130/orinoco/httputils"
	"github.com/vail130/orinoco/sieve"
	"github.com/vail130/orinoco/stringutils"
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

type TriggerRequest struct {
	Event   string  `json:"event"`
	Trigger Trigger `json:"trigger"`
	Value   float64 `json:"value"`
}

func triggerEvent(event string, trigger Trigger, metricValue float64) {
	fmt.Println("Trigger", event, trigger.Endpoint, metricValue)

	triggerRequest := TriggerRequest{
		event,
		trigger,
		metricValue,
	}
	jsonData, err := json.Marshal(triggerRequest)
	if err != nil {
		log.Fatal(err.Error())
	} else {
		http.Post(trigger.Endpoint, "application/json", bytes.NewBuffer(jsonData))
	}
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
		data, err := httputils.GetDataFromUrl(url)
		if err != nil {
			data = make([]byte, 0)
		}

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
