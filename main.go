package main

import (
	"os"

	"fmt"

	"github.com/makii42/gottaw/config"
	"github.com/makii42/gottaw/defaults"
	"github.com/makii42/gottaw/output"
	"github.com/makii42/gottaw/watch"
	"gopkg.in/urfave/cli.v1"
	"log"
	"reflect"
)

var (
	Trace  bool
	Quiet  bool
	Output *output.Output
	Config *config.Config
)

type ActionFactoryFunc func(*config.Config, *output.Output) cli.ActionFunc

func main() {
	app := cli.NewApp()
	app.Name = "gotta watch"
	app.Usage = "Run command(s) when files in the folder change."
	app.EnableBashCompletion = true
	if wa, ok := watch.WatchCmd.Action.(ActionFactoryFunc); ok {
		app.Action = inject(wa)
		log.Printf("Action: %#v", app.Action)
	} else {
		app.Action = watch.WatchCmd.Action
	}
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
			Name:        "t, trace",
			Usage:       "Log more details",
			Destination: &Trace,
		},
		cli.BoolFlag{
			Name:        "q, quiet",
			Usage:       "Logs nothing but pass-through output and success messages",
			Destination: &Quiet,
		},
		cli.BoolFlag{
			Name:  "g, growl",
			Usage: "Notify OS via growl about pipeline result",
		},
	}
	app.Commands = []cli.Command{
		filterCmd(watch.WatchCmd),
		filterCmd(watch.OneRunCmd),
		defaults.DefaultsCmd,
	}

	app.Run(os.Args)
}

func filterCmd(cmd cli.Command) cli.Command {
	if cmd.Action != nil {
		log.Printf("%s: testing for factory", cmd.Name)
		// this is a hack and MIGHT not work in future versions
		// right now, the Command.Action field is of type interface{},
		// which allows us to to check for type ActionFactoryFunc and
		// wrap-and-swap the actual action that is returned by it.
		// If that code is deprecated, we revisit this.
		if af, ok := cmd.Action.(ActionFactoryFunc); ok {
			log.Printf("%s: replacing action for with injector(1)", cmd.Name)
			cmd.Action = inject(af)
		} else if af, ok := cmd.Action.(func(*config.Config, *output.Output) cli.ActionFunc); ok {
			log.Printf("%s: replacing action for with injector (2)", cmd.Name)
		    cmd.Action = inject(af)
		} else {
			t := reflect.TypeOf(cmd.Action)
			log.Printf("%s: NOPE!!! It's a %v", cmd.Name, t)
		}
	}
	return cmd
}

func inject(f ActionFactoryFunc) cli.ActionFunc {

	return func(ctx *cli.Context) error {

		// setup Config
		Config = config.Setup(ctx.GlobalString("config"))

		// set up logging
		var lvl output.Level
		if Trace && Quiet {
			return fmt.Errorf("please specify either --trace or --quiet")
		} else if Trace {
			lvl = output.L_TRACE
		} else if Quiet {
			lvl = output.L_QUIET
		} else {
			lvl = output.L_NOTICE
		}
		Output = output.NewLog(lvl, Config)

		return f(Config, Output)(ctx)
	}
}

