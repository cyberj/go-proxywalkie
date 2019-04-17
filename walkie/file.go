package walkie

import (
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"time"
)

// A File in the path
type File struct {
	Name   string    `json:"name"`
	Mtime  time.Time `json:"mtime"`
	Size   int64     `json:"size"`
	SHA256 string    `json:"sha256"`
	file   os.FileInfo
}

//
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
