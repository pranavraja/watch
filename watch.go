package main

import (
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/howeyc/fsnotify"
)

type rules []string

func (r rules) Match(filename string) bool {
	for _, rule := range r {
		if matched, _ := path.Match(rule, `/`+filename); matched {
			return true
		}
	}
	return false
}

type RecursiveWatcher struct {
	*fsnotify.Watcher
	ignore rules
}

func NewRecursiveWatcher(path string, ignore rules) (watcher RecursiveWatcher, err error) {
	w, err := fsnotify.NewWatcher()
	if err != nil {
		return
	}
	watcher = RecursiveWatcher{w, ignore}
	err = watcher.WatchRecursive(path)
	if err != nil {
		return
	}
	return watcher, nil
}

func (watcher RecursiveWatcher) WatchRecursive(base string) (err error) {
	err = watcher.Watch(base)
	filepath.Walk(base, func(path string, info os.FileInfo, e error) error {
		if strings.HasPrefix(path, ".") {
			return nil
		}
		if watcher.ignore.Match(path) {
			return filepath.SkipDir
		}
		if info.IsDir() {
			err = watcher.Watch(path)
		}
		return nil
	})
	return
}

func (watcher RecursiveWatcher) Handle(ev *fsnotify.FileEvent) {
	if ev.IsCreate() {
		info, err := os.Stat(ev.Name)
		if err == nil {
			if info.IsDir() && !watcher.ignore.Match(ev.Name) {
				watcher.WatchRecursive(ev.Name)
			}
		}
	} else if ev.IsDelete() {
		watcher.RemoveWatch(ev.Name)
	}
}

func (watcher RecursiveWatcher) Next() (*fsnotify.FileEvent, error) {
	select {
	case event := <-watcher.Event:
		watcher.Handle(event)
		return event, nil
	case err := <-watcher.Error:
		return nil, err
	}
}
