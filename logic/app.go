package logic

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/hailongz/kk-lib/duktape"
)

import "C"

type IProtocol func(app IApp, ctx duktape.Context)

type IApp interface {
	Path() string
	Exec(path string, protocols []IProtocol) error
	ExecCode(name string, code string, protocols []IProtocol) error
	Clear()
	ServeHTTP(resp http.ResponseWriter, req *http.Request)
	SessionKey() string
	SessionMaxAge() int
}

var g_protocols = []IProtocol{}

func AddProtocol(protocol IProtocol) {
	g_protocols = append(g_protocols, protocol)
}

type App struct {
	fs            http.Handler
	dir           string
	cached        map[string]string
	lock          sync.RWMutex
	sessionKey    string
	sessionMaxAge int
}

func NewApp(dir string, cached bool, sessionKey string, SessionMaxAge int) IApp {
	v := App{}
	v.fs = http.FileServer(http.Dir(dir))
	v.dir = dir
	v.sessionKey = sessionKey
	v.sessionMaxAge = SessionMaxAge
	if cached {
		v.cached = map[string]string{}
	}
	return &v
}

func (A *App) Clear() {
	if A.cached != nil {
		A.lock.Lock()
		A.cached = map[string]string{}
		A.lock.Unlock()
	}
}

func (A *App) SessionKey() string {
	return A.sessionKey
}

func (A *App) SessionMaxAge() int {
	return A.sessionMaxAge
}

func (A *App) Path() string {
	return A.dir
}

func (A *App) ExecCode(name string, code string, protocols []IProtocol) error {

	var err error

	scope := duktape.NewScope()
	ctx := duktape.New(scope)

	for _, protocol := range g_protocols {
		protocol(A, ctx)
	}

	if protocols != nil {
		for _, protocol := range protocols {
			protocol(A, ctx)
		}
	}

	if A.cached == nil {
		duktape.PushGlobalObject(ctx)
		duktape.PushString(ctx, "debug")
		duktape.PushBoolean(ctx, true)
		duktape.PutProp(ctx, -3)
		duktape.Pop(ctx, 1)
	}

	duktape.PushString(ctx, name)
	duktape.Compile(ctx, name, code)

	if duktape.IsFunction(ctx, -1) {
		err = duktape.Call(ctx, 0)
	}

	duktape.Pop(ctx, 1)

	duktape.Recycle(ctx)

	return err
}

func (A *App) Exec(path string, protocols []IProtocol) error {

	var code string
	var hasCode bool

	if A.cached != nil {
		A.lock.RLock()
		code, hasCode = A.cached[path]
		A.lock.RUnlock()
	}

	if !hasCode {

		fd, err := os.Open(filepath.Join(A.dir, path))

		if err != nil {
			return err
		}

		defer fd.Close()

		b, err := ioutil.ReadAll(fd)

		if err != nil {
			return err
		}

		code = string(b)

		if A.cached != nil {
			A.lock.Lock()
			A.cached[path] = code
			A.lock.Unlock()
		}

	}

	return A.ExecCode(path, code, protocols)
}

func newSessionId() string {
	v := time.Now()
	m := md5.New()
	rand.Seed(int64(v.Nanosecond()))
	m.Write([]byte(fmt.Sprintf("kk_%d-%d..", v.Nanosecond(), rand.Int())))
	return hex.EncodeToString(m.Sum(nil))
}

