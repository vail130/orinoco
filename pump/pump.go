package pump

import (
	"time"

	"github.com/vail130/orinoco/sieve"
	"github.com/vail130/orinoco/timeutils"
)

func forkStream(streams [](map[string]string), eventArray []sieve.Event) {
	for i := 0; i < len(streams); i++ {
		switch {
		case streams[i]["type"] == "stdout":
			stream := StdoutStream{}
			go stream.Process(eventArray)
		case streams[i]["type"] == "log":
			stream := LogStream{
				streams[i]["path"],
			}
			go stream.Process(eventArray)
		case streams[i]["type"] == "http":
			stream := HTTPStream{
				streams[i]["url"],
			}
			go stream.Process(eventArray)
		case streams[i]["type"] == "s3":
			stream := S3Stream{
				streams[i]["region"],
				streams[i]["bucket"],
				streams[i]["prefix"],
			}
			go stream.Process(eventArray)
		}
	}
}

func Pump(minBatchSize int, maxBatchDelay int, streams [](map[string]string), eventChannel chan *sieve.Event) {
	var start time.Time
	var now time.Time
	for {
		start = timeutils.UtcNow()
		eventArray := make([]sieve.Event, 0)
		for event := range eventChannel {
			eventArray = append(eventArray, *event)
			now = timeutils.UtcNow()
			if now.Sub(start) >= time.Duration(maxBatchDelay) || len(eventChannel) >= minBatchSize {
				forkStream(streams, eventArray)
				start = now
				break
			}
		}
	}
}
