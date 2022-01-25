package forward

import (
	"fmt"
	"github.com/v2fly/v2ray-core/v5/common/buf"
	"strings"
)

var httpMethods = []string{"GET", "POST", "PUT"}

type readError struct {
	error
}

func (e readError) Error() string {
	return e.error.Error()
}

func (e readError) Inner() error {
	return e.error
}

type writeError struct {
	error
}

func (e writeError) Error() string {
	return e.error.Error()
}

func (e writeError) Inner() error {
	return e.error
}

func isHTTPMethod(s string) bool {
	upper := strings.ToUpper(s)
	for _, method := range httpMethods {
		if upper == method {
			return true
		}
	}
	return false
}

func copy(reader buf.Reader, writer buf.Writer, host string) error {
	for {
		buffer, err := reader.ReadMultiBuffer()
		if !buffer.IsEmpty() {
			bufStr := buffer.String()
			bufArr := strings.Split(bufStr, " ")
			if len(bufArr) > 1 && isHTTPMethod(bufArr[0]) && strings.HasPrefix(bufArr[1], "/") {
				// TODO lack situation of https
				bufArr[1] = fmt.Sprintf("http://%s%s", host, bufArr[1])
			}
			newBufStr := ""
			for _, s := range bufArr {
				newBufStr += s
				newBufStr += " "
			}

			newB, err := buf.ReadFrom(strings.NewReader(newBufStr))
			if err != nil {
				return readError{err}
			}
			buffer = newB
			if wErr := writer.WriteMultiBuffer(buffer); wErr != nil {
				return writeError{wErr}
			}
		}

		if err != nil {
			return readError{err}
		}
	}
	return nil
}
