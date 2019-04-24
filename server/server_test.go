package proxy

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"time"

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
func TestServer(t *testing.T) {
	require := require.New(t)

	// logrus.SetLevel(logrus.DebugLevel)
	var err error

	require.NoError(clean())
	defer clean()

	// parent_dir := getTestAssetsDir()
	testdir := getTestDir()
	// synced_dir := filepath.Join(parent_dir, "synced_dir")

	woriginal, err := walkie.NewWalkie(testdir)
	require.NoError(err)
	require.NoError(woriginal.Explore())

	server, err := NewServer(testdir)
	require.NoError(err)

	ts := httptest.NewServer(server)
	defer ts.Close()

	require.NotNil(server.dircache)
	require.NotEqual(0, len(server.dircache))

	// Get Index
	time.Sleep(1 * time.Second)
	res, err := http.Get(ts.URL)
	require.NoError(err)
	require.Equal(200, res.StatusCode)
	exporteddir := &walkie.Directory{}

	require.NoError(json.NewDecoder(res.Body).Decode(exporteddir))
	res.Body.Close()
	require.True(woriginal.Directory.DeepEquals(*exporteddir))

	// Test without gzip compression
	tr := &http.Transport{
		DisableCompression: true,
	}
	client := &http.Client{Transport: tr}
	req, err := http.NewRequest("GET", ts.URL, nil)
	require.NoError(err)
	res, err = client.Do(req)
	require.NoError(err)
	exporteddir = &walkie.Directory{}
	require.NoError(json.NewDecoder(res.Body).Decode(exporteddir))
	res.Body.Close()
	require.True(woriginal.Directory.DeepEquals(*exporteddir))

}
