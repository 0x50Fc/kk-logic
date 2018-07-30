package duktape

import (
	"unsafe"
	"fmt"
	"errors"
	"log"
)

/*
#cgo CFLAGS: -std=c99 -Os -fomit-frame-pointer -fstrict-aliasing
#cgo linux LDFLAGS: -lm
#cgo freebsd LDFLAGS: -lm
#include "duk_config.h"
#include "duktape.h"
#include "duk.h"
extern duk_ret_t go_ObjectDealloc(duk_context * ctx);
extern duk_ret_t go_FunctionCall(duk_context * ctx);
*/
import "C"

type Context *C.struct_Context

type Scope struct {
	objects map[string]interface{}
	autoId uint64
}

func NewScope() *Scope {
	return &Scope{map[string]interface{}{},0}
}

func (O *Scope) newObject(object interface{}) string {
	iid := O.autoId + 1
	O.autoId = iid
	id := fmt.Sprintf("%d",iid)
 	O.objects[id] = object
	return id
}

func (O *Scope) removeObject(id string) {
	delete(O.objects,id)
}


func New(scope *Scope) Context {

	ctx := C.NewContext()
	ctx.scope = unsafe.Pointer(scope)

	C.duk_push_global_object(ctx.ctx);
	key := C.CString("__Scope");
	C.duk_push_string(ctx.ctx,key);
	C.duk_push_pointer(ctx.ctx,nil)
	C.duk_put_prop(ctx.ctx,C.duk_idx_t(-3));
	C.duk_pop(ctx.ctx);
	C.free(unsafe.Pointer(key))
	return ctx
}

func Recycle(ctx Context) {
	C.RecycleContext(ctx)
}

func getScope(ctx * C.duk_context) *Scope {

	var v * Scope = nil

	C.duk_push_global_object(ctx);
	key := C.CString("__Context");
	C.duk_push_string(ctx,key);
	C.free(unsafe.Pointer(key))
	C.duk_get_prop(ctx,C.duk_idx_t(-2));

	if(C.duk_is_pointer(ctx,C.duk_idx_t(-1)) != C.duk_bool_t(0)) {
		ctx := (Context)( C.duk_to_pointer(ctx,C.duk_idx_t(-1)))
		v = (*Scope)(ctx.scope)
	}

	C.duk_pop_n(ctx,2);

	return v
}

func PushGlobalObject(ctx Context) {
	C.duk_push_global_object(ctx.ctx);
}

func  PushObject(ctx Context) {
	C.duk_push_object(ctx.ctx);
}

func  PushArray(ctx Context) {
	C.duk_push_array(ctx.ctx);
}

func  PushInt(ctx Context,value int) {
	C.duk_push_int(ctx.ctx,C.duk_int_t(value));
}

func  PushUndefined(ctx Context) {
	C.duk_push_undefined(ctx.ctx);
}

func  PushNull(ctx Context) {
	C.duk_push_null(ctx.ctx);
}

func  PushNumber(ctx Context,value float64) {
	C.duk_push_number(ctx.ctx,C.duk_double_t(value));
}

func  PushBoolean(ctx Context,value bool) {
	if(value) {
		C.duk_push_boolean(ctx.ctx,C.duk_bool_t(1));
	} else {
		C.duk_push_boolean(ctx.ctx,C.duk_bool_t(0));
	}
}

func  PushString(ctx Context,value string) {
	key := C.CString(value)
	C.duk_push_string(ctx.ctx,key);
	C.free(unsafe.Pointer(key))
}

func PushGoObject(ctx Context,object interface{}) {
	
	if(object == nil) {
		C.duk_push_undefined(ctx.ctx)
		return;
	}

	scope := getScope(ctx.ctx)

	if scope == nil {
		C.duk_push_undefined(ctx.ctx)
		return
	}

	id := scope.newObject(object)
	C.duk_push_object(ctx.ctx)
	
	PushString(ctx,"__id")
	PushString(ctx,id)
	C.duk_put_prop(ctx.ctx,-3)
	
	C.duk_push_c_function(ctx.ctx,C.duk_c_function(C.go_ObjectDealloc),C.DUK_VARARGS);
	C.duk_set_finalizer(ctx.ctx,-2)
}

func  PushGoFunction(ctx Context,fn func() int) {

	if(fn == nil) {
		C.duk_push_undefined(ctx.ctx)
		return;
	}

	scope := getScope(ctx.ctx)

	if scope == nil {
		C.duk_push_undefined(ctx.ctx)
		return
	}

	id := scope.newObject(fn)
	C.duk_push_c_function(ctx.ctx,C.duk_c_function(C.go_FunctionCall),C.DUK_VARARGS);
	
	PushString(ctx,"__id")
	PushString(ctx,id)
	C.duk_put_prop(ctx.ctx,-3)

	C.duk_push_c_function(ctx.ctx,C.duk_c_function(C.go_ObjectDealloc),C.DUK_VARARGS);
	C.duk_set_finalizer(ctx.ctx,-2)
}

