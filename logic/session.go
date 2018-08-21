package logic

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"math/rand"
	"net/http"
	"time"
)

type ISession interface {
	Http(resp http.ResponseWriter, req *http.Request) string
}

type Session struct {
	key    string
	maxAge int
}

func NewSession(key string, maxAge int) *Session {
	return &Session{key, maxAge}
}

func newSessionId() string {
	atime := time.Now()
	rand.Seed(atime.UnixNano())
	m := md5.New()
	m.Write([]byte(fmt.Sprintf("kk_%d_%d..&*(", atime.UnixNano(), rand.Int())))
	return hex.EncodeToString(m.Sum(nil))
}

func (S *Session) Http(resp http.ResponseWriter, req *http.Request) string {

	cookie, err := req.Cookie(S.key)

	if err == nil && cookie != nil {
		return cookie.Value
	}

	cookie = &http.Cookie{}
	cookie.Name = S.key
	cookie.Value = newSessionId()
	cookie.HttpOnly = true
	cookie.Path = "/"
	cookie.MaxAge = S.maxAge

	http.SetCookie(resp, cookie)

	return cookie.Value
}
