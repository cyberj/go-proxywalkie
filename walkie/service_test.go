package walkie

import (
	"encoding/json"
	"testing"

	"github.com/cyberj/go-proxywalkie/testutils"
	"github.com/stretchr/testify/require"
)

// Test unexistant path
func TestWalkingUnknownPath(t *testing.T) {
	require := require.New(t)
	testdir := "./carmen/sandiego"

	_, err := NewWalkie(testdir)
	require.Error(err)
	// t.Fail()
}

func BenchmarkWalking(b *testing.B) {
	// run the Fib function b.N times

	testdirs, err := testutils.NewTestDir()
	if err != nil {
		b.Error(err)
	}
	defer testdirs.Clean()

	w, err := NewWalkie(testdirs.TestassetsDir)
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
	testdirs, err := testutils.NewTestDir()
	require.NoError(err)
	defer testdirs.Clean()

	w, err := NewWalkie(testdirs.TestassetsDir)
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

	require.Len(w.directories, 4)
	require.Len(w.files, 11)

	// fmt.Printf("%s", data)
	// t.Fail()
}
