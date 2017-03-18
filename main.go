package main

import (
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"strings"

	"github.com/fsnotify/fsnotify"
	"github.com/ianschenck/envflag"
	"github.com/urfave/cli"
	"gopkg.in/yaml.v2"
)

const (
	DEFAULT_FILE_NAME = "./.gottaw.yml"
)

func main() {
	envflag.Parse()
	app := cli.NewApp()

	app.Name = "gotta watch"
	app.Usage = "Run command(s) when files in the folder change."
	app.Action = WatchIt
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
			Value: "800ms",
			Usage: "Delay of the pipeline action after event",
		},
	}

	app.Run(os.Args)
}

// WatchIt does the work
func WatchIt(c *cli.Context) error {
	var cfg *Config

	if c.IsSet("config") {
		if parsed, err := parseConfig(c.String("config")); err != nil {
			panic(err)
		} else {
			cfg = parsed
		}
	} else if _, err := os.Stat(DEFAULT_FILE_NAME); err == nil {
		log.Printf("Using default config from %s", DEFAULT_FILE_NAME)
		if parsed, err := parseConfig(DEFAULT_FILE_NAME); err != nil {
			panic(err)
		} else {
			cfg = parsed
		}
	}
	delay, err := time.ParseDuration(c.String("delay"))
	if err != nil {
		panic(err)
	}
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		panic(err)
	}
	defer watcher.Close()

	watchlist := make(Watchlist)

	done := make(chan bool)
	go func() {
		var timer *time.Timer
		for {
			select {
			case ev := <-watcher.Events:
				if ev.Op&fsnotify.Chmod == fsnotify.Chmod {
					continue
				}
				if ignore, err := isIgnored(ev.Name, cfg); err != nil {
					panic(err)
				} else if ignore {
					continue
				}
				log.Printf("event: %#v", ev)
				if timer == nil {
					timer = time.AfterFunc(delay, createAction(ev, cfg.Pipeline, func() {
						timer = nil
					}))
				} else if timer != nil {
					timer.Reset(delay)
				}

			case err := <-watcher.Errors:
				log.Println("error:", err)
			}
		}

	}()

	if f, err := os.Stat(c.String("folder")); err != nil {
		panic(err)
	} else if err := watchDirRecursive(f.Name(), watcher, watchlist, cfg); err != nil {
		panic(err)
	}
	log.Printf("Watchlist those folders: %v", watchlist)
	<-done

	return nil
}

func watchDirRecursive(dir string, watcher *fsnotify.Watcher, watchlist Watchlist, cfg *Config) error {
	watchlist[dir] = true
	if err := watcher.Add(dir); err != nil {
		return err
	}
	var recorder filepath.WalkFunc = func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			if ignored, err := isIgnored(path, cfg); err != nil {
				return err
			} else if !ignored {
				watchlist[path] = true
				if err := watcher.Add(path); err != nil {
					return err
				}
			}
		}
		return nil
	}
	err := filepath.Walk(dir, recorder)
	return err
}

func createAction(ev fsnotify.Event, pipeline []string, cleanup func()) func() {
	return func() {
		for i, commandStr := range pipeline {
			log.Printf(">> (%d) '%s'\n", i, commandStr)
			elements := strings.Split(commandStr, " ")
			commandStr, elements = elements[0], elements[1:]
			cmd := exec.Command(commandStr, elements...)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			if err := cmd.Run(); err != nil {
				log.Printf("<< (%d) ERROR!!! \n", i)
				break
			}
			log.Printf("<< (%d) done\n", i)
		}
		cleanup()
	}
}

func isIgnored(f string, cfg *Config) (bool, error) {
	for _, exclude := range cfg.Excludes {
		if ignore, err := filepath.Match(exclude, f); err != nil {
			return false, err
		} else if ignore {
			//log.Printf("ignoring %s because of %s", f, exclude)
			return true, nil
		}
	}
	//log.Printf("not ignoring %s for %v", f, cfg.Excludes)
	return false, nil
}

func parseConfig(cfgFile string) (*Config, error) {
	var cfg Config
	source, err := ioutil.ReadFile(cfgFile)
	if err != nil {
		return nil, err
	}
	if err := yaml.Unmarshal(source, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

// Config is the root config object
type Config struct {
	Excludes []string `yaml:"excludes"`
	Pipeline []string `yaml:"pipeline"`
}

type Watchlist map[string]bool
