package sieve

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

func (stream *StdoutStream) Process(streamName string, data []byte) {
	fmt.Println(string(data))
}

func (stream *LogStream) Process(streamName string, data []byte) {
	os.MkdirAll(path.Dir(stream.Path), 0666)
	logFile := filepath.Join(stream.Path, stringutils.Concat(streamName, ".log"))
	file, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln(err)
		return
	}

	data = sliceutils.ConcatByteSlices(data, []byte("\n"))
	file.Write(data)
	file.Close()
}

func (stream *HTTPStream) Process(streamName string, data []byte) {
	httputils.PostDataToUrl(stream.URL, "application/json", data)
}

func (stream *S3Stream) Process(streamName string, data []byte) {
	auth, err := aws.EnvAuth()
	if err != nil {
		log.Fatalln(err)
		return
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
		return
	}

	base32UUID = strings.TrimSuffix(base32UUID, "======")
	objectKey := stringutils.Concat(prefix, unixTimeStamp, "_", base32UUID, ".gz")

	compressedData := compressutils.GzipData(data)

	err = bucket.Put(objectKey, compressedData, "binary/octet-stream", s3.Private)
	if err != nil {
		log.Fatalln(err)
		return
	}
}
