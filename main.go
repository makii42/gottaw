package main

import (
	"os"

	"github.com/makii42/gottaw/config"
	"github.com/makii42/gottaw/defaults"
	"github.com/makii42/gottaw/output"
	"github.com/makii42/gottaw/watch"
	"gopkg.in/urfave/cli.v1"
)


func main() {
	app := cli.NewApp()
	app.Name = "gotta watch"
	app.Usage = "Run command(s) when files in the folder change."
	app.EnableBashCompletion = true
	app.Action = watch.WatchCmd.Action
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "f, folder",
			Value:  ".",
			Usage:  "Folder to watch for changes",
			EnvVar: "GOTTAW_FOLDER",
		},
		cli.StringFlag{
			Name:  "c, config",
			Value: "./.gottaw.yml",
			Usage: "Config file to read",
			Destination: &config.File,
		},
		cli.StringFlag{
			Name:  "d, delay",
			Value: "1500ms",
			Usage: "Delay of the pipeline action after event",
		},
		cli.BoolFlag{
			Name:        "t, trace",
			Usage:       "Log more details",
			Destination: &output.Trace,
		},
		cli.BoolFlag{
			Name:        "q, quiet",
			Usage:       "Logs nothing but pass-through output and success messages",
			Destination: &output.Quiet,
		},
		cli.BoolFlag{
			Name:  "g, growl",
			Usage: "Notify OS via growl about pipeline result",
		},
	}
	app.Commands = []cli.Command{
		watch.WatchCmd,
		watch.OneRunCmd,
		defaults.DefaultsCmd,
	}

	app.Run(os.Args)
}