func PushBytes(ctx Context,bytes []byte) {

	size := C.duk_size_t(len(bytes))

	b := C.duk_push_buffer_raw(ctx.ctx,size,0)

	C.memcpy(b,C.CBytes(bytes),C.size_t(size))

	C.duk_push_buffer_object(ctx.ctx,-1,0,size,C.DUK_BUFOBJ_UINT8ARRAY);

	C.duk_remove(ctx.ctx,-2)

}

func Remove(ctx Context,idx int) {
	C.duk_remove(ctx.ctx,C.duk_idx_t(idx))
}

func Dup(ctx Context,idx int) {
	C.duk_dup(ctx.ctx,C.duk_idx_t(idx))
}

func  ToInt(ctx Context,idx int) int {
	return int(C.duk_to_int(ctx.ctx,C.duk_idx_t(idx)));
}

func  ToNumber(ctx Context,idx int) float64 {
	return float64(C.duk_to_number(ctx.ctx,C.duk_idx_t(idx)));
}

func  ToBoolean(ctx Context,idx int) bool {
	return C.duk_to_boolean(ctx.ctx,C.duk_idx_t(idx)) != C.duk_bool_t(0);
}

func  ToString(ctx Context,idx int) string {
	s := C.duk_to_string(ctx.ctx,C.duk_idx_t(idx));
	if s != nil {
		return C.GoString(s);
	}
	return ""
}

func  ToBytes(ctx Context,idx int) []byte {
	n := C.duk_size_t(0);
	b := C.duk_get_buffer_data(ctx.ctx,C.duk_idx_t(idx),&n);
	return C.GoBytes(b,C.int(n))
}

func  ToGoObject(ctx Context,idx int) interface{} {
	
	scope := getScope(ctx.ctx)

	if(scope == nil) {
		return nil
	}

	PushString(ctx,"__id");
	GetProp(ctx,idx -1);
	
	if IsString(ctx,-1) {

		id := ToString(ctx,-1)
		Pop(ctx,1)
		v ,ok := scope.objects[id]
		if ok {
			return v
		}
	}

	Pop(ctx,1)

	return nil
}


func  PutProp(ctx Context,idx int) {
	C.duk_put_prop(ctx.ctx,C.duk_idx_t(idx))
}

func  GetProp(ctx Context,idx int) {
	C.duk_get_prop(ctx.ctx,C.duk_idx_t(idx))
}

func  IsUndefined(ctx Context,idx int) bool {
	return C.duk_is_undefined(ctx.ctx,C.duk_idx_t(idx)) != C.duk_bool_t(0);
}

func  IsNull(ctx Context,idx int) bool {
	return C.duk_is_null(ctx.ctx,C.duk_idx_t(idx)) != C.duk_bool_t(0);
}

func  IsObject(ctx Context,idx int) bool {
	return C.duk_is_object(ctx.ctx,C.duk_idx_t(idx)) != C.duk_bool_t(0);
}

func  IsArray(ctx Context,idx int) bool {
	return C.duk_is_array(ctx.ctx,C.duk_idx_t(idx)) != C.duk_bool_t(0);
}

func  IsString(ctx Context,idx int) bool {
	return C.duk_is_string(ctx.ctx,C.duk_idx_t(idx)) != C.duk_bool_t(0);
}

func  IsNumber(ctx Context,idx int) bool {
	return C.duk_is_number(ctx.ctx,C.duk_idx_t(idx)) != C.duk_bool_t(0);
}

func  IsBoolean(ctx Context,idx int) bool {
	return C.duk_is_boolean(ctx.ctx,C.duk_idx_t(idx)) != C.duk_bool_t(0);
}

func  IsBytes(ctx Context,idx int) bool {
	return C.duk_is_buffer_data(ctx.ctx,C.duk_idx_t(idx)) != C.duk_bool_t(0);
}

func  IsFunction(ctx Context,idx int) bool {
	return C.duk_is_function(ctx.ctx,C.duk_idx_t(idx)) != C.duk_bool_t(0);
}

func Call(ctx Context,n int) error {
	
	if C.duk_pcall(ctx.ctx,C.duk_idx_t(n)) == C.DUK_EXEC_SUCCESS  {
		return nil
	}

	PushString(ctx,"lineNumber");
	GetProp(ctx,-2)
	line := ToInt(ctx,-1);
	Pop(ctx,1)

	PushString(ctx,"stack");
	GetProp(ctx,-2)
	stack := ToString(ctx,-1);
	Pop(ctx,1)

	PushString(ctx,"fileName");
	GetProp(ctx,-2)
	fileName := ToString(ctx,-1);
	Pop(ctx,1)

	return errors.New(fmt.Sprintf("%s(%d): %s",fileName,line,stack))
}

func  Pop(ctx Context,n int) {
	C.duk_pop_n(ctx.ctx,C.duk_idx_t(n))
}

func  Enum(ctx Context,idx int) {
	C.duk_enum(ctx.ctx,C.duk_idx_t(idx),C.DUK_ENUM_INCLUDE_SYMBOLS)
}

func  Next(ctx Context,idx int) bool {
	return C.duk_next(ctx.ctx,C.duk_idx_t(idx),1) != C.duk_bool_t(0)
}

