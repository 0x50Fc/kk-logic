package lib

import (
	"github.com/hailongz/kk-lib/dynamic"
	"github.com/hailongz/kk-logic/logic"
)

type ScriptLogic struct {
	logic.Logic
}

func (L *ScriptLogic) Exec(ctx logic.IContext, app logic.IApp) error {

	L.Logic.Exec(ctx, app)

	code := L.Get(ctx, app, "code")

	done := "done"

	if code != nil {

		err := ctx.Call(dynamic.StringValue(code, ""), L.Tag, func(name string) {
			done = name
		})

		if err != nil {
			return L.Error(ctx, app, err)
		}
	}

	return L.Done(ctx, app, done)
}

func init() {
	logic.Openlib("kk.Logic.Script", func(object interface{}) logic.ILogic {
		v := ScriptLogic{}
		v.Init(object)
		return &v
	})
}
