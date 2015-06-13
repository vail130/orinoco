package compressutils

import (
	"bytes"
	"compress/gzip"
)

func GzipData(data []byte) []byte {
	var compressedData bytes.Buffer
	w := gzip.NewWriter(&compressedData)
	w.Write(data)

	// Close to ensure write flushes data to buffer
	w.Close()

	return compressedData.Bytes()
}
