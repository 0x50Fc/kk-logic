package logic

import (
	"log"
	"net/url"
	"strings"

	"github.com/hailongz/kk-lib/dynamic"
	"gopkg.in/yaml.v2"
)

type Ali struct {
	scheme   string
	host     string
	basePath string
}

func NewAli(baseURL string) *Ali {

	v := Ali{}
	u, _ := url.Parse(baseURL)
	v.scheme = u.Scheme
	v.host = u.Host
	v.basePath = u.Path

	if strings.HasSuffix(v.basePath, "/") {
		v.basePath = v.basePath[0 : len(v.basePath)-1]
	}

	return &v
}

func (S *Ali) getType(stype string) string {
	switch stype {
	case "int", "long", "integer":
		return "integer"
	case "float", "double", "number":
		return "number"
	case "bool", "boolean":
		return "boolean"
	case "file":
		return "file"
	}
	return "string"
}

func (S *Ali) Object(store IStore) interface{} {

	basePath := S.basePath

	if basePath == "" {
		basePath = "/"
	}

	v := map[string]interface{}{
		"swagger":  "2.0",
		"basePath": basePath,
		"info": map[string]interface{}{
			"title":   "kk-logic",
			"version": "1.0",
		},
		"schemes": []interface{}{
			"https",
		},
	}

	paths := map[string]interface{}{}

	v["paths"] = paths

	store.Walk(func(path string) {

		app, err := store.Get(path)

		if err != nil {
			log.Println("[ERROR]", path, err)
			return
		}

		name := "/" + path[0:len(path)-4] + "json"

		input := dynamic.Get(app.Object(), "input")

		contentType := dynamic.StringValue(dynamic.Get(app.Object(), "contentType"), "application/json")
		method := dynamic.StringValue(dynamic.Get(input, "method"), "POST")

		parameters := []interface{}{}

		dynamic.Each(dynamic.Get(input, "fields"), func(key interface{}, field interface{}) bool {

			parameters = append(parameters, map[string]interface{}{
				"name":        dynamic.StringValue(dynamic.Get(field, "name"), ""),
				"description": dynamic.StringValue(dynamic.Get(field, "title"), ""),
				"in":          "formData",
				"type":        S.getType(dynamic.StringValue(dynamic.Get(field, "type"), "string")),
				"pattern":     dynamic.StringValue(dynamic.Get(field, "pattern"), ""),
				"required":    dynamic.BooleanValue(dynamic.Get(field, "required"), false),
			})

			return true
		})

		produces := []interface{}{"application/json"}

		if contentType != "application/json" {
			produces = append(produces, contentType)
		}

		consumes := []interface{}{"application/x-www-form-urlencoded"}

		if method == "POST" {
			consumes = append(produces, "multipart/form-data")
		}

		id := strings.Replace(strings.Replace(name, "/", "_", -1), ".", "_", -1)

		object := map[string]interface{}{
			"x-aliyun-apigateway-paramater-handling": "MAPPING",
			"x-aliyun-apigateway-auth-type":          "ANONYMOUS",
			"x-aliyun-apigateway-backend": map[string]interface{}{
				"type":    "HTTP",
				"address": S.scheme + "://" + S.host,
				"path":    S.basePath + name,
				"method":  strings.ToLower(method),
				"timeout": 10000,
			},
			"schemes": []interface{}{
				"https",
			},
			"operationId": id,
			"summary":     dynamic.StringValue(dynamic.Get(app.Object(), "title"), ""),
			"produces":    produces,
			"consumes":    consumes,
			"parameters":  parameters,
			"responses": map[string]interface{}{
				"200": map[string]interface{}{
					"description": "OK",
				},
			},
		}

		if method == "POST" {
			paths[name] = map[string]interface{}{"post": object}
		} else {
			paths[name] = map[string]interface{}{"get": object}
		}

	})

	return v
}

func (S *Ali) Marshal(store IStore) ([]byte, error) {
	return yaml.Marshal(S.Object(store))
}
