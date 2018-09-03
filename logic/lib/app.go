package lib

import (
	"github.com/hailongz/kk-lib/dynamic"
	"github.com/hailongz/kk-logic/logic"
)

type AppLogic struct {
	logic.Logic
	app  logic.IApp
	done string
}

func init() {
	logic.Openlib("kk.Logic.App", func(object interface{}) logic.ILogic {
		v := AppLogic{}
		v.Init(object)
		return &v
	})
}

func (L *AppLogic) Exec(ctx logic.IContext, app logic.IApp) error {
	L.Logic.Exec(ctx, app)

	if L.app == nil {

		var err error = nil
		path := dynamic.StringValue(L.Get(ctx, app, "path"), "")
		L.app, err = app.Store().Get(path)
		if err != nil {
			return err
		}

		L.app.Each(func(name string, v logic.ILogic) bool {

			outlet, ok := v.(*OutletLogic)

			if ok {
				outlet.OnOutlet = func(ctx logic.IContext, app logic.IApp) error {
					L.done = name
					return nil
				}
			}

			return true
		})
	}

	if L.app != nil {

		L.done = "done"

		params := L.Get(ctx, app, logic.ParamsKeys[0])
		output := ctx.Get(logic.OutputKeys)
		ctx.Begin()
		ctx.Set(logic.ParamsKeys, params)
		ctx.Set(logic.ResultKeys, logic.Nil)
		ctx.Set(logic.OutputKeys, output)
		err := L.app.Exec(ctx, "in")
		v := ctx.Get(logic.ResultKeys)
		if v == logic.Nil {
			v = nil
		}
		ctx.End()

		if err != nil {
			if L.Has("error") {
				ctx.Set(logic.ErrorKeys, logic.GetErrorObject(err))
				return L.Done(ctx, app, "error")
			}
			return err
		}

		if v != nil {
			ctx.Set(logic.ResultKeys, v)
		}

		return L.Done(ctx, app, L.done)
	}

	return L.Done(ctx, app, "done")
}
