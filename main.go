package main

import (
	"fmt"
	"github.com/aishuchen/goctl-swagger/render"
	"github.com/urfave/cli/v2"
	"os"
	"runtime"
)

var (
	version  = "0.0.1"
	commands = []*cli.Command{
		{
			Name:  "swagger",
			Usage: "generates swagger json file",
			//Action: action.Generator,
			Action: render.Do,
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
					Name:  "target",
					Usage: "swagger save file name",
				},
				&cli.StringFlag{
					Name:  "schemes",
					Usage: "swagger support schemes: http, https, ws, wss",
				},
			},
		},
		{
			Name:   "version",
			Action: printVersion,
		},
	}
)

func main() {
	app := cli.NewApp()
	app.Usage = "a plugin of goctl to generate swagger json file"
	app.Version = fmt.Sprintf("%s %s/%s", version, runtime.GOOS, runtime.GOARCH)
	app.Commands = commands
	if err := app.Run(os.Args); err != nil {
		fmt.Printf("goctl-swagger: %+v\n", err)
	}
}

func printVersion(ctx *cli.Context) error {
	fmt.Println("goctl-swagger version", version)
	return nil
}
