package walkie

import (
	"os"
	"path/filepath"

	"github.com/Sirupsen/logrus"
)

// Walkie service
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
	err = filepath.Walk(w.path, func(path string, info os.FileInfo, werr error) error {

		if err != nil {
			logrus.Debugf("prevent panic by handling failure accessing a path %q: %v\n", path, err)
			return nil
		}
		return nil
	})
	return
}
