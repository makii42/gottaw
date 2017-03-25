package defaults

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	c "github.com/makii42/gottaw/config"
	"github.com/makii42/gottaw/output"
	"github.com/urfave/cli"
)

var guesser DefaultGuesser
var l *output.Logger

func init() {
	l = output.NewLogger(&c.Config{})
	guesser = DefaultGuesser{
		GolangDefault{},
		NodeYarnDefault{},
		NodeNpmDefault{},
		JavaMavenDefault{},
	}
}

var DefaultsCmd = cli.Command{
	Name:   "defaults",
	Usage:  "Prints and optionally writes the defaults for a folder",
	Action: defaults,
	Flags:  []cli.Flag{},
}

func defaults(cli *cli.Context) {
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
		def := guesser.Find(rootDir)
		if def != nil {
			l.Successf("🎯  Identified default %s\n", def.Name())
		} else {
			l.Errorf("🚫  No known default matched contents of %s\n", rootDir)
			fmt.Printf("\nAvailable defaults are:\n\n")
			for _, def := range guesser {
				fmt.Printf("  - %s\n", def.Name())
			}
			fmt.Println("\nFeel free to contribute your default at https://github.com/makii42/gottaw")
		}
	}

}

func GuessDefault(path string) Default {
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
