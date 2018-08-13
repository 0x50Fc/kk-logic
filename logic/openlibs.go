package logic

import (
	"log"

	"github.com/hailongz/kk-lib/duktape"
	"github.com/hailongz/kk-logic/assert"
)

func openlib(app IApp, ctx duktape.Context, name string) {

	data, ok := assert.Get(name)

	if ok {
		duktape.Compile(ctx, name, string(data))
		if duktape.IsFunction(ctx, -1) {
			err := duktape.Call(ctx, 0)
			if err != nil {
				log.Println("[LOGIC] [ERROR]", err)
			}
		}
		duktape.Pop(ctx, 1)
	}

}

func init() {

	AddProtocol(func(app IApp, ctx duktape.Context) {

		openlib(app, ctx, "require.js")
		openlib(app, ctx, "kk.js")
		openlib(app, ctx, "kk-date.js")
		openlib(app, ctx, "kk-logic-var.js")
		openlib(app, ctx, "kk-logic-http.js")
		openlib(app, ctx, "kk-logic-each.js")
		openlib(app, ctx, "kk-logic-outlet.js")
		openlib(app, ctx, "kk-logic-app.js")

	})

}
