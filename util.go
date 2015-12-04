package bme280

import "encoding/binary"
import "bytes"

func convert(b []byte, data interface{}) error {
	buf := bytes.NewReader(b)
	return binary.Read(buf, binary.LittleEndian, data)
}
