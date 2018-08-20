package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/hailongz/kk-logic/logic"
)

func main() {

	prefix := "/"
	dir := "."
	port := 8080
	cached := true
	sessionKey := "kk"
	sessionMaxAge := 1800
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

			} else if v == "-d" {

				cached = false

				i += 1

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

			}
			i += 1
		}

	}

	app := logic.NewApp(dir, cached, sessionKey, sessionMaxAge)

	log.Println("PORT: ", port)
	log.Println("CACHED: ", cached)
	log.Println("ROOT: ", dir)
	log.Println("PREFIX: ", prefix)

	http.Handle(prefix, app)

	srv := &http.Server{
		Addr: fmt.Sprintf(":%d", port),
	}

	log.Panic(srv.ListenAndServe())

}
