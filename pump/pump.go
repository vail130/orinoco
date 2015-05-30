package pump

import (
	"bufio"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/vail130/orinoco/httputils"
	"github.com/vail130/orinoco/stringutils"
)

func sendEventOverHttp(url string, data []byte) {
	_, err := httputils.PostDataToUrl(url, "application/json", data)
	if err != nil {
		log.Fatalln(err)
	}
}

func consumeLogs(logPath string, url string) {
	now := time.Now()
	unixTimeStamp := strconv.FormatInt(now.Unix(), 10)

	consumingPath := stringutils.Concat(logPath, ".", unixTimeStamp, ".consuming")
	err := os.Rename(logPath, consumingPath)
	if err != nil {
		return
	}

	if file, err := os.OpenFile(consumingPath, os.O_RDONLY, 0666); err == nil {
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			messageData := scanner.Bytes()
			if string(messageData) != "null" {
				sendEventOverHttp(url, messageData)
			}
		}

		if err = scanner.Err(); err != nil {
			log.Fatalln(err)
		}

		file.Close()
	}

	consumedPath := stringutils.Concat(logPath, ".", unixTimeStamp, ".consumed")
	os.Rename(consumingPath, consumedPath)
}

func Pump(logPath string, url string) {
	for {
		consumeLogs(logPath, url)
	}
}
