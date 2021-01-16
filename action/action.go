package action

import (
	"encoding/json"
	"github.com/urfave/cli/v2"
	"github.com/zeromicro/goctl-swagger/generate"
	"io/ioutil"
	"os"
)

func Generator(ctx *cli.Context) error {
	pkg := ctx.String("package")
	fileName := ctx.String("filename")
	std, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		return err
	}

	var plugin generate.Plugin
	plugin.ParentPackage = pkg
	plugin.FileName = fileName

	if len(plugin.FileName) == 0 {
		plugin.FileName = "rest.swagger.json"
	}
	err = json.Unmarshal(std, &plugin)
	if err != nil {
		return err
	}

	return generate.Do(plugin)
}
