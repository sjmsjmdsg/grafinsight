package main

import (
	"fmt"
	"os"
	"runtime"

	"github.com/fatih/color"
	"github.com/openinsight-project/grafinsight/pkg/cmd/grafinsight-cli/commands"
	"github.com/openinsight-project/grafinsight/pkg/cmd/grafinsight-cli/logger"
	"github.com/openinsight-project/grafinsight/pkg/cmd/grafinsight-cli/services"
	"github.com/openinsight-project/grafinsight/pkg/cmd/grafinsight-cli/utils"
	"github.com/urfave/cli/v2"
)

var version = "master"

func main() {
	setupLogging()

	app := &cli.App{
		Name: "GrafInsight CLI",
		Authors: []*cli.Author{
			{
				Name:  "GrafInsight Project",
				Email: "hello@grafinsight.com",
			},
		},
		Version: version,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "pluginsDir",
				Usage:   "Path to the Grafinsight plugin directory",
				Value:   utils.GetGrafinsightPluginDir(runtime.GOOS),
				EnvVars: []string{"GF_PLUGIN_DIR"},
			},
			&cli.StringFlag{
				Name:    "repo",
				Usage:   "URL to the plugin repository",
				Value:   "https://grafinsight.com/api/plugins",
				EnvVars: []string{"GF_PLUGIN_REPO"},
			},
			&cli.StringFlag{
				Name:    "pluginUrl",
				Usage:   "Full url to the plugin zip file instead of downloading the plugin from grafinsight.com/api",
				Value:   "",
				EnvVars: []string{"GF_PLUGIN_URL"},
			},
			&cli.BoolFlag{
				Name:  "insecure",
				Usage: "Skip TLS verification (insecure)",
			},
			&cli.BoolFlag{
				Name:  "debug, d",
				Usage: "Enable debug logging",
			},
			&cli.StringFlag{
				Name:  "configOverrides",
				Usage: "Configuration options to override defaults as a string. e.g. cfg:default.paths.log=/dev/null",
			},
			&cli.StringFlag{
				Name:  "homepath",
				Usage: "Path to Grafinsight install/home path, defaults to working directory",
			},
			&cli.StringFlag{
				Name:  "config",
				Usage: "Path to config file",
			},
		},
		Commands:        commands.Commands,
		CommandNotFound: cmdNotFound,
	}

	app.Before = func(c *cli.Context) error {
		services.Init(version, c.Bool("insecure"))
		return nil
	}

	if err := app.Run(os.Args); err != nil {
		logger.Errorf("%s: %s %s\n", color.RedString("Error"), color.RedString("✗"), err)
		os.Exit(1)
	}
}

func setupLogging() {
	for _, f := range os.Args {
		if f == "-d" || f == "--debug" || f == "-debug" {
			logger.SetDebug(true)
		}
	}
}

func cmdNotFound(c *cli.Context, command string) {
	fmt.Printf(
		"%s: '%s' is not a %s command. See '%s --help'.\n",
		c.App.Name,
		command,
		c.App.Name,
		os.Args[0],
	)
	os.Exit(1)
}
