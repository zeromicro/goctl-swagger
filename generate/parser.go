package generate

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"unsafe"

	"github.com/zeromicro/go-zero/tools/goctl/api/spec"
	"github.com/zeromicro/go-zero/tools/goctl/plugin"
)

var strColon = []byte(":")

const (
	validateKey     = "validate"
	defaultOption   = "default"
	stringOption    = "string"
	optionalOption  = "optional"
	omitemptyOption = "omitempty"
	optionsOption   = "options"
	rangeOption     = "range"
	exampleOption   = "example"
	optionSeparator = "|"
	equalToken      = "="
	atRespDoc       = "@respdoc-"
)

func parseRangeOption(option string) (float64, float64, bool) {
	const str = "\\[([+-]?\\d+(\\.\\d+)?):([+-]?\\d+(\\.\\d+)?)\\]"
	result := regexp.MustCompile(str).FindStringSubmatch(option)
	if len(result) != 5 {
		return 0, 0, false
	}

	min, err := strconv.ParseFloat(result[1], 64)
	if err != nil {
		return 0, 0, false
	}

	max, err := strconv.ParseFloat(result[3], 64)
	if err != nil {
		return 0, 0, false
	}

	if max < min {
		return min, min, true
	}
	return min, max, true
}

func applyGenerate(p *plugin.Plugin, host string, basePath string, schemes string) (*swaggerObject, error) {
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

	if len(schemes) > 0 {
		supportedSchemes := []string{"http", "https", "ws", "wss"}
		ss := strings.Split(schemes, ",")
		for i := range ss {
			scheme := ss[i]
			scheme = strings.TrimSpace(scheme)
			if !contains(supportedSchemes, scheme) {
				log.Fatalf("unsupport scheme: [%s], only support [http, https, ws, wss]", scheme)
			}
			ss[i] = scheme
		}
		s.Schemes = ss
	}
	s.SecurityDefinitions = swaggerSecurityDefinitionsObject{}
	newSecDefValue := swaggerSecuritySchemeObject{}
	newSecDefValue.Name = "Authorization"
	newSecDefValue.Description = "Enter JWT Bearer token **_only_**"
	newSecDefValue.Type = "apiKey"
	newSecDefValue.In = "header"
	s.SecurityDefinitions["apiKey"] = newSecDefValue

	// s.Security = append(s.Security, swaggerSecurityRequirementObject{"apiKey": []string{}})

	requestResponseRefs := refMap{}
	renderServiceRoutes(p.Api.Service, p.Api.Service.Groups, s.Paths, requestResponseRefs)
	m := messageMap{}

	renderReplyAsDefinition(s.Definitions, m, p.Api.Types, requestResponseRefs)

	return &s, nil
}

func renderServiceRoutes(service spec.Service, groups []spec.Group, paths swaggerPathsObject, requestResponseRefs refMap) {
	for _, group := range groups {
		for _, route := range group.Routes {
			path := group.GetAnnotation("prefix") + route.Path
			if path[0] != '/' {
				path = "/" + path
			}
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
				for _, member := range defineStruct.Members {
					if hasHeaderParameters(member) {
						parameters = parseHeader(member, parameters)
					}
				}
				if strings.ToUpper(route.Method) == http.MethodGet {
					for _, member := range defineStruct.Members {
						if hasPathParameters(member) || hasHeaderParameters(member) {
							continue
						}
						if embedStruct, isEmbed := member.Type.(spec.DefineStruct); isEmbed {
							for _, m := range embedStruct.Members {
								parameters = append(parameters, renderStruct(m))
							}
							continue
						}
						parameters = append(parameters, renderStruct(member))
					}
				} else {

					reqRef := fmt.Sprintf("#/definitions/%s", route.RequestType.Name())

					if len(route.RequestType.Name()) > 0 {
						schema := swaggerSchemaObject{
							schemaCore: schemaCore{
								Ref: reqRef,
							},
						}

						parameter := swaggerParameterObject{
							Name:     "body",
							In:       "body",
							Required: true,
							Schema:   &schema,
						}

						doc := strings.Join(route.RequestType.Documents(), ",")
						doc = strings.Replace(doc, "//", "", -1)

						if doc != "" {
							parameter.Description = doc
						}

						parameters = append(parameters, parameter)
					}
				}
			}

			pathItemObject, ok := paths[path]
			if !ok {
				pathItemObject = swaggerPathItemObject{}
			}

			desc := "A successful response."
			respSchema := schemaCore{}
			// respRef := swaggerSchemaObject{}
			if route.ResponseType != nil && len(route.ResponseType.Name()) > 0 {
				if strings.HasPrefix(route.ResponseType.Name(), "[]") {

					refTypeName := strings.Replace(route.ResponseType.Name(), "[", "", 1)
					refTypeName = strings.Replace(refTypeName, "]", "", 1)

					respSchema.Type = "array"
					respSchema.Items = &swaggerItemsObject{Ref: fmt.Sprintf("#/definitions/%s", refTypeName)}
				} else {
					respSchema.Ref = fmt.Sprintf("#/definitions/%s", route.ResponseType.Name())
				}
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
							schemaCore: respSchema,
						},
					},
				},
			}

			if defineStruct, ok := route.RequestType.(spec.DefineStruct); ok {
				for _, member := range defineStruct.Members {
					if member.IsFormMember() {
						operationObject.Consumes = []string{"multipart/form-data"}
						break
					}
				}
			}

			for _, v := range route.Doc {
				markerIndex := strings.Index(v, atRespDoc)
				if markerIndex >= 0 {
					l := strings.Index(v, "(")
					r := strings.Index(v, ")")
					code := strings.TrimSpace(v[markerIndex+len(atRespDoc) : l])
					var comment string
					commentIndex := strings.Index(v, "//")
					if commentIndex > 0 {
						comment = strings.TrimSpace(strings.Trim(v[commentIndex+2:], "*/"))
					}
					content := strings.TrimSpace(v[l+1 : r])
					if strings.Index(v, ":") > 0 {
						lines := strings.Split(content, "\n")
						kv := make(map[string]string, len(lines))
						for _, line := range lines {
							sep := strings.Index(line, ":")
							key := strings.TrimSpace(line[:sep])
							value := strings.TrimSpace(line[sep+1:])
							kv[key] = value
						}
						kvByte, err := json.Marshal(kv)
						if err != nil {
							continue
						}
						operationObject.Responses[code] = swaggerResponseObject{
							Description: comment,
							Schema: swaggerSchemaObject{
								schemaCore: schemaCore{
									Example: string(kvByte),
								},
							},
						}
					} else if len(content) > 0 {
						operationObject.Responses[code] = swaggerResponseObject{
							Description: comment,
							Schema: swaggerSchemaObject{
								schemaCore: schemaCore{
									Ref: fmt.Sprintf("#/definitions/%s", content),
								},
							},
						}
					}
				}
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

			if group.Annotation.Properties["jwt"] != "" {
				operationObject.Security = &[]swaggerSecurityRequirementObject{{"apiKey": []string{}}}
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
			case http.MethodPatch:
				pathItemObject.Patch = operationObject
			}

			paths[path] = pathItemObject
		}
	}
}

