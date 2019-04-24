package proxy

import (
	"encoding/json"
	"io/ioutil"
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

// Test Server info
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
	client_nogzip := &http.Client{Transport: tr}
	req, err := http.NewRequest("GET", ts.URL, nil)
	require.NoError(err)
	res, err = client_nogzip.Do(req)
	require.NoError(err)
	exporteddir = &walkie.Directory{}
	require.NoError(json.NewDecoder(res.Body).Decode(exporteddir))
	res.Body.Close()
	require.True(woriginal.Directory.DeepEquals(*exporteddir))

	// Test Head file
	client := &http.Client{}
	req, err = http.NewRequest("HEAD", ts.URL+"/folder1/file_xx", nil)
	require.NoError(err)
	res, err = client.Do(req)
	require.Equal(404, res.StatusCode)
	req, err = http.NewRequest("HEAD", ts.URL+"/folder1/file_1a", nil)
	require.NoError(err)
	res, err = client.Do(req)
	require.NoError(err)
	require.Equal(200, res.StatusCode)
	require.Equal("05eae7dd459fc32142c65246877d9625f51bcec8a48e79432936227637d170af", res.Header.Get("X-ProxyWalkie-Hash"))
	require.Equal("8", res.Header.Get("X-ProxyWalkie-Size"))
	require.Equal("2019-04-23 10:53:19.475828971 +0200 CEST", res.Header.Get("X-ProxyWalkie-Mtime"))

	// Test get file
	req, err = http.NewRequest("GET", ts.URL+"/folder1/file_xx", nil)
	require.NoError(err)
	res, err = client.Do(req)
	require.Equal(404, res.StatusCode)
	req, err = http.NewRequest("GET", ts.URL+"/folder1/file_1a", nil)
	require.NoError(err)
	res, err = client.Do(req)
	require.NoError(err)
	require.Equal(200, res.StatusCode)
	require.Equal("05eae7dd459fc32142c65246877d9625f51bcec8a48e79432936227637d170af", res.Header.Get("X-ProxyWalkie-Hash"))
	require.Equal("8", res.Header.Get("X-ProxyWalkie-Size"))
	require.Equal("2019-04-23 10:53:19.475828971 +0200 CEST", res.Header.Get("X-ProxyWalkie-Mtime"))

	b, err := ioutil.ReadAll(res.Body)
	require.NoError(err)
	res.Body.Close()
	require.Equal("file_1a\n", string(b))

}
