package walkie

import (
	"github.com/Sirupsen/logrus"
	"github.com/fsnotify/fsnotify"

	"path/filepath"
)

// init_notify (re) create a watcher
func (w *Walkie) notify_init() error {

	if w.watcher != nil {
		w.watcher.Close()
	}

	watcher, err := fsnotify.NewWatcher()

	w.watcher = watcher

	if err != nil {
		return err
	}

	go w.notify_loop()

	return nil
}

// Start Watcher loop
func (w *Walkie) notify_loop() {

	// done := make(chan bool)

	for {
		// logrus.Debug("Notify loop")
		select {
		case event, ok := <-w.watcher.Events:
			if !ok {
				return
			}
			logrus.Debugf("event:%s", event)
			if event.Op&fsnotify.Write == fsnotify.Write {
				logrus.Debugf("modified file:%s", event.Name)
			}
		case err, ok := <-w.watcher.Errors:
			if !ok {
				return
			}
			logrus.Errorf("error:%s", err)
		}
	}

	// <-done
}

// Start Watcher loop
func (w *Walkie) notify_all() {

	// done := make(chan bool)
	err := w.watcher.Add(w.path)
	if err != nil {
		logrus.Errorf("Can't watch self directory %s", w.path)
	}

	for path := range w.directories {
		err = w.watcher.Add(filepath.Join(w.path, path))
		if err != nil {
			logrus.Errorf("Can't watch directory %s", filepath.Join(w.path, path))
		}
	}
	for path := range w.files {
		err = w.watcher.Add(filepath.Join(w.path, path))
		if err != nil {
			logrus.Errorf("Can't watch file %s", filepath.Join(w.path, path))
		}
	}

	// <-done
}

// Start Watcher loop
func (w *Walkie) Watch() {

	w.notify_all()
}
