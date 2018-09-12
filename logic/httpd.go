package logic

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/hailongz/kk-lib/dynamic"
)

func HandlerFunc(store IStore, session ISession, maxMemory int64) func(resp http.ResponseWriter, req *http.Request) {

	fs := http.FileServer(http.Dir(store.Dir()))

	return func(resp http.ResponseWriter, req *http.Request) {

		if strings.HasSuffix(req.URL.Path, ".json") {

			path := req.URL.Path[0:len(req.URL.Path)-4] + "yaml"

			app, err := store.Get(path)

			if err == nil {

				ctx := NewContext()

				defer ctx.Recycle()

				var inputData interface{} = nil

				{

					if req.Method == "POST" {

						ctype := req.Header.Get("Content-Type")

						if strings.Contains(ctype, "text/json") || strings.Contains(ctype, "application/json") {
							b, err := ioutil.ReadAll(req.Body)
							if err != nil {
								json.Unmarshal(b, &inputData)
							}
							req.Body.Close()
						} else if strings.Contains(ctype, "text/xml") || strings.Contains(ctype, "text/plain") {
							b, err := ioutil.ReadAll(req.Body)
							if err != nil {
								ctx.Set(ContentKeys, string(b))
							}
							req.Body.Close()
						} else if strings.Contains(ctype, "multipart/form-data") {
							inputData = map[string]interface{}{}
							req.ParseMultipartForm(maxMemory)
							if req.MultipartForm != nil {
								for key, values := range req.MultipartForm.Value {
									dynamic.Set(inputData, key, values[0])
								}
								for key, values := range req.MultipartForm.File {
									dynamic.Set(inputData, key, values[0])
								}
							}
						} else {

							inputData = map[string]interface{}{}

							req.ParseForm()

							for key, values := range req.Form {
								dynamic.Set(inputData, key, values[0])
							}

						}

					} else {
						inputData = map[string]interface{}{}
						req.ParseForm()

						for key, values := range req.Form {
							dynamic.Set(inputData, key, values[0])
						}
					}

					ctx.Set(InputKeys, inputData)

				}

				{
					var ip = req.Header.Get("X-CLIENT-IP")

					if ip == "" {
						ip = req.Header.Get("X-Real-IP")
					}

					if ip == "" {
						ip = req.RemoteAddr
					}

					ip = strings.Split(ip, ":")[0]

					ctx.Set(ClientIpKeys, ip)
				}

				{
					ctx.Set([]string{"referer"}, req.Referer())
					ctx.Set(RefererKeys, req.UserAgent())
					if strings.HasPrefix(req.Proto, "HTTPS/") {
						ctx.Set([]string{"url"}, "https://"+req.Host+req.RequestURI)
						ctx.Set([]string{"scheme"}, "https")
					} else {
						ctx.Set([]string{"url"}, "http://"+req.Host+req.RequestURI)
						ctx.Set([]string{"scheme"}, "http")
					}
				}

				{
					headers := map[string]interface{}{}

					for key, values := range req.Header {
						headers[key] = values[0]
					}

					ctx.Set(HeadersKeys, headers)
				}

				ctx.Set(URLKeys, map[string]interface{}{
					"path":     req.URL.Path,
					"host":     req.URL.Host,
					"hostname": req.URL.Hostname(),
					"port":     req.URL.Port(),
					"fragment": req.URL.Fragment,
					"scheme":   req.URL.Scheme,
				})

				ctx.Set(MethodKeys, req.Method)
				ctx.Set(SessionIdKeys, session.Http(resp, req))
				ctx.Set(OutputKeys, map[string]interface{}{})

				err = app.Exec(ctx, "in")

				if err != nil {

					{
						r, ok := err.(*Redirect)

						if ok {
							resp.Header().Set("Location", r.URL)
							resp.WriteHeader(302)
							return
						}
					}

					{
						r, ok := err.(*Content)

						if ok {

							for key, vs := range r.Header {
								resp.Header()[key] = vs
							}

							resp.Header().Set("Content-Type", r.ContentType)

							resp.Write(r.Content)

							return
						}
					}

					b, _ := json.Marshal(GetErrorObject(err))
					resp.Header().Set("Content-Type", "application/json; charset=utf-8")
					resp.Write(b)
					return
				}

				{
					v := ctx.Get(ViewKeys)
					if v != nil {
						view, ok := v.(*View)
						if ok {
							resp.Header().Set("Content-Type", view.ContentType)
							if view.Headers != nil {
								for key, value := range view.Headers {
									resp.Header()[key] = []string{value}
								}
							}
							resp.WriteHeader(200)
							resp.Write(view.Content)
							return
						}
					}
				}

				{
					v := ctx.Get(OutputKeys)
					b, _ := json.Marshal(v)
					resp.Header().Set("Content-Type", "application/json; charset=utf-8")
					resp.Write(b)
					return
				}

			}

			if err != nil {
				resp.WriteHeader(404)
				resp.Write([]byte(err.Error()))
			}

		} else if strings.HasSuffix(req.URL.Path, ".yaml") || strings.HasSuffix(req.URL.Path, ".yml") {
			resp.WriteHeader(404)
		} else if req.URL.Path == "/*.swagger" {

			s := NewSwagger(fmt.Sprintf("http://%s/", req.Host))

			b, err := s.Marshal(store)

			if err != nil {
				resp.WriteHeader(404)
				resp.Write([]byte(err.Error()))
			} else {
				resp.Header().Add("Content-Type", "text/yaml; charset=utf-8")
				resp.Write(b)
			}
		} else if req.URL.Path == "/*.ali" {

			s := NewAli(req.URL.Query().Get("vpc"))

			b, err := s.Marshal(store)

			if err != nil {
				resp.WriteHeader(404)
				resp.Write([]byte(err.Error()))
			} else {
				resp.Header().Add("Content-Type", "text/yaml; charset=utf-8")
				resp.Write(b)
			}
		} else {
			fs.ServeHTTP(resp, req)
		}
	}

}
