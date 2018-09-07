package lib

import (
	"github.com/hailongz/kk-lib/dynamic"
	"github.com/hailongz/kk-logic/logic"
)

type ContentLogic struct {
	logic.Logic
}

func (L *ContentLogic) Exec(ctx logic.IContext, app logic.IApp) error {

	L.Logic.Exec(ctx, app)

	contentType := dynamic.StringValue(L.Get(ctx, app, "contentType"), "")
	content := L.Get(ctx, app, "content")

	v := logic.Content{
		ContentType: contentType,
	}

	if content == nil {
		v.Content = []byte{}
	}

	if v.Content == nil {
		b, ok := content.([]byte)
		if ok {
			v.Content = b
		}
	}

	if v.Content == nil {
		b, ok := content.(string)
		if ok {
			v.Content = []byte(b)
		}
	}

	if v.Content == nil {
		v.Content = []byte{}
	}

	dynamic.Each(L.Get(ctx, app, "headers"), func(key interface{}, value interface{}) bool {

		v.Header[dynamic.StringValue(key, "")] = []string{dynamic.StringValue(value, "")}

		return true
	})

	return &v

}

func init() {
	logic.Openlib("kk.Logic.Content", func(object interface{}) logic.ILogic {
		v := ContentLogic{}
		v.Init(object)
		return &v
	})
}
