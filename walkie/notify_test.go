package walkie

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/Sirupsen/logrus"
	"github.com/stretchr/testify/require"
)

func clean() (err error) {
	// require := require.New(t)
	// synced_dir := filepath.Join(getTestAssetsDir(), "synced_dir")
	parent_dir := getTestAssetsDir()
	synced_dir := filepath.Join(parent_dir, "synced_dir")

	err = os.RemoveAll(synced_dir)
	if err != nil {
		return
	}
	// os.RemoveAll(synced_dir)
	err = os.MkdirAll(synced_dir, 0755)
	if err != nil {
		return
	}

	return

}

// Test File deletion
func TestNotify(t *testing.T) {
	require := require.New(t)

	logrus.SetLevel(logrus.DebugLevel)
	logrus.Debugf("TestNotify")

	defer logrus.Debugf("TestNotify - End")
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
	wresult.Watch()

	// Create two useless files/dir
	require.NoError(os.MkdirAll(filepath.Join(synced_dir, "useless_dir"), 0755))
	_, err = os.Create(filepath.Join(synced_dir, "useless_file"))
	require.NoError(err)

	// Re-check diff again
	require.Len(wresult.Directory.Files, 1)
	require.Len(wresult.Directory.Directories, 1)
	require.Len(wresult.files, 1)
	require.Len(wresult.directories, 1)
}
