package defaults

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"io/ioutil"

	c "github.com/makii42/gottaw/config"
	"github.com/makii42/gottaw/output"
	"github.com/urfave/cli"
)

var DefaultsCmd = cli.Command{
	Name:   "defaults",
	Usage:  "Prints and optionally writes the defaults for a folder",
	Action: defaults,
	Flags: []cli.Flag{
		cli.BoolFlag{
			Name:  "w, write",
			Usage: "Writes default config to set configuration file or default location.",
		},
	},
}

func defaults(cli *cli.Context) {
	conf := &c.Config{}
	trace := cli.GlobalBool("trace")
	var l *output.Logger
	if trace {
		l = output.NewLogger(output.TRACE, conf)
	} else {
		l = output.NewLogger(output.NOTICE, conf)
	}
	configFile, err := filepath.Abs(cli.GlobalString("config"))
	if err != nil {
		panic(err)
	}
	file, err := os.Stat(configFile)
	if err == nil && file.Mode().IsRegular() {
		// file exists
		fmt.Printf("Config file exists: %s\n", configFile)
	}
	// err != nil assumes file does not exist.
	// checking if FOLDER exists
	rootDir := filepath.Dir(configFile)
	dir, err := os.Stat(rootDir)
	if err != nil {
		panic(err)
	}
	if dir.IsDir() {
		l.Noticef("🔬  evaluating %s\n", rootDir)
		def := GuessDefault(rootDir, l)
		if def != nil {
			l.Successf("🎯  Identified default %s\n", def.Name())
			if cli.Bool("write") {
				data, err := c.SerializeConfig(def.Config())
				if err != nil {
					log.Fatalf("error serializing default: %s", err)
				}
				fmt.Printf(
					"Default config for %s:\n===\n%s===\nWrite to '%s'? [y/N] ",
					def.Name(),
					string(data),
					cli.GlobalString("config"),
				)
				reader := bufio.NewReader(os.Stdin)
				input, _ := reader.ReadString('\n')
				input = strings.Trim(input, " \n")
				if strings.ToLower(input) == "y" {
					err := ioutil.WriteFile(cli.GlobalString("config"), data, 0660)
					if err != nil {
						log.Fatalf("error writing '%s': %s", cli.GlobalString("config"), err)
					}
					l.Successf("✅  Okay!\n")
				} else {
					l.Noticef("🌮  Okay, never mind!\n")
				}
			}
		} else {
			l.Errorf("🚫  No known default matched contents of %s\n", rootDir)
			fmt.Println("\nFeel free to contribute your default at https://github.com/makii42/gottaw")
		}
	}
}

func GuessDefault(path string, l *output.Logger) Default {
	util := newDefaultsUtil(l)

	guesser := DefaultGuesser{
		NewGolangDefault(util),
		NewNodeYarnDefault(util),
		NewNodeNpmDefault(util),
		NewJavaMavenDefault(util),
	}

	workdir, err := filepath.Abs(path)
	if err != nil {
		log.Fatal(err)
	}
	def := guesser.Find(workdir)
	return def
}

type DefaultGuesser []Default

func (d DefaultGuesser) Find(dir string) Default {
	for _, candidate := range d {
		if candidate.Test(dir) {
			return candidate
		}
	}
	return nil
}

type Default interface {
	Name() string
	Test(dir string) bool
	Config() *c.Config
}
