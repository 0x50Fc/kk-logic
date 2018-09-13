package logic

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/go-yaml/yaml"
)

type IStore interface {
	Dir() string
	Get(path string) (IApp, error)
	Walk(fn func(path string))
	Recycle()
}

type storeApp struct {
	app     IApp
	endTime time.Time
	modTime time.Time
}

type MemStore struct {
	dir     string
	app     map[string]*storeApp
	expires time.Duration
	lock    sync.RWMutex
}

func NewMemStore(dir string, expires time.Duration) *MemStore {
	v := MemStore{}
	v.dir = dir
	v.app = map[string]*storeApp{}
	v.expires = expires
	return &v
}

func (S *MemStore) Recycle() {
	S.lock.Lock()
	defer S.lock.Unlock()
	for _, a := range S.app {
		a.app.Recycle()
	}
}

func (S *MemStore) Dir() string {
	return S.dir
}

func (S *MemStore) Get(path string) (IApp, error) {

	atime := time.Now()

	S.lock.RLock()

	v, ok := S.app[path]

	S.lock.RUnlock()

	if ok && atime.Before(v.endTime) {
		return v.app, nil
	}

	p := filepath.Join(S.dir, path)

	st, err := os.Stat(p)

	if err != nil {
		return nil, err
	}

	if ok {
		if v.modTime.Equal(st.ModTime()) {
			v.endTime = atime.Add(S.expires)
			return v.app, nil
		}
	}

	fd, err := os.Open(p)

	if err != nil {
		return nil, err
	}

	defer fd.Close()

	b, err := ioutil.ReadAll(fd)

	if err != nil {
		return nil, err
	}

	var object interface{} = nil

	if strings.HasSuffix(p, ".json") {
		err = json.Unmarshal(b, &object)
		if err != nil {
			return nil, err
		}
	} else if strings.HasSuffix(p, ".yaml") {
		err = yaml.Unmarshal(b, &object)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, errors.New("No implemented " + filepath.Ext(p))
	}

	app := NewApp(object, S, path)

	vv := &storeApp{}
	vv.app = app
	vv.modTime = st.ModTime()
	vv.endTime = atime.Add(S.expires)

	S.lock.Lock()
	S.app[path] = vv
	S.lock.Unlock()

	if v != nil {
		v.app.Recycle()
	}

	return app, nil
}

func (S *MemStore) Walk(fn func(path string)) {

	filepath.Walk(S.dir, func(path string, info os.FileInfo, err error) error {

		if strings.HasSuffix(path, ".yaml") {
			p, _ := filepath.Rel(S.dir, path)
			fn(p)
		}

		return nil
	})
}
