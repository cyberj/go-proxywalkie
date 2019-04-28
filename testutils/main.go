package testutils

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

type TestDir struct {
	ParentDir     string
	TestassetsDir string
	SyncedDir     string
}

// Create a new file testing structure
// Be sure to clean before leaving !
func NewTestDir() (td *TestDir, err error) {

	td = &TestDir{}

	// Create dir in
	td.ParentDir, err = ioutil.TempDir("", "prxwlk")
	td.createStructure()

	// create syncedDir
	td.SyncedDir = filepath.Join(td.ParentDir, "synced_dir")
	err = os.MkdirAll(td.SyncedDir, 0755)
	if err != nil {
		return
	}

	return

}

func (td *TestDir) Clean() (err error) {

	if td.ParentDir == "" {
		return
	}

	err = os.RemoveAll(td.ParentDir)
	if err != nil {
		return
	}

	return
}

// Prepare structure in parent_dir
func (td *TestDir) createStructure() (err error) {

	completetestdir := filepath.Join(td.ParentDir, "complete_test")
	td.TestassetsDir = completetestdir
	folder1 := filepath.Join(completetestdir, "folder1")
	folder2 := filepath.Join(completetestdir, "folder2")
	folder21 := filepath.Join(folder2, "folder_21")
	folder22 := filepath.Join(folder2, "folder_22")

	dirs := []string{completetestdir, folder1, folder2, folder21, folder22}

	// Create all directories
	for _, val := range dirs {
		os.MkdirAll(val, 0755)
		if err != nil {
			return
		}
	}

	files := []fdir{
		fdir{folder1, "file_1a"},
		fdir{folder1, "file_1b"},
		fdir{folder1, "file_dup"},
		fdir{folder21, "file_21a"},
		fdir{folder21, "file_21b"},
		fdir{folder21, "file_dup"},
		fdir{folder22, "file_22a"},
		fdir{folder22, "file_22b"},
		fdir{folder2, "file_2a"},
		fdir{folder2, "file_2b"},
		fdir{completetestdir, "README.md"},
	}

	for _, val := range files {
		err = val.create()
		if err != nil {
			return
		}
	}

	return

}

/**

Tests utils
*/

type fdir struct {
	directory string
	name      string
}

func (fd fdir) create() (err error) {
	f, err := os.Create(filepath.Join(fd.directory, fd.name))
	if err != nil {
		return
	}
	fmt.Fprintln(f, fd.name)
	err = f.Close()
	if err != nil {
		return
	}
	return
}
