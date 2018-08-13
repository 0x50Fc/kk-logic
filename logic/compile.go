package logic

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/hailongz/kk-lib/duktape"
)

func init() {

	AddProtocol(func(app IApp, ctx duktape.Context) {

		duktape.PushGlobalObject(ctx)

		duktape.PushString(ctx, "compile")
		duktape.PushGoFunction(ctx, func() int {

			var path string
			var prefix string
			var suffix string

			var top = duktape.GetTop(ctx)

			if top > 0 && duktape.IsString(ctx, -top) {

				path = duktape.ToString(ctx, -top)

				if top > 1 && duktape.IsString(ctx, -top+1) {
					prefix = duktape.ToString(ctx, -top+1)
				}

				if top > 2 && duktape.IsString(ctx, -top+2) {
					suffix = duktape.ToString(ctx, -top+2)
				}

				fd, err := os.Open(filepath.Join(app.Path(), path))

				if err != nil {
					duktape.PushString(ctx, err.Error())
					return 1
				}

				defer fd.Close()

				b, err := ioutil.ReadAll(fd)

				if err != nil {
					duktape.PushString(ctx, err.Error())
					return 1
				}

				duktape.Compile(ctx, path, fmt.Sprintf("%s%s%s", prefix, string(b), suffix))

				return 1
			}

			return 0
		})

		duktape.PutProp(ctx, -3)
		duktape.Pop(ctx, 1)
	})

}
