package proxy

import (
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"

	server "github.com/cyberj/go-proxywalkie/server"
	"github.com/cyberj/go-proxywalkie/testutils"
	"github.com/cyberj/go-proxywalkie/walkie"
	"github.com/stretchr/testify/require"
)

// Test File deletion
func TestSync(t *testing.T) {
	require := require.New(t)

	// logrus.SetLevel(logrus.DebugLevel)
	var err error

	testdirs, err := testutils.NewTestDir()
	require.NoError(err)
	defer testdirs.Clean()

	woriginal, err := walkie.NewWalkie(testdirs.TestassetsDir)
	require.NoError(err)
	require.NoError(woriginal.Explore())

	srv, err := server.NewServer(testdirs.TestassetsDir)
	require.NoError(err)

	ts := httptest.NewServer(srv)
	defer ts.Close()

	// ts.URL

	proxy, err := NewProxy(testdirs.SyncedDir, ts.URL)
	require.NoError(err)

	// grok grok...
	for !proxy.Ready() {
	}

	// directories created
	require.Len(proxy.walkiedir.ListDirs(), 4)

	require.False(proxy.checkFile("/totototo"))
	require.False(proxy.checkFile("folder1/file_1a"))
	require.False(proxy.checkFile("folder2/folder_22/file_22b"))

	originalFile := woriginal.Directory.Directories["folder2"].Directories["folder_22"].Files["file_22b"]
	require.NotNil(originalFile)

	f, err := os.Open(filepath.Join(testdirs.TestassetsDir, "folder2", "folder_22", "file_22b"))
	require.NoError(err)
	defer f.Close()

	err = proxy.getFile("folder2/folder_22/file_22b")
	require.NoError(err)

	require.True(proxy.checkFile("folder2/folder_22/file_22b"))
}

// Test File deletion
func TestSyncClean(t *testing.T) {
	require := require.New(t)

	// logrus.SetLevel(logrus.DebugLevel)
	var err error

	testdirs, err := testutils.NewTestDir()
	require.NoError(err)
	defer testdirs.Clean()

	woriginal, err := walkie.NewWalkie(testdirs.TestassetsDir)
	require.NoError(err)
	require.NoError(woriginal.Explore())

	srv, err := server.NewServer(testdirs.TestassetsDir)
	require.NoError(err)

	ts := httptest.NewServer(srv)
	defer ts.Close()

	uselessfile_path := filepath.Join(testdirs.SyncedDir, "useless_file")
	f, err := os.Create(uselessfile_path)
	f.Close()

	proxy, err := NewProxy(testdirs.SyncedDir, ts.URL)
	require.NoError(err)

	_, err = os.Stat(uselessfile_path)
	require.NoError(err)

	// grok grok...
	proxy.Stop()

	proxy.Clean = true
	require.NoError(proxy.Run())

	// for !proxy.running {
	// }

	// File must no exists
	_, err = os.Stat(uselessfile_path)
	require.Error(err)
	require.True(os.IsNotExist(err))

}

// Test File deletion
func TestSyncBackground(t *testing.T) {
	require := require.New(t)

	// logrus.SetLevel(logrus.DebugLevel)
	var err error

	testdirs, err := testutils.NewTestDir()
	require.NoError(err)
	defer testdirs.Clean()

	woriginal, err := walkie.NewWalkie(testdirs.TestassetsDir)
	require.NoError(err)
	require.NoError(woriginal.Explore())

	srv, err := server.NewServer(testdirs.TestassetsDir)
	require.NoError(err)

	ts := httptest.NewServer(srv)
	defer ts.Close()

	// uselessfile_path := filepath.Join(synced_dir, "useless_file")
	// _, err = os.Create(uselessfile_path)

	proxy, err := NewProxyParams(testdirs.SyncedDir, ts.URL, 10*time.Minute, false, true)
	require.NoError(err)

	// Wait sync
	time.Sleep(500 * time.Millisecond)
	require.Len(proxy.walkiedir.ListFiles(), 11)
	require.Len(proxy.walkiedir.ListDirs(), 4)

}
