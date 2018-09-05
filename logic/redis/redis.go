package http

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/hailongz/kk-lib/dynamic"
	"github.com/hailongz/kk-logic/logic"
	"gopkg.in/redis.v5"
)

type Redis struct {
	client *redis.Client
	prefix string
}

type RedisOpenLogic struct {
	logic.Logic
}

func (L *RedisOpenLogic) Exec(ctx logic.IContext, app logic.IApp) error {

	L.Logic.Exec(ctx, app)

	prefix := dynamic.StringValue(L.Get(ctx, app, "prefix"), "")
	name := dynamic.StringValue(L.Get(ctx, app, "name"), "redis")
	addr := dynamic.StringValue(L.Get(ctx, app, "addr"), "127.0.0.1:6379")
	pwd := dynamic.StringValue(L.Get(ctx, app, "password"), "")
	db := dynamic.IntValue(L.Get(ctx, app, "db"), 0)

	v := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: pwd,     // no password set
		DB:       int(db), // use default DB
	})

	_, err := v.Ping().Result()

	if err != nil {
		if L.Has("error") {
			ctx.Set(logic.ErrorKeys, logic.GetErrorObject(err))
			return L.Done(ctx, app, "error")
		}
		return err
	}

	r := Redis{v, prefix}

	ctx.SetGlobal(name, &r)

	ctx.AddRecycle(func() {
		v.Close()
	})

	return L.Done(ctx, app, "done")
}

type RedisGetLogic struct {
	logic.Logic
}

func (L *RedisGetLogic) Exec(ctx logic.IContext, app logic.IApp) error {

	L.Logic.Exec(ctx, app)

	name := dynamic.StringValue(L.Get(ctx, app, "name"), "redis")
	ttype := dynamic.StringValue(L.Get(ctx, app, "type"), "text")
	key := dynamic.StringValue(L.Get(ctx, app, "key"), "")

	v := ctx.Get([]string{name})

	if v == nil {
		return L.Error(ctx, app, logic.NewError(logic.ERROR_UNKNOWN, fmt.Sprintf("未找到 Redis [%s]", name)))
	}

	r, ok := v.(*Redis)

	if !ok {
		return L.Error(ctx, app, logic.NewError(logic.ERROR_UNKNOWN, fmt.Sprintf("未找到 Redis [%s]", name)))
	}

	vv, err := r.client.Get(key).Result()

	if err != nil {
		return L.Error(ctx, app, err)
	}

	if ttype == "json" {
		v = nil
		err = json.Unmarshal([]byte(vv), &v)
		if err != nil {
			return L.Error(ctx, app, err)
		}
		ctx.Set(logic.ResultKeys, v)
	} else {
		ctx.Set(logic.ResultKeys, vv)
	}

	return L.Done(ctx, app, "done")
}

type RedisSetLogic struct {
	logic.Logic
}

func (L *RedisSetLogic) Exec(ctx logic.IContext, app logic.IApp) error {
	L.Logic.Exec(ctx, app)

	name := dynamic.StringValue(L.Get(ctx, app, "name"), "redis")
	ttype := dynamic.StringValue(L.Get(ctx, app, "type"), "text")
	key := dynamic.StringValue(L.Get(ctx, app, "key"), "")
	expires := dynamic.IntValue(L.Get(ctx, app, "expires"), 0)
	value := L.Get(ctx, app, "value")

	v := ctx.Get([]string{name})

	if v == nil {
		return L.Done(ctx, app, "done")
	}

	r, ok := v.(*Redis)

	if !ok {
		return L.Done(ctx, app, "done")
	}

	var vv string = ""

	if ttype == "json" {
		b, _ := json.Marshal(value)
		vv = string(b)
	} else {
		vv = dynamic.StringValue(value, "")
	}

	_, err := r.client.Set(key, vv, time.Duration(expires)*time.Second).Result()

	if err != nil {
		return L.Error(ctx, app, err)
	}

	return L.Done(ctx, app, "done")
}

type RedisDelLogic struct {
	logic.Logic
}

func (L *RedisDelLogic) Exec(ctx logic.IContext, app logic.IApp) error {

	L.Logic.Exec(ctx, app)

	name := dynamic.StringValue(L.Get(ctx, app, "name"), "redis")
	key := dynamic.StringValue(L.Get(ctx, app, "key"), "")

	v := ctx.Get([]string{name})

	if v == nil {
		return L.Done(ctx, app, "done")
	}

	r, ok := v.(*Redis)

	if !ok {
		return L.Done(ctx, app, "done")
	}

	_, err := r.client.Del(key).Result()

	if err != nil {
		return L.Error(ctx, app, err)
	}

	return L.Done(ctx, app, "done")
}

func init() {
	logic.Openlib("kk.Logic.Redis.Open", func(object interface{}) logic.ILogic {
		v := RedisOpenLogic{}
		v.Init(object)
		return &v
	})
	logic.Openlib("kk.Logic.Redis.Get", func(object interface{}) logic.ILogic {
		v := RedisGetLogic{}
		v.Init(object)
		return &v
	})
	logic.Openlib("kk.Logic.Redis.Set", func(object interface{}) logic.ILogic {
		v := RedisSetLogic{}
		v.Init(object)
		return &v
	})
	logic.Openlib("kk.Logic.Redis.Del", func(object interface{}) logic.ILogic {
		v := RedisDelLogic{}
		v.Init(object)
		return &v
	})
}
