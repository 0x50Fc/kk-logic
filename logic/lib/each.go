package lib

import (
	"reflect"
	"strings"

	"github.com/hailongz/kk-lib/dynamic"
	"github.com/hailongz/kk-logic/logic"
)

type EachLogic struct {
	logic.Logic
}

func setValue(v map[string]interface{}, key string, value interface{}) {
	if key == "_" {
		dynamic.Each(value, func(k interface{}, vv interface{}) bool {
			v[dynamic.StringValue(k, "")] = vv
			return true
		})
	} else if value != nil {
		v[key] = value
	}
}

func (L *EachLogic) item(ctx logic.IContext, app logic.IApp, object interface{}, fields interface{}) interface{} {

	if fields == nil {
		return object
	}

	v := map[string]interface{}{}

	dynamic.Each(fields, func(key interface{}, value interface{}) bool {

		skey := dynamic.StringValue(key, "")

		{
			s, ok := value.(string)
			if ok {
				if strings.HasPrefix(s, "=") {
					ctx.Begin()
					ctx.Set(logic.ObjectKeys, object)
					value = L.EvaluateValue(ctx, app, value, object)
					ctx.End()
				} else {
					value = dynamic.Get(object, s)
				}

				setValue(v, skey, value)

				return true
			}
		}

		value = L.EvaluateValue(ctx, app, value, object)

		setValue(v, skey, value)

		return true
	})

	return v
}

func (L *EachLogic) Exec(ctx logic.IContext, app logic.IApp) error {

	L.Logic.Exec(ctx, app)

	key := dynamic.StringValue(L.Get(ctx, app, "key"), "")
	value := L.Get(ctx, app, "value")
	fields := dynamic.Get(L.Object(), "fields")

	var keys []string = nil

	if key == "" {
		keys = logic.RefererKeys
	} else {
		keys = strings.Split(key, ".")
	}

	if value == nil {
		ctx.Set(keys, nil)
		return L.Done(ctx, app, "done")
	}

	{
		a, ok := value.([]interface{})
		if ok {
			v := []interface{}{}
			for i, vv := range a {
				ctx.Set(logic.KeyKeys, i)
				vv = L.item(ctx, app, vv, fields)
				if vv != nil {
					v = append(v, vv)
				}
			}
			ctx.Set(keys, v)
			return L.Done(ctx, app, "done")
		}
	}

	{
		switch reflect.ValueOf(value).Kind() {
		case reflect.Slice:
			v := []interface{}{}
			dynamic.Each(value, func(key interface{}, vv interface{}) bool {
				ctx.Set(logic.KeyKeys, key)
				vv = L.item(ctx, app, vv, fields)
				if vv != nil {
					v = append(v, vv)
				}
				return true
			})
			ctx.Set(keys, v)
			return L.Done(ctx, app, "done")
		}
	}

	v := L.item(ctx, app, value, fields)

	ctx.Set(keys, v)

	return L.Done(ctx, app, "done")
}

func init() {
	logic.Openlib("kk.Logic.Each", func(object interface{}) logic.ILogic {
		v := EachLogic{}
		v.Init(object)
		return &v
	})
}
