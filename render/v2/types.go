package v2

import (
	"bytes"
	"encoding/json"
	"gopkg.in/yaml.v3"
)

// Swagger 2.0
// Swagger list the resource
type Swagger struct {
	Swagger             string               `json:"swagger,omitempty" yaml:"swagger,omitempty"`
	Info                Information          `json:"info" yaml:"info"`
	Host                string               `json:"host,omitempty" yaml:"host,omitempty"`
	BasePath            string               `json:"basePath,omitempty" yaml:"basePath,omitempty"`
	Schemes             []string             `json:"schemes,omitempty" yaml:"schemes,omitempty"`
	Consumes            []string             `json:"consumes,omitempty" yaml:"consumes,omitempty"`
	Produces            []string             `json:"produces,omitempty" yaml:"produces,omitempty"`
	Paths               map[string]*Path     `json:"paths" yaml:"paths"`
	Definitions         map[string]*Schema   `json:"definitions,omitempty" yaml:"definitions,omitempty"`
	SecurityDefinitions map[string]*Security `json:"securityDefinitions,omitempty" yaml:"securityDefinitions,omitempty"`
	Security            map[string][]string  `json:"security,omitempty" yaml:"security,omitempty"`
	Tags                []*Tag               `json:"tags,omitempty" yaml:"tags,omitempty"`
	ExternalDocs        *ExternalDocs        `json:"externalDocs,omitempty" yaml:"externalDocs,omitempty"`
}

func (s *Swagger) EncodeJSON() ([]byte, error) {
	var formatted bytes.Buffer
	enc := json.NewEncoder(&formatted)
	enc.SetIndent("", "  ")

	if err := enc.Encode(s); err != nil {
		return nil, err
	}
	return formatted.Bytes(), nil
}

func (s *Swagger) EncodeYAML() ([]byte, error) {
	return yaml.Marshal(s)
}

// Information Provides metadata about the API. The metadata can be used by the clients if needed.
type Information struct {
	Title          string   `json:"title" yaml:"title"`
	Description    string   `json:"description,omitempty" yaml:"description,omitempty"`
	Version        string   `json:"version" yaml:"version"`
	TermsOfService string   `json:"termsOfService,omitempty" yaml:"termsOfService,omitempty"`
	Contact        *Contact `json:"contact,omitempty" yaml:"contact,omitempty"`
	License        *License `json:"license,omitempty" yaml:"license,omitempty"`
}

// Contact information for the exposed API.
type Contact struct {
	Name  string `json:"name,omitempty" yaml:"name,omitempty"`
	URL   string `json:"url,omitempty" yaml:"url,omitempty"`
	Email string `json:"email,omitempty" yaml:"email,omitempty"`
}

// License information for the exposed API.
type License struct {
	Name string `json:"name,omitempty" yaml:"name,omitempty"`
	URL  string `json:"url,omitempty" yaml:"url,omitempty"`
}

// Path Describes the operations available on a single path.
type Path struct {
	Ref     string     `json:"$ref,omitempty" yaml:"$ref,omitempty"`
	Get     *Operation `json:"get,omitempty" yaml:"get,omitempty"`
	Put     *Operation `json:"put,omitempty" yaml:"put,omitempty"`
	Post    *Operation `json:"post,omitempty" yaml:"post,omitempty"`
	Delete  *Operation `json:"delete,omitempty" yaml:"delete,omitempty"`
	Options *Operation `json:"options,omitempty" yaml:"options,omitempty"`
	Head    *Operation `json:"head,omitempty" yaml:"head,omitempty"`
	Patch   *Operation `json:"patch,omitempty" yaml:"patch,omitempty"`
}

// Operation Describes a single API operation on a path.
type Operation struct {
	Tags        []string             `json:"tags,omitempty" yaml:"tags,omitempty"`
	Summary     string               `json:"summary,omitempty" yaml:"summary,omitempty"`
	Description string               `json:"description,omitempty" yaml:"description,omitempty"`
	OperationID string               `json:"operationId,omitempty" yaml:"operationId,omitempty"`
	Consumes    []string             `json:"consumes,omitempty" yaml:"consumes,omitempty"`
	Produces    []string             `json:"produces,omitempty" yaml:"produces,omitempty"`
	Schemes     []string             `json:"schemes,omitempty" yaml:"schemes,omitempty"`
	Parameters  []*Parameter         `json:"parameters,omitempty" yaml:"parameters,omitempty"`
	Responses   map[string]*Response `json:"responses,omitempty" yaml:"responses,omitempty"`
	Deprecated  bool                 `json:"deprecated,omitempty" yaml:"deprecated,omitempty"`
}

// Parameter Describes a single operation parameter.
type Parameter struct {
	In          string          `json:"in,omitempty" yaml:"in,omitempty"`
	Name        string          `json:"name,omitempty" yaml:"name,omitempty"`
	Description string          `json:"description,omitempty" yaml:"description,omitempty"`
	Required    bool            `json:"required,omitempty" yaml:"required,omitempty"`
	Schema      *Schema         `json:"schema,omitempty" yaml:"schema,omitempty"`
	Type        string          `json:"type,omitempty" yaml:"type,omitempty"`
	Format      string          `json:"format,omitempty" yaml:"format,omitempty"`
	Items       *ParameterItems `json:"items,omitempty" yaml:"items,omitempty"`
	Default     interface{}     `json:"default,omitempty" yaml:"default,omitempty"`
}