func renderStruct(member spec.Member) swaggerParameterObject {
	tempKind := swaggerMapTypes[strings.Replace(member.Type.Name(), "[]", "", -1)]

	ftype, format, ok := primitiveSchema(tempKind, member.Type.Name())
	if !ok {
		ftype = tempKind.String()
		format = "UNKNOWN"
	}
	sp := swaggerParameterObject{In: "query", Type: ftype, Format: format, Schema: &swaggerSchemaObject{}}

	for _, tag := range member.Tags() {
		if tag.Key == validateKey {
			continue
		}

		sp.Name = tag.Name
		if len(tag.Options) == 0 {
			sp.Required = true
			continue
		}

		required := true
		for _, option := range tag.Options {
			if strings.HasPrefix(option, optionsOption) {
				segs := strings.SplitN(option, equalToken, 2)
				if len(segs) == 2 {
					sp.Enum = strings.Split(segs[1], optionSeparator)
				}
			}

			if strings.HasPrefix(option, rangeOption) {
				segs := strings.SplitN(option, equalToken, 2)
				if len(segs) == 2 {
					min, max, ok := parseRangeOption(segs[1])
					if ok {
						sp.Schema.Minimum = min
						sp.Schema.Maximum = max
					}
				}
			}

			if strings.HasPrefix(option, defaultOption) {
				segs := strings.Split(option, equalToken)
				if len(segs) == 2 {
					sp.Default = segs[1]
				}
			} else if strings.HasPrefix(option, optionalOption) || strings.HasPrefix(option, omitemptyOption) {
				required = false
			}

			if strings.HasPrefix(option, exampleOption) {
				segs := strings.Split(option, equalToken)
				if len(segs) == 2 {
					sp.Example = segs[1]
				}
			}
		}
		sp.Required = required
	}

	if len(member.Comment) > 0 {
		sp.Description = strings.TrimLeft(member.Comment, "//")
	}

	return sp
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
			if hasPathParameters(member) || hasHeaderParameters(member) {
				continue
			}
			kv := keyVal{Value: schemaOfField(member)}
			kv.Key = member.Name
			if tag, err := member.GetPropertyName(); err == nil {
				kv.Key = tag
			}
			if kv.Key == "" {
				memberStruct, _ := member.Type.(spec.DefineStruct)
				for _, m := range memberStruct.Members {
					if hasHeaderParameters(m) || hasPathParameters(m) {
						continue
					}

					mkv := keyVal{
						Value: schemaOfField(m),
						Key:   m.Name,
					}

					if tag, err := m.GetPropertyName(); err == nil {
						mkv.Key = tag
					}
					if schema.Properties == nil {
						schema.Properties = &swaggerSchemaObjectProperties{}
					}
					*schema.Properties = append(*schema.Properties, mkv)
				}
				continue
			}
			if schema.Properties == nil {
				schema.Properties = &swaggerSchemaObjectProperties{}
			}
			*schema.Properties = append(*schema.Properties, kv)

			for _, tag := range member.Tags() {
				if tag.Key == validateKey {
					continue
				}
				if len(tag.Options) == 0 {
					if !contains(schema.Required, tag.Name) && tag.Name != "required" {
						schema.Required = append(schema.Required, tag.Name)
					}
					continue
				}

				required := true
				for _, option := range tag.Options {
					// case strings.HasPrefix(option, defaultOption):
					// case strings.HasPrefix(option, optionsOption):

					if strings.HasPrefix(option, optionalOption) || strings.HasPrefix(option, omitemptyOption) {
						required = false
					}
				}

				if required && !contains(schema.Required, tag.Name) {
					schema.Required = append(schema.Required, tag.Name)
				}
			}
		}

		d[i2.Name()] = schema
	}
}

