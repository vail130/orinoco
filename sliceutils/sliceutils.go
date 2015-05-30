package sliceutils

func ConcatByteSlices(byteSlices ...[]byte) []byte {
	newByteSlice := make([]byte, 0)
	for i := 0; i < len(byteSlices); i++ {
		newByteSlice = append(newByteSlice, byteSlices[i]...)
	}
	return newByteSlice
}
