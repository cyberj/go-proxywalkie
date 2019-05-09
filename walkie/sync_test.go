package walkie

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/Sirupsen/logrus"
	"github.com/cyberj/go-proxywalkie/testutils"
	"github.com/stretchr/testify/require"
)

// Test Directory copy
func TestCopydir(t *testing.T) {
	require := require.New(t)

	var err error

	testdirs, err := testutils.NewTestDir()
	require.NoError(err)
	defer testdirs.Clean()

	woriginal, err := NewWalkie(testdirs.TestassetsDir)
	require.NoError(err)
	require.NoError(woriginal.Explore())

	wresult, err := NewWalkie(testdirs.SyncedDir)
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

	require.NoError(os.MkdirAll(filepath.Join(testdirs.SyncedDir, "useless_dir"), 0755))
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

	testdirs, err := testutils.NewTestDir()
	require.NoError(err)
	defer testdirs.Clean()

	woriginal, err := NewWalkie(testdirs.TestassetsDir)
	require.NoError(err)
	require.NoError(woriginal.Explore())

	wresult, err := NewWalkie(testdirs.SyncedDir)
	require.NoError(err)
	require.NoError(wresult.Explore())

	// require.NoError(woriginal.Directory.CopyDir(synced_dir))
	require.NoError(os.MkdirAll(filepath.Join(testdirs.SyncedDir, "useless_dir"), 0755))
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

	logrus.SetLevel(logrus.DebugLevel)
	logrus.Debug("Start test")
	var err error

	testdirs, err := testutils.NewTestDir()
	require.NoError(err)
	defer testdirs.Clean()

	woriginal, err := NewWalkie(testdirs.TestassetsDir)
	require.NoError(err)
	require.NoError(woriginal.Explore())

	wresult, err := NewWalkie(testdirs.SyncedDir)
	require.NoError(err)
	require.NoError(wresult.Explore())

	// require.NoError(woriginal.Directory.CopyDir(synced_dir))
	// require.NoError(os.MkdirAll(filepath.Join(synced_dir, "useless_dir"), 0755))
	f, err := os.Create(filepath.Join(testdirs.SyncedDir, "useless_file"))
	f.Close()
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

// Test Directory Sync
func TestSyncdir(t *testing.T) {
	require := require.New(t)

	var err error

	testdirs, err := testutils.NewTestDir()
	require.NoError(err)
	defer testdirs.Clean()

	require.NoError(os.MkdirAll(filepath.Join(testdirs.SyncedDir, "useless_dir"), 0755))

	woriginal, err := NewWalkie(testdirs.TestassetsDir)
	require.NoError(err)
	require.NoError(woriginal.Explore())

	wresult, err := NewWalkie(testdirs.SyncedDir)
	require.NoError(err)
	require.NoError(wresult.Explore())

	require.Len(wresult.Directory.Directories, 1)
	toadd, toremove := wresult.Directory.DiffDir(*woriginal.Directory)
	require.Len(toadd, 4)
	require.Len(toremove, 1)

	sourcedir := *woriginal.Directory

	add, del, err := wresult.SyncDir(sourcedir)
	require.NoError(err)
	require.Equal(4, add)
	require.Equal(1, del)

	// Reexplore
	require.Len(wresult.Directory.Directories, 2)

	// Re-check diff again
	toadd, toremove = woriginal.Directory.DiffDir(*wresult.Directory)
	require.Len(toadd, 0)
	require.Len(toremove, 0)

}

// Test Directory Sync
func TestCreateFile(t *testing.T) {
	require := require.New(t)

	var err error

	testdirs, err := testutils.NewTestDir()
	require.NoError(err)
	defer testdirs.Clean()

	woriginal, err := NewWalkie(testdirs.TestassetsDir)
	require.NoError(err)
	require.NoError(woriginal.Explore())

	wresult, err := NewWalkie(testdirs.SyncedDir)
	require.NoError(err)
	require.NoError(wresult.Explore())

	// Sync dirs
	sourcedir := *woriginal.Directory
	add, del, err := wresult.SyncDir(sourcedir)
	require.NoError(err)
	require.Equal(4, add)
	require.Equal(0, del)

	originalFile := woriginal.Directory.Directories["folder2"].Directories["folder_22"].Files["file_22b"]
	require.NotNil(originalFile)

	f, err := os.Open(filepath.Join(testdirs.TestassetsDir, "folder2", "folder_22", "file_22b"))
	require.NoError(err)
	defer f.Close()
	require.NoError(wresult.UpdateOrCreateFile("folder2/folder_22/file_22b", f, *originalFile))

	require.Len(wresult.files, 1)
	newFile := wresult.Directory.Directories["folder2"].Directories["folder_22"].Files["file_22b"]
	require.NotNil(newFile)

	require.Equal("04337b307b9fe41137554ae2b1fddf1f0c6eb344fcfded725d971037d97311d4", newFile.SHA256)
	require.Equal(originalFile.SHA256, newFile.SHA256)
	require.Equal(originalFile.Mtime, newFile.Mtime)
	require.True(newFile.Equals(*originalFile))

	require.NoError(wresult.Explore())

	require.Len(wresult.files, 1)
	require.Equal("04337b307b9fe41137554ae2b1fddf1f0c6eb344fcfded725d971037d97311d4", wresult.Directory.Directories["folder2"].Directories["folder_22"].Files["file_22b"].SHA256)
	newFile = wresult.Directory.Directories["folder2"].Directories["folder_22"].Files["file_22b"]
	require.True(newFile.Equals(*originalFile))
}
