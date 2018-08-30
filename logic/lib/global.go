package lib

import (
	"github.com/hailongz/kk-lib/dynamic"
	"github.com/hailongz/kk-logic/logic"
)

type GlobalLogic struct {
	logic.Logic
}

func (L *GlobalLogic) Exec(ctx logic.IContext, app logic.IApp) error {

	L.Logic.Exec(ctx, app)

	key := L.Get(ctx, app, "name")
	value := L.Get(ctx, app, "value")

	if key != nil {
		skey := dynamic.StringValue(key, "")
		ctx.SetGlobal(skey, value)
	}

	return L.Done(ctx, app, "done")
}

func init() {
	logic.Openlib("kk.Logic.Global", func(object interface{}) logic.ILogic {
		v := GlobalLogic{}
		v.Init(object)
		return &v
	})
}
