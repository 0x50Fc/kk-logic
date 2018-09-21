package lib

import (
	"log"

	"github.com/hailongz/kk-logic/logic"
)

type DumpLogic struct {
	logic.Logic
}

func (L *DumpLogic) Exec(ctx logic.IContext, app logic.IApp) error {

	L.Logic.Exec(ctx, app)

	log.Println("[DUMP]", ctx.Get(nil))

	return L.Done(ctx, app, "done")
}

func init() {
	logic.Openlib("kk.Logic.Dump", func(object interface{}) logic.ILogic {
		v := DumpLogic{}
		v.Init(object)
		return &v
	})
}
