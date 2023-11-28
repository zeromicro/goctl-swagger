package v2

import (
	"github.com/aishuchen/goctl-swagger/render/types"
	"github.com/zeromicro/go-zero/tools/goctl/api/parser"
	"github.com/zeromicro/go-zero/tools/goctl/plugin"
	"os"
	"testing"
)

func TestRenderer_Render(t *testing.T) {
	apiPath := os.Getenv("GOCTL_API_PATH")
	renderer := &Renderer{}
	apiSpec, err := parser.Parse(apiPath)
	if err != nil {
		t.Fatal(err)
	}
	plg := &plugin.Plugin{
		Api: apiSpec,
	}
	swagger, err := renderer.Render(plg, types.Option{})
	if err != nil {
		t.Fatal(err)
	}
	data, err := swagger.EncodeJSON()
	if err != nil {
		t.Fatal(err)
	}
	outPath := os.Getenv("SWAGGER_OUT_PATH")
	if err := os.WriteFile(outPath, data, 0644); err != nil {
		t.Fatal(err)
	}
}
