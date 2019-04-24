package proxy

// Test File deletion
// func TestServer(t *testing.T) {
// 	require := require.New(t)
//
// 	logrus.SetLevel(logrus.DebugLevel)
// 	var err error
//
// 	require.NoError(clean())
// 	defer clean()
//
// 	// parent_dir := getTestAssetsDir()
// 	testdir := getTestDir()
// 	// synced_dir := filepath.Join(parent_dir, "synced_dir")
//
// 	woriginal, err := walkie.NewWalkie(testdir)
// 	require.NoError(err)
// 	require.NoError(woriginal.Explore())
//
// 	proxy, err := NewProxy(testdir)
// 	require.NoError(err)
//
// 	ts := httptest.NewServer(proxy.Router())
// 	defer ts.Close()
//
// 	res, err := http.Get(ts.URL)
// 	require.NoError(err)
//
// 	exporteddir := &walkie.Directory{}
// 	require.NoError(json.NewDecoder(res.Body).Decode(exporteddir))
//
// 	res.Body.Close()
//
// 	require.True(woriginal.Directory.DeepEquals(*exporteddir))
//
// }
