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
	Get(path string) (IApp, error)
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
	lock    sync.Mutex
}

func NewMemStore(dir string, expires time.Duration) *MemStore {
	v := MemStore{}
	v.dir = dir
	v.app = map[string]*storeApp{}
	v.expires = expires
	return &v
}

func (S *MemStore) Get(path string) (IApp, error) {

	atime := time.Now()

	S.lock.Lock()

	v, ok := S.app[path]

	S.lock.Unlock()

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

	v = &storeApp{}
	v.app = app
	v.modTime = st.ModTime()
	v.endTime = atime.Add(S.expires)

	S.lock.Lock()
	S.app[path] = v
	S.lock.Unlock()

	return app, nil
}

type FileStore struct {
	dir string
}

func NewFileStore(dir string) *FileStore {
	v := FileStore{}
	v.dir = dir
	return &v
}

func (S *FileStore) Get(path string) (IApp, error) {

	p := filepath.Join(S.dir, path)

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

	return NewApp(object, S, path), nil
}
