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

// Diff for 2 directories structure
func (d Directory) DiffDir(ref Directory) (toadd, toremove []string) {

	my_dirs := d.getSubdirs()
	ref_dirs := ref.getSubdirs()

	// copy(toadd, my_dirs)
	// copy(toremove, ref_dirs)
	var found bool

	// Check for missing dirs
	for k := range ref_dirs {
		found = false
		for t := range my_dirs {
			if k == t {
				found = true
				break
			}
		}
		if !found {
			toadd = append(toadd, k)
		}
	}

	// Check for useless dirs
	for k := range my_dirs {
		found = false
		for t := range ref_dirs {
			if k == t {
				found = true
				break
			}
		}
		if !found {
			toremove = append(toremove, k)
		}
	}

	return
}

// Get Subdirectories list
func (d Directory) getSubfiles() (files map[string]*File) {

	files = map[string]*File{}

	for k, v := range d.Files {
		files[k] = v
	}

	for _, dir := range d.Directories {

		dirfiles := dir.getSubfiles()

		for k, v := range dirfiles {
			files[filepath.ToSlash(filepath.Join(dir.Name, filepath.FromSlash(k)))] = v
		}
	}

	return
}

// Get all files list
func (d Directory) ListFiles() (files map[string]*File) {
	return d.getSubfiles()
}

// Get Subdirectories list
func (d *Directory) getSubdirs() (directories map[string]*Directory) {

	directories = map[string]*Directory{}

	for _, v := range d.Directories {
		directories[v.Name] = v

		for subdirpath, subdir := range v.getSubdirs() {
			directories[filepath.ToSlash(filepath.Join(v.Name, filepath.FromSlash(subdirpath)))] = subdir
		}
	}

	return
}

// Diff for 2 directories structure : files only
func (d Directory) DiffFiles(ref Directory) (toadd, toremove []string) {

	my_files := d.getSubfiles()
	ref_files := ref.getSubfiles()

	// copy(toadd, my_dirs)
	// copy(toremove, ref_dirs)
	var equal bool

	// Check for missing or incorrect files
	for k, v := range ref_files {
		equal = false
		for _, t := range my_files {
			if v == t {

				equal = v.Equals(*t)
				break
			}
		}

		// Not found : we add
		if !equal {
			toadd = append(toadd, k)
		}
	}

	var found bool
	// Check for useless files
	for k, v := range my_files {
		found = false
		for _, t := range ref_files {
			if v == t {
				found = true
				break
			}
		}
		if !found {
			toremove = append(toremove, k)
		}
	}

	return
}
