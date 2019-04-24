package proxy

import (
	"net/http/httptest"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	server "github.com/cyberj/go-proxywalkie/server"
	"github.com/cyberj/go-proxywalkie/walkie"
	"github.com/stretchr/testify/require"
)

func getTestDir() string {
	return filepath.Join(getTestAssetsDir(), "complete_test")
}

func getTestAssetsDir() string {
	_, filename, _, _ := runtime.Caller(0)
	testdir := filepath.Join(filepath.Dir(filename), "..", "tests_assets")
	return testdir
}

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
func TestCache(t *testing.T) {
	require := require.New(t)

	// logrus.SetLevel(logrus.DebugLevel)
	var err error

	require.NoError(clean())
	defer require.NoError(clean())

	parent_dir := getTestAssetsDir()
	testdir := getTestDir()
	synced_dir := filepath.Join(parent_dir, "synced_dir")

	woriginal, err := walkie.NewWalkie(testdir)
	require.NoError(err)
	require.NoError(woriginal.Explore())

	srv, err := server.NewServer(testdir)
	require.NoError(err)

	ts := httptest.NewServer(srv)
	defer ts.Close()

	// ts.URL

	proxy, err := NewProxy(synced_dir, ts.URL)
	require.NoError(err)

	proxy.Ready()

	// directories created
	require.Len(proxy.walkiedir.ListDirs(), 4)

	require.False(proxy.checkFile("/totototo"))
	require.False(proxy.checkFile("folder1/file_1a"))
	require.False(proxy.checkFile("folder2/folder_22/file_22b"))

	originalFile := woriginal.Directory.Directories["folder2"].Directories["folder_22"].Files["file_22b"]
	require.NotNil(originalFile)

	f, err := os.Open(filepath.Join(testdir, "folder2", "folder_22", "file_22b"))
	require.NoError(err)
	defer f.Close()

	err = proxy.getFile("folder2/folder_22/file_22b")
	require.NoError(err)

	require.True(proxy.checkFile("folder2/folder_22/file_22b"))

	// require.False(proxy.checkFile("folder1/file_1a"))
}
