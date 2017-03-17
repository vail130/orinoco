package sieve

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

func deleteObsoleteDateKeysForTime(t time.Time) {
	// Delete date keys more than 24 hours old
	startTime := t.AddDate(0, 0, -1)
	for dateKey, tt := range dateKeyMap {
		if tt.Year() < startTime.Year() ||
				(tt.Year() == startTime.Year() && tt.Month() < startTime.Month()) ||
				(tt.Year() == startTime.Year() && tt.Month() == startTime.Month() && tt.Day() < startTime.Day()) {
			delete(dateKeyMap, dateKey)
		}
	}
}

func getStreamSummary(now time.Time, stream string, streamMap map[string](map[string]int)) *StreamSummary {
	timeUnits := map[string]time.Duration{
		"hour":   time.Hour,
		"minute": time.Minute,
		"second": time.Second,
	}

	keyFormats := map[string]string{
		"hour":   hourKeyFormat,
		"minute": minuteKeyFormat,
		"second": secondKeyFormat,
	}

	valuesToDate := map[string]int{
		"hour":   0,
		"minute": 0,
		"second": 0,
	}

	trailingCounts := map[string]int{
		"hour":   0,
		"minute": 0,
		"second": 0,
	}

	changePerPeriod := map[string]int{
		"hour":   0,
		"minute": 0,
		"second": 0,
	}

	var timeKey1 string
	var timeKey2 string
	var timeKey3 string

	if timeMap, streamExists := streamMap[stream]; streamExists {
		for period := range trailingCounts {
			if periodToDate, timeKeyExists := timeMap[now.Format(keyFormats[period])]; timeKeyExists {
				valuesToDate[period] = periodToDate
			}

			onePeriodAgo := now.Add(-1 * timeUnits[period])
			timeKey1 = onePeriodAgo.Format(keyFormats[period])
			timeValue1, timeKey1Exists := timeMap[timeKey1]
			if timeKey1Exists {
				trailingCounts[period] += timeValue1
			}

			twoPeriodsAgo := now.Add(-2 * timeUnits[period])
			timeKey2 = twoPeriodsAgo.Format(keyFormats[period])
			timeValue2, timeKey2Exists := timeMap[timeKey2]
			if timeKey2Exists {
				trailingCounts[period] += timeValue2
			}

			threePeriodsAgo := now.Add(-3 * timeUnits[period])
			timeKey3 = threePeriodsAgo.Format(keyFormats[period])
			timeValue3, timeKey3Exists := timeMap[timeKey3]
			if timeKey3Exists {
				trailingCounts[period] += timeValue3
			}

			if timeKey1Exists {
				changePerPeriod[period] = timeValue1
				if timeKey2Exists {
					changePerPeriod[period] = timeValue1 - timeValue2
				}
			}
		}
	}

	projectedThisMinute := float32(valuesToDate["minute"]) / float32(now.Second()+1) * float32(60)
	projectedThisHour := float32(valuesToDate["hour"]) / float32(now.Minute()+1) * float32(60)

	return &StreamSummary{
		stream,
		now.Format(time.RFC3339),
		valuesToDate["second"],
		valuesToDate["minute"],
		valuesToDate["hour"],
		projectedThisMinute,
		projectedThisHour,
		float32(trailingCounts["second"]) / float32(3.0),
		float32(trailingCounts["minute"]) / float32(3.0),
		float32(trailingCounts["hour"]) / float32(3.0),
		changePerPeriod["second"],
		changePerPeriod["minute"],
		changePerPeriod["hour"],
	}
}

func DeleteAllStreamsHandler(w http.ResponseWriter, r *http.Request) {
	dateMap = make(map[string](map[string](map[string]int)))
	dateKeyMap = make(map[string]time.Time)
	w.WriteHeader(http.StatusNoContent)
}

func GetAllStreamSummaries(now time.Time, streamMap map[string](map[string]int)) []StreamSummary {
	var streamSummaries []StreamSummary
	for stream := range streamMap {
		streamSummaries = append(streamSummaries, *getStreamSummary(now, stream, streamMap))
	}
	return streamSummaries
}

func GetAllStreamsHandler(w http.ResponseWriter, r *http.Request) {
	now := GetTimestampForRequest(r.URL.Query(), nil)
	streamMap := GetStreamMapForTime(now)

	var jsonData []byte

	if len(streamMap) > 0 {
		streamSummaries := GetAllStreamSummaries(now, streamMap)

		var err error
		jsonData, err = json.Marshal(streamSummaries)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		jsonData = []byte("null")
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}

func GetStreamHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	stream := vars["stream"]

	now := GetTimestampForRequest(r.URL.Query(), nil)
	streamMap := GetStreamMapForTime(now)

	var jsonData []byte

	if _, ok := streamMap[stream]; ok {
		streamSummary := *getStreamSummary(now, stream, streamMap)

		var err error
		jsonData, err = json.Marshal(streamSummary)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		jsonData = []byte("null")
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}

func GetStreamMapForTime(t time.Time) map[string](map[string]int) {
	dateKey := t.Format(dateKeyFormat)
	streamMap, ok := dateMap[dateKey]
	if ok == false {
		dateMap[dateKey] = make(map[string](map[string]int))
		streamMap = dateMap[dateKey]
		dateKeyMap[dateKey] = t
	}

	return streamMap
}
