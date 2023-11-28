package v2

import (
	"github.com/aishuchen/goctl-swagger/render/types"
	"github.com/zeromicro/go-zero/tools/goctl/api/spec"
	"net/http"
)

func renderParameters(obj spec.DefineStruct, method string) []*Parameter {
	var (
		parameters []*Parameter
		body       *spec.DefineStruct // request body
	)

	for _, field := range obj.Members {
		if field.Name == "" { // 匿名字段一定是结构体
			parameters = append(parameters, renderParameters(field.Type.(spec.DefineStruct), method)...)
			continue
		}
		tags := field.Tags()
		if len(tags) == 0 {
			continue
		}
		tag := lookupGozeroTag(tags)
		if tag == nil {
			continue
		}
		var param Parameter
		if method == http.MethodGet { // Don't support body parameters in GET method
			switch tag.Key {
			case types.PathTagKey:
				param.In = "path"
			case types.FormTagKey:
				param.In = "query"
			default:
				continue // request parameter only support path and query parameters, todo: support header and form parameters
			}
		} else {
			switch tag.Key {
			case types.PathTagKey:
				param.In = "path"
			case types.FormTagKey:
				param.In = "query"
			case types.JsonTagKey:
				body = &obj
				continue
			default:
				continue // just support application/json yet. TODO support multipart/form-data
			}
		}
		param.Name = tag.Name
		param.Description = parseComment(field.Comment)
		param.Required = !isOptionalTag(tag)
		param.Type, param.Format = rawTypeFormat(field.Type.Name())

		parameters = append(parameters, &param)
	}
	if body != nil {
		parameters = append(parameters, renderRequestBody(*body))
	}
	return parameters
}

func renderRequestBody(body spec.DefineStruct) *Parameter {
	param := Parameter{
		Name: "body",
		In:   "body",
		Schema: &Schema{
			Type: "object",
		},
	}

	name, schema := renderSchema(body)
	ref := registerModel(name, schema)
	param.Schema.Ref = ref
	return &param
}
