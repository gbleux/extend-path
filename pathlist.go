package main

import (
	"os"
	"path/filepath"
	"strings"
)

type Directory string

/**
 * verify each candidate is part of the target array. the order
 * of appearence is not checked, only the presence. this method
 * returns false if at least a single candidate is NOT contained
 * in the target slice.
 */
func verifyEntries(candidates []Directory, target []Directory) bool {
	var filtered = filterDirectories(candidates, target)

	return 0 == len(filtered)
}

/**
 * remove all candidates from the target slice. the returned array
 * will contain none of the candidates, even if they occur multiple
 * times.
 */
func deleteEntries(candidates []Directory, target []Directory) []Directory {
	var filtered = filterDirectories(target, candidates)

	return filtered
}

func appendEntries(appendence []Directory, directories []Directory, relocate bool) []Directory {
	var dirs1, dirs2 = splitDirectories(appendence, directories, relocate)

	return append(dirs1, dirs2...)
}

func prependEntries(prependence []Directory, directories []Directory, relocate bool) []Directory {
	var dirs1, dirs2 = splitDirectories(prependence, directories, relocate)

	return append(dirs2, dirs1...)
}

func directoriesContains(directories []Directory, directory Directory) bool {
	for _, dir := range directories {
		if directoryEquals(dir, directory) {
			return true
		}
	}

	return false
}

/**
 * directory
 */
func directoryEquals(dir1 Directory, dir2 Directory) bool {
	var cmp1 = strings.TrimRight(string(dir1), "\\/")
	var cmp2 = strings.TrimRight(string(dir2), "\\/")

	return cmp1 == cmp2
}

/**
 * remove all entries from the directories, which are part of the filter.
 *
 * example:
 *     var directories = []Directory{"a", "b", "d"}
 *     var filter = []Directory{"b", "c"}
 *     filterDirectories(directories, filter) // ["a", "d"]
 */
func filterDirectories(directories []Directory, filter []Directory) []Directory {
	var filtered = []Directory{}

	for _, directory := range directories {
		if false == directoriesContains(filter, directory) {
			filtered = append(filtered, directory)
		}
	}

	return filtered
}

func splitDirectories(needles []Directory, haystack []Directory, relocate bool) (directories1 []Directory, directories2 []Directory) {
	if relocate {
		directories1 = filterDirectories(haystack, needles)
		directories2 = needles
	} else {
		directories2 = filterDirectories(needles, haystack)
		directories1 = haystack
	}

	return
}

func resolveDirectories(directories []string, excludeInvalid bool) []Directory {
	var filtered = []Directory{}
	var err error = nil
	var absolute string
	var handler os.FileInfo

	for _, directory := range directories {
		absolute, err = filepath.Abs(directory)

		if nil != err {
			continue
		}

		if excludeInvalid {
			handler, err = os.Stat(absolute)

			if nil != err || false == handler.IsDir() {
				continue
			}
		}

		filtered = append(filtered, Directory(absolute))
	}

	return filtered
}

func splitDirectoryPath(directories string) []Directory {
	var splitlist = filepath.SplitList(directories)
	var converted = make([]Directory, len(splitlist))

	for index, dir := range splitlist {
		converted[index] = Directory(dir)
	}

	return converted
}