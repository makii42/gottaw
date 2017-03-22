package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"strings"

	c "github.com/fatih/color"
	"github.com/fsnotify/fsnotify"
	"github.com/urfave/cli"
	"gopkg.in/yaml.v2"
)

var success, notices, triggers, errors *c.Color

func init() {
	success = c.New(c.FgGreen)
	notices = c.New(c.FgBlue)
	triggers = c.New(c.FgYellow)
	errors = c.New(c.FgHiRed)
}

func main() {
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
			Value: "1500ms",
			Usage: "Delay of the pipeline action after event",
		},
	}

	app.Run(os.Args)
}

// WatchIt does the work
func WatchIt(c *cli.Context) error {
	var cfg *Config

	configFile := c.String("config")
	if parsed, err := parseConfig(configFile); err != nil {
		panic(err)
	} else if parsed == nil {
		panic(fmt.Errorf("🚨  parsed config is empty: '%s'", configFile))
	} else {
		cfg = parsed
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
				triggers.Printf("🔎  change detected: %s\n", ev.Name)
				if timer == nil {
					timer = time.AfterFunc(delay, createAction(ev, cfg.Pipeline, func() {
						timer = nil
					}))
				} else if timer != nil {
					timer.Reset(delay)
				}

			case err := <-watcher.Errors:
				errors.Printf("error: %v\n", err)
			}
		}

	}()

	if f, err := os.Stat(c.String("folder")); err != nil {
		panic(err)
	} else if err := watchDirRecursive(f.Name(), watcher, watchlist, cfg); err != nil {
		panic(err)
	}
	notices.Printf("watching %d folders.\n", len(watchlist))
	<-done

	return nil
}

func watchDirRecursive(dir string, watcher *fsnotify.Watcher, watchlist Watchlist, cfg *Config) error {
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
		start := time.Now()
		for i, commandStr := range pipeline {
			elements := strings.Split(commandStr, " ")
			command, elements := elements[0], elements[1:]
			cmd := exec.Command(command, elements...)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			err := cmd.Start()
			pid := cmd.Process.Pid
			if err != nil {
				errors.Printf("🚨  (%d@%d) ERROR starting '%s': %v", i, pid, commandStr, err)
				break
			}
			notices.Printf("♻️  (%d@%d) started '%s'\n", i, pid, commandStr)
			if err := cmd.Wait(); err != nil {
				if exitErr, ok := err.(*exec.ExitError); ok {
					errors.Printf("🚨  (%d@%d) ERROR: %#v", i, pid, exitErr.ProcessState)
				} else {
					errors.Printf("🚨  (%d@%d) ERROR: %#v \n", i, pid, err)
				}
			}

			notices.Printf("♻️  (%d@%d) done\n", i, pid)
		}
		dur := time.Since(start)
		success.Printf("✅  Pipeline done after %s\n", dur.String())
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
