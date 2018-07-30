package logic

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	URL "net/url"
	"strings"

	"github.com/hailongz/kk-logic/duktape"
)

var ca *x509.CertPool

func init() {

	ca = x509.NewCertPool()
	ca.AppendCertsFromPEM(pemCerts)

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig:   &tls.Config{RootCAs: ca},
			DisableKeepAlives: false,
		},
	}

	AddProtocol(func(app IApp, ctx duktape.Context) {

		duktape.PushGlobalObject(ctx)

		duktape.PushString(ctx, "http")
		duktape.PushObject(ctx)

		duktape.PushString(ctx, "send")
		duktape.PushGoFunction(ctx, func() int {

			top := duktape.GetTop(ctx)

			if top > 0 && duktape.IsObject(ctx, -top) {
				method := "GET"
				url := ""
				ttype := "application/x-www-form-urlencoded"
				dataType := "json"

				duktape.PushString(ctx, "method")
				duktape.GetProp(ctx, -top-1)
				if duktape.IsString(ctx, -1) {
					method = duktape.ToString(ctx, -1)
				}
				duktape.Pop(ctx, 1)

				duktape.PushString(ctx, "url")
				duktape.GetProp(ctx, -top-1)
				if duktape.IsString(ctx, -1) {
					url = duktape.ToString(ctx, -1)
				}
				duktape.Pop(ctx, 1)

				duktape.PushString(ctx, "type")
				duktape.GetProp(ctx, -top-1)
				if duktape.IsString(ctx, -1) {
					ttype = duktape.ToString(ctx, -1)
				}
				duktape.Pop(ctx, 1)

				duktape.PushString(ctx, "dataType")
				duktape.GetProp(ctx, -top-1)
				if duktape.IsString(ctx, -1) {
					dataType = duktape.ToString(ctx, -1)
				}
				duktape.Pop(ctx, 1)

				var req *http.Request = nil
				var err error = nil
				if method == "POST" {

					b := bytes.NewBuffer(nil)
					idx := 0

					if ttype == "application/x-www-form-urlencoded" {
						duktape.PushString(ctx, "data")
						duktape.GetProp(ctx, -top-1)
						if duktape.IsObject(ctx, -1) {
							duktape.Enum(ctx, -1)
							for duktape.Next(ctx, -1) {
								if idx != 0 {
									b.WriteString("&")
								}
								idx = idx + 1
								b.WriteString(duktape.ToString(ctx, -2))
								b.WriteString("=")
								if duktape.IsString(ctx, -1) {
									b.WriteString(URL.QueryEscape(duktape.ToString(ctx, -1)))
								} else if duktape.IsBoolean(ctx, -1) {
									if duktape.ToBoolean(ctx, -1) {
										b.WriteString("true")
									} else {
										b.WriteString("false")
									}
								} else if duktape.IsNumber(ctx, -1) {
									b.WriteString(fmt.Sprintf("%g", duktape.ToNumber(ctx, -1)))
								} else {
									v, _ := json.Marshal(duktape.ToValue(ctx, -1))
									b.WriteString(URL.QueryEscape(string(v)))
								}
								duktape.Pop(ctx, 2)
							}
							duktape.Pop(ctx, 1)
						}
						duktape.Pop(ctx, 1)
					} else if ttype == "application/json" {
						duktape.PushString(ctx, "data")
						duktape.GetProp(ctx, -top-1)
						v := duktape.JsonEncode(ctx, -1)
						b.WriteString(v)
						duktape.Pop(ctx, 1)
					} else {
						duktape.PushString(ctx, "data")
						duktape.GetProp(ctx, -top-1)
						if duktape.IsString(ctx, -1) {
							b.WriteString(duktape.ToString(ctx, -1))
						} else if duktape.IsBytes(ctx, -1) {
							b.Write(duktape.ToBytes(ctx, -1))
						}
						duktape.Pop(ctx, 1)
					}

					req, err = http.NewRequest(method, url, b)

				} else {

					b := bytes.NewBuffer(nil)

					b.WriteString(url)

					idx := 0

					if strings.HasSuffix(url, "?") {
						idx = 0
					} else if strings.Contains(url, "?") {
						idx = 1
					} else {
						idx = -1
					}

					duktape.PushString(ctx, "data")
					duktape.GetProp(ctx, -top-1)
					if duktape.IsObject(ctx, -1) {
						duktape.Enum(ctx, -1)
						for duktape.Next(ctx, -1) {
							if idx == -1 {
								b.WriteString("?")
								idx = 0
							}
							if idx != 0 {
								b.WriteString("&")
							}
							idx = idx + 1
							b.WriteString(duktape.ToString(ctx, -2))
							b.WriteString("=")
							if duktape.IsString(ctx, -1) {
								b.WriteString(URL.QueryEscape(duktape.ToString(ctx, -1)))
							} else if duktape.IsBoolean(ctx, -1) {
								if duktape.ToBoolean(ctx, -1) {
									b.WriteString("true")
								} else {
									b.WriteString("false")
								}
							} else if duktape.IsNumber(ctx, -1) {
								b.WriteString(fmt.Sprintf("%g", duktape.ToNumber(ctx, -1)))
							} else {
								v, _ := json.Marshal(duktape.ToValue(ctx, -1))
								b.WriteString(URL.QueryEscape(string(v)))
							}
							duktape.Pop(ctx, 2)
						}
						duktape.Pop(ctx, 1)
					}
					duktape.Pop(ctx, 1)

					req, err = http.NewRequest(method, b.String(), nil)
				}

				if err != nil {
					duktape.PushThis(ctx)
					duktape.PushString(ctx, "errmsg")
					duktape.PushString(ctx, err.Error())
					duktape.PutProp(ctx, -3)
					duktape.Pop(ctx, 1)
					duktape.PushBoolean(ctx, false)
					return 1
				}

				duktape.PushString(ctx, "headers")
				duktape.GetProp(ctx, -top-1)
				if duktape.IsObject(ctx, -1) {
					duktape.Enum(ctx, -1)
					for duktape.Next(ctx, -1) {
						key := duktape.ToString(ctx, -2)

						if duktape.IsString(ctx, -1) {
							req.Header[key] = []string{duktape.ToString(ctx, -1)}
						} else if duktape.IsBoolean(ctx, -1) {
							if duktape.ToBoolean(ctx, -1) {
								req.Header[key] = []string{"true"}
							} else {
								req.Header[key] = []string{"false"}
							}
						} else if duktape.IsNumber(ctx, -1) {
							req.Header[key] = []string{fmt.Sprintf("%g", duktape.ToNumber(ctx, -1))}
						}
						duktape.Pop(ctx, 2)
					}
					duktape.Pop(ctx, 1)
				}
				duktape.Pop(ctx, 1)

				resp, err := client.Do(req)

				if err != nil {
					duktape.PushThis(ctx)
					duktape.PushString(ctx, "errmsg")
					duktape.PushString(ctx, err.Error())
					duktape.PutProp(ctx, -3)
					duktape.Pop(ctx, 1)
					duktape.PushBoolean(ctx, false)
					return 1
				}

				if resp.StatusCode == 200 {

					b, err := ioutil.ReadAll(resp.Body)

					defer resp.Body.Close()

					if err != nil {
						duktape.PushThis(ctx)
						duktape.PushString(ctx, "errmsg")
						duktape.PushString(ctx, err.Error())
						duktape.PutProp(ctx, -3)
						duktape.Pop(ctx, 1)
						duktape.PushBoolean(ctx, false)
						return 1
					}

					if dataType == "json" {

						duktape.PushString(ctx, string(b))
						duktape.JsonDecode(ctx, -1)

						return 1
					} else if dataType == "text" {
						duktape.PushString(ctx, string(b))
						return 1
					} else {
						duktape.PushBytes(ctx, b)
						return 1
					}

				} else {
					duktape.PushThis(ctx)
					duktape.PushString(ctx, "errmsg")
					duktape.PushString(ctx, fmt.Sprintf("[%d] %s", resp.StatusCode, resp.Status))
					duktape.PutProp(ctx, -3)
					duktape.Pop(ctx, 1)
					duktape.PushBoolean(ctx, false)
					return 1
				}

			}

			return 0
		})

		duktape.PutProp(ctx, -3)

		duktape.PutProp(ctx, -3)

		duktape.Pop(ctx, 1)
	})

}
