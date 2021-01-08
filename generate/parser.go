package generate

import (
	"bytes"
	"fmt"
	"github.com/tal-tech/go-zero/tools/goctl/api/spec"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"unsafe"
)

var (
	strColon = []byte(":")
)

func applyGenerate(p Plugin) (*swaggerObject, error) {

	s := swaggerObject{
		Swagger:           "2.0",
		Schemes:           []string{"http", "https"},
		Consumes:          []string{"application/json"},
		Produces:          []string{"application/json"},
		Paths:             make(swaggerPathsObject),
		Definitions:       make(swaggerDefinitionsObject),
		StreamDefinitions: make(swaggerDefinitionsObject),
		Info: swaggerInfoObject{
			Title:   p.Api.Info.Title,
			Version: p.Api.Info.Version,
		},
	}

	s.SecurityDefinitions = swaggerSecurityDefinitionsObject{}
	newSecDefValue := swaggerSecuritySchemeObject{}
	newSecDefValue.Name = "Authorization"
	newSecDefValue.Description = "Enter JWT Bearer token **_only_**"
	newSecDefValue.Type = "apiKey"
	newSecDefValue.In = "header"
	s.SecurityDefinitions["apiKey"] = newSecDefValue

	requestResponseRefs := refMap{}
	renderServiceRoutes(p.Api.Service, p.Api.Service.Groups, s.Paths, requestResponseRefs)
	m := messageMap{}
	renderReplyAsDefinition(s.Definitions, m, p.Api.Types, requestResponseRefs)
	return &s, nil
}

func renderServiceRoutes(service spec.Service, groups []spec.Group, paths swaggerPathsObject, requestResponseRefs refMap) {

	for _, group := range groups {

		for _, route := range group.Routes {
			path := route.Path
			parameters := swaggerParametersObject{}
			if countParams(path) > 0 {
				p := strings.Split(path, "/")
				for i := range p {
					part := p[i]
					if strings.Contains(part, ":") {
						key := strings.TrimPrefix(p[i], ":")
						path = strings.Replace(path, fmt.Sprintf(":%s", key), fmt.Sprintf("{%s}", key), 1)
						parameters = append(parameters, swaggerParameterObject{
							Name:     key,
							In:       "path",
							Required: true,
							Type:     "string",
						})
					}
				}
			}

			reqRef := fmt.Sprintf("#/definitions/%s", route.RequestType.Name)
			if len(route.RequestType.Name) > 0 {
				var schema = swaggerSchemaObject{
					schemaCore: schemaCore{
						Ref: reqRef,
					},
				}
				parameters = append(parameters, swaggerParameterObject{
					Name:     "body",
					In:       "body",
					Required: true,
					Schema:   &schema,
				})
			}

			pathItemObject, ok := paths[path]
			if !ok {
				pathItemObject = swaggerPathItemObject{}
			}

			desc := "A successful response."
			respRef := fmt.Sprintf("#/definitions/%s", route.ResponseType.Name)
			if len(route.ResponseType.Name) < 1 {
				respRef = ""
			}

			tags := service.Name
			if group.Annotations != nil && len(group.Annotations) > 0 {
				if groupName, ok := group.Annotations[0].Properties["group"]; ok {
					tags = groupName
				}
			}

			operationObject := &swaggerOperationObject{
				Tags:       []string{tags},
				Parameters: parameters,
				Responses: swaggerResponsesObject{
					"200": swaggerResponseObject{
						Description: desc,
						Schema: swaggerSchemaObject{
							schemaCore: schemaCore{
								Ref: respRef,
							},
						},
					},
				},
			}

			// set OperationID
			for _, annotation := range route.Annotations {
				if annotation.Name == "handler" {
					operationObject.OperationID = annotation.Value
				}
			}

			for _, param := range operationObject.Parameters {
				if param.Schema != nil && param.Schema.Ref != "" {
					requestResponseRefs[param.Schema.Ref] = struct{}{}
				}
			}

			if len(route.Annotations) > 0 {
				operationObject.Summary, _ = strconv.Unquote(route.Annotations[0].Properties["summary"])
				operationObject.Description, _ = strconv.Unquote(route.Annotations[0].Properties["description"])
			}

			switch strings.ToUpper(route.Method) {
			case http.MethodGet:
				pathItemObject.Get = operationObject
			case http.MethodPost:
				pathItemObject.Post = operationObject
			case http.MethodDelete:
				pathItemObject.Delete = operationObject
			case http.MethodPut:
				pathItemObject.Put = operationObject
			}

			paths[path] = pathItemObject
		}
	}

}

