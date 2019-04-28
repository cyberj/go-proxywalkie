package proxy

import (
	"net/http/httptest"
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"time"

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
func TestSync(t *testing.T) {
	require := require.New(t)

	// logrus.SetLevel(logrus.DebugLevel)
	var err error

	require.NoError(clean())

	defer clean()

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

	f, err := os.Open(filepath.Join(testdir, "folder2", "folder_22", "file_22b"))
	require.NoError(err)
	defer f.Close()

	err = proxy.getFile("folder2/folder_22/file_22b")
	require.NoError(err)

	require.True(proxy.checkFile("folder2/folder_22/file_22b"))

	// require.False(proxy.checkFile("folder1/file_1a"))
}

// Test File deletion
func TestSyncClean(t *testing.T) {
	require := require.New(t)

	// logrus.SetLevel(logrus.DebugLevel)
	var err error

	require.NoError(clean())

	defer clean()

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

	uselessfile_path := filepath.Join(synced_dir, "useless_file")
	_, err = os.Create(uselessfile_path)

	proxy, err := NewProxy(synced_dir, ts.URL)
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

	require.NoError(clean())

	defer clean()

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

	// uselessfile_path := filepath.Join(synced_dir, "useless_file")
	// _, err = os.Create(uselessfile_path)

	proxy, err := NewProxyParams(synced_dir, ts.URL, 10*time.Minute, false, true)
	require.NoError(err)

	// Wait sync
	time.Sleep(500 * time.Millisecond)
	require.Len(proxy.walkiedir.ListFiles(), 11)
	require.Len(proxy.walkiedir.ListDirs(), 4)

}
