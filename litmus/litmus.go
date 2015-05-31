package litmus

import (
	"encoding/json"
	"io/ioutil"
	"log"
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
	Stream    string `yaml:"stream"`
	Metric    string `yaml:"metric"`
	Condition string `yaml:"condition"`
	Endpoint  string `yaml:"endpoint"`
}

type Config struct {
	Url      string             `yaml:"url"`
	Triggers map[string]Trigger `yaml:"triggers"`
}

type TriggerRequest struct {
	Stream  string      `json:"stream"`
	Trigger Trigger     `json:"trigger"`
	Value   interface{} `json:"value"`
}

func triggerStream(stream string, trigger Trigger, metricValue interface{}) {
	triggerRequest := TriggerRequest{
		stream,
		trigger,
		metricValue,
	}

	jsonData, err := json.Marshal(triggerRequest)
	if err != nil {
		log.Fatalln(err.Error())
		return
	}

	httputils.PostDataToUrl(trigger.Endpoint, "application/json", jsonData)
}

func evaluateTriggerForStreamSummary(stream string, trigger Trigger, streamSummary sieve.StreamSummary, conditionRegexp *regexp.Regexp) {
	if trigger.Stream == "*" || streamSummary.Stream == trigger.Stream {
		fieldName := stringutils.UnderscoreToTitle(trigger.Metric)
		reflectedValue := reflect.ValueOf(streamSummary)
		metricInterface := reflectedValue.FieldByName(fieldName).Interface()
		metricValueType := reflect.TypeOf(metricInterface).Name()

		// TODO Refactor this duplicate code

		if metricValueType == "int" {
			metricValue := int64(reflectedValue.FieldByName(fieldName).Int())

			matches := conditionRegexp.FindAllStringSubmatch(trigger.Condition, -1)
			num, err := strconv.ParseInt(matches[0][2], 10, 64)
			if err != nil {
				return
			}

			condition := matches[0][1]
			if (condition == "==" && metricValue == num) ||
				(condition == ">" && metricValue > num) ||
				(condition == ">=" && metricValue >= num) ||
				(condition == "<" && metricValue < num) ||
				(condition == "<=" && metricValue <= num) {
				go triggerStream(stream, trigger, metricValue)
			}

		} else if metricValueType == "float" {
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
				go triggerStream(stream, trigger, metricValue)
			}

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

		var streamSummaries []sieve.StreamSummary
		err = json.Unmarshal(data, &streamSummaries)
		if err != nil {
			log.Fatalln(err)
		}

		for stream, trigger := range triggerMap {
			for i := 0; i < len(streamSummaries); i++ {
				evaluateTriggerForStreamSummary(stream, trigger, streamSummaries[i], conditionRegexp)
			}
		}

		time.Sleep(time.Second)
	}
}

func Litmus(url string, configPath string) {
	configData, err := ioutil.ReadFile(configPath)
	if err != nil {
		log.Fatal(err)
	}

	var config Config
	yaml.Unmarshal(configData, &config)

	if url == "" {
		url = config.Url
	}

	monitorSieve(stringutils.Concat(url, "/streams"), config.Triggers)
}
