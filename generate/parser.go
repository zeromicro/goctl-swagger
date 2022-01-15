package generate

import (
	"bytes"
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"unsafe"

	"github.com/tal-tech/go-zero/tools/goctl/api/spec"
	plugin2 "github.com/tal-tech/go-zero/tools/goctl/plugin"
)

var strColon = []byte(":")

const (
	defaultOption   = "default"
	stringOption    = "string"
	optionalOption  = "optional"
	optionsOption   = "options"
	rangeOption     = "range"
	optionSeparator = "|"
	equalToken      = "="
)

func applyGenerate(p *plugin2.Plugin, host string, basePath string) (*swaggerObject, error) {
	title, _ := strconv.Unquote(p.Api.Info.Properties["title"])
	version, _ := strconv.Unquote(p.Api.Info.Properties["version"])
	desc, _ := strconv.Unquote(p.Api.Info.Properties["desc"])

	s := swaggerObject{
		Swagger:           "2.0",
		Schemes:           []string{"http", "https"},
		Consumes:          []string{"application/json"},
		Produces:          []string{"application/json"},
		Paths:             make(swaggerPathsObject),
		Definitions:       make(swaggerDefinitionsObject),
		StreamDefinitions: make(swaggerDefinitionsObject),
		Info: swaggerInfoObject{
			Title:       title,
			Version:     version,
			Description: desc,
		},
	}
	if len(host) > 0 {
		s.Host = host
	}
	if len(basePath) > 0 {
		s.BasePath = basePath
	}

	s.SecurityDefinitions = swaggerSecurityDefinitionsObject{}
	newSecDefValue := swaggerSecuritySchemeObject{}
	newSecDefValue.Name = "Authorization"
	newSecDefValue.Description = "Enter JWT Bearer token **_only_**"
	newSecDefValue.Type = "apiKey"
	newSecDefValue.In = "header"
	s.SecurityDefinitions["apiKey"] = newSecDefValue

	s.Security = append(s.Security, swaggerSecurityRequirementObject{"apiKey": []string{}})

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

						spo := swaggerParameterObject{
							Name:     key,
							In:       "path",
							Required: true,
							Type:     "string",
						}

						// extend the comment functionality
						// to allow query string parameters definitions
						// EXAMPLE:
						// @doc(
						// 	summary: "Get Cart"
						// 	description: "returns a shopping cart if one exists"
						// 	customerId: "customer id"
						// )
						//
						// the format for a parameter is
						// paramName: "the param description"
						//

						prop := route.AtDoc.Properties[key]
						if prop != "" {
							// remove quotes
							spo.Description = strings.Trim(prop, "\"")
						}

						parameters = append(parameters, spo)
					}
				}
			}
			if defineStruct, ok := route.RequestType.(spec.DefineStruct); ok {
				if strings.ToUpper(route.Method) == http.MethodGet {
					for _, member := range defineStruct.Members {
						if strings.Contains(member.Tag, "path") {
							continue
						}
						tempKind := swaggerMapTypes[strings.Replace(member.Type.Name(), "[]", "", -1)]

						ftype, format, ok := primitiveSchema(tempKind, member.Type.Name())
						if !ok {
							ftype = tempKind.String()
							format = "UNKNOWN"
						}
						sp := swaggerParameterObject{In: "query", Type: ftype, Format: format}

						for _, tag := range member.Tags() {
							sp.Name = tag.Name
							if len(tag.Options) == 0 {
								sp.Required = true
								continue
							}
							for _, option := range tag.Options {
								if strings.HasPrefix(option, defaultOption) {
									segs := strings.Split(option, equalToken)
									if len(segs) == 2 {
										sp.Default = segs[1]
									}
								} else if !strings.HasPrefix(option, optionalOption) {
									sp.Required = true
								}
							}
						}

						if len(member.Comment) > 0 {
							sp.Description = strings.TrimLeft(member.Comment, "//")
						}

						parameters = append(parameters, sp)
					}
				} else {

					reqRef := fmt.Sprintf("#/definitions/%s", route.RequestType.Name())

					if len(route.RequestType.Name()) > 0 {
						schema := swaggerSchemaObject{
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
				}
			}

			pathItemObject, ok := paths[path]
			if !ok {
				pathItemObject = swaggerPathItemObject{}
			}

			desc := "A successful response."
			respRef := ""
			if route.ResponseType != nil && len(route.ResponseType.Name()) > 0 {
				respRef = fmt.Sprintf("#/definitions/%s", route.ResponseType.Name())
			}
			tags := service.Name
			if value := group.GetAnnotation("group"); len(value) > 0 {
				tags = value
			}

			if value := group.GetAnnotation("swtags"); len(value) > 0 {
				tags = value
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
			operationObject.OperationID = route.Handler

			for _, param := range operationObject.Parameters {
				if param.Schema != nil && param.Schema.Ref != "" {
					requestResponseRefs[param.Schema.Ref] = struct{}{}
				}
			}
			operationObject.Summary = strings.ReplaceAll(route.JoinedDoc(), "\"", "")

			if len(route.AtDoc.Properties) > 0 {
				operationObject.Description, _ = strconv.Unquote(route.AtDoc.Properties["description"])
			}

			operationObject.Description = strings.ReplaceAll(operationObject.Description, "\"", "")

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
		defineStruct, _ := i2.(spec.DefineStruct)

		schema.Title = defineStruct.Name()

		for _, member := range defineStruct.Members {
			kv := keyVal{Value: schemaOfField(member)}
			kv.Key = member.Name
			if tag, err := member.GetPropertyName(); err == nil {
				kv.Key = tag
			}
			if schema.Properties == nil {
				schema.Properties = &swaggerSchemaObjectProperties{}
			}
			*schema.Properties = append(*schema.Properties, kv)

			for _, tag := range member.Tags() {
				if len(tag.Options) == 0 {
					if !contains(schema.Required, tag.Name) && tag.Name != "required" {
						schema.Required = append(schema.Required, tag.Name)
					}
					continue
				}
				for _, option := range tag.Options {
					switch {
					case !strings.HasPrefix(option, optionalOption):
						if !contains(schema.Required, tag.Name) {
							schema.Required = append(schema.Required, tag.Name)
						}
						// case strings.HasPrefix(option, defaultOption):
						// case strings.HasPrefix(option, optionsOption):
					}
				}
			}
		}

		d[i2.Name()] = schema
	}
}

func schemaOfField(member spec.Member) swaggerSchemaObject {
	ret := swaggerSchemaObject{}

	var core schemaCore

	kind := swaggerMapTypes[member.Type.Name()]
	var props *swaggerSchemaObjectProperties

	comment := member.GetComment()
	comment = strings.Replace(comment, "//", "", -1)

	switch ft := kind; ft {
	case reflect.Invalid: //[]Struct 也有可能是 Struct
		// []Struct
		// map[ArrayType:map[Star:map[StringExpr:UserSearchReq] StringExpr:*UserSearchReq] StringExpr:[]*UserSearchReq]
		refTypeName := strings.Replace(member.Type.Name(), "[", "", 1)
		refTypeName = strings.Replace(refTypeName, "]", "", 1)
		refTypeName = strings.Replace(refTypeName, "*", "", 1)
		refTypeName = strings.Replace(refTypeName, "{", "", 1)
		refTypeName = strings.Replace(refTypeName, "}", "", 1)
		// interface

		if refTypeName == "interface" {
			core = schemaCore{Type: "object"}
		} else if refTypeName == "mapstringstring" {
			core = schemaCore{Type: "object"}
		} else {
			core = schemaCore{
				Ref: "#/definitions/" + refTypeName,
			}
		}
	case reflect.Slice:
		tempKind := swaggerMapTypes[strings.Replace(member.Type.Name(), "[]", "", -1)]
		ftype, format, ok := primitiveSchema(tempKind, member.Type.Name())

		if ok {
			core = schemaCore{Type: ftype, Format: format}
		} else {
			core = schemaCore{Type: ft.String(), Format: "UNKNOWN"}
		}
	default:
		ftype, format, ok := primitiveSchema(ft, member.Type.Name())
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
		if strings.HasPrefix(member.Type.Name(), "[]") {
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
		if strings.HasPrefix(member.Type.Name(), "map") {
			fmt.Println("暂不支持map类型")
		}
	default:
		ret = swaggerSchemaObject{
			schemaCore: core,
			Properties: props,
		}
	}
	ret.Description = comment

	for _, tag := range member.Tags() {
		if len(tag.Options) == 0 {
			continue
		}
		for _, option := range tag.Options {
			switch {
			case strings.HasPrefix(option, defaultOption):
				segs := strings.Split(option, equalToken)
				if len(segs) == 2 {
					ret.Default = segs[1]
				}
			case strings.HasPrefix(option, optionsOption):

			}
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
	case reflect.Float32:
		return "number", "float", true
	case reflect.Float64:
		return "number", "double", true
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

func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}
