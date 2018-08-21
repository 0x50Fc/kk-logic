package lib

import (
	"strings"

	"github.com/hailongz/kk-lib/dynamic"
	"github.com/hailongz/kk-logic/logic"
)

type EachLogic struct {
	logic.Logic
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

				v[skey] = value
				return true
			}
		}

		value = L.EvaluateValue(ctx, app, value, object)

		v[skey] = value

		return true
	})

	return v
}

func (L *EachLogic) Exec(ctx logic.IContext, app logic.IApp) error {

	L.Logic.Exec(ctx, app)

	value := L.Get(ctx, app, "value")
	fields := dynamic.Get(L.Object(), "fields")

	if value == nil {
		ctx.Set(logic.ResultKeys, nil)
		return L.Done(ctx, app, "done")
	}

	{
		a, ok := value.([]interface{})
		if ok {
			v := []interface{}{}
			for i, vv := range a {
				ctx.Set(logic.KeyKeys, i)
				vv := L.item(ctx, app, vv, fields)
				if vv != nil {
					v = append(v, vv)
				}
			}
			ctx.Set(logic.ResultKeys, v)
			return L.Done(ctx, app, "done")
		}
	}

	{
		a, ok := value.(map[string]interface{})
		if ok {
			v := map[string]interface{}{}
			for key, vv := range a {
				ctx.Set(logic.KeyKeys, key)
				vv := L.item(ctx, app, vv, fields)
				if vv != nil {
					v[key] = vv
				}
			}
			ctx.Set(logic.ResultKeys, v)
			return L.Done(ctx, app, "done")
		}
	}

	return L.Done(ctx, app, "done")
}

func init() {
	logic.Openlib("kk.Logic.Each", func(object interface{}) logic.ILogic {
		v := EachLogic{}
		v.Init(object)
		return &v
	})
}
