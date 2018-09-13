package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/hailongz/kk-logic/logic"
	_ "github.com/hailongz/kk-logic/logic/captcha"
	_ "github.com/hailongz/kk-logic/logic/http"
	_ "github.com/hailongz/kk-logic/logic/lib"
	_ "github.com/hailongz/kk-logic/logic/oss"
	_ "github.com/hailongz/kk-logic/logic/redis"
)

func main() {

	prefix := "/"
	dir := "."
	port := 8080
	sessionKey := "kk"
	sessionMaxAge := 1800
	maxMemory := int64(4096000)

	{
		i := 1
		n := len(os.Args)

		for i < n {
			v := os.Args[i]

			if v == "-p" && i+1 < n {

				ii, _ := strconv.ParseInt(os.Args[i+1], 10, 64)

				port = int(ii)

				i += 2
				continue

			} else if v == "-r" && i+1 < n {

				dir = os.Args[i+1]

				i += 2
				continue

			} else if v == "--prefix" && i+1 < n {

				prefix = os.Args[i+1]

				i += 2
				continue

			} else if v == "--sessionKey" && i+1 < n {

				sessionKey = os.Args[i+1]

				i += 2
				continue

			} else if v == "--sessionMaxAge" && i+1 < n {

				sessionMaxAge, _ = strconv.Atoi(os.Args[i+1])

				i += 2
				continue

			} else if v == "--maxMemory" && i+1 < n {

				maxMemory, _ = strconv.ParseInt(os.Args[i+1], 10, 64)

				i += 2
				continue

			}
			i += 1
		}

	}

	var store logic.IStore = logic.NewMemStore(dir, 6*time.Second)

	session := logic.NewSession(sessionKey, sessionMaxAge)

	log.Println("PORT: ", port)
	log.Println("ROOT: ", dir)
	log.Println("PREFIX: ", prefix)
	log.Println("maxMemory: ", maxMemory)

	http.HandleFunc(prefix, logic.HandlerFunc(store, session, maxMemory))

	srv := &http.Server{
		Addr: fmt.Sprintf(":%d", port),
	}

	log.Panic(srv.ListenAndServe())

}
