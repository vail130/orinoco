package pump

import (
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"

	"gopkg.in/amz.v1/aws"
	"gopkg.in/amz.v1/s3"

	"github.com/vail130/orinoco/compressutils"
	"github.com/vail130/orinoco/httputils"
	"github.com/vail130/orinoco/sieve"
	"github.com/vail130/orinoco/sliceutils"
	"github.com/vail130/orinoco/stringutils"
	"github.com/vail130/orinoco/timeutils"
)

type StdoutStream struct{}
type LogStream struct {
	Path string
}
type HTTPStream struct {
	URL string
}

type S3Stream struct {
	Region string
	Bucket string
	Prefix string
}

func indexEventsByStream(eventArray []sieve.Event) map[string]([]byte) {
	streamMap := make(map[string]([]byte))

	for i := 0; i < len(eventArray); i++ {
		event := eventArray[i]

		if _, ok := streamMap[event.Stream]; !ok {
			streamMap[event.Stream] = make([]byte, 0)
		}

		streamMap[event.Stream] = sliceutils.ConcatByteSlices(streamMap[event.Stream], event.Data, []byte("\n"))
	}

	return streamMap
}

func (stream *StdoutStream) Process(eventArray []sieve.Event) {
	outputString := ""
	for i := 0; i < len(eventArray); i++ {
		event := eventArray[i]
		outputString = stringutils.Concat(string(event.Data), "\n")
	}

	fmt.Print(outputString)
}

func (stream *LogStream) Process(eventArray []sieve.Event) {
	streamMap := indexEventsByStream(eventArray)

	for streamName, data := range streamMap {
		os.MkdirAll(path.Dir(stream.Path), 0666)
		logFile := filepath.Join(stream.Path, stringutils.Concat(streamName, ".log"))
		file, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			// TODO: Should this be fatal? Probably not. If so, detect earlier!
			log.Fatalln(err)
			continue
		}

		file.Write(data)
		file.Close()
	}
}

func (stream *HTTPStream) Process(eventArray []sieve.Event) {
	// TODO: Batch (?)
	for i := 0; i < len(eventArray); i++ {
		event := eventArray[i]
		url := stream.URL
		if strings.Contains(url, "{{stream}}") {
			url = strings.Replace(url, "{{stream}}", event.Stream, -1)
		}

		_, err := httputils.PostDataToUrl(url, "application/json", event.Data)
		if err != nil {
			// TODO: Log to error log? Retry?
		}
	}
}

func (stream *S3Stream) Process(eventArray []sieve.Event) {
	streamMap := indexEventsByStream(eventArray)

	for streamName, data := range streamMap {
		auth, err := aws.EnvAuth()
		if err != nil {
			log.Fatalln(err)
			continue
		}
		s3Instance := s3.New(auth, aws.Regions[stream.Region])
		bucket := s3Instance.Bucket(stream.Bucket)

		prefix := stream.Prefix
		now := timeutils.UtcNow()
		if strings.Contains(prefix, "{{stream}}") {
			prefix = strings.Replace(prefix, "{{stream}}", streamName, -1)
		}
		if strings.Contains(prefix, "{{year}}") {
			prefix = strings.Replace(prefix, "{{year}}", now.Format("2006"), -1)
		}
		if strings.Contains(prefix, "{{month}}") {
			prefix = strings.Replace(prefix, "{{month}}", now.Format("01"), -1)
		}
		if strings.Contains(prefix, "{{day}}") {
			prefix = strings.Replace(prefix, "{{day}}", now.Format("02"), -1)
		}

		unixTimeStamp := strconv.FormatInt(now.Unix(), 10)
		base32UUID, err := stringutils.GetBase32UUID()
		if err != nil {
			log.Fatalln(err)
			continue
		}

		base32UUID = strings.TrimSuffix(base32UUID, "======")
		objectKey := stringutils.Concat(prefix, unixTimeStamp, "_", base32UUID, ".gz")

		compressedData := compressutils.GzipData(data)

		err = bucket.Put(objectKey, compressedData, "binary/octet-stream", s3.Private)
		if err != nil {
			log.Fatalln(err)
			continue
		}
	}
}
