package logic

import (
	"errors"
	"fmt"
	"log"
	"reflect"
	"strconv"

	"github.com/hailongz/kk-lib/dynamic"
	"gopkg.in/olebedev/go-duktape.v3"
)

var InputKeys = []string{"input"}
var ObjectKeys = []string{"object"}
var ResultKeys = []string{"result"}
var OutputKeys = []string{"output"}
var ViewKeys = []string{"view"}
var SessionIdKeys = []string{"sessionId"}
var UrlKeys = []string{"url"}
var URLKeys = []string{"URL"}
var ClientIpKeys = []string{"clientIp"}
var ContentKeys = []string{"content"}
var RefererKeys = []string{"referer"}
var UserAgentKeys = []string{"userAgent"}
var HeadersKeys = []string{"headers"}
var ParamsKeys = []string{"params"}
var ErrorKeys = []string{"error"}
var KeyKeys = []string{"key"}
var MethodKeys = []string{"method"}

const (
	DUK_DEFPROP_WRITABLE          = (1 << 0)
	DUK_DEFPROP_ENUMERABLE        = (1 << 1)
	DUK_DEFPROP_CONFIGURABLE      = (1 << 2)
	DUK_DEFPROP_HAVE_WRITABLE     = (1 << 3)
	DUK_DEFPROP_HAVE_ENUMERABLE   = (1 << 4)
	DUK_DEFPROP_HAVE_CONFIGURABLE = (1 << 5)
	DUK_DEFPROP_HAVE_VALUE        = (1 << 6)
	DUK_DEFPROP_HAVE_GETTER       = (1 << 7)
	DUK_DEFPROP_HAVE_SETTER       = (1 << 8)
	DUK_DEFPROP_FORCE             = (1 << 9)
)

type Function int64

type IContext interface {
	Begin()
	End()
	Get(keys []string) interface{}
	Set(keys []string, value interface{})
	SetGlobal(key string, value interface{})
	Evaluate(evaluateCode string, name string) interface{}
	Call(evaluateCode string, name string, done func(name string)) error
	Recycle()
	AddRecycle(fn func())
}

type Context struct {
	jsContext *duktape.Context
	scopes    []interface{}
	recycles  []func()
}

func NewContext() IContext {
	v := Context{}
	v.jsContext = duktape.New()
	v.scopes = []interface{}{map[string]interface{}{}}
	v.recycles = []func(){}
	return &v
}

func (C *Context) AddRecycle(fn func()) {
	C.recycles = append(C.recycles, fn)
}

func (C *Context) Begin() {
	v := C.scopes[len(C.scopes)-1]
	object := map[string]interface{}{}
	dynamic.Each(v, func(key interface{}, value interface{}) bool {
		dynamic.Set(object, dynamic.StringValue(key, ""), value)
		return true
	})
	C.scopes = append(C.scopes, object)
}

func (C *Context) End() {
	C.scopes = C.scopes[0 : len(C.scopes)-1]
}

func (C *Context) Get(keys []string) interface{} {

	if keys == nil || len(keys) == 0 {
		return C.scopes[len(C.scopes)-1]
	}

	key := keys[0]
	i := len(C.scopes) - 1

	var object interface{} = nil

	for i >= 0 {

		object = C.scopes[i]

		if dynamic.Get(object, key) != nil {
			break
		}

		i = i - 1
	}

	return dynamic.GetWithKeys(object, keys)
}

func (C *Context) SetGlobal(key string, value interface{}) {
	object := C.scopes[0]
	dynamic.Set(object, key, value)
}

func (C *Context) Set(keys []string, value interface{}) {
	object := C.scopes[len(C.scopes)-1]
	dynamic.SetWithKeys(object, keys, value)
}

