package logic

import (
	"log"
	"time"

	"github.com/go-redis/redis"
	"github.com/hailongz/kk-lib/duktape"
)

func init() {

	AddProtocol(func(app IApp, ctx duktape.Context) {

		duktape.PushGlobalObject(ctx)

		duktape.PushString(ctx, "Redis")
		duktape.PushObject(ctx)

		duktape.PushString(ctx, "get")
		duktape.PushGoFunction(ctx, func() int {

			top := duktape.GetTop(ctx)

			if top > 0 && duktape.IsString(ctx, -top) {

				key := duktape.ToString(ctx, -top)

				duktape.PushThis(ctx)

				v := duktape.ToGoObject(ctx, -1)

				duktape.Pop(ctx, 1)

				if v != nil {

					cli := v.(*redis.Client)

					if cli != nil {

						v, err := cli.Get(key).Result()

						if err != nil {
							duktape.PushBoolean(ctx, false)
							log.Println("[REDIS] [ERROR]", err.Error())
							return 1
						} else {
							duktape.PushString(ctx, v)
							return 1
						}
					}
				}

			}

			return 0
		})
		duktape.PutProp(ctx, -3)

		duktape.PushString(ctx, "set")
		duktape.PushGoFunction(ctx, func() int {

			top := duktape.GetTop(ctx)

			if top > 1 && duktape.IsString(ctx, -top) && duktape.IsString(ctx, -top+1) {

				key := duktape.ToString(ctx, -top)
				value := duktape.ToString(ctx, -top+1)
				exp := time.Duration(0)

				if top > 2 && duktape.IsNumber(ctx, -top+2) {
					exp = time.Millisecond * time.Duration(duktape.ToInt(ctx, -top+2))
				}

				duktape.PushThis(ctx)

				v := duktape.ToGoObject(ctx, -1)

				duktape.Pop(ctx, 1)

				if v != nil {

					cli := v.(*redis.Client)

					if cli != nil {

						_, err := cli.Set(key, value, exp).Result()

						if err != nil {
							duktape.PushBoolean(ctx, false)
							log.Println("[REDIS] [ERROR]", err.Error())
							return 1
						} else {
							duktape.PushBoolean(ctx, true)
							return 1
						}
					}
				}

			}

			return 0
		})
		duktape.PutProp(ctx, -3)

		duktape.PushString(ctx, "del")
		duktape.PushGoFunction(ctx, func() int {
			top := duktape.GetTop(ctx)

			if top > 0 && duktape.IsString(ctx, -top) {

				key := duktape.ToString(ctx, -top)

				duktape.PushThis(ctx)

				v := duktape.ToGoObject(ctx, -1)

				duktape.Pop(ctx, 1)

				if v != nil {

					cli := v.(*redis.Client)

					if cli != nil {

						_, err := cli.Del(key).Result()

						if err != nil {
							duktape.PushBoolean(ctx, false)
							log.Println("[REDIS] [ERROR]", err.Error())
							return 1
						} else {
							duktape.PushBoolean(ctx, true)
							return 1
						}
					}
				}

			}

			return 0
		})
		duktape.PutProp(ctx, -3)

		duktape.PutProp(ctx, -3)

		duktape.PushString(ctx, "redis")
		duktape.PushObject(ctx)

		duktape.PushString(ctx, "open")
		duktape.PushGoFunction(ctx, func() int {

			addr := ""
			pwd := ""
			db := 0

			top := duktape.GetTop(ctx)

			if top > 0 && duktape.IsString(ctx, -top) {
				addr = duktape.ToString(ctx, -top)
			}

			if top > 1 && duktape.IsString(ctx, -top+1) {
				pwd = duktape.ToString(ctx, -top+1)
			}

			if top > 2 && duktape.IsNumber(ctx, -top+2) {
				db = duktape.ToInt(ctx, -top)
			}

			client := redis.NewClient(&redis.Options{
				Addr:     addr,
				Password: pwd,     // no password set
				DB:       int(db), // use default DB
			})

			_, err := client.Ping().Result()

			if err == nil {
				duktape.PushGoObject(ctx, client)
				duktape.PushGlobalObject(ctx)
				duktape.PushString(ctx, "Redis")
				duktape.GetProp(ctx, -2)
				duktape.SetPrototype(ctx, -3)
				duktape.Pop(ctx, 1)
				return 1
			} else {
				duktape.PushBoolean(ctx, false)
				log.Println("[REDIS] [ERROR]", err.Error())
				return 1
			}

		})
		duktape.PutProp(ctx, -3)

		duktape.PutProp(ctx, -3)

		duktape.Pop(ctx, 1)
	})

}
