package logic

import (
	"log"
	"net/url"
	"strings"

	"github.com/hailongz/kk-lib/dynamic"
	"gopkg.in/yaml.v2"
)

func init() {
	SetGlobalIgnoreKey("contentType")
}

type Swagger struct {
	scheme   string
	host     string
	basePath string
}

func NewSwagger(baseURL string) *Swagger {

	v := Swagger{}
	u, _ := url.Parse(baseURL)
	v.scheme = u.Scheme
	v.host = u.Host
	v.basePath = u.Path

	if v.basePath != "/" {
		if strings.HasSuffix(v.basePath, "/") {
			v.basePath = v.basePath[0 : len(v.basePath)-1]
		}
	}

	return &v
}

func (S *Swagger) getType(stype string) string {
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

func (S *Swagger) Object(store IStore) interface{} {

	v := map[string]interface{}{
		"swagger":  "2.0",
		"host":     S.host,
		"basePath": S.basePath,
		"info": map[string]interface{}{
			"title":   "kk-logic",
			"version": "1.0",
		},
		"schemes": []interface{}{
			S.scheme,
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

		in := "query"

		if method == "POST" {
			in = "formData"
		}

		parameters := []interface{}{}

		dynamic.Each(dynamic.Get(input, "fields"), func(key interface{}, field interface{}) bool {

			parameters = append(parameters, map[string]interface{}{
				"name":        dynamic.StringValue(dynamic.Get(field, "name"), ""),
				"description": dynamic.StringValue(dynamic.Get(field, "title"), ""),
				"in":          in,
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
			consumes = append(consumes, "multipart/form-data")
		}

		object := map[string]interface{}{
			"summary":    dynamic.StringValue(dynamic.Get(app.Object(), "title"), ""),
			"produces":   produces,
			"consumes":   consumes,
			"parameters": parameters,
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

func (S *Swagger) Marshal(store IStore) ([]byte, error) {
	return yaml.Marshal(S.Object(store))
}
