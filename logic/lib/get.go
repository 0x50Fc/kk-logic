package lib

import (
	"strings"

	"github.com/hailongz/kk-lib/dynamic"
	"github.com/hailongz/kk-lib/json"
	"github.com/hailongz/kk-logic/logic"
	"gopkg.in/yaml.v2"
)

type GetLogic struct {
	logic.Logic
}

func (L *GetLogic) Exec(ctx logic.IContext, app logic.IApp) error {

	L.Logic.Exec(ctx, app)

	key := dynamic.StringValue(L.Get(ctx, app, "key"), "")
	sourceType := dynamic.StringValue(L.Get(ctx, app, "sourceType"), "")

	var keys []string = nil

	if key == "" {
		keys = logic.RefererKeys
	} else {
		keys = strings.Split(key, ".")
	}

	var err error = nil
	var v []byte = nil

	path := L.Get(ctx, app, "path")

	if path != nil {

		s := app.Store()

		if s != nil {
			v, err = s.GetContent(dynamic.StringValue(path, ""))
			if err != nil {
				return L.Error(ctx, app, err)
			}
		}
	}

	if v == nil {
		ctx.Set(keys, nil)
		return L.Done(ctx, app, "done")
	}

	switch sourceType {
	case "json":
		var value interface{} = nil
		err := json.Unmarshal(v, &value)
		if err != nil {
			return L.Error(ctx, app, err)
		}
		ctx.Set(keys, value)
		break
	case "yaml":
		var value interface{} = nil
		err := yaml.Unmarshal(v, &value)
		if err != nil {
			return L.Error(ctx, app, err)
		}
		ctx.Set(keys, value)
		break
	case "text":
		ctx.Set(keys, string(v))
		break
	default:
		ctx.Set(keys, v)
		break
	}

	return L.Done(ctx, app, "done")
}

func init() {
	logic.Openlib("kk.Logic.Get", func(object interface{}) logic.ILogic {
		v := GetLogic{}
		v.Init(object)
		return &v
	})
}