func hasPathParameters(member spec.Member) bool {
	for _, tag := range member.Tags() {
		if tag.Key == "path" {
			return true
		}
	}

	return false
}

func hasHeaderParameters(member spec.Member) bool {
	for _, tag := range member.Tags() {
		if tag.Key == "header" {
			return true
		}
	}
	return false
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
		} else if strings.HasPrefix(refTypeName, "[]") {
			core = schemaCore{Type: "array"}

			tempKind := swaggerMapTypes[strings.Replace(refTypeName, "[]", "", -1)]
			ftype, format, ok := primitiveSchema(tempKind, refTypeName)
			if ok {
				core.Items = &swaggerItemsObject{Type: ftype, Format: format}
			} else {
				core.Items = &swaggerItemsObject{Type: ft.String(), Format: "UNKNOWN"}
			}

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
				segs := strings.SplitN(option, equalToken, 2)
				if len(segs) == 2 {
					ret.Enum = strings.Split(segs[1], optionSeparator)
				}
			case strings.HasPrefix(option, rangeOption):
				segs := strings.SplitN(option, equalToken, 2)
				if len(segs) == 2 {
					min, max, ok := parseRangeOption(segs[1])
					if ok {
						ret.Minimum = min
						ret.Maximum = max
					}
				}
			case strings.HasPrefix(option, exampleOption):
				segs := strings.Split(option, equalToken)
				if len(segs) == 2 {
					ret.Example = segs[1]
				}
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
	case reflect.Uint:
		return "integer", "uint32", true
	case reflect.Int8:
		return "integer", "int8", true
	case reflect.Uint8:
		return "integer", "uint8", true
	case reflect.Int16:
		return "integer", "int16", true
	case reflect.Uint16:
		return "integer", "uin16", true
	case reflect.Int64:
		return "integer", "int64", true
	case reflect.Uint64:
		return "integer", "uint64", true
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

func parseHeader(m spec.Member, parameters []swaggerParameterObject) []swaggerParameterObject {

	tempKind := swaggerMapTypes[strings.Replace(m.Type.Name(), "[]", "", -1)]
	ftype, format, ok := primitiveSchema(tempKind, m.Type.Name())
	if !ok {
		ftype = tempKind.String()
		format = "UNKNOWN"
	}
	sp := swaggerParameterObject{In: "header", Type: ftype, Format: format, Schema: &swaggerSchemaObject{}}

	for _, tag := range m.Tags() {
		sp.Name = tag.Name
		if len(tag.Options) == 0 {
			sp.Required = true
			continue
		}

		required := true
		for _, option := range tag.Options {
			if strings.HasPrefix(option, optionsOption) {
				segs := strings.SplitN(option, equalToken, 2)
				if len(segs) == 2 {
					sp.Enum = strings.Split(segs[1], optionSeparator)
				}
			}

			if strings.HasPrefix(option, rangeOption) {
				segs := strings.SplitN(option, equalToken, 2)
				if len(segs) == 2 {
					min, max, ok := parseRangeOption(segs[1])
					if ok {
						sp.Schema.Minimum = min
						sp.Schema.Maximum = max
					}
				}
			}

			if strings.HasPrefix(option, defaultOption) {
				segs := strings.Split(option, equalToken)
				if len(segs) == 2 {
					sp.Default = segs[1]
				}
			} else if strings.HasPrefix(option, optionalOption) || strings.HasPrefix(option, omitemptyOption) {
				required = false
			}

			if strings.HasPrefix(option, exampleOption) {
				segs := strings.Split(option, equalToken)
				if len(segs) == 2 {
					sp.Example = segs[1]
				}
			}
		}
		sp.Required = required
	}
	sp.Description = strings.TrimLeft(m.Comment, "//")
	if m.Name == "" {
		memberDefineStruct, ok := m.Type.(spec.DefineStruct)
		if !ok {
			return parameters
		}
		for _, cm := range memberDefineStruct.Members {
			if hasHeaderParameters(cm) {
				parameters = parseHeader(cm, parameters)
			}
		}
	}
	return append(parameters, sp)
}
