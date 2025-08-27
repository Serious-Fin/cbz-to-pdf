package main

import (
	"slices"
	"testing"
)

func TestIsCBZFile(t *testing.T) {
	testFailed := false
	files := []string{
		"foo.cbz",
		".cbz",
		"hello world.cbz",
		"test.CBZ",
		"test.CbZ",
		"test.CBz",
		"test.cBZ",
		"test (2020 edition).cbz",
		"test.",
		"test.txt",
		" .aa",
		" .cbz",
		"foo.cbz.txt",
		".cbz.txt",
	}

	want := []string{
		".cbz",
		" .cbz",
		"foo.cbz",
		"hello world.cbz",
		"test.CBZ",
		"test.CbZ",
		"test.CBz",
		"test.cBZ",
		"test (2020 edition).cbz",
	}

	got := make([]string, 0)

	for _, filename := range files {
		if isCbzFile(filename) {
			got = append(got, filename)
		}
	}

	for _, filename := range want {
		if !slices.Contains(got, filename) {
			t.Logf("missing filename \"%s\"", filename)
			testFailed = true
		}
	}

	for _, filename := range got {
		if !slices.Contains(want, filename) {
			t.Logf("file \"%s\" should not be cbz", filename)
			testFailed = true
		}
	}

	if testFailed {
		t.Fail()
	}
}
