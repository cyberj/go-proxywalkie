package walkie

import (
	"os"
	"path/filepath"

	"github.com/Sirupsen/logrus"
)

// Copy directories, return number of created directories
func (w *Walkie) CopyDir(ref Directory) (nb int, err error) {

	toadd, _ := w.Directory.DiffDir(ref)

	for _, v := range toadd {
		path := filepath.Join(w.path, v)
		logrus.Debugf("create %s", path)
		err2 := os.MkdirAll(path, 0755)
		if err2 != nil {
			logrus.Errorf("CopyDir/MkdirAll : %s", err2)
			if nb != 0 {
				w.Explore()
			}
			return
		} else {
			nb++
		}
	}

	logrus.Debugf("CopyDir copied nb=%d", nb)
	if nb != 0 {
		w.Explore()
	}
	return
}

// Delete directories, return number of changed files
func (w *Walkie) CleanDir(ref Directory) (nb int, err error) {

	_, toremove := w.Directory.DiffDir(ref)

	for _, v := range toremove {
		path := filepath.Join(w.path, v)
		logrus.Debugf("Delete %s", path)
		err2 := os.RemoveAll(path)
		if err2 != nil {
			if nb != 0 {
				w.Explore()
			}
			return
		} else {
			nb++
		}
	}

	logrus.Debugf("CleanDir deleted nb=%d", nb)
	if nb != 0 {
		w.Explore()
	}
	return
}

// Synchronize directory
func (w *Walkie) SyncDir(ref Directory) (add, del int, err error) {
	add, err = w.CopyDir(ref)
	if err != nil {
		return
	}
	del, err = w.CleanDir(ref)
	if err != nil {
		return
	}

	return
}

// Delete files, return number of changed files
func (w *Walkie) CleanFiles(ref Directory) (nb int, err error) {

	_, toremove := w.Directory.DiffFiles(ref)

	for _, v := range toremove {
		path := filepath.Join(w.path, v)
		logrus.Debugf("Delete %s", path)
		err = os.RemoveAll(path)
		if err != nil {
			if nb != 0 {
				w.Explore()
			}
			logrus.Debugf("CleanFiles got error %s when trying to delete %s", err, path)
			return
		} else {
			nb++
		}
	}

	logrus.Debugf("CleanFiles deleted nb=%d", nb)
	if nb != 0 {
		w.Explore()
	}
	return
}
