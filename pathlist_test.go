package main

import (
	"path/filepath"
	"strings"

	"testing"
)

var (
	// shared test path
	PATH = []Directory{"/home/bin", "/bin", "/sbin"}
	// directory NOT part of path
	MYBIN = Directory("/opt/bin")
)

func TestAppendEmpty(t *testing.T) {
	var path = pathFactory()
	var entries = []Directory{}
	var expected = path
	var actual = appendEntries(entries, path, true)

	assertArrayEquals(t, expected, actual)
}

func TestAppendRelocate(t *testing.T) {
	var path = pathFactory()
	var entries = []Directory{MYBIN, path[1]}
	var pathdel = append([]Directory{path[0]}, path[2:]...)
	var expected = append(pathdel, entries...)
	var actual = appendEntries(entries, path, true)

	assertArrayEquals(t, expected, actual)
}

func TestAppendNoRelocate(t *testing.T) {
	var path = pathFactory()
	var entries = []Directory{MYBIN, path[1]}
	var expected = append(path, MYBIN)
	var actual = appendEntries(entries, path, false)

	assertArrayEquals(t, expected, actual)
}

func TestPrependEmpty(t *testing.T) {
	var path = pathFactory()
	var entries = []Directory{}
	var expected = path
	var actual = prependEntries(entries, path, true)

	assertArrayEquals(t, expected, actual)
}

func TestPrependRelocate(t *testing.T) {
	var path = pathFactory()
	var entries = []Directory{MYBIN, path[1]}
	var pathdel = append([]Directory{path[0]}, path[2:]...)
	var expected = append(entries, pathdel...)
	var actual = prependEntries(entries, path, true)

	assertArrayEquals(t, expected, actual)
}

func TestPrependNoRelocate(t *testing.T) {
	var path = pathFactory()
	var entries = []Directory{MYBIN, path[1]}
	var expected = append([]Directory{MYBIN}, path...)
	var actual = prependEntries(entries, path, false)

	assertArrayEquals(t, expected, actual)
}

func TestDeleteAffected(t *testing.T) {
	var path = pathFactory()
	var entries = []Directory{MYBIN, path[1]}
	var expected = append([]Directory{path[0]}, path[2:]...)
	var actual = deleteEntries(entries, path)

	assertArrayEquals(t, expected, actual)
}

func TestDeleteDupes(t *testing.T) {
	var path = pathFactory()
	var entries = []Directory{MYBIN}
	var target = []Directory{MYBIN, path[0], MYBIN}
	var expected = []Directory{path[0]}
	var actual = deleteEntries(entries, target)

	assertArrayEquals(t, expected, actual)
}

func TestDeleteEmpty(t *testing.T) {
	var path = pathFactory()
	var entries = []Directory{}
	var expected = path
	var actual = deleteEntries(entries, path)

	assertArrayEquals(t, expected, actual)
}

func TestDeleteUnaffected(t *testing.T) {
	var path = pathFactory()
	var entries = []Directory{MYBIN, "/opt/bin"}
	var expected = path
	var actual = deleteEntries(entries, path)

	assertArrayEquals(t, expected, actual)
}

func TestVerifyPartialFail(t *testing.T) {
	var path = pathFactory()
	var input = append(path, MYBIN)

	assertVerification(t, false, input, path)
}

func TestVerifyFail(t *testing.T) {
	var path = pathFactory()
	var input = []Directory{MYBIN}

	assertVerification(t, false, input, path)
}

func TestVerifyEmptySuccess(t *testing.T) {
	var path = pathFactory()
	var input = []Directory{}

	assertVerification(t, true, input, path)
}

func TestVerifyPartialSuccess(t *testing.T) {
	var path = pathFactory()
	var input = path[1:]

	assertVerification(t, true, input, path)
}

func TestVerifySuccess(t *testing.T) {
	var path = pathFactory()
	var input = path

	assertVerification(t, true, input, path)
}

func TestDirectoriesContains(t *testing.T) {
	var path = pathFactory()

	assertContains(t, true, path[0], path)
	assertContains(t, true, path[0]+"/", path)

	assertContains(t, false, MYBIN, path)
	assertContains(t, false, "/"+path[0], path)
}

