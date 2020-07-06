package utils

import (
	"bytes"
	"io"
	"net/mail"
)

func ReadMail(r io.Reader) (*mail.Message, error) {
	var buffer *bytes.Buffer

	buffer = bytes.NewBuffer(nil)
	io.Copy(buffer, r)
	return mail.ReadMessage(buffer)
}