func  JsonEncode(ctx Context,idx int) string{
	v := C.duk_json_encode(ctx.ctx,C.duk_idx_t(idx))
	return C.GoString(v);
}

func  JsonDecode(ctx Context,idx int) {
	C.duk_json_decode(ctx.ctx,C.duk_idx_t(idx))
}

func  PushValue(ctx Context,value interface{}) {
	if value == nil {
		PushUndefined(ctx)
		return;
	} 

	{
		v,ok := value.(bool)
		if ok {
			PushBoolean(ctx,v)
			return;
		}
	}

	{
		v,ok := value.(int)
		if ok {
			PushInt(ctx,v)
			return;
		}
	}

	{
		v,ok := value.(float64)
		if ok {
			PushNumber(ctx,v)
			return;
		}
	}

	{
		v,ok := value.(string)
		if ok {
			PushString(ctx,v)
			return;
		}
	}

	{
		v,ok := value.([]byte)
		if ok {
			PushBytes(ctx,v)
			return;
		}
	}

	{
		m,ok := value.(map[string]interface{})
		if ok {
			PushObject(ctx,)
			for key,v := range m {
				PushString(ctx,key)
				PushValue(ctx,v)
				PutProp(ctx,-3)
			}
 			return;
		}
	}

	{
		a,ok := value.([]interface{})
		if ok {
			PushObject(ctx)
			for i,v := range a {
				PushInt(ctx,i)
				PushValue(ctx,v)
				PutProp(ctx,-3)
			}
 			return;
		}
	}

	PushUndefined(ctx)
}

func  ToValue(ctx Context,idx int) interface{} {
	
	if IsNumber(ctx,idx) {
		v := ToNumber(ctx,idx)
		iv := int64(v)
		if float64(iv) ==v {
			return iv
		}
		return v
	} else if IsBoolean(ctx,idx) {
		return ToBoolean(ctx,idx)
	} else if IsString(ctx,idx) {
		return ToString(ctx,idx)
	} else if IsArray(ctx,idx) {
		a := []interface{}{}
		Enum(ctx,idx)
		for Next(ctx,idx) {
			a = append(a,ToValue(ctx,-1))
			Pop(ctx,2)
		}
		Pop(ctx,1)
		return a
	} else if IsObject(ctx,idx) {
		
		v := ToGoObject(ctx,idx);
		
		if v != nil {
			return v
		}

		m := map[string]interface{}{}

		Enum(ctx,idx)
		for Next(ctx,idx) {
			key := ToString(ctx,-2);
			m[key] = ToValue(ctx,-1)
			Pop(ctx,2)
		}
		Pop(ctx,1)

		return m
	}

	return nil
}

func  Compile(ctx Context,name string,code string)  {
	PushString(ctx,name)
	v := C.CString(code)
	C.duk_compile_raw(ctx.ctx,v,0,1 | C.DUK_COMPILE_NOSOURCE | C.DUK_COMPILE_STRLEN)
	C.free(unsafe.Pointer(v))
}

func  GetTop(ctx Context) int {
	return int(C.duk_get_top(ctx.ctx))
}

func GetHeapptr(ctx Context,idx int) unsafe.Pointer {
	return C.duk_get_heapptr(ctx.ctx,C.duk_idx_t(idx));
}

func PushHeapptr(ctx Context, heapptr unsafe.Pointer) {
	C.duk_push_heapptr(ctx.ctx,heapptr);
}

func PushThis(ctx Context) {
	C.duk_push_this(ctx.ctx);
}

//export go_ObjectDealloc
func go_ObjectDealloc(ctx * C.duk_context) C.duk_ret_t {

	scope := getScope(ctx)

	if scope != nil {

		key := C.CString("__id");

		C.duk_push_string(ctx,key)
		C.duk_get_prop(ctx,-2)

		if C.duk_is_string(ctx,-1) != C.duk_bool_t(0) {
			id := C.duk_to_string(ctx,-1)
			if id != nil {
				scope.removeObject(C.GoString(id))
				log.Println("[DEALLOC]",id)
			}
		}

		C.duk_pop(ctx);

		C.free(unsafe.Pointer(key))
		
	}
	
	return C.duk_ret_t(0)
}


//export go_FunctionCall
func go_FunctionCall(ctx * C.duk_context) C.duk_ret_t {

	scope := getScope(ctx)

	if scope != nil {

		var fn func() int = nil

		C.duk_push_current_function(ctx);

		key := C.CString("__id");

		C.duk_push_string(ctx,key)

		C.duk_get_prop(ctx,-2)

		if C.duk_is_string(ctx,-1) != C.duk_bool_t(0) {
			id := C.duk_to_string(ctx,-1)
			if id != nil {
				fn = scope.objects[C.GoString(id)].(func() int)
			}
		} 

		C.duk_pop_n(ctx,2);

		C.free(unsafe.Pointer(key))
		
		if fn != nil {
			n := fn()
			return C.duk_ret_t(n)
		}
	}

	return C.duk_ret_t(0)
}