func TestFilterDirectories(t *testing.T) {
	var path = pathFactory()
	var actual = filterDirectories(path, []Directory{path[0]})

	assertArrayEquals(t, path[1:], actual)
}

func TestFilterDirectoriesNone(t *testing.T) {
	var path = pathFactory()
	var actual = filterDirectories(path, []Directory{MYBIN})

	assertArrayEquals(t, path, actual)
}

func TestFilterDirectoriesEmpty(t *testing.T) {
	var path = pathFactory()
	var actual = filterDirectories(path, []Directory{})

	assertArrayEquals(t, path, actual)
}

func TestFilterDirectoriesAll(t *testing.T) {
	var path = pathFactory()
	var actual = filterDirectories(path, path)

	assertArrayEquals(t, []Directory{}, actual)
}

func TestFilterDirectoriesDoc(t *testing.T) {
	var directories = []Directory{"a", "b", "d"}
	var filter = []Directory{"b", "c"}
	var actual = filterDirectories(directories, filter)

	assertArrayEquals(t, []Directory{"a", "d"}, actual)
}

func TestSplitDirectoriesRelocate(t *testing.T) {
	var needles = []Directory{"b", "d"}
	var haystack = []Directory{"a", "b", "c"}
	var actual1, actual2 = splitDirectories(needles, haystack, true)

	assertArrayEquals(t, []Directory{"a", "c"}, actual1)
	assertArrayEquals(t, []Directory{"b", "d"}, actual2)
}

func TestSplitDirectoriesNoRelocate(t *testing.T) {
	var needles = []Directory{"b", "d"}
	var haystack = []Directory{"a", "b", "c"}
	var actual1, actual2 = splitDirectories(needles, haystack, false)

	assertArrayEquals(t, []Directory{"a", "b", "c"}, actual1)
	assertArrayEquals(t, []Directory{"d"}, actual2)
}

func TestDirectoryEquals(t *testing.T) {
	var path = pathFactory()

	assertDirectoryEquals(t, true, MYBIN, MYBIN)
	assertDirectoryEquals(t, true, MYBIN+"/", MYBIN)
	assertDirectoryEquals(t, true, MYBIN, MYBIN+"/")
	assertDirectoryEquals(t, true, MYBIN+"//", MYBIN)

	assertDirectoryEquals(t, false, "//"+MYBIN, MYBIN)
	assertDirectoryEquals(t, false, "/usr"+MYBIN, MYBIN)
	assertDirectoryEquals(t, false, MYBIN+"/vendor", MYBIN)
	assertDirectoryEquals(t, false, path[0], MYBIN)
}

func TestSplitDirectoryPath(t *testing.T) {
	var input = strings.Join([]string{"a", "b", "c"}, string(filepath.ListSeparator))
	var expected = []Directory{"a", "b", "c"}
	var actual = splitDirectoryPath(input)

	assertArrayEquals(t, expected, actual)
}

func assertVerification(t *testing.T, expected bool, directories []Directory, path []Directory) {
	if expected != verifyEntries(directories, path) {
		t.Errorf("Verification result mismatch for %v in %v.\nExpected: %v",
			directories, path, expected)
	}
}

func assertContains(t *testing.T, expected bool, dir Directory, path []Directory) {
	if expected != directoriesContains(path, dir) {
		t.Errorf("Invalid directory lookup result for %s.\nExpected: %v",
			dir, expected)
	}
}

func assertDirectoryEquals(t *testing.T, expected bool, dir1 Directory, dir2 Directory) {
	if expected != directoryEquals(dir1, dir2) {
		t.Errorf("Invalid directory equality result for %s and %s.\nExpected %v",
			dir1, dir2, expected)
	}
}

func assertArrayEquals(t *testing.T, expected []Directory, actual []Directory) {
	if len(actual) != len(expected) {
		t.Fatalf("Array length mismatch.\nExpected: %d: %v,\nbut was: %d: %v",
			len(expected), expected, len(actual), actual)
	}

	for i := 0; i < len(expected); i++ {
		if actual[i] != expected[i] {
			t.Errorf("Array item mismatch.\nExpected: %s,\nbut was: %s",
				expected[i], actual[i])
		}
	}
}

/**
 * copy method for PATH to avoid side-effects
 */
func pathFactory() []Directory {
	artifact := make([]Directory, len(PATH))

	copy(artifact, PATH)

	return artifact
}
