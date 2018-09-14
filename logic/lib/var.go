package lib

import (
	"strings"

	"github.com/hailongz/kk-lib/dynamic"
	"github.com/hailongz/kk-lib/json"
	"github.com/hailongz/kk-logic/logic"
	"gopkg.in/yaml.v2"
)

type VarLogic struct {
	logic.Logic
}

func (L *VarLogic) Exec(ctx logic.IContext, app logic.IApp) error {

	L.Logic.Exec(ctx, app)

	key := L.Get(ctx, app, "key")
	value := L.Get(ctx, app, "value")
	stype := dynamic.StringValue(L.Get(ctx, app, "type"), "")

	if key != nil {

		switch stype {
		case "json":
			b, err := json.Marshal(value)
			if err != nil {
				return L.Error(ctx, app, err)
			}
			value = string(b)
			break
		case "yaml":
			b, err := yaml.Marshal(value)
			if err != nil {
				return L.Error(ctx, app, err)
			}
			value = string(b)
			break
		}
		skey := dynamic.StringValue(key, "")
		ctx.Set(strings.Split(skey, "."), value)
	}

	return L.Done(ctx, app, "done")
}

func init() {
	logic.Openlib("kk.Logic.Var", func(object interface{}) logic.ILogic {
		v := VarLogic{}
		v.Init(object)
		return &v
	})
}
