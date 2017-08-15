package watch

import (
	c "github.com/makii42/gottaw/config"
	"gopkg.in/fsnotify.v1"
)

// Tracker keeps track of file system changes, and let's you control
// what it monitors.
type Tracker interface {
	Tracked() []string
	Add(p string) error
	IsTracked(p string) bool
	Remove(path string)
	Events() chan fsnotify.Event
	Errors() chan error
	Close() error
}

// NewTracker returns a new tracker based on the config passed.
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

// FSNotifyTracker is a tracker that uses FSNotify.
type FSNotifyTracker struct {
	tracked map[string]bool
	watcher *fsnotify.Watcher
	cfg     *c.Config
}

// Tracked returns all paths that this tracker will report changes on.
func (t *FSNotifyTracker) Tracked() []string {
	var keys []string
	for k := range t.tracked {
		keys = append(keys, k)
	}
	return keys
}

// Add adds a path to the list of tracked paths
func (t *FSNotifyTracker) Add(path string) error {
	t.tracked[path] = true
	if err := t.watcher.Add(path); err != nil {
		return err
	}
	return nil
}

// IsTracked determines whether the passed path is tracked (true) or not (false).
func (t *FSNotifyTracker) IsTracked(path string) bool {
	_, ok := t.tracked[path]
	return ok
}

// Remove removes the handed path from the tracker.
func (t *FSNotifyTracker) Remove(path string) {
	if _, ok := t.tracked[path]; ok {
		t.watcher.Remove(path)
		delete(t.tracked, path)
	}
}

// Events returns a channel that will deliver all events issued by this tracker.
func (t *FSNotifyTracker) Events() chan fsnotify.Event {
	return t.watcher.Events
}

// Errors returns a channel that will return all errors issued by this tracker.
func (t *FSNotifyTracker) Errors() chan error {
	return t.watcher.Errors
}

// Close stops all watching and event publication of this tracker.
func (t *FSNotifyTracker) Close() error {
	return t.watcher.Close()
}
