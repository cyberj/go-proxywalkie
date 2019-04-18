package walkie

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

// Test Directory walking
func TestSyncDirectories(t *testing.T) {

	require := require.New(t)

	require.NoError(clean())
	defer clean()

	parent_dir := getTestAssetsDir()
	testdir := getTestDir()
	synced_dir := filepath.Join(parent_dir, "synced_dir")

	woriginal, err := NewWalkie(testdir)
	require.NoError(err)
	err = woriginal.Explore()
	require.NoError(err)

	wresult, err := NewWalkie(synced_dir)
	require.NoError(err)
	err = wresult.Explore()
	require.NoError(err)

	// require.False(woriginal.Directory.Equals(*wresult.Directory))
	//
	//
	//
	// require.True(woriginal.Directory.Equals(*wresult.Directory))
}
