package logic

import (
	"log"
	"strings"

	"github.com/hailongz/kk-lib/dynamic"
)

var globalLogicCreator = map[string]LogicCreator{}

type LogicFunction func(ctx IContext, app IApp) error
type LogicCreator func(object interface{}) ILogic

type ILogic interface {
	Object() interface{}
	EvaluateValue(ctx IContext, app IApp, value interface{}, object interface{}) interface{}
	Get(ctx IContext, app IApp, key string) interface{}
	Exec(ctx IContext, app IApp) error
	On(name string, value interface{})
	Done(ctx IContext, app IApp, name string) error
	Error(ctx IContext, app IApp, err error) error
	Has(name string) bool
}

type Logic struct {
	Tag    string
	object interface{}
	on     map[string]interface{}
}

func NewLogic(class string, object interface{}) ILogic {
	fn, ok := globalLogicCreator[class]
	if ok {
		return fn(object)
	}
	v := Logic{}
	v.Init(object)
	return &v
}

func Openlib(class string, creator LogicCreator) {
	globalLogicCreator[class] = creator
	// log.Printf("[OPENLIB] [%s] [OK]\n", class)
}

func init() {
	Openlib("kk.Logic", func(object interface{}) ILogic {
		v := Logic{}
		v.Init(object)
		return &v
	})
}

func (L *Logic) Init(object interface{}) {
	L.object = object
	L.on = map[string]interface{}{}
	L.Tag = dynamic.StringValue(dynamic.Get(L.object, "$id"), "")

	if L.Tag == "" {
		L.Tag = dynamic.StringValue(dynamic.Get(L.object, "$tag"), "")
	}

	if L.Tag == "" {
		L.Tag = dynamic.StringValue(dynamic.Get(L.object, "$class"), "")
	}

	dynamic.Each(object, func(key interface{}, value interface{}) bool {

		skey := dynamic.StringValue(key, "")

		if strings.HasPrefix(skey, "on") {
			class := dynamic.StringValue(dynamic.Get(value, "$class"), "")
			if class != "" {
				L.On(skey[2:], NewLogic(class, value))
			} else {
				L.On(skey[2:], value)
			}
		}

		return true
	})
}

func (L *Logic) Object() interface{} {
	return L.object
}

func (L *Logic) EvaluateValue(ctx IContext, app IApp, value interface{}, object interface{}) interface{} {

	if value == nil {
		return nil
	}

	{
		s, ok := value.(string)
		if ok {
			if strings.HasPrefix(s, "=") {
				return ctx.Evaluate(s[1:], L.Tag)
			}
			return s
		}
	}

	{
		s, ok := value.(ILogic)

		if ok {
			ctx.Begin()
			ctx.Set(ObjectKeys, object)
			ctx.Set(ResultKeys, Nil)
			s.Exec(ctx, app)
			v := ctx.Get(ResultKeys)
			if v == Nil {
				v = nil
			}
			ctx.End()
			return v
		}
	}

	{
		s, ok := value.([]interface{})

		if ok {
			v := []interface{}{}
			for _, i := range s {
				vv := L.EvaluateValue(ctx, app, i, object)
				if vv != nil {
					v = append(v, vv)
				}
			}
			return v
		}
	}

	{
		s := dynamic.StringValue(dynamic.Get(value, "$class"), "")

		if s != "" {
			logic := NewLogic(s, value)
			return L.EvaluateValue(ctx, app, logic, object)
		}
	}

	{
		s, ok := value.(map[string]interface{})

		if ok {
			v := map[string]interface{}{}
			for key, i := range s {
				vv := L.EvaluateValue(ctx, app, i, object)
				if vv != nil {
					v[key] = vv
				}
			}
			return v
		}
	}

	{
		s, ok := value.(map[interface{}]interface{})

		if ok {
			v := map[string]interface{}{}
			for key, i := range s {
				vv := L.EvaluateValue(ctx, app, i, object)
				if vv != nil {
					v[dynamic.StringValue(key, "")] = vv
				}
			}
			return v
		}
	}

	return value
}

func (L *Logic) Get(ctx IContext, app IApp, key string) interface{} {
	return L.EvaluateValue(ctx, app, dynamic.Get(L.object, key), L.object)
}

func (L *Logic) Exec(ctx IContext, app IApp) error {
	log.Println("[EXEC]", app.Path(), ">>", L.Tag)
	return nil
}

func (L *Logic) On(name string, value interface{}) {
	L.on[name] = value
}

func (L *Logic) Has(name string) bool {
	_, ok := L.on[name]
	return ok
}

func (L *Logic) Error(ctx IContext, app IApp, err error) error {
	if L.Has("error") {
		ctx.Set(ErrorKeys, GetErrorObject(err))
		return L.Done(ctx, app, "error")
	}
	return err
}

func (L *Logic) Done(ctx IContext, app IApp, name string) error {

	log.Println("[DONE]", app.Path(), ">>", L.Tag, ">>", name)

	fn, ok := L.on[name]

	if ok {

		{
			s, ok := fn.(string)

			if ok {
				return app.Exec(ctx, s)
			}
		}

		{
			s, ok := fn.(LogicFunction)

			if ok {
				return s(ctx, app)
			}
		}

		{
			s, ok := fn.(ILogic)
			if ok {
				return s.Exec(ctx, app)
			}
		}
	}

	return nil

}
