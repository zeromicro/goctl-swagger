package action

import (
	"encoding/json"
	"github.com/shyandsy/goctl-swagger/generate"
	"github.com/urfave/cli/v2"
	"github.com/zeromicro/go-zero/tools/goctl/api/parser"
	"github.com/zeromicro/go-zero/tools/goctl/config"
	"github.com/zeromicro/go-zero/tools/goctl/plugin"
	"github.com/zeromicro/go-zero/tools/goctl/util/pathx"
	"path/filepath"
)

func Generator(ctx *cli.Context) error {
	fileName := ctx.String("filename")
	apiFileName := ctx.String("api")
	generatedDir := ctx.String("dir")

	if len(fileName) == 0 {
		fileName = "rest.swagger.json"
	}

	content, err := prepareArgs(apiFileName, generatedDir)
	if err != nil {
		panic(err)
	}

	p, err := NewPlugin(content)
	if err != nil {
		return err
	}
	basepath := ctx.String("basepath")
	host := ctx.String("host")
	return generate.Do(fileName, host, basepath, p)
}

func NewPlugin(content []byte) (*plugin.Plugin, error) {
	var plugin plugin.Plugin

	var info struct {
		ApiFilePath string
		Style       string
		Dir         string
	}
	err := json.Unmarshal(content, &info)
	if err != nil {
		return nil, err
	}

	plugin.ApiFilePath = info.ApiFilePath
	plugin.Style = info.Style
	plugin.Dir = info.Dir
	api, err := parser.Parse(info.ApiFilePath)
	if err != nil {
		return nil, err
	}

	plugin.Api = api
	return &plugin, nil
}

func prepareArgs(apiPath string, generatedDir string) ([]byte, error) {
	var transferData plugin.Plugin
	if len(apiPath) > 0 && pathx.FileExists(apiPath) {
		api, err := parser.Parse(apiPath)
		if err != nil {
			return nil, err
		}
		transferData.Api = api
	}
	absApiFilePath, err := filepath.Abs(apiPath)
	if err != nil {
		return nil, err
	}
	transferData.ApiFilePath = absApiFilePath
	dirAbs, err := filepath.Abs(generatedDir)
	if err != nil {
		return nil, err
	}

	transferData.Dir = dirAbs
	transferData.Style = config.DefaultFormat //VarStringStyle
	content, err := json.Marshal(transferData)
	if err != nil {
		return nil, err
	}
	return content, nil
}
