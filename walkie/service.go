package walkie

import (
	"os"
	"path/filepath"

	"github.com/Sirupsen/logrus"
)

// Walkie service
// We expect Directory to always be synced
type Walkie struct {
	path string

	Directory *Directory
}

// NewWalkie create new Walkie service
func NewWalkie(path string) (walkie *Walkie, err error) {
	walkie = &Walkie{
		path: path,
	}

	_, err = os.Stat(path)
	if err != nil {
		logrus.Errorf("Error opening path %s : %s", path, err)
		return
	}

	return
}

// Explore Directory to make filemap
func (w *Walkie) Explore() (err error) {

	dirlist := map[string]*Directory{}

	err = filepath.Walk(w.path, func(path string, info os.FileInfo, werr error) error {
		// logrus.Info(path)

		// realpath := strings.TrimPrefix(path, w.path)
		// logrus.Info(realpath)

		if info.IsDir() {
			d, err2 := NewDirectory(info)
			if err2 != nil {
				logrus.Debugf("NewDirectory err on path %q: %v\n", path, err2)
				return err
			}

			dirlist[path] = d
			if path != w.path {
				dirlist[filepath.Dir(path)].Directories[d.Name] = d
			}

			// data, _ := json.Marshal(d)
			// fmt.Printf("%s", data)

		} else {
			f, err2 := NewFile(path, info)
			if err2 != nil {
				logrus.Debugf("NewFile err on path %q: %v\n", path, err2)
				return err
			}

			dirlist[filepath.Dir(path)].Files[f.Name] = f
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

	w.Directory = dirlist[w.path]
	// data2, _ := json.Marshal(w.Directory)
	// fmt.Printf("%s", data2)

	return
}
