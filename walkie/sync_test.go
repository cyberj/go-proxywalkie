package walkie

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

// Test Directory copy
func TestCopydir(t *testing.T) {
	require := require.New(t)

	var err error

	require.NoError(clean())
	defer clean()

	// logrus.SetLevel(logrus.DebugLevel)

	parent_dir := getTestAssetsDir()
	testdir := getTestDir()
	synced_dir := filepath.Join(parent_dir, "synced_dir")

	woriginal, err := NewWalkie(testdir)
	require.NoError(err)
	require.NoError(woriginal.Explore())

	wresult, err := NewWalkie(synced_dir)
	require.NoError(err)
	require.NoError(wresult.Explore())

	require.Len(wresult.Directory.Directories, 0)
	toadd, toremove := wresult.Directory.DiffDir(*woriginal.Directory)
	require.Len(toadd, 4)
	require.Len(toremove, 0)

	sourcedir := *woriginal.Directory

	nb, err := wresult.CopyDir(sourcedir)
	require.NoError(err)
	require.Equal(4, nb)

	// Reexplore
	require.Len(wresult.Directory.Directories, 2)

	// Re-check diff again
	toadd, toremove = woriginal.Directory.DiffDir(*wresult.Directory)
	require.Len(toadd, 0)
	require.Len(toremove, 0)

	require.NoError(os.MkdirAll(filepath.Join(synced_dir, "useless_dir"), 0755))
	// Re-check diff again
	require.NoError(wresult.Explore())
	toadd, toremove = wresult.Directory.DiffDir(*woriginal.Directory)
	require.Len(toadd, 0)
	require.Len(toremove, 1)
}

// Test Directory erase
func TestDeletedir(t *testing.T) {
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

	// require.NoError(woriginal.Directory.CopyDir(synced_dir))
	require.NoError(os.MkdirAll(filepath.Join(synced_dir, "useless_dir"), 0755))
	require.NoError(wresult.Explore())

	sourcedir := *woriginal.Directory

	nb, err := wresult.CleanDir(sourcedir)
	require.NoError(err)
	require.Equal(1, nb)

	// Re-check diff again
	require.Len(wresult.Directory.Files, 0)
	require.Len(wresult.Directory.Directories, 0)
	_, toremove := wresult.Directory.DiffDir(*woriginal.Directory)
	require.Len(toremove, 0)
}

// Test File deletion
func TestDeleteFile(t *testing.T) {
	require := require.New(t)

	// logrus.SetLevel(logrus.DebugLevel)

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

	// require.NoError(woriginal.Directory.CopyDir(synced_dir))
	// require.NoError(os.MkdirAll(filepath.Join(synced_dir, "useless_dir"), 0755))
	_, err = os.Create(filepath.Join(synced_dir, "useless_file"))
	require.NoError(err)
	require.NoError(wresult.Explore())

	sourcedir := *woriginal.Directory

	nb, err := wresult.CleanFiles(sourcedir)
	require.NoError(err)
	require.Equal(1, nb)

	// Re-check diff again
	require.Len(wresult.Directory.Files, 0)
	require.Len(wresult.Directory.Directories, 0)
	_, toremove := wresult.Directory.DiffFiles(*woriginal.Directory)
	require.Len(toremove, 0)
}