func renderReplyAsDefinition(d swaggerDefinitionsObject, m messageMap, p []spec.Type, refs refMap) {
	for _, i2 := range p {
		schema := swaggerSchemaObject{
			schemaCore: schemaCore{
				Type: "object",
			},
		}
		schema.Title = i2.Name

		for _, member := range i2.Members {
			kv := keyVal{Value: schemaOfField(member)}
			kv.Key = member.Name
			if tag, err := member.GetPropertyName(); err == nil {
				kv.Key = tag
			}
			if schema.Properties == nil {
				schema.Properties = &swaggerSchemaObjectProperties{}
			}
			*schema.Properties = append(*schema.Properties, kv)
		}
		d[i2.Name] = schema
	}

}

func schemaOfField(member spec.Member) swaggerSchemaObject {
	ret := swaggerSchemaObject{}

	var core schemaCore

	kind := swaggerMapTypes[member.Type]
	var props *swaggerSchemaObjectProperties

	switch ft := kind; ft {
	case reflect.Invalid: //[]Struct 也有可能是 Struct
		// []Struct
		refTypeName := strings.Replace(member.Type, "[", "", 1)
		refTypeName = strings.Replace(refTypeName, "]", "", 1)
		core = schemaCore{
			Ref: "#/definitions/" + refTypeName,
		}
	default:
		ftype, format, ok := primitiveSchema(ft, member.Type)
		if ok {
			core = schemaCore{Type: ftype, Format: format}
		} else {
			core = schemaCore{Type: ft.String(), Format: "UNKNOWN"}
		}
	}

	switch ft := kind; ft {
	case reflect.Slice:
		ret = swaggerSchemaObject{
			schemaCore: schemaCore{
				Type:  "array",
				Items: (*swaggerItemsObject)(&core),
			},
		}
	case reflect.Invalid:
		// 判断是否数组
		if strings.HasPrefix(member.Type, "[]") {
			ret = swaggerSchemaObject{
				schemaCore: schemaCore{
					Type:  "array",
					Items: (*swaggerItemsObject)(&core),
				},
			}
		} else {

			ret = swaggerSchemaObject{
				schemaCore: core,
				Properties: props,
			}
		}
	default:
		ret = swaggerSchemaObject{
			schemaCore: core,
			Properties: props,
		}
	}

	return ret
}

// https://swagger.io/specification/ Data Types
func primitiveSchema(kind reflect.Kind, t string) (ftype, format string, ok bool) {
	switch kind {
	case reflect.Int:
		return "integer", "int32", true
	case reflect.Int64:
		return "integer", "int64", true
	case reflect.Bool:
		return "boolean", "boolean", true
	case reflect.String:
		return "string", "", true
	case reflect.Slice:
		return strings.Replace(t, "[]", "", -1), "", true
	default:
		return "", "", false
	}
}

// StringToBytes converts string to byte slice without a memory allocation.
func stringToBytes(s string) (b []byte) {
	return *(*[]byte)(unsafe.Pointer(
		&struct {
			string
			Cap int
		}{s, len(s)},
	))
}

func countParams(path string) uint16 {
	var n uint16
	s := stringToBytes(path)
	n += uint16(bytes.Count(s, strColon))
	return n
}
