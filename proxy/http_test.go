package proxy

import (
	"net/http"
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
func TestHTTPBackground(t *testing.T) {
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

	proxy, err := NewProxyParams(testdirs.SyncedDir, ts.URL, 10*time.Minute, false, false)
	require.NoError(err)

	tsproxy := httptest.NewServer(proxy.Router())
	defer tsproxy.Close()

	time.Sleep(500 * time.Millisecond)

	client := &http.Client{}
	req, err := http.NewRequest("GET", tsproxy.URL+"/folder1/file_1a", nil)
	require.NoError(err)
	res, err := client.Do(req)
	require.NoError(err)
	require.Equal(200, res.StatusCode)
	require.Equal("", res.Header.Get("X-ProxyWalkie-Cached"))

	res, err = client.Do(req)
	require.NoError(err)
	require.Equal(200, res.StatusCode)
	require.Equal("true", res.Header.Get("X-ProxyWalkie-Cached"))

}

// Test File deletion
func TestHTTPOffline(t *testing.T) {
	require := require.New(t)

	// logrus.SetLevel(logrus.DebugLevel)
	var err error

	testdirs, err := testutils.NewTestDir()
	require.NoError(err)
	defer testdirs.Clean()

	uselessfile_path := filepath.Join(testdirs.SyncedDir, "useless_file")
	f, err := os.Create(uselessfile_path)
	f.Close()

	// uselessfile_path := filepath.Join(synced_dir, "useless_file")
	// _, err = os.Create(uselessfile_path)

	proxy, err := NewProxyParams(testdirs.SyncedDir, "https://127.0.0.1:6666", 10*time.Minute, false, false)
	require.NoError(err)

	tsproxy := httptest.NewServer(proxy.Router())
	defer tsproxy.Close()

	time.Sleep(500 * time.Millisecond)

	client := &http.Client{}
	req, err := http.NewRequest("GET", tsproxy.URL+"/useless_file", nil)
	require.NoError(err)
	res, err := client.Do(req)
	require.NoError(err)
	require.Equal(200, res.StatusCode)
	require.Equal("true", res.Header.Get("X-ProxyWalkie-Cached"))
}
