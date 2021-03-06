//go:generate swag init --dir ./ --generalInfo routes/routes.go --propertyStrategy snakecase --output ./routes/docs

// Package main freedom is my stock bot
package main

import (
	"fmt"
	"os"
	"time"

	"github.com/283713406/freedom/cmds"
	"github.com/283713406/freedom/logging"
	"github.com/283713406/freedom/version"
	"github.com/spf13/viper"
	"github.com/urfave/cli/v2"
)

var (
	// DefaultLoglevel 日志级别默认值
	DefaultLoglevel = "info"
	// ProcessorOptions 要启动运行的进程可选项
	ProcessorOptions = []string{cmds.ProcessorChecker, cmds.ProcessorExportor}
)

func init() {
	viper.SetDefault("app.chan_size", 50)
}

func main() {
	app := cli.NewApp()
	app.EnableBashCompletion = true
	app.Name = "freedom"
	app.Usage = "mamba 的股票筛选和检测程序"
	app.UsageText = `该程序不构成任何投资建议，程序只是个人辅助工具，具体分析仍然需要自己判断。`

	app.Version = version.Version
	app.Compiled = time.Now()
	app.Authors = []*cli.Author{
		{
			Name:  "mamba",
			Email: "283713406@qq.com",
		},
	}
	app.Copyright = "(c) 2021 mamba"

	cli.VersionFlag = &cli.BoolFlag{
		Name:    "version",
		Aliases: []string{"v"},
		Usage:   "show the version",
	}

	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:        "loglevel",
			Aliases:     []string{"l"},
			Value:       DefaultLoglevel,
			Usage:       "cmd 日志级别 [debug|info|warn|error]",
			EnvVars:     []string{"FREEDOM_CMD_LOGLEVEL"},
			DefaultText: DefaultLoglevel,
		},
	}
	app.BashComplete = func(c *cli.Context) {
		if c.NArg() > 0 {
			return
		}
		for _, i := range ProcessorOptions {
			fmt.Println(i)
		}
	}

	app.Commands = append(app.Commands, cmds.CommandExportor())
	app.Commands = append(app.Commands, cmds.CommandChecker())

	if err := app.Run(os.Args); err != nil {
		logging.Fatal(nil, err.Error())
	}

}
