package action

import (
	plugin2 "github.com/tal-tech/go-zero/tools/goctl/plugin"
	"github.com/urfave/cli/v2"
	"github.com/zeromicro/goctl-swagger/generate"
)

func Generator(ctx *cli.Context) error {

	fileName := ctx.String("filename")

	if len(fileName) == 0 {
		fileName = "rest.swagger.json"
	}

	p, err := plugin2.NewPlugin()
	if err != nil {
		return err
	}
	return generate.Do(fileName, p)
}
