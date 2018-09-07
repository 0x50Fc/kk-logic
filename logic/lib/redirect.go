package lib

import (
	"github.com/hailongz/kk-lib/dynamic"
	"github.com/hailongz/kk-logic/logic"
)

type RedirectLogic struct {
	logic.Logic
}

func (L *RedirectLogic) Exec(ctx logic.IContext, app logic.IApp) error {

	L.Logic.Exec(ctx, app)

	url := dynamic.StringValue(L.Get(ctx, app, "url"), "")

	return logic.NewRedirect(url)

}

func init() {
	logic.Openlib("kk.Logic.Redirect", func(object interface{}) logic.ILogic {
		v := RedirectLogic{}
		v.Init(object)
		return &v
	})
}
