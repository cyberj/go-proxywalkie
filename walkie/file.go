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

	err := f.Compare(x)
	if err != nil {
		return false
	}

	return true
}

// Compare check if File f == file x and return an error
func (f File) Compare(x File) (err error) {

	if f.Name == x.Name && f.Mtime.Round(time.Second).Equal(x.Mtime.Round(time.Second)) && f.Size == x.Size && f.SHA256 == x.SHA256 {
		// if f.Name == x.Name && f.Mtime.Equal(x.Mtime) && f.Size == x.Size && f.SHA256 == x.SHA256 {
		return
	}

	return FileCompareError{origin: f, other: x}
}

// copy fileinfo
func (f File) copy() (newfile *File) {

	file := File(f)
	newfile = &file

	return
}
