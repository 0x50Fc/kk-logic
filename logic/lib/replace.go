package lib

import (
	"regexp"
	"strings"

	"github.com/hailongz/kk-lib/dynamic"
	"github.com/hailongz/kk-logic/logic"
)

type ReplaceLogic struct {
	logic.Logic
}

func (L *ReplaceLogic) Exec(ctx logic.IContext, app logic.IApp) error {

	L.Logic.Exec(ctx, app)

	key := dynamic.StringValue(L.Get(ctx, app, "key"), "")
	value := dynamic.StringValue(L.Get(ctx, app, "value"), "")
	items := L.Get(ctx, app, "items")

	var keys []string = nil

	if key == "" {
		keys = logic.RefererKeys
	} else {
		keys = strings.Split(key, ".")
	}

	var err error = nil

	dynamic.Each(items, func(key interface{}, item interface{}) bool {

		var pattern *regexp.Regexp = nil

		pattern, err = regexp.Compile(dynamic.StringValue(dynamic.Get(item, "pattern"), ""))

		if err != nil {
			return false
		}

		value = pattern.ReplaceAllString(value, dynamic.StringValue(dynamic.Get(item, "value"), ""))

		return true
	})

	if err != nil {
		return L.Error(ctx, app, err)
	}

	ctx.Set(keys, value)

	return L.Done(ctx, app, "done")
}

func init() {
	logic.Openlib("kk.Logic.Replace", func(object interface{}) logic.ILogic {
		v := ReplaceLogic{}
		v.Init(object)
		return &v
	})
}
