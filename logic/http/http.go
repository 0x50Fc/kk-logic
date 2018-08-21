package http

import (
	"github.com/hailongz/kk-lib/dynamic"
	"github.com/hailongz/kk-lib/http"
	"github.com/hailongz/kk-logic/logic"
)

type HttpLogic struct {
	logic.Logic
}

func (L *HttpLogic) Exec(ctx logic.IContext, app logic.IApp) error {

	L.Logic.Exec(ctx, app)

	options := http.Options{}

	options.Method = dynamic.StringValue(L.Get(ctx, app, "method"), "GET")
	options.Url = dynamic.StringValue(L.Get(ctx, app, "url"), "")
	options.Type = dynamic.StringValue(L.Get(ctx, app, "type"), http.OptionTypeUrlencode)
	options.ResponseType = dynamic.StringValue(L.Get(ctx, app, "dataType"), http.OptionResponseTypeJson)
	options.Data = L.Get(ctx, app, "data")

	v, err := http.Send(&options)

	if err != nil {
		if L.Has("error") {
			ctx.Set(logic.ErrorKeys, logic.GetErrorObject(err))
			return L.Done(ctx, app, "error")
		}
		return err
	}

	ctx.Set(logic.ResultKeys, v)

	return L.Done(ctx, app, "done")
}

func init() {
	logic.Openlib("kk.Logic.Http", func(object interface{}) logic.ILogic {
		v := HttpLogic{}
		v.Init(object)
		return &v
	})
}
