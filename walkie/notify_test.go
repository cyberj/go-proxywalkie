package walkie

import (
	"os"
	"path/filepath"
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

	// _, err := NewWalkie(testdir)
	// require.Error(err)
	// t.Fail()
}
