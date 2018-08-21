package lib

import (
	"strings"

	"github.com/hailongz/kk-lib/dynamic"
	"github.com/hailongz/kk-logic/logic"
)

type VarLogic struct {
	logic.Logic
}

func (L *VarLogic) Exec(ctx logic.IContext, app logic.IApp) error {

	L.Logic.Exec(ctx, app)

	key := L.Get(ctx, app, "key")
	value := L.Get(ctx, app, "value")

	if key != nil {
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
