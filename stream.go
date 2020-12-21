package begundal

import (
	"bytes"
	"strings"
)

type streamInBetween struct {
	startData  bool
	initiate   strings.Builder
	temp       *bytes.Buffer
	startPoint string
	endPoint   rune
}

func newStreamInBetween(startPoint string, endPoint rune) *streamInBetween {
	return &streamInBetween{
		temp:       bytes.NewBuffer(nil),
		startPoint: startPoint,
		endPoint:   endPoint,
	}
}

func (e *streamInBetween) Read(p []byte) (n int, err error) {
	return e.temp.Read(p)
}

func (e *streamInBetween) Write(p []byte) (n int, err error) {
	for _, b := range p {
		if e.startData {
			// Stop early since we hit the closing quote
			if b == byte(e.endPoint) {
				return 0, nil
			}
			err = e.temp.WriteByte(b)
			if err != nil {
				return 0, err
			}
		} else {
			err = e.initiate.WriteByte(b)
			if err != nil {
				return 0, err
			}
			dd := e.initiate.String()
			if strings.Contains(dd, e.startPoint) {
				e.startData = true
			}
		}
		n++
	}
	return
}
