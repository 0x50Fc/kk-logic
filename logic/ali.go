package logic

import (
	"log"
	"strings"

	"github.com/hailongz/kk-lib/dynamic"
	"gopkg.in/yaml.v2"
)

type Ali struct {
	vpc string
}

func NewAli(vpc string) *Ali {

	v := Ali{}
	v.vpc = vpc

	return &v
}

func (S *Ali) getType(stype string) string {
	switch stype {
	case "int", "long", "integer":
		return "integer"
	case "float", "double", "number":
		return "number"
	case "bool", "boolean":
		return "string"
	case "file":
		return "file"
	}
	return "string"
}

func (S *Ali) getFormat(stype string) string {
	switch stype {
	case "int", "integer":
		return "int32"
	case "long":
		return "int64"
	case "float", "double", "number":
		return "number"
	case "bool", "boolean":
		return "Boolean"
	case "file":
		return "file"
	}
	return "string"
}

func (S *Ali) Object(store IStore) interface{} {

	v := map[string]interface{}{
		"swagger":  "2.0",
		"basePath": "/",
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

		path = path[0 : len(path)-4]

		name := "/" + path + "json"

		input := dynamic.Get(app.Object(), "input")

		contentType := dynamic.StringValue(dynamic.Get(app.Object(), "contentType"), "application/json")
		method := dynamic.StringValue(dynamic.Get(input, "method"), "POST")

		in := "query"

		if method == "POST" {
			in = "formData"
		}

		handling := "MAPPING"

		parameters := []interface{}{}

		dynamic.Each(dynamic.Get(input, "fields"), func(key interface{}, field interface{}) bool {

			stype := dynamic.StringValue(dynamic.Get(field, "type"), "string")

			parameters = append(parameters, map[string]interface{}{
				"name":        dynamic.StringValue(dynamic.Get(field, "name"), ""),
				"description": dynamic.StringValue(dynamic.Get(field, "title"), ""),
				"in":          in,
				"type":        S.getType(stype),
				"format":      S.getFormat(stype),
				"pattern":     dynamic.StringValue(dynamic.Get(field, "pattern"), ""),
				"required":    dynamic.BooleanValue(dynamic.Get(field, "required"), false),
			})

			if stype == "file" {
				handling = "PASSTHROUGH"
			}

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

		if handling == "PASSTHROUGH" {
			parameters = parameters[0:0]
		}

		id := strings.Replace(strings.Replace(strings.Replace(path, "/", "_", -1), ".", "_", -1), "-", "_", -1) + "json"

		object := map[string]interface{}{
			"x-aliyun-apigateway-paramater-handling": handling,
			"x-aliyun-apigateway-auth-type":          "ANONYMOUS",
			"x-aliyun-apigateway-backend": map[string]interface{}{
				"type":          "HTTP-VPC",
				"vpcAccessName": S.vpc,
				"path":          name,
				"method":        strings.ToLower(method),
				"timeout":       10000,
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
