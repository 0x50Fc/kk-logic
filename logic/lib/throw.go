package lib

import (
	"github.com/hailongz/kk-lib/dynamic"
	"github.com/hailongz/kk-logic/logic"
)

type ThrowLogic struct {
	logic.Logic
}

func (L *ThrowLogic) Exec(ctx logic.IContext, app logic.IApp) error {

	L.Logic.Exec(ctx, app)

	errno := L.Get(ctx, app, "errno")
	errmsg := L.Get(ctx, app, "errmsg")

	return logic.NewError(int(dynamic.IntValue(errno, logic.ERROR_UNKNOWN)), dynamic.StringValue(errmsg, "未知错误"))

}

func init() {
	logic.Openlib("kk.Logic.Throw", func(object interface{}) logic.ILogic {
		v := ThrowLogic{}
		v.Init(object)
		return &v
	})
}
