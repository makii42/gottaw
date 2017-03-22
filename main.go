package main

import (
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
var cfg Config

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
	delay := setup(c)
	tracker := NewTracker(&cfg)
	defer tracker.Close()

	trackingRoot, err := filepath.Abs(c.String("folder"))
	if err != nil {
		errors.Printf("🚨  problem with your folder: '%s'", c.String("folder"))
		panic(err)
	}
	if _, err := os.Stat(trackingRoot); err != nil {
		errors.Printf("🚨  problem with your folder: '%s'", c.String("folder"))
	}

	done := make(chan bool)
	go func() {
		var timer *time.Timer
		for {
			select {
			case ev := <-tracker.Events():
				if ev.Op&fsnotify.Chmod == fsnotify.Chmod { // couldn't care less
					continue
				}
				if isIgnored(ev.Name, &cfg) {
					continue
				}
				if ev.Op&fsnotify.Create == fsnotify.Create {
					if newFile, err := os.Stat(ev.Name); err != nil {
						panic(err)
					} else if newFile.IsDir() {
						tracker.Add(ev.Name)
						triggers.Printf(
							"🔭  added '%s', now watching %d folders\n",
							ev.Name,
							len(tracker.Tracked()),
						)
					}
				} else if ev.Op&fsnotify.Remove == fsnotify.Remove {
					if tracker.IsTracked(ev.Name) {
						tracker.Remove(ev.Name)
						triggers.Printf(
							"🔭  removed '%s', now watching %d folders\n",
							ev.Name,
							len(tracker.Tracked()),
						)
					}
				} else if ev.Op&fsnotify.Write == fsnotify.Write && cfg.File == ev.Name {
					parseConfig(cfg.File)
					triggers.Printf("🛠  reloaded config '%s'\n", cfg.File)
				} else {
					triggers.Printf("🔎  change detected: %s\n", ev.Name)
				}
				if timer == nil {
					timer = time.AfterFunc(delay, executePipeline(cfg.Pipeline, func() {
						timer = nil
					}))
				} else if timer != nil {
					triggers.Printf("🔎  even more changes detected: %s\n", ev.Name)
					if !timer.Stop() {
						<-timer.C
					}
					timer.Reset(delay)
				}

			case err := <-tracker.Errors():
				errors.Printf("error: %v\n", err)
			}
		}
	}()

	if err := watchDirRecursive(trackingRoot, tracker, &cfg); err != nil {
		panic(err)
	}
	notices.Printf("🔭  watching %d folders. %#v\n", len(tracker.Tracked()), tracker.Tracked())
	executePipeline(cfg.Pipeline, func() {})()
	<-done
	return nil
}

func watchDirRecursive(dir string, tracker *Tracker, cfg *Config) error {
	var recorder filepath.WalkFunc = func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			if ignored := isIgnored(path, cfg); err != nil {
				return err
			} else if !ignored {
				tracker.Add(path)
			} else {
				return filepath.SkipDir
			}
		}
		return nil
	}
	err := filepath.Walk(dir, recorder)
	return err
}

func executePipeline(pipeline []string, cleanup func()) func() {
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
				break
			}

			notices.Printf("♻️  (%d@%d) done\n", i, pid)
		}
		dur := time.Since(start)
		success.Printf("✅  Pipeline done after %s\n", dur.String())
		cleanup()
	}
}

func isIgnored(f string, cfg *Config) bool {
	f, err := filepath.Abs(f)
	if err != nil {
		errors.Printf("🚨  Please check your excludes in your config: '%s'", f)
		panic(err)
	}
	for _, exclude := range cfg.Excludes {
		absExclude, err := filepath.Abs(exclude)
		if err != nil {
			errors.Printf("🚨  Please check your excludes in your config: '%s'", exclude)
			panic(err)
		}
		if ignore, err := filepath.Match(absExclude, f); err != nil {
			errors.Printf("🚨  Please check your excludes in your config: '%s'", exclude)
			panic(err)
		} else if ignore {
			return true
		}
	}
	return false
}

func parseConfig(cfgFile string) error {
	source, err := ioutil.ReadFile(cfgFile)
	if err != nil {
		return err
	}
	if err := yaml.Unmarshal(source, &cfg); err != nil {
		return err
	}
	cfg.File = cfgFile
	return nil
}

// Config is the root config object
type Config struct {
	File     string
	Excludes []string `yaml:"excludes"`
	Pipeline []string `yaml:"pipeline"`
}

func setup(c *cli.Context) time.Duration {
	configFile, err := filepath.Abs(c.String("config"))
	if err != nil {
		panic(err)
	}
	err = parseConfig(configFile)
	if err != nil {
		panic(err)
	}
	delay, err := time.ParseDuration(c.String("delay"))
	if err != nil {
		panic(err)
	}
	return delay
}