func pushValue(jsContext *duktape.Context, value interface{}) {

	if value == nil {
		jsContext.PushUndefined()
		return
	}

	v := reflect.ValueOf(value)

	switch v.Kind() {
	case reflect.Int, reflect.Int16, reflect.Int8, reflect.Int64, reflect.Int32:
		vv := v.Int()
		iv := int(vv)
		if vv == int64(iv) {
			jsContext.PushInt(iv)
		} else {
			jsContext.PushString(strconv.FormatInt(vv, 10))
		}
		return
	case reflect.Uint, reflect.Uint16, reflect.Uint8, reflect.Uint64, reflect.Uint32:
		vv := v.Uint()
		iv := uint(vv)
		if vv == uint64(iv) {
			jsContext.PushUint(iv)
		} else {
			jsContext.PushString(strconv.FormatUint(vv, 10))
		}
		return
	case reflect.Bool:
		jsContext.PushBoolean(v.Bool())
		return
	case reflect.Float32, reflect.Float64:
		jsContext.PushNumber(v.Float())
		return
	case reflect.String:
		jsContext.PushString(v.String())
		return
	}

	{
		a, ok := value.([]interface{})
		if ok {
			jsContext.PushArray()
			for i, v := range a {
				jsContext.PushInt(i)
				pushValue(jsContext, v)
				jsContext.PutProp(-3)
			}
			return
		}
	}

	pushObject(jsContext, value)
}

func pushObject(jsContext *duktape.Context, object interface{}) {

	jsContext.PushObject()

	dynamic.Each(object, func(key interface{}, value interface{}) bool {

		jsContext.PushString(dynamic.StringValue(key, ""))

		jsContext.PushGoFunction(func(jsContext *duktape.Context) int {
			pushValue(jsContext, value)
			return 1
		})

		jsContext.DefProp(-3, uint(DUK_DEFPROP_HAVE_GETTER|DUK_DEFPROP_HAVE_ENUMERABLE|DUK_DEFPROP_ENUMERABLE))

		return true
	})

}

func pushContext(jsContext *duktape.Context, ctx IContext) {

	jsContext.PushObject()

	jsContext.PushString("get")
	jsContext.PushGoFunction(func(jsContext *duktape.Context) int {

		top := jsContext.GetTop()

		if top > 0 {

			keys := []string{}

			if jsContext.IsArray(-top) {
				jsContext.Enum(-1, duktape.EnumArrayIndicesOnly)
				for jsContext.Next(-1, true) {
					keys = append(keys, dynamic.StringValue(toValue(jsContext, -1), ""))
					jsContext.Pop2()
				}
				jsContext.Pop()
			} else if jsContext.IsString(-top) {
				keys = append(keys, jsContext.ToString(-top))
			}

			v := ctx.Get(keys)

			pushValue(jsContext, v)

			return 1
		}

		return 0
	})
	jsContext.PutProp(-3)

	jsContext.PushString("set")
	jsContext.PushGoFunction(func(jsContext *duktape.Context) int {

		top := jsContext.GetTop()

		if top > 0 {

			keys := []string{}

			if jsContext.IsArray(-top) {
				jsContext.Enum(-1, duktape.EnumArrayIndicesOnly)
				for jsContext.Next(-1, true) {
					keys = append(keys, dynamic.StringValue(toValue(jsContext, -1), ""))
					jsContext.Pop2()
				}
				jsContext.Pop()
			} else if jsContext.IsString(-top) {
				keys = append(keys, jsContext.ToString(-top))
			}

			var v interface{} = nil

			if top > 1 {
				v = toValue(jsContext, -top+1)
			}

			ctx.Set(keys, v)

		}

		return 0
	})
	jsContext.PutProp(-3)

}