func (A *App) ServeHTTP(resp http.ResponseWriter, req *http.Request) {

	if strings.HasSuffix(req.URL.Path, ".json") {

		path := req.URL.Path[0 : len(req.URL.Path)-2]

		err := A.Exec(path, []IProtocol{func(app IApp, ctx duktape.Context) {

			duktape.PushGlobalObject(ctx)

			duktape.PushString(ctx, "header")
			duktape.PushGoFunction(ctx, func() int {

				top := duktape.GetTop(ctx)

				if top > 1 && duktape.IsString(ctx, -top) && duktape.IsString(ctx, -top+1) {
					resp.Header()[duktape.ToString(ctx, -top)] = []string{duktape.ToString(ctx, -top+1)}
				}

				return 0
			})

			duktape.PutProp(ctx, -3)

			duktape.PushString(ctx, "echo")
			duktape.PushGoFunction(ctx, func() int {

				top := duktape.GetTop(ctx)
				for i := 0; i < top; i++ {
					if duktape.IsString(ctx, -top+i) {
						resp.Write([]byte(duktape.ToString(ctx, -top+i)))
					} else if duktape.IsBytes(ctx, -top+i) {
						resp.Write(duktape.ToBytes(ctx, -top+i))
					}
				}

				return 0
			})

			duktape.PutProp(ctx, -3)

			duktape.PushString(ctx, "_REQUEST")
			duktape.PushObject(ctx)

			duktape.PushString(ctx, "url")
			duktape.PushString(ctx, req.URL.String())
			duktape.PutProp(ctx, -3)

			duktape.PushString(ctx, "method")
			duktape.PushString(ctx, req.Method)
			duktape.PutProp(ctx, -3)

			duktape.PushString(ctx, "path")
			duktape.PushString(ctx, req.URL.Path)
			duktape.PutProp(ctx, -3)

			duktape.PushString(ctx, "protocol")
			duktape.PushString(ctx, req.URL.Scheme)
			duktape.PutProp(ctx, -3)

			duktape.PushString(ctx, "hostname")
			duktape.PushString(ctx, req.URL.Hostname())
			duktape.PutProp(ctx, -3)

			duktape.PutProp(ctx, -3)

			duktape.PushString(ctx, "_HEADER")
			duktape.PushObject(ctx)

			for key, vs := range req.Header {
				duktape.PushString(ctx, key)
				duktape.PushString(ctx, vs[0])
				duktape.PutProp(ctx, -3)
			}

			duktape.PutProp(ctx, -3)

			duktape.PushString(ctx, "_COOKIE")
			duktape.PushObject(ctx)

			sessionKey := app.SessionKey()
			sessionId := ""

			for _, cookie := range req.Cookies() {
				duktape.PushString(ctx, cookie.Name)
				duktape.PushString(ctx, cookie.Value)
				duktape.PutProp(ctx, -3)
				if cookie.Name == sessionKey {
					sessionId = cookie.Value
				}
			}

			if sessionId == "" {
				sessionId = newSessionId()
				cookie := http.Cookie{Name: sessionKey, Value: sessionId, Path: "/", HttpOnly: true, MaxAge: app.SessionMaxAge()}
				http.SetCookie(resp, &cookie)
			}

			duktape.PushString(ctx, "_SESSIONID")
			duktape.PushObject(ctx)
			duktape.PushString(ctx, sessionId)
			duktape.PutProp(ctx, -3)

			duktape.PutProp(ctx, -3)

			duktape.PushString(ctx, "_GET")
			duktape.PushObject(ctx)

			for key, vs := range req.URL.Query() {
				duktape.PushString(ctx, key)
				duktape.PushString(ctx, vs[0])
				duktape.PutProp(ctx, -3)
			}

			duktape.PutProp(ctx, -3)

			if req.Method == "POST" {

				duktape.PushString(ctx, "_POST")
				duktape.PushObject(ctx)

				for key, vs := range req.PostForm {
					duktape.PushString(ctx, key)
					duktape.PushString(ctx, vs[0])
					duktape.PutProp(ctx, -3)
				}

				duktape.PutProp(ctx, -3)

			}

			duktape.Pop(ctx, 1)
		}})

		if err != nil {
			resp.WriteHeader(500)
			resp.Write([]byte(err.Error()))
		}

	} else if strings.HasSuffix(req.URL.Path, ".js") {
		resp.WriteHeader(404)
	} else {
		A.fs.ServeHTTP(resp, req)
	}

}
