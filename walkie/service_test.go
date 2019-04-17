package walkie

import (
	"encoding/json"
	"path/filepath"
	"runtime"
	"testing"
	"time"

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

// Test unexistant path
func TestWalkingUnknownPath(t *testing.T) {
	require := require.New(t)
	testdir := "./carmen/sandiego"

	_, err := NewWalkie(testdir)
	require.Error(err)
	// t.Fail()
}

func clean(t *testing.T) {
	require := require.New(t)
	testdir := "./carmen/sandiego"

	_, err := NewWalkie(testdir)
	require.Error(err)
	// t.Fail()
}

func BenchmarkWalinkg(b *testing.B) {
	// run the Fib function b.N times
	testdir := getTestDir()
	w, err := NewWalkie(testdir)
	if err != nil {
		b.Fatal(err)
	}

	for n := 0; n < b.N; n++ {
		err = w.Explore()
		if err != nil {
			b.Fatal(err)
		}
	}
}

// Test Directory walking
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

	_, err = json.MarshalIndent(w.Directory, "", "  ")
	require.NoError(err)

	// fmt.Printf("%s", data)
	// t.Fail()
}

// Test Directory walking
func TestCompare(t *testing.T) {
	require := require.New(t)
	testdir := getTestDir()

	w, err := NewWalkie(testdir)
	require.NoError(err)
	err = w.Explore()
	require.NoError(err)

	f1 := w.Directory.Directories["folder1"].Files["file_1a"]
	f2 := w.Directory.Directories["folder1"].Files["file_1b"]

	require.True(f1.Equals(*f1))
	require.False(f1.Equals(*f2))

	d1 := w.Directory.Directories["folder1"]
	d2 := w.Directory.Directories["folder2"]

	// Shallow test
	require.True(d1.Equals(*d1))
	require.False(d1.Equals(*d2))

	// Do a copy
	d3 := w.Directory.Directories["folder1"].copy()
	// d3.Files["file_1a"] = d3.Files["file_1a"]
	require.True(d1.DeepEquals(*d3))

	// Change a bit then test
	d3.Files["file_1a"] = d3.Files["file_1a"].copy()
	d3.Files["file_1a"].Mtime = time.Now()
	//
	require.NotEqual(d3.Files["file_1a"].Mtime, d1.Files["file_1a"].Mtime)
	require.False(d1.DeepEquals(*d3))
}
