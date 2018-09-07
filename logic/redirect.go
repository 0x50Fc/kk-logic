package logic

import (
	"fmt"
)

type Redirect struct {
	URL string
}

func (E *Redirect) Error() string {
	return fmt.Sprintf("[Location] %s", E.URL)
}

func NewRedirect(url string) *Redirect {
	return &Redirect{URL: url}
}
