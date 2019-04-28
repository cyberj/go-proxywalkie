package walkie

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/Sirupsen/logrus"
	"github.com/cyberj/go-proxywalkie/testutils"
	"github.com/stretchr/testify/require"
)

// Test File deletion
func TestNotify(t *testing.T) {
	require := require.New(t)

	// logrus.SetLevel(logrus.DebugLevel)
	logrus.Debugf("TestNotify")

	defer logrus.Debugf("TestNotify - End")
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
	wresult.Watch()

	// Create two useless files/dir
	require.NoError(os.MkdirAll(filepath.Join(testdirs.SyncedDir, "useless_dir"), 0755))
	_, err = os.Create(filepath.Join(testdirs.SyncedDir, "useless_file"))
	require.NoError(err)

	t.SkipNow()
	// Re-check diff again
	require.Len(wresult.Directory.Files, 1)
	require.Len(wresult.Directory.Directories, 1)
	require.Len(wresult.files, 1)
	require.Len(wresult.directories, 1)
}
