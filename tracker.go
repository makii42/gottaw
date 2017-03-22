package main

import (
	"github.com/fsnotify/fsnotify"
)

type Tracker struct {
	tracked map[string]bool
	watcher *fsnotify.Watcher
	cfg     *Config
}

func NewTracker(cfg *Config) *Tracker {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		panic(err)
	}
	return &Tracker{
		tracked: make(map[string]bool),
		watcher: watcher,
		cfg:     cfg,
	}
}

func (t *Tracker) Tracked() []string {
	var keys []string
	for k := range t.tracked {
		keys = append(keys, k)
	}
	return keys
}

func (t *Tracker) Add(path string) error {
	t.tracked[path] = true
	if err := t.watcher.Add(path); err != nil {
		return err
	}
	return nil
}

func (t *Tracker) IsTracked(path string) bool {
	_, ok := t.tracked[path]
	return ok
}

func (t *Tracker) Remove(path string) {
	if _, ok := t.tracked[path]; ok {
		t.watcher.Remove(path)
		delete(t.tracked, path)
	}
}

func (t *Tracker) Events() chan fsnotify.Event {
	return t.watcher.Events
}

func (t *Tracker) Errors() chan error {
	return t.watcher.Errors
}

func (t *Tracker) Close() error {
	return t.watcher.Close()
}
