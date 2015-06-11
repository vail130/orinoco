package sieve

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"path"
	"path/filepath"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/mux"

	"github.com/vail130/orinoco/httputils"
	"github.com/vail130/orinoco/sliceutils"
	"github.com/vail130/orinoco/stringutils"
	"github.com/vail130/orinoco/timeutils"
)

type StreamSummary struct {
	Stream                   string  `json:"stream"`
	Timestamp                string  `json:"timestamp"`
	SecondToDate             int     `json:"second_to_date"`
	MinuteToDate             int     `json:"minute_to_date"`
	HourToDate               int     `json:"hour_to_date"`
	ProjectedThisMinute      float32 `json:"projected_this_minute"`
	ProjectedThisHour        float32 `json:"projected_this_hour"`
	TrailingAveragePerSecond float32 `json:"trailing_average_per_second"`
	TrailingAveragePerMinute float32 `json:"trailing_average_per_minute"`
	TrailingAveragePerHour   float32 `json:"trailing_average_per_hour"`
	ChangePerSecond          int     `json:"change_per_second"`
	ChangePerMinute          int     `json:"change_per_minute"`
	ChangePerHour            int     `json:"change_per_hour"`
}

type Event struct {
	Stream    string `json:"stream"`
	Timestamp string `json:"timestamp"`
	Data      string `json:"data"`
}

type StdoutStream struct{}
type LogStream struct {
	LogDir string
}
type HTTPStream struct {
	URL string
}

//type WebsocketStream struct {
//	URL string
//}

// Reference time for formats: Mon Jan 2 15:04:05 -0700 MST 2006
const dateKeyFormat = "2006-01-02"
const hourKeyFormat = "2006-01-02-15"
const minuteKeyFormat = "2006-01-02-15-04"
const secondKeyFormat = "2006-01-02-15-04-05"

var dateMap = make(map[string](map[string](map[string]int)))
var dateKeyMap = make(map[string]time.Time)

func GetTimestampForRequest(queryValues url.Values, data []byte) time.Time {
	if !isTestEnv {
		return timeutils.UtcNow()
	}

	if timestampString, ok := queryValues["timestamp"]; ok {
		if timestamp, err := time.Parse(time.RFC3339, timestampString[0]); err == nil {
			return timestamp.UTC()
		}
	}

	if len(data) > 0 {
		var f interface{}
		err := json.Unmarshal(data, &f)
		if err != nil {
			return timeutils.UtcNow()
		}

		if timestampString, ok := f.(map[string]interface{})["timestamp"]; ok {
			if timestamp, err := time.Parse(time.RFC3339, timestampString.(string)); err == nil {
				return timestamp
			}
		}
	}

	return timeutils.UtcNow()
}

func trackStreamForTime(stream string, t time.Time) {
	streamMap := GetStreamMapForTime(t)
	deleteObsoleteDateKeysForTime(t)

	dateMap, ok := streamMap[stream]
	if ok == false {
		streamMap[stream] = make(map[string]int)
		dateMap = streamMap[stream]
	}

	timeKeys := []string{
		t.Format(hourKeyFormat),
		t.Format(minuteKeyFormat),
		t.Format(secondKeyFormat),
	}

	for i := 0; i < len(timeKeys); i++ {
		if _, ok := dateMap[timeKeys[i]]; ok == false {
			dateMap[timeKeys[i]] = 0
		}
		dateMap[timeKeys[i]] += 1
	}
}

func (stream *StdoutStream) Process(wg *sync.WaitGroup, streamName string, data []byte) {
	defer wg.Done()
	fmt.Println(string(data))
}

func (stream *LogStream) Process(wg *sync.WaitGroup, streamName string, data []byte) {
	defer wg.Done()
	
	os.MkdirAll(path.Dir(stream.LogDir), 0666)
	logFile := filepath.Join(stream.LogDir, stringutils.Concat(streamName, ".log"))
	file, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln(err)
		return
	}

	data = sliceutils.ConcatByteSlices(data, []byte("\n"))
	file.Write(data)
	file.Close()
}

func (stream *HTTPStream) Process(wg *sync.WaitGroup, streamName string, data []byte) {
	defer wg.Done()
	httputils.PostDataToUrl(stream.URL, "application/json", data)
}

//func (stream *WebsocketStream) Process(wg *sync.WaitGroup, streamName string, data []byte) {
//	defer wg.Done()

//}

func PostStreamHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	streamName := vars["stream"]

	defer r.Body.Close()
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		data = make([]byte, 0)
	}

	t := GetTimestampForRequest(r.URL.Query(), data)
	trackStreamForTime(streamName, t)

	var wg sync.WaitGroup
	for i := 0; i < len(configuredStreams); i++ {
		switch {

		case configuredStreams[i] == "-":
			stream := StdoutStream{}
			wg.Add(1)
			go stream.Process(&wg, streamName, data)

		case strings.HasPrefix(configuredStreams[i], "/"):
			stream := LogStream{
				configuredStreams[i],
			}
			wg.Add(1)
			go stream.Process(&wg, streamName, data)

		case strings.HasPrefix(configuredStreams[i], "http"):
			stream := HTTPStream{
				configuredStreams[i],
			}
			wg.Add(1)
			go stream.Process(&wg, streamName, data)

			//		case strings.HasPrefix(configuredStreams[i], "ws"):
			//			stream := WebsocketStream{
			//				configuredStreams[i],
			//			}
			//			wg.Add(1)
			//			go stream.Process(&wg, streamName, data)

			// TODO
			// Add other streams

		}
	}
	wg.Wait()
	w.WriteHeader(http.StatusOK)
}