// A limited subset of JSON-Schema's items object. It is used by parameter definitions that are not located in "body".
// http://swagger.io/specification/#itemsObject
type ParameterItems struct {
	Type             string          `json:"type,omitempty" yaml:"type,omitempty"`
	Format           string          `json:"format,omitempty" yaml:"format,omitempty"`
	Items            *ParameterItems `json:"items,omitempty" yaml:"items,omitempty"` //Required if type is "array". Describes the type of items in the array.
	CollectionFormat string          `json:"collectionFormat,omitempty" yaml:"collectionFormat,omitempty"`
	Default          string          `json:"default,omitempty" yaml:"default,omitempty"`
}

// Schema Object allows the definition of input and output data types.
type Schema struct {
	Ref                  string             `json:"$ref,omitempty" yaml:"$ref,omitempty"`
	Title                string             `json:"title,omitempty" yaml:"title,omitempty"`
	Description          string             `json:"description,omitempty" yaml:"description,omitempty"`
	Default              interface{}        `json:"default,omitempty" yaml:"default,omitempty"`
	Type                 string             `json:"type,omitempty" yaml:"type,omitempty"`
	Example              string             `json:"example,omitempty" yaml:"example,omitempty"`
	Required             []string           `json:"required,omitempty" yaml:"required,omitempty"`
	Format               string             `json:"format,omitempty" yaml:"format,omitempty"`
	ReadOnly             bool               `json:"readOnly,omitempty" yaml:"readOnly,omitempty"`
	Items                *Schema            `json:"items,omitempty" yaml:"items,omitempty"`
	AllOf                []*Schema          `json:"allOf,omitempty" yaml:"allOf,omitempty"`
	Properties           map[string]*Schema `json:"properties,omitempty" yaml:"properties,omitempty"`
	AdditionalProperties *Schema            `json:"additionalProperties,omitempty" yaml:"additionalProperties,omitempty"`

	required bool // 当认为 Schema 是 Property 时使用
}

// Property are taken from the JSON Schema definition but their definitions were adjusted to the Swagger Specification
//type Property struct {
//	Ref                  string               `json:"$ref,omitempty" yaml:"$ref,omitempty"`
//	Title                string               `json:"title,omitempty" yaml:"title,omitempty"`
//	Description          string               `json:"description,omitempty" yaml:"description,omitempty"`
//	Default              interface{}          `json:"default,omitempty" yaml:"default,omitempty"`
//	Type                 string               `json:"type,omitempty" yaml:"type,omitempty"`
//	Example              string               `json:"example,omitempty" yaml:"example,omitempty"`
//	Required             []string             `json:"required,omitempty" yaml:"required,omitempty"`
//	Format               string               `json:"format,omitempty" yaml:"format,omitempty"`
//	ReadOnly             bool                 `json:"readOnly,omitempty" yaml:"readOnly,omitempty"`
//	Properties           map[string]*Property `json:"properties,omitempty" yaml:"properties,omitempty"`
//	Items                *Property            `json:"items,omitempty" yaml:"items,omitempty"`
//	AdditionalProperties *Property            `json:"additionalProperties,omitempty" yaml:"additionalProperties,omitempty"`
//}

// Response as they are returned from executing this operation.
type Response struct {
	Description string  `json:"description" yaml:"description"`
	Schema      *Schema `json:"schema,omitempty" yaml:"schema,omitempty"`
	Ref         string  `json:"$ref,omitempty" yaml:"$ref,omitempty"`
}

// Security Allows the definition of a security scheme that can be used by the operations
type Security struct {
	Type             string            `json:"type,omitempty" yaml:"type,omitempty"` // Valid values are "basic", "apiKey" or "oauth2".
	Description      string            `json:"description,omitempty" yaml:"description,omitempty"`
	Name             string            `json:"name,omitempty" yaml:"name,omitempty"`
	In               string            `json:"in,omitempty" yaml:"in,omitempty"`     // Valid values are "query" or "header".
	Flow             string            `json:"flow,omitempty" yaml:"flow,omitempty"` // Valid values are "implicit", "password", "application" or "accessCode".
	AuthorizationURL string            `json:"authorizationUrl,omitempty" yaml:"authorizationUrl,omitempty"`
	TokenURL         string            `json:"tokenUrl,omitempty" yaml:"tokenUrl,omitempty"`
	Scopes           map[string]string `json:"scopes,omitempty" yaml:"scopes,omitempty"` // The available scopes for the OAuth2 security scheme.
}

// Tag Allows adding meta data to a single tag that is used by the Operation Object
type Tag struct {
	Name         string        `json:"name,omitempty" yaml:"name,omitempty"`
	Description  string        `json:"description,omitempty" yaml:"description,omitempty"`
	ExternalDocs *ExternalDocs `json:"externalDocs,omitempty" yaml:"externalDocs,omitempty"`
}

// ExternalDocs include Additional external documentation
type ExternalDocs struct {
	Description string `json:"description,omitempty" yaml:"description,omitempty"`
	URL         string `json:"url,omitempty" yaml:"url,omitempty"`
}
