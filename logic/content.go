package logic

import (
	"fmt"
	"net/http"
)

type Content struct {
	ContentType string
	Content     []byte
	Header      http.Header
}

func (E *Content) Error() string {
	return fmt.Sprintf("[%s] %d", E.ContentType, len(E.Content))
}

func NewContent(contentType string, content []byte, keyValue ...string) *Content {
	v := Content{}
	v.ContentType = contentType
	v.Content = content
	v.Header = http.Header{}

	if keyValue != nil {

		key := ""

		for i, vv := range keyValue {
			if i%2 == 0 {
				key = vv
			} else {
				v.Header[key] = []string{vv}
			}
		}
	}

	return &v
}
