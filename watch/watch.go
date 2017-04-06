package watch

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/makii42/gottaw/config"
	"github.com/makii42/gottaw/daemon"
	"github.com/makii42/gottaw/output"
	"gopkg.in/fsnotify.v1"
	"gopkg.in/urfave/cli.v1"
)

var log *output.Logger
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
	log = output.NewLogger(output.TRACE, watchCfg)
	defer tracker.Close()

	trackingRoot, err := filepath.Abs(c.String("folder"))
	if err != nil {
		log.Errorf("🚨  problem with your folder: '%s'", c.String("folder"))
		panic(err)
	}
	if _, err := os.Stat(trackingRoot); err != nil {
		log.Errorf("🚨  problem with your folder: '%s'", c.String("folder"))
	}
	var serverd daemon.Daemon

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
				} else if ev.Op&fsnotify.Write == fsnotify.Write && watchCfg.File == ev.Name {
					watchCfg, _ := config.ParseConfig(watchCfg.File)
					log.Triggerf("🛠  reloaded config '%s'\n", watchCfg.File)
					go WatchIt(c)
					break WatchLoop
				} else if timer == nil {
					log.Triggerf("🔎  change detected: %s\n", ev.Name)
				}

				if timer != nil {
					log.Triggerf("🔎  even more changes detected: %s\n", ev.Name)
					timer.Reset(delay)
				} else {
					pipelineFunc := executePipeline(func() {
						if serverd != nil {
							if err := serverd.Stop(); err != nil {
								panic(err)
							}
						}
					}, watchCfg.Pipeline, func() {
						timer = nil
						if serverd != nil {
							serverd.Start()
						}
					})
					timer = time.AfterFunc(delay, pipelineFunc)
				}

			case err := <-tracker.Errors():
				log.Errorf("error: %v\n", err)
			}
		}
	}()

	if err := watchDirRecursive(trackingRoot, tracker, watchCfg); err != nil {
		panic(err)
	}
	log.Noticef("🔭  watching %d folder(s). %s\n", len(tracker.Tracked()), tracker.Tracked())
	if watchCfg.Server != "" {
		serverd = daemon.NewDaemon(log, watchCfg.Server)
	}
	executePipeline(nil, watchCfg.Pipeline, func() {
		if serverd != nil {
			if err := serverd.Start(); err != nil {
				panic(err)
			}
		}
	})()
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

func executePipeline(preProcess func(), pipeline []string, postProcess func()) func() {
	return func() {
		start := time.Now()
		if preProcess != nil {
			preProcess()
		}
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
				log.Errorf("🚨  (%d@%d) ERROR starting '%s': %v", i, pid, commandStr, err)
				return
			}
			log.Noticef("♻️  (%d@%d) started '%s'\n", i, pid, commandStr)
			if err := cmd.Wait(); err != nil {
				log.Errorf("🚨  (%d@%d) ERROR: %s \n", i, pid, err)
				return
			}

			log.Noticef("♻️  (%d@%d) done\n", i, pid)
		}
		dur := time.Since(start)
		log.Successf("✅  Pipeline done after %s\n", dur.String())
		if postProcess != nil {
			postProcess()
		}
	}
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
