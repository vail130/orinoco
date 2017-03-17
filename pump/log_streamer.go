package pump

import (
	"bufio"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/vail130/orinoco/stringutils"
)

type LogStreamer struct {
	StreamName    string
	SaveLogFiles  bool
	LogPath       string
	ConsumingPath string
}

type StreamHandler func(string, []byte)

func (logStreamer *LogStreamer) cleanupLogFile() error {
	if !logStreamer.SaveLogFiles {
		err := os.Remove(logStreamer.ConsumingPath)
		if err != nil {
			return err
		}
	}

	return nil
}

func (logStreamer *LogStreamer) prepareLogFile() error {
	now := time.Now()
	unixTimeStamp := strconv.FormatInt(now.Unix(), 10)
	base64UUID, err := stringutils.GetBase64UUID()
	if err != nil {
		log.Fatalln(err)
		return err
	}

	uniquePath := stringutils.Concat(logStreamer.LogPath, ".", unixTimeStamp, ".", base64UUID)
	logStreamer.ConsumingPath = stringutils.Concat(uniquePath, ".consuming")

	err = os.Rename(logStreamer.LogPath, logStreamer.ConsumingPath)
	if err != nil {
		return err
	}

	return nil
}

func (logStreamer *LogStreamer) ConsumeLogs(streamHandler StreamHandler) {
	err := logStreamer.prepareLogFile()
	if err != nil {
		if !os.IsNotExist(err) {
			log.Fatalln(err)
		}
		return
	}

	file, err := os.OpenFile(logStreamer.ConsumingPath, os.O_RDONLY, 0666)
	if err != nil {
		log.Fatalln(err)
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		messageData := scanner.Bytes()
		if string(messageData) != "null" {
			streamHandler(logStreamer.StreamName, messageData)
		}
	}

	if err = scanner.Err(); err != nil {
		log.Fatalln(err)
	}

	file.Close()

	logStreamer.cleanupLogFile()
}
