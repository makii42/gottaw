package watch

import (
	"os"
	"path/filepath"
	"time"

	"github.com/makii42/gottaw/config"
	"github.com/makii42/gottaw/daemon"
	"github.com/makii42/gottaw/pipeline"
	"gopkg.in/fsnotify.v1"
	"gopkg.in/urfave/cli.v1"
	"github.com/makii42/gottaw/output"
	"fmt"
)

var WatchCmd = cli.Command{
	Name:   "watch",
	Usage:  "starts watching folder(s)",
	Action: watchIt,
	Flags:  []cli.Flag{},
}

var log output.Logger

// WatchIt does the work
func watchIt(c *cli.Context) error {
	cfg := config.Load()
	_log, err := output.NewLog(cfg)
	log = _log
	if err != nil {
		return err
	}
	delay, err := time.ParseDuration(c.GlobalString("delay"))
	if err != nil {
		return err
	}
	tracker := NewTracker(cfg)
	defer tracker.Close()

	trackingRoot, err := filepath.Abs(c.String("folder"))
	if err != nil {
		log.Errorf("🚨  problem with your folder: '%s'", c.String("folder"))
		return err
	}
	if _, err := os.Stat(trackingRoot); err != nil {
		log.Errorf("🚨  problem with your folder: '%s'", c.String("folder"))
		return fmt.Errorf("error while accessing tracking root %s", trackingRoot)
	}
	var serverd daemon.Daemon

	done := make(chan bool)
	go func() {
		var timer *time.Timer
		for {
			select {
			case ev := <-tracker.Events():
				if ev.Op&fsnotify.Chmod == fsnotify.Chmod || isIgnored(ev.Name, cfg) { // couldn't care less
					continue
				}
				if ev.Op&fsnotify.Create == fsnotify.Create {
					if newFile, err := os.Stat(ev.Name); err == nil {
						if newFile.IsDir() {
							tracker.Add(ev.Name)
							log.Triggerf(
								"🔭  added '%s', now watching %d folders\n",
								ev.Name,
								len(tracker.Tracked()),
							)
						}
					}
				} else if ev.Op&fsnotify.Remove == fsnotify.Remove && tracker.IsTracked(ev.Name) {
					tracker.Remove(ev.Name)
					log.Triggerf(
						"🔭  removed '%s', now watching %d folders\n",
						ev.Name,
						len(tracker.Tracked()),
					)
				} else if ev.Op&fsnotify.Write == fsnotify.Write && cfg.GetConfigFile() == ev.Name {
					cfg.Reload()
					log.Triggerf("🛠  reloaded config '%s'\n", cfg.GetConfigFile())
				} else if timer == nil {
					log.Triggerf("🔎  change detected: %s\n", ev.Name)
				}

				if timer != nil {
					log.Triggerf("🔎  even more changes detected: %s\n", ev.Name)
					timer.Reset(delay)
				} else {
					pl := pipeline.NewPipeline(func() {
						if serverd != nil {
							if err := serverd.Stop(); err != nil {
								panic(err)
							}
						}
					}, log, cfg.Pipeline, func() {
						timer = nil
						if serverd != nil {
							serverd.Start()
						}
					})
					timer = time.AfterFunc(delay, pl.Executor())
				}

			case err := <-tracker.Errors():
				log.Errorf("error: %v\n", err)
			}
		}
	}()

	if err := watchDirRecursive(trackingRoot, tracker, cfg); err != nil {
		return err
	}
	log.Noticef("🔭  watching %d folder(s). %s\n", len(tracker.Tracked()), tracker.Tracked())
	if cfg.Server != "" {
		serverd = daemon.NewDaemon(cfg.Server)
	}
	pl := pipeline.NewPipeline(nil, log, cfg.Pipeline, func() {
		if serverd != nil {
			if err := serverd.Start(); err != nil {
				panic(err)
			}
		}
	})
	pl.Executor()()
	<-done
	return nil
}

func watchDirRecursive(dir string, t Tracker, cfg *config.Config) error {
	var recorder filepath.WalkFunc = func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			if isIgnored(path, cfg) {
				return filepath.SkipDir
			}
			t.Add(path)
		}
		return nil
	}
	err := filepath.Walk(dir, recorder)
	return err
}

func isIgnored(f string, cfg *config.Config) bool {
	f, err := filepath.Abs(f)
	if err != nil {
		log.Errorf("🚨  Please check your excludes in your config: '%s'", f)
		panic(err)
	}
	for _, exclude := range cfg.Excludes {
		wd := "."
		if cfg.WorkingDirectory != "" {
			wd = cfg.WorkingDirectory
		}
		ude, err := filepath.Abs(filepath.Join(wd, exclude))
		if err != nil {
			log.Errorf("🚨  Please check your excludes in your config: '%s'", ude)
			panic(err)
		}
		if ignore, err := filepath.Match(ude, f); err != nil {
			log.Errorf("🚨  Please check your excludes in your config: '%s'", ude)
			panic(err)
		} else if ignore {
			return true
		}
	}
	return false
}
