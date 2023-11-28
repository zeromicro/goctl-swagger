package v2

import (
	"github.com/aishuchen/goctl-swagger/render/types"
	"github.com/zeromicro/go-zero/tools/goctl/plugin"
	"strings"
)

type Renderer struct {
}

func (r *Renderer) Render(plg *plugin.Plugin, opt types.Option) (types.Swagger, error) {
	registerTypes(plg.Api.Types)
	var contact *Contact
	if len(plg.Api.Info.Email) > 0 || len(plg.Api.Info.Author) > 0 {
		contact = &Contact{
			Name:  plg.Api.Info.Author,
			Email: plg.Api.Info.Email,
		}
	}
	info := Information{
		Title:       strings.Trim(plg.Api.Info.Properties["title"], `"`),
		Description: plg.Api.Info.Properties["desc"],
		Version:     plg.Api.Info.Properties["version"],
		Contact:     contact,
	}
	paths := renderPaths(plg.Api.Service)
	swagger := &Swagger{
		Swagger:  "2.0",
		Info:     info,
		Consumes: []string{"application/json"},
		Produces: []string{"application/json"},
		Paths:    paths,
		Schemes:  opt.Schemes,
	}
	swagger.Definitions = models
	return swagger, nil
}
