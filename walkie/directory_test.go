package walkie

import (
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// Test Directory walking
func TestCompare(t *testing.T) {
	require := require.New(t)
	testdir := getTestDir()

	w, err := NewWalkie(testdir)
	require.NoError(err)
	err = w.Explore()
	require.NoError(err)

	dirlist := w.Directory.getSubdirs()

	require.Len(dirlist, 4)
	// dirlist = w.Directory.getSubdirsOnly()
	// require.Len(dirlist, 4)

	// data, err := json.MarshalIndent(dirlist, "", "  ")
	// require.NoError(err)
	//
	// fmt.Printf("%s", data)
	// t.Fail()

	f1 := w.Directory.Directories["folder1"].Files["file_1a"]
	f2 := w.Directory.Directories["folder1"].Files["file_1b"]

	require.True(f1.Equals(*f1))
	require.False(f1.Equals(*f2))

	d1 := w.Directory.Directories["folder1"]
	d2 := w.Directory.Directories["folder2"]

	// Shallow test
	require.True(d1.Equals(*d1))
	require.False(d1.Equals(*d2))

	// Do a copy
	d3 := w.Directory.Directories["folder1"].copy()
	// d3.Files["file_1a"] = d3.Files["file_1a"]
	require.True(d1.DeepEquals(*d3))
	require.True(d3.DeepEquals(*d1))

	// Change a bit then test
	d3.Files["file_1a"] = d3.Files["file_1a"].copy()
	d3.Files["file_1a"].Mtime = time.Now()
	//
	require.NotEqual(d3.Files["file_1a"].Mtime, d1.Files["file_1a"].Mtime)
	require.False(d1.DeepEquals(*d3))
	require.False(d3.DeepEquals(*d1))

	d4 := w.Directory.Directories["folder2"].copy()
	require.True(d2.DeepEquals(*d4))
	require.True(d4.DeepEquals(*d2))
	d4.Directories["folder_22"].Files["file_22b"].SHA256 = "Hello"
	require.False(d2.DeepEquals(*d4))
	require.False(d4.DeepEquals(*d2))

}

// Test file list
func TestGetSubfiles(t *testing.T) {
	require := require.New(t)

	var err error

	testdir := getTestDir()

	woriginal, err := NewWalkie(testdir)
	require.NoError(err)
	require.NoError(woriginal.Explore())

	files := woriginal.Directory.getSubfiles()
	require.Len(files, 11)

	files = woriginal.Directory.ListFiles()
	require.Len(files, 11)

	keys := []string{}
	for k := range files {
		keys = append(keys, k)
	}
	require.Contains(keys, filepath.Join("folder1", "file_1a"))
	require.Contains(keys, filepath.Join("folder2", "folder_21", "file_21b"))

	dirnb, filenb := woriginal.Stat()
	require.Equal(4, dirnb)
	require.Equal(11, filenb)

	files_pub := woriginal.ListFiles()
	require.Len(files_pub, 11)

}

// Test file list
func TestDiffFiles(t *testing.T) {
	require := require.New(t)

	var err error
	require.NoError(clean())
	defer clean()

	parent_dir := getTestAssetsDir()
	testdir := getTestDir()
	synced_dir := filepath.Join(parent_dir, "synced_dir")

	woriginal, err := NewWalkie(testdir)
	require.NoError(err)
	require.NoError(woriginal.Explore())

	wresult, err := NewWalkie(synced_dir)
	require.NoError(err)
	require.NoError(wresult.Explore())

	toadd, toremove := woriginal.Directory.DiffFiles(*woriginal.Directory)
	require.Len(toadd, 0)
	require.Len(toremove, 0)

	toadd, toremove = wresult.Directory.DiffFiles(*woriginal.Directory)
	require.Len(toadd, 11)
	require.Len(toremove, 0)

}
