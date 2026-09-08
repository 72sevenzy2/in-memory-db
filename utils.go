package root

import (
	"unsafe"
)

// small zero copy utils for string/byte conversions

func StringToByte(str string) []byte {
	return unsafe.Slice(unsafe.StringData(str), len(str))
}

func ByteToString(b []byte) string {
	return unsafe.String(unsafe.SliceData(b), len(b))
}
