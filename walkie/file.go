package walkie

import (
	"hash"
	"os"
	"time"
)

// A File in the path
type File struct {
	Name   string    `json:"name"`
	Mtime  time.Time `json:"mtime"`
	Size   int64     `json:"size"`
	SHA256 hash.Hash `json:"sha256"`
	file   *os.FileInfo
}

// A Directory with files
type Directory struct {
	Name        string    `json:"name"`
	Mtime       time.Time `json:"mtime"`
	file        *os.FileInfo
	Files       []*File      `json:"files"`
	Directories []*Directory `json:"directories"`
}
