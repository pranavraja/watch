package main

import (
	"github.com/howeyc/fsnotify"
	"os"
	"path/filepath"
	"strings"
)

type RecursiveWatcher struct {
	*fsnotify.Watcher
}

func NewRecursiveWatcher(path string) (watcher RecursiveWatcher, err error) {
	w, err := fsnotify.NewWatcher()
	if err != nil {
		return
	}
	watcher = RecursiveWatcher{w}
	err = watcher.WatchRecursive(path)
	if err != nil {
		return
	}
	return watcher, nil
}

func (watcher RecursiveWatcher) WatchRecursive(path string) (err error) {
	err = watcher.Watch(path)
	filepath.Walk(path, func(path string, info os.FileInfo, e error) error {
		if strings.HasPrefix(path, ".") {
			return nil
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
			if info.IsDir() {
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
