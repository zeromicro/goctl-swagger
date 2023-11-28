package render

import (
	"bytes"
	"github.com/aishuchen/goctl-swagger/render/types"
	"github.com/aishuchen/goctl-swagger/render/v2"
	"github.com/urfave/cli/v2"
	"github.com/zeromicro/go-zero/tools/goctl/plugin"
	"os"
	"strings"
)

func Render(plg *plugin.Plugin, opt types.Option) error {
	var renderer types.Renderer
	if opt.Version == "" {
		opt.Version = "2.0"
	}
	if opt.RenderType == "" {
		opt.RenderType = "json"
	}
	if opt.Target == "" {
		opt.Target = "swagger.json"
	}
	if len(opt.Schemes) == 0 {
		opt.Schemes = []string{"http"}
	}
	renderer = getRenderer(opt.Version)
	swagger, err := renderer.Render(plg, opt)
	if err != nil {
		return err
	}
	switch opt.RenderType {
	case "json":
		return renderJson(swagger, opt.Target)
	case "yaml":
		return renderYaml(swagger, opt.Target)
	}
	return nil
}

func getRenderer(version string) types.Renderer {
	switch version {
	case "2.0":
		return &v2.Renderer{}
	default:
		return &v2.Renderer{}
	}
}

func renderJson(swag types.Swagger, target string) error {
	content, err := swag.EncodeJSON()
	if err != nil {
		return err
	}
	buf := bytes.NewBuffer(content)
	if err = writeFile(target, buf); err != nil {
		return err
	}
	return nil
}

func renderYaml(swag types.Swagger, target string) error {
	content, err := swag.EncodeYAML()
	if err != nil {
		return err
	}
	buf := bytes.NewBuffer(content)
	if err = writeFile(target, buf); err != nil {
		return err
	}
	return nil
}

func writeFile(path string, buffer *bytes.Buffer) error {
	return os.WriteFile(path, buffer.Bytes(), 0666)
}

func Do(ctx *cli.Context) error {
	target := ctx.String("target")
	basepath := ctx.String("basepath")
	host := ctx.String("host")
	scheme := ctx.String("schemes")
	var schemes []string
	if len(scheme) > 0 {
		schemes = strings.Split(scheme, ",")
	}
	opt := types.Option{
		Host:       host,
		BasePath:   basepath,
		Schemes:    schemes,
		Target:     target,
		Version:    "2.0",  // TODO: can be configurable
		RenderType: "json", // TODO: can be configurable
	}
	p, err := plugin.NewPlugin()
	if err != nil {
		return err
	}
	return Render(p, opt)
}
