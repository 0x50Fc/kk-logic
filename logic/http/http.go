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

	checkType := dynamic.StringValue(L.Get(ctx, app, "checkType"), "errno")

	options := http.Options{}

	options.Method = dynamic.StringValue(L.Get(ctx, app, "method"), "GET")
	options.Url = dynamic.StringValue(L.Get(ctx, app, "url"), "")
	options.Type = dynamic.StringValue(L.Get(ctx, app, "type"), http.OptionTypeUrlencode)
	options.ResponseType = dynamic.StringValue(L.Get(ctx, app, "dataType"), http.OptionResponseTypeJson)
	options.Data = L.Get(ctx, app, "data")

	v, err := http.Send(&options)

	if err != nil {
		return L.Error(ctx, app, err)
	}

	if checkType == "errno" {
		if dynamic.Get(v, "errno") != nil {
			errno := int(dynamic.IntValue(dynamic.Get(v, "errno"), logic.ERROR_UNKNOWN))
			errmsg := dynamic.StringValue(dynamic.Get(v, "errmsg"), "未知错误")
			return L.Error(ctx, app, logic.NewError(errno, errmsg))
		}
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
