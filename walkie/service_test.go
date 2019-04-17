package walkie

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/require"
)

func getTestDir() string {
	_, filename, _, _ := runtime.Caller(0)
	// fmt.Println("Current test filename: " + filename)
	// fmt.Println("Current test dir: " + filepath.Dir(filename))

	testdir := filepath.Join(filepath.Dir(filename), "..", "tests_assets", "complete_test")
	// fmt.Println("Target test dir: " + testdir)
	return testdir
}

func TestWalkingUnknownPath(t *testing.T) {
	require := require.New(t)
	testdir := "./carmen/sandiego"

	_, err := NewWalkie(testdir)
	require.Error(err)
	// t.Fail()
}

func TestWalking(t *testing.T) {
	require := require.New(t)
	testdir := getTestDir()

	w, err := NewWalkie(testdir)
	require.NoError(err)
	err = w.Explore()
	require.NoError(err)

	require.Equal("complete_test", w.Directory.Name)
	require.Len(w.Directory.Directories, 2)
	require.Len(w.Directory.Files, 1)
	require.Equal("05eae7dd459fc32142c65246877d9625f51bcec8a48e79432936227637d170af", w.Directory.Directories["folder1"].Files["file_1a"].SHA256)
	require.Equal("04337b307b9fe41137554ae2b1fddf1f0c6eb344fcfded725d971037d97311d4", w.Directory.Directories["folder2"].Directories["folder_22"].Files["file_22b"].SHA256)

	data, err := json.MarshalIndent(w.Directory, "", "  ")
	require.NoError(err)

	fmt.Printf("%s", data)
	// t.Fail()
}
