package watch

import (
	c "github.com/makii42/gottaw/config"
	"gopkg.in/fsnotify.v1"
)

type Tracker interface {
	Tracked() []string
	Add(p string) error
	IsTracked(p string) bool
	Remove(path string)
	Events() chan fsnotify.Event
	Errors() chan error
	Close() error
}

func NewTracker(cfg *c.Config) Tracker {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		panic(err)
	}
	return &FSNotifyTracker{
		tracked: make(map[string]bool),
		watcher: watcher,
		cfg:     cfg,
	}
}

type FSNotifyTracker struct {
	tracked map[string]bool
	watcher *fsnotify.Watcher
	cfg     *c.Config
}

func (t *FSNotifyTracker) Tracked() []string {
	var keys []string
	for k := range t.tracked {
		keys = append(keys, k)
	}
	return keys
}

func (t *FSNotifyTracker) Add(path string) error {
	t.tracked[path] = true
	if err := t.watcher.Add(path); err != nil {
		return err
	}
	return nil
}

func (t *FSNotifyTracker) IsTracked(path string) bool {
	_, ok := t.tracked[path]
	return ok
}

func (t *FSNotifyTracker) Remove(path string) {
	if _, ok := t.tracked[path]; ok {
		t.watcher.Remove(path)
		delete(t.tracked, path)
	}
}

func (t *FSNotifyTracker) Events() chan fsnotify.Event {
	return t.watcher.Events
}

func (t *FSNotifyTracker) Errors() chan error {
	return t.watcher.Errors
}

func (t *FSNotifyTracker) Close() error {
	return t.watcher.Close()
}
