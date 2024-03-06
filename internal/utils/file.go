package utils

import "bytes"

func Peek(buf *bytes.Buffer, b []byte) (int, error) {
	buf2 := bytes.NewBuffer(buf.Bytes())
	return buf2.Read(b)
}
