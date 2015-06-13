package compressutils

import (
	"bytes"
	"compress/gzip"
)

func GzipData(data []byte) []byte {
	var compressedData bytes.Buffer
	w := gzip.NewWriter(&compressedData)
	defer w.Close()
	w.Write(data)
	return compressedData
}