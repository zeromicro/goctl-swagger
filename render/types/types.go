package types

import (
	"github.com/zeromicro/go-zero/tools/goctl/plugin"
)

type Option struct {
	Host       string
	BasePath   string
	Schemes    []string
	RenderType string
	Version    string
	Target     string
}

type Renderer interface {
	Render(plg *plugin.Plugin, opt Option) (Swagger, error)
}

type Swagger interface {
	EncodeJSON() ([]byte, error)
	EncodeYAML() ([]byte, error)
}

// Go zero defined tag keys.
const (
	JsonTagKey        = "json"
	FormTagKey        = "form"
	PathTagKey        = "path"
	HeaderTagKey      = "header"
	DefaultSummaryKey = "summary"
)
