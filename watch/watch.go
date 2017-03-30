package watch

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/makii42/gottaw/config"
	"github.com/makii42/gottaw/output"
	"github.com/urfave/cli"
)

var l *output.Logger
var watchCfg *config.Config

var WatchCmd = cli.Command{
	Name:   "watch",
	Usage:  "starts watching folder(s)",
	Action: WatchIt,
	Flags:  []cli.Flag{},
}

// WatchIt does the work
func WatchIt(c *cli.Context) error {

	var delay time.Duration
	watchCfg, delay = config.Setup(c)
	tracker := NewTracker(watchCfg)
	l = output.NewLogger(output.TRACE, watchCfg)
	defer tracker.Close()

	trackingRoot, err := filepath.Abs(c.String("folder"))
	if err != nil {
		l.Errorf("🚨  problem with your folder: '%s'", c.String("folder"))
		panic(err)
	}
	if _, err := os.Stat(trackingRoot); err != nil {
		l.Errorf("🚨  problem with your folder: '%s'", c.String("folder"))
	}

	done := make(chan bool)
	go func() {
		var timer *time.Timer
	WatchLoop:
		for {
			select {
			case ev := <-tracker.Events():
				if ev.Op&fsnotify.Chmod == fsnotify.Chmod { // couldn't care less
					continue
				}
				if isIgnored(ev.Name, watchCfg) {
					continue
				}
				if ev.Op&fsnotify.Create == fsnotify.Create {
					if newFile, err := os.Stat(ev.Name); err == nil {
						if newFile.IsDir() {
							tracker.Add(ev.Name)
							l.Triggerf(
								"🔭  added '%s', now watching %d folders\n",
								ev.Name,
								len(tracker.Tracked()),
							)
						}
					}
				} else if ev.Op&fsnotify.Remove == fsnotify.Remove && tracker.IsTracked(ev.Name) {
					tracker.Remove(ev.Name)
					l.Triggerf(
						"🔭  removed '%s', now watching %d folders\n",
						ev.Name,
						len(tracker.Tracked()),
					)
				} else if ev.Op&fsnotify.Write == fsnotify.Write && watchCfg.File == ev.Name {
					watchCfg, _ := config.ParseConfig(watchCfg.File)
					l.Triggerf("🛠  reloaded config '%s'\n", watchCfg.File)
					go WatchIt(c)
					break WatchLoop
				} else {
					l.Triggerf("🔎  change detected: %s\n", ev.Name)
				}

				if timer != nil {
					l.Triggerf("🔎  even more changes detected: %s\n", ev.Name)
					timer.Reset(delay)
				} else {
					timer = time.AfterFunc(delay, executePipeline(watchCfg.Pipeline, func() {
						timer = nil
					}))
				}

			case err := <-tracker.Errors():
				l.Errorf("error: %v\n", err)
			}
		}
	}()

	if err := watchDirRecursive(trackingRoot, tracker, watchCfg); err != nil {
		panic(err)
	}
	l.Noticef("🔭  watching %d folder(s). %s\n", len(tracker.Tracked()), tracker.Tracked())
	executePipeline(watchCfg.Pipeline, func() {})()
	<-done
	return nil
}

func watchDirRecursive(dir string, tracker *Tracker, cfg *config.Config) error {
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
			if watchCfg.WorkingDirectory != "" {
				cmd.Dir = watchCfg.WorkingDirectory
			}
			err := cmd.Start()
			pid := cmd.Process.Pid
			if err != nil {
				l.Errorf("🚨  (%d@%d) ERROR starting '%s': %v", i, pid, commandStr, err)
				return
			}
			l.Noticef("♻️  (%d@%d) started '%s'\n", i, pid, commandStr)
			if err := cmd.Wait(); err != nil {
				l.Errorf("🚨  (%d@%d) ERROR: %s \n", i, pid, err)
				return
			}

			l.Noticef("♻️  (%d@%d) done\n", i, pid)
		}
		dur := time.Since(start)
		l.Successf("✅  Pipeline done after %s\n", dur.String())
		cleanup()
	}
}

func isIgnored(f string, cfg *config.Config) bool {
	f, err := filepath.Abs(f)
	if err != nil {
		l.Errorf("🚨  Please check your excludes in your config: '%s'", f)
		panic(err)
	}
	for _, exclude := range cfg.Excludes {
		ude, err := filepath.Abs(exclude)
		if err != nil {
			l.Errorf("🚨  Please check your excludes in your config: '%s'", exclude)
			panic(err)
		}
		if ignore, err := filepath.Match(ude, f); err != nil {
			l.Errorf("🚨  Please check your excludes in your config: '%s'", exclude)
			panic(err)
		} else if ignore {
			return true
		}
	}
	return false
}
