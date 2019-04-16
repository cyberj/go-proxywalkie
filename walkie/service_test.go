package walkie

import (
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
	// t.Fail()
}
