package logic

import (
	"log"

	"github.com/hailongz/kk-lib/dynamic"
)

var globalIgnoreKeys = map[string]bool{"title": true, "type": true, "meta": true, "input": true}

type IApp interface {
	Path() string
	Store() IStore
	Object() interface{}
	Exec(ctx IContext, name string) error
	Each(fn func(name string, logic ILogic) bool)
	Recycle()
}

type App struct {
	path   string
	store  IStore
	object interface{}
	logics map[string]interface{}
}

func SetGlobalIgnoreKey(keys ...string) {
	for _, key := range keys {
		globalIgnoreKeys[key] = true
	}
}

func NewApp(object interface{}, store IStore, path string) *App {
	a := App{}
	a.path = path
	a.store = store
	a.object = object
	a.logics = map[string]interface{}{}

	dynamic.Each(object, func(key interface{}, value interface{}) bool {

		skey := dynamic.StringValue(key, "")

		if globalIgnoreKeys[skey] {
			return true
		}

		{
			s, ok := value.(string)
			if ok {
				a.logics[skey] = s
				return true
			}
		}

		{
			class := dynamic.StringValue(dynamic.Get(value, "$class"), "")

			if class != "" {
				a.logics[skey] = NewLogic(class, value)
				return true
			}
		}

		return true
	})
	return &a
}

func (A *App) Store() IStore {
	return A.store
}

func (A *App) Object() interface{} {
	return A.object
}

func (A *App) Path() string {
	return A.path
}

func (A *App) Exec(ctx IContext, name string) error {

	log.Printf("[APP] [EXEC] %s >> %s", A.path, name)

	v, ok := A.logics[name]

	if !ok {
		return nil
	}

	{
		s, ok := v.(string)
		if ok {
			return A.Exec(ctx, s)
		}
	}

	{
		s, ok := v.(ILogic)

		if ok {
			return s.Exec(ctx, A)
		}
	}

	return nil
}

func (A *App) Each(fn func(name string, logic ILogic) bool) {

	for name, v := range A.logics {
		{
			vv, ok := v.(ILogic)
			if ok {
				if !fn(name, vv) {
					break
				}
			}
		}
	}

}

func (A *App) Recycle() {

	for _, v := range A.logics {
		{
			vv, ok := v.(ILogic)
			if ok {
				vv.Recycle()
			}
		}
	}

}
