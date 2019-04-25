package walkie

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"

	"github.com/Sirupsen/logrus"
	"github.com/fsnotify/fsnotify"
)

var Version = "X.X.X"

// Walkie service
// We expect Directory to always be synced
type Walkie struct {
	path string

	Directory *Directory

	// hashmap
	directories map[string]*Directory
	files       map[string]*File

	watcher *fsnotify.Watcher

	mu sync.RWMutex
}

// NewWalkie create new Walkie service
func NewWalkie(path string) (walkie *Walkie, err error) {
	dirPath, err := filepath.Abs(path)
	if err != nil {
		return
	}

	walkie = &Walkie{
		path: dirPath,

		Directory: &Directory{
			Name: "",
			// Mtime:       info.ModTime(),
			Files:       map[string]*File{},
			Directories: map[string]*Directory{},
		},

		directories: map[string]*Directory{},
		files:       map[string]*File{},
	}

	_, err = os.Stat(path)
	if err != nil {
		logrus.Errorf("Error opening path %s : %s", path, err)
		return
	}

	err = walkie.notify_init()
	if err != nil {
		logrus.Errorf("Error creating watcher : %s", err)
		return
	}
	return
}

// Explore Directory to make filemap
func (w *Walkie) Explore() (err error) {

	dirlist := map[string]*Directory{}

	err = filepath.Walk(w.path, func(path string, info os.FileInfo, werr error) error {
		logrus.Info(path)

		// realpath := strings.TrimPrefix(path, w.path)
		// logrus.Info(realpath)
		// logrus.Error(w.path)

		if werr != nil {
			logrus.Errorf("walk error %s %s", w.path, werr)
			return nil
		}

		if info.IsDir() {
			d, err2 := NewDirectory(info)
			if err2 != nil {
				logrus.Debugf("NewDirectory err on path %q: %v\n", path, err2)
				return err
			}

			dirlist[filepath.ToSlash(path)] = d
			if path != w.path {
				dirlist[filepath.ToSlash(filepath.Dir(path))].Directories[d.Name] = d
			}

			// data, _ := json.Marshal(d)
			// fmt.Printf("%s", data)

		} else {
			f, err2 := NewFile(path, info)
			if err2 != nil {
				logrus.Debugf("NewFile err on path %q: %v\n", path, err2)
				return err
			}

			dirlist[filepath.ToSlash(filepath.Dir(path))].Files[f.Name] = f
			// data, _ := json.Marshal(f)
			// fmt.Printf("%s", data)
		}

		if err != nil {
			logrus.Debugf("prevent panic by handling failure accessing a path %q: %v\n", path, err)
			return nil
		}
		return nil
	})

	// data2, _ := json.Marshal(dirlist)
	// fmt.Printf("%s", data2)

	w.mu.Lock()
	w.Directory = dirlist[filepath.ToSlash(w.path)]
	w.files = w.Directory.getSubfiles()
	w.directories = w.Directory.getSubdirs()
	w.mu.Unlock()
	// data2, _ := json.Marshal(w.Directory)
	// fmt.Printf("%s", data2)

	return
}

// Close Watcher
func (w *Walkie) Close() {

	w.watcher.Close()

	return
}

// Stats
func (w *Walkie) Stat() (nbdir, nbfiles int) {
	w.mu.RLock()
	nbdir = len(w.directories)
	nbfiles = len(w.files)
	w.mu.RUnlock()

	return
}

// Get files directory list
func (w *Walkie) ListFiles() (files []string) {

	for k := range w.files {
		files = append(files, k)
	}

	return
}

// Get Subdirectories list
func (w *Walkie) ListDirs() (files []string) {

	for k := range w.directories {
		files = append(files, k)
	}

	return
}

// Get a file
func (w *Walkie) GetFile(path string) (file File, found bool) {

	f, ok := w.files[path]
	if !ok {
		return
	}

	return *f, true
}

// Get a directory
func (w *Walkie) GetDir(path string) (dir Directory, found bool) {
	w.mu.RLock()

	d, ok := w.directories[path]
	if !ok {
		return
	}
	w.mu.RUnlock()

	return *d, true
}

// Create or update a file (use slashes as path)
func (w *Walkie) UpdateOrCreateFile(path string, r io.Reader, original_file File) (err error) {

	var directory *Directory
	var file *File

	w.mu.Lock()
	defer w.mu.Unlock()

	name := filepath.Base(path)
	dir := filepath.Dir(path)

	if dir == "." {
		directory = w.Directory
	} else {
		var ok bool

		directory, ok = w.directories[filepath.ToSlash(dir)]
		if !ok {
			return fmt.Errorf("UpdateOrCreateFile : Directory not found")
		}
	}

	systempath := filepath.Join(w.path, filepath.FromSlash(path))
	// _, exists := w.files[path]

	// Create or update file
	f, err := os.OpenFile(systempath, os.O_CREATE|os.O_WRONLY, 0644)
	defer f.Close()
	if err != nil {
		return
	}

	// Copy file
	_, err = io.Copy(f, r)
	if err != nil {
		return
	}

	// Close it
	if err = f.Close(); err != nil {
		return
	}

	// Change mtime
	err = os.Chtimes(systempath, original_file.Mtime, original_file.Mtime)
	if err != nil {
		return
	}

	// Get info and create file
	info, err := os.Stat(systempath)
	if err != nil {
		return
	}
	file, err = NewFile(systempath, info)
	if err != nil {
		return
	}

	// Last check
	if file.SHA256 != original_file.SHA256 {
		return fmt.Errorf("Different SHA : newfile=%s expected=%s", file.SHA256, original_file.SHA256)
	}

	directory.Files[name] = file
	w.files[path] = file

	return
}
