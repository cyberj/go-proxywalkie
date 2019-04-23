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
	toadd, toremove := woriginal.Directory.DiffDir(*wresult.Directory)
	require.Len(toadd, 4)
	require.Len(toremove, 0)

	err = woriginal.Directory.CopyDir(synced_dir)
	require.NoError(err)

	// Reexplore
	require.NoError(wresult.Explore())
	require.Len(wresult.Directory.Directories, 2)

	// Re-check diff again
	toadd, toremove = woriginal.Directory.DiffDir(*wresult.Directory)
	require.Len(toadd, 0)
	require.Len(toremove, 0)

	require.NoError(os.MkdirAll(filepath.Join(synced_dir, "useless_dir"), 0755))
	// Re-check diff again
	require.NoError(wresult.Explore())
	toadd, toremove = woriginal.Directory.DiffDir(*wresult.Directory)
	require.Len(toadd, 0)
	require.Len(toremove, 1)
}

// Test Directory copy
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

	nb, err := woriginal.Directory.CleanDir(synced_dir, *wresult.Directory)
	require.NoError(err)
	require.Equal(1, nb)

	// Re-check diff again
	require.NoError(wresult.Explore())
	require.Len(wresult.Directory.Files, 0)
	require.Len(wresult.Directory.Directories, 0)
	_, toremove := woriginal.Directory.DiffDir(*wresult.Directory)
	require.Len(toremove, 0)
}

// Test Directory walking
func TestSyncDirectories(t *testing.T) {

	require := require.New(t)

	require.NoError(clean())
	defer clean()

	parent_dir := getTestAssetsDir()
	testdir := getTestDir()
	synced_dir := filepath.Join(parent_dir, "synced_dir")

	woriginal, err := NewWalkie(testdir)
	require.NoError(err)
	err = woriginal.Explore()
	require.NoError(err)

	wresult, err := NewWalkie(synced_dir)
	require.NoError(err)
	err = wresult.Explore()
	require.NoError(err)

	// require.False(woriginal.Directory.Equals(*wresult.Directory))
	//
	//
	//
	// require.True(woriginal.Directory.Equals(*wresult.Directory))
}
