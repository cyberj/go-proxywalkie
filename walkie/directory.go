package walkie

import (
	"os"
	"path/filepath"
)

// A Directory with files
type Directory struct {
	Name string `json:"name"`
	// Mtime       time.Time             `json:"mtime"`
	Files       map[string]*File      `json:"files"`
	Directories map[string]*Directory `json:"directories"`
	file        os.FileInfo
}

func NewDirectory(info os.FileInfo) (d *Directory, err error) {
	d = &Directory{
		Name: info.Name(),
		// Mtime:       info.ModTime(),
		file:        info,
		Files:       map[string]*File{},
		Directories: map[string]*Directory{},
	}
	return
}

// Equals check recursively if Directory f == Directory x
func (d Directory) Equals(x Directory) bool {

	// check basic info
	if d.Name != x.Name {
		return false
	}

	return true
}

// Equals check recursively if Directory f == Directory x
func (d Directory) DeepEquals(x Directory) bool {

	if !d.Equals(x) {
		return false
	}

	// check files by len
	if len(d.Files) != len(x.Files) {
		return false
	}
	// Then one by one
	var xf *File
	var exists bool
	for path, f := range d.Files {

		xf, exists = x.Files[path]
		if !exists {
			return false
		}

		if !f.Equals(*xf) {
			return false
		}

	}
	// check directories by len
	if len(d.Directories) != len(x.Directories) {
		return false
	}
	// Then one by one
	var xdir *Directory
	for path, d := range d.Directories {

		xdir, exists = x.Directories[path]
		if !exists {
			return false
		}

		if !d.DeepEquals(*xdir) {
			return false
		}

	}

	return true
}

// Deep copy a directory
func (d Directory) copy() (newdir *Directory) {

	directory := Directory(d)
	newdir = &directory

	newdir.Directories = map[string]*Directory{}
	newdir.Files = map[string]*File{}

	for k, v := range d.Files {
		newdir.Files[k] = v.copy()
	}

	for k, v := range d.Directories {
		newdir.Directories[k] = v.copy()
	}
	return
}

// CopyDir Copy directory structure from this dir
func (d Directory) CopyDir(path string) (err error) {

	directories := d.getSubdirsOnly()

	for _, v := range directories {
		// logrus.Infof("create %s", filepath.Join(path, v))
		err = os.MkdirAll(filepath.Join(path, v), 0755)
		if err != nil {
			return
		}
	}

	return
}

// Diff for 2 directories structure
func (d Directory) DiffDir(target Directory) (toadd, toremove []string) {

	my_dirs := d.getSubdirsOnly()
	target_dirs := target.getSubdirsOnly()

	// copy(toadd, my_dirs)
	// copy(toremove, target_dirs)
	var found bool

	// Check for missing dirs
	for _, v := range my_dirs {
		found = false
		for _, t := range target_dirs {
			if v == t {
				found = true
				break
			}
		}
		if !found {
			toadd = append(toadd, v)
		}
	}

	// Check for useless dirs
	for _, v := range target_dirs {
		found = false
		for _, t := range my_dirs {
			if v == t {
				found = true
				break
			}
		}
		if !found {
			toremove = append(toremove, v)
		}
	}

	return
}

// Get Subdirectories list
func (d Directory) getSubdirs() (directories []string) {

	directories = []string{d.Name}

	for _, v := range d.Directories {
		for _, subdir := range v.getSubdirs() {
			directories = append(directories, filepath.Join(d.Name, subdir))
		}
	}

	return
}

// Get Subdirectories list without itself
func (d Directory) getSubdirsOnly() (directories []string) {

	directories = []string{}

	for _, v := range d.Directories {
		for _, subdir := range v.getSubdirs() {
			directories = append(directories, subdir)
		}
	}

	return
}
