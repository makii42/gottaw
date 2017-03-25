package main

import (
	"os"

	"github.com/makii42/gottaw/defaults"
	"github.com/makii42/gottaw/watch"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "gotta watch"
	app.Usage = "Run command(s) when files in the folder change."
	app.Action = watch.WatchIt
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
		},
		cli.StringFlag{
			Name:  "d, delay",
			Value: "1500ms",
			Usage: "Delay of the pipeline action after event",
		},
		cli.BoolFlag{
			Name:  "g, growl",
			Usage: "Notify OS via growl about pipeline result",
		},
	}
	app.Commands = []cli.Command{
		watch.WatchCmd,
		defaults.DefaultsCmd,
	}
	app.Run(os.Args)
}
