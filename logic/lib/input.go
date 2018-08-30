package lib

import (
	"fmt"

	"regexp"

	"github.com/hailongz/kk-lib/dynamic"
	"github.com/hailongz/kk-logic/logic"
)

/**
 * 验证输入参数
 */
type InputLogic struct {
	logic.Logic
}

func (L *InputLogic) Exec(ctx logic.IContext, app logic.IApp) error {

	L.Logic.Exec(ctx, app)

	var err error = nil

	inputData := ctx.Get([]string{"input"})
	input := dynamic.Get(app.Object(), "input")

	method := dynamic.Get(input, "method")

	if method != nil && dynamic.StringValue(method, "") != dynamic.StringValue(ctx.Get(logic.MethodKeys), "GET") {
		err = logic.NewError(logic.ERROR_INPUT, "不支持的方法")
	}

	if err == nil {
		dynamic.Each(dynamic.Get(input, "fields"), func(key interface{}, field interface{}) bool {

			name := dynamic.StringValue(dynamic.Get(field, "name"), "")
			errno := int(dynamic.IntValue(dynamic.Get(field, "errno"), logic.ERROR_INPUT))
			errmsg := dynamic.StringValue(dynamic.Get(field, "errmsg"), fmt.Sprintf("参数错误 %s", name))
			required := dynamic.BooleanValue(dynamic.Get(field, "required"), false)
			pattern := dynamic.StringValue(dynamic.Get(field, "pattern"), "")

			v := dynamic.Get(inputData, name)

			if required && dynamic.IsEmpty(v) {
				err = logic.NewError(errno, errmsg)
				return false
			}

			if pattern != "" {
				var r *regexp.Regexp = nil
				r, err = regexp.Compile(pattern)
				if err != nil {
					return false
				}
				if !r.MatchString(dynamic.StringValue(v, "")) {
					err = logic.NewError(errno, errmsg)
					return false
				}
			}

			return true
		})
	}

	if err != nil {
		return L.Error(ctx, app, err)
	}

	return L.Done(ctx, app, "done")
}

func init() {
	logic.Openlib("kk.Logic.Input", func(object interface{}) logic.ILogic {
		v := InputLogic{}
		v.Init(object)
		return &v
	})
}
