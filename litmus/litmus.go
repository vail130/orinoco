package litmus

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"log"
	"reflect"
	"regexp"
	"strconv"
	"time"

	"github.com/vail130/orinoco/httputils"
	"github.com/vail130/orinoco/sieve"
	"github.com/vail130/orinoco/stringutils"
	"github.com/vail130/orinoco/timeutils"
)

type Trigger struct {
	Stream    string `yaml:"stream"`
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
	Stream  string      `json:"stream"`
	Trigger Trigger     `json:"trigger"`
	Value   interface{} `json:"value"`
}

var triggers []Trigger

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

	_, err = httputils.PostDataToUrl(trigger.Endpoint, "application/json", jsonData)
	if err != nil {
		// TODO: Retry or error log?
		log.Print(err)
	}
}

func evaluateTriggerForStreamSummary(trigger Trigger, streamSummary sieve.StreamSummary) {
	conditionRegexp, _ := regexp.Compile(`([=<>]+)([0-9.]+)`)

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
				go triggerStream(streamSummary.Stream, trigger, metricValue)
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
				go triggerStream(streamSummary.Stream, trigger, metricValue)
			}

		}
	}
}

func evaluateTriggers(now time.Time) {
	streamMap := sieve.GetStreamMapForTime(now)
	streamSummaries := sieve.GetAllStreamSummaries(now, streamMap)

	for i := 0; i < len(triggers); i++ {
		for j := 0; j < len(streamSummaries); j++ {
			evaluateTriggerForStreamSummary(triggers[i], streamSummaries[j])
		}
	}
}

func PutEvaluateLitmusTriggersHandler(w http.ResponseWriter, r *http.Request) {
	// This is a test endpoint for doing deterministic integration tests
	defer r.Body.Close()
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		data = make([]byte, 0)
	}

	t := sieve.GetTimestampForRequest(r.URL.Query(), data)
	evaluateTriggers(t)

	w.WriteHeader(http.StatusOK)
}

func Litmus(configTriggers []Trigger) {
	triggers = configTriggers
	for {
		time.Sleep(time.Second)
		now := timeutils.UtcNow()
		evaluateTriggers(now)
	}
}
