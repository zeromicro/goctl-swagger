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
	std, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		return err
	}

	var plugin generate.Plugin
	plugin.ParentPackage = pkg

	err = json.Unmarshal(std, &plugin)
	if err != nil {
		return err
	}

	return generate.Do(plugin)
}
