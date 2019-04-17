package walkie

import (
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"time"
)

// type Filer interface  {
//     stats() (name string, mtime, time.Time, size int64, sha256 string)
// }

// A File in the path
type File struct {
	Name   string    `json:"name"`
	Mtime  time.Time `json:"mtime"`
	Size   int64     `json:"size"`
	SHA256 string    `json:"sha256"`
	file   os.FileInfo
}

// NewFile and compute Hash
func NewFile(path string, info os.FileInfo) (f *File, err error) {

	// Open file for hashing
	osfile, err := os.Open(path)
	if err != nil {
		return
	}
	hasher := sha256.New()

	defer osfile.Close()
	if _, err2 := io.Copy(hasher, osfile); err2 != nil {
		return nil, err2
	}

	f = &File{
		Name:   info.Name(),
		Mtime:  info.ModTime(),
		Size:   info.Size(),
		SHA256: fmt.Sprintf("%x", hasher.Sum(nil)),
		file:   info,
	}
	return
}

// Equals check if File f == file x
func (f File) Equals(x File) bool {

	if f.Name == x.Name && f.Mtime == x.Mtime && f.Size == x.Size && f.SHA256 == x.SHA256 {
		return true
	}
	return false
}

// copy fileinfo
func (f File) copy() (newfile *File) {

	file := File(f)
	newfile = &file

	return
}

// A Directory with files
type Directory struct {
	Name        string                `json:"name"`
	Mtime       time.Time             `json:"mtime"`
	Files       map[string]*File      `json:"files"`
	Directories map[string]*Directory `json:"directories"`
	file        os.FileInfo
}

func NewDirectory(info os.FileInfo) (d *Directory, err error) {
	d = &Directory{
		Name:        info.Name(),
		Mtime:       info.ModTime(),
		file:        info,
		Files:       map[string]*File{},
		Directories: map[string]*Directory{},
	}
	return
}

// Equals check recursively if Directory f == Directory x
func (d Directory) Equals(x Directory) bool {

	// check basic info
	if d.Name != x.Name || d.Mtime != x.Mtime {
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
