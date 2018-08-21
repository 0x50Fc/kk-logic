package http

import (
	"bytes"

	"github.com/dchest/captcha"
	"github.com/hailongz/kk-lib/dynamic"
	"github.com/hailongz/kk-logic/logic"
)

type CaptchaLogic struct {
	logic.Logic
}

func (L *CaptchaLogic) Exec(ctx logic.IContext, app logic.IApp) error {

	L.Logic.Exec(ctx, app)

	length := dynamic.IntValue(L.Get(ctx, app, "length"), 4)
	width := dynamic.IntValue(L.Get(ctx, app, "width"), 200)
	height := dynamic.IntValue(L.Get(ctx, app, "height"), 80)

	id := captcha.NewLen(int(length))
	b := bytes.NewBuffer(nil)
	captcha.WriteImage(b, id, int(width), int(height))

	ctx.Set(logic.ViewKeys, &logic.View{Content: b.Bytes(), ContentType: "image/png"})
	ctx.Set(logic.ResultKeys, id)

	return L.Done(ctx, app, "done")
}

type CaptchaVerifyLogic struct {
	logic.Logic
}

func (L *CaptchaVerifyLogic) Exec(ctx logic.IContext, app logic.IApp) error {

	L.Logic.Exec(ctx, app)

	id := dynamic.StringValue(L.Get(ctx, app, "id"), "")
	code := dynamic.StringValue(L.Get(ctx, app, "code"), "")
	errno := dynamic.IntValue(L.Get(ctx, app, "errno"), logic.ERROR_UNKNOWN)
	errmsg := dynamic.StringValue(L.Get(ctx, app, "errmsg"), "验证码错误")

	if captcha.VerifyString(id, code) {
		return L.Done(ctx, app, "done")
	}

	return L.Error(ctx, app, logic.NewError(int(errno), errmsg))
}

func init() {
	logic.Openlib("kk.Logic.Captcha", func(object interface{}) logic.ILogic {
		v := CaptchaLogic{}
		v.Init(object)
		return &v
	})
	logic.Openlib("kk.Logic.Captcha.Verify", func(object interface{}) logic.ILogic {
		v := CaptchaVerifyLogic{}
		v.Init(object)
		return &v
	})
}