func toValue(jsContext *duktape.Context, idx int) interface{} {
	if jsContext.IsNumber(idx) {
		v := jsContext.ToNumber(idx)
		iv := int64(v)
		if v == float64(iv) {
			return iv
		}
		return v
	} else if jsContext.IsBoolean(idx) {
		return jsContext.ToBoolean(idx)
	} else if jsContext.IsString(idx) {
		return jsContext.ToString(idx)
	} else if jsContext.IsArray(idx) {
		v := []interface{}{}
		jsContext.Enum(idx, duktape.EnumArrayIndicesOnly)
		for jsContext.Next(-1, true) {
			vv := toValue(jsContext, -1)
			if vv != nil {
				v = append(v, vv)
			}
			jsContext.Pop2()
		}
		jsContext.Pop()
		return v
	} else if jsContext.IsObject(idx) {

		v := map[string]interface{}{}

		jsContext.Enum(idx, duktape.EnumIncludeInternal)

		for jsContext.Next(-1, true) {
			key := toValue(jsContext, -2)
			vv := toValue(jsContext, -1)
			if key != nil && vv != nil {
				v[dynamic.StringValue(key, "")] = vv
			}
			jsContext.Pop2()
		}

		jsContext.Pop()

		return v
	}

	return nil
}

func getErrorString(jsContext *duktape.Context, idx int) string {

	fileName := ""
	lineNumber := 0
	stack := ""

	jsContext.GetPropString(idx, "fileName")
	fileName = jsContext.ToString(-1)
	jsContext.Pop()

	jsContext.GetPropString(idx, "lineNumber")
	lineNumber = jsContext.ToInt(-1)
	jsContext.Pop()

	jsContext.GetPropString(idx, "stack")
	stack = jsContext.ToString(-1)
	jsContext.Pop()

	return fmt.Sprintf("%s(%d): %s\n", fileName, lineNumber, stack)
}

func dumpError(jsContext *duktape.Context, tag string, idx int) {

	log.Printf("%s %s\n", tag, getErrorString(jsContext, idx))

}

func (C *Context) Call(evaluateCode string, name string, done func(name string)) error {

	var err error = nil

	C.jsContext.PushString(name)
	C.jsContext.CompileStringFilename(0, fmt.Sprintf("(function(ctx,done){ %s })", evaluateCode))

	if C.jsContext.IsFunction(-1) {

		if C.jsContext.Pcall(0) == duktape.ExecSuccess {

			if C.jsContext.IsFunction(-1) {

				pushContext(C.jsContext, C)

				C.jsContext.PushGoFunction(func(jsContext *duktape.Context) int {
					top := jsContext.GetTop()

					if top > 0 && jsContext.IsString(-top) {
						done(jsContext.ToString(-top))
					}

					return 0
				})

				if C.jsContext.Pcall(2) == duktape.ExecSuccess {

				} else {
					dumpError(C.jsContext, "[CONTEXT] [CALL]", -1)
					err = errors.New(getErrorString(C.jsContext, -1))
				}

			}

		} else {
			dumpError(C.jsContext, "[CONTEXT] [CALL]", -1)
			err = errors.New(getErrorString(C.jsContext, -1))
		}

	}

	C.jsContext.Pop()

	return err
}

func (C *Context) Evaluate(evaluateCode string, name string) interface{} {

	var v interface{} = nil

	C.jsContext.PushString(name)
	C.jsContext.CompileStringFilename(0, "(function(object,evaluate){ var _G; with (object) { _G = eval('(' + evaluate + ')'); } return _G;})")

	if C.jsContext.IsFunction(-1) {

		if C.jsContext.Pcall(0) == duktape.ExecSuccess {

			pushObject(C.jsContext, C.scopes[len(C.scopes)-1])
			C.jsContext.PushString(evaluateCode)

			if C.jsContext.Pcall(2) == duktape.ExecSuccess {
				v = toValue(C.jsContext, -1)
			} else {
				dumpError(C.jsContext, "[CONTEXT] [Evaluate]", -1)
			}

		} else {
			dumpError(C.jsContext, "[CONTEXT] [Evaluate]", -1)
		}
	}

	C.jsContext.Pop()

	return v

}

func (C *Context) Recycle() {
	C.jsContext.DestroyHeap()
	C.jsContext = nil
	C.scopes = nil
	if C.recycles != nil {
		for _, fn := range C.recycles {
			fn()
		}
		C.recycles = nil
	}
}
