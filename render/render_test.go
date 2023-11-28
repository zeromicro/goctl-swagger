package render

import (
	"github.com/aishuchen/goctl-swagger/render/types"
	"github.com/zeromicro/go-zero/tools/goctl/api/parser"
	"github.com/zeromicro/go-zero/tools/goctl/plugin"
	"os"
	"testing"
)

func TestRender20(t *testing.T) {
	apiPath := os.Getenv("GOCTL_API_PATH")
	apiSpec, err := parser.Parse(apiPath)
	if err != nil {
		t.Fatal(err)
	}
	plg := &plugin.Plugin{
		Api: apiSpec,
	}
	outPath := os.Getenv("SWAGGER_OUT_PATH")
	opt := types.Option{
		Target:     outPath,
		Version:    "2.0",
		RenderType: "json",
	}
	if err := Render(plg, opt); err != nil {
		t.Fatal(err)
	}
}
