package main

import (
	"fmt"
	"os"
	"runtime"

	"github.com/urfave/cli/v2"
	"github.com/zeromicro/goctl-swagger/action"
)

var (
	version  = "20220621"
	commands = []*cli.Command{
		{
			Name:   "swagger",
			Usage:  "generates swagger.json",
			Action: action.Generator,
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:  "host",
					Usage: "api request address",
				},
				&cli.StringFlag{
					Name:  "basepath",
					Usage: "url request prefix",
				},
				&cli.StringFlag{
					Name:  "filename",
					Usage: "swagger save file name",
				},
			},
		},
	}
)

func main() {
	app := cli.NewApp()
	app.Usage = "a plugin of goctl to generate swagger.json"
	app.Version = fmt.Sprintf("%s %s/%s", version, runtime.GOOS, runtime.GOARCH)
	app.Commands = commands
	if err := app.Run(os.Args); err != nil {
		fmt.Printf("goctl-swagger: %+v\n", err)
	}
}
