package main

import (
	"bytes"
	"testing"
)

const testFullResult = "├───\033[1;34mproject\033[0;m\n" +
	"│	├───\033[1;32mfile.txt\033[0;m \033[1;35m(19b)\033[0;m\n" +
	"│	└───\033[1;32mgopher.png\033[0;m \033[1;35m(70372b)\033[0;m\n" +
	"├───\033[1;34mstatic\033[0;m\n" +
	"│	├───\033[1;34ma_lorem\033[0;m\n" +
	"│	│	├───\033[1;32mdolor.txt\033[0;m \033[1;35m(empty)\033[0;m\n" +
	"│	│	├───\033[1;32mgopher.png\033[0;m \033[1;35m(70372b)\033[0;m\n" +
	"│	│	└───\033[1;34mipsum\033[0;m\n" +
	"│	│		└───\033[1;32mgopher.png\033[0;m \033[1;35m(70372b)\033[0;m\n" +
	"│	├───\033[1;34mcss\033[0;m\n" +
	"│	│	└───\033[1;32mbody.css\033[0;m \033[1;35m(28b)\033[0;m\n" +
	"│	├───\033[1;32mempty.txt\033[0;m \033[1;35m(empty)\033[0;m\n" +
	"│	├───\033[1;34mhtml\033[0;m\n" +
	"│	│	└───\033[1;32mindex.html\033[0;m \033[1;35m(57b)\033[0;m\n" +
	"│	├───\033[1;34mjs\033[0;m\n" +
	"│	│	└───\033[1;32msite.js\033[0;m \033[1;35m(10b)\033[0;m\n" +
	"│	└───\033[1;34mz_lorem\033[0;m\n" +
	"│		├───\033[1;32mdolor.txt\033[0;m \033[1;35m(empty)\033[0;m\n" +
	"│		├───\033[1;32mgopher.png\033[0;m \033[1;35m(70372b)\033[0;m\n" +
	"│		└───\033[1;34mipsum\033[0;m\n" +
	"│			└───\033[1;32mgopher.png\033[0;m \033[1;35m(70372b)\033[0;m\n" +
	"├───\033[1;34mzline\033[0;m\n" +
	"│	├───\033[1;32mempty.txt\033[0;m \033[1;35m(empty)\033[0;m\n" +
	"│	└───\033[1;34mlorem\033[0;m\n" +
	"│		├───\033[1;32mdolor.txt\033[0;m \033[1;35m(empty)\033[0;m\n" +
	"│		├───\033[1;32mgopher.png\033[0;m \033[1;35m(70372b)\033[0;m\n" +
	"│		└───\033[1;34mipsum\033[0;m\n" +
	"│			└───\033[1;32mgopher.png\033[0;m \033[1;35m(70372b)\033[0;m\n" +
	"└───\033[1;32mzzfile.txt\033[0;m \033[1;35m(empty)\033[0;m\n"

func TestRecursiveTreeFull(t *testing.T) {
	out := new(bytes.Buffer)
	err := dirTree(out, "testdata", true)
	if err != nil {
		t.Errorf("test for OK Failed - error")
	}
	result := out.String()
	if result != testFullResult {
		t.Errorf("test for OK Failed - results not match\nGot:\n%v\nExpected:\n%v", result, testFullResult)
	}
}

func TestIterativeTreeFull(t *testing.T) {
	out := new(bytes.Buffer)
	err := dirTreeIterative(out, "testdata", true)
	if err != nil {
		t.Errorf("test for OK Failed - error")
	}
	result := out.String()
	if result != testFullResult {
		t.Errorf("test for OK Failed - results not match\nGot:\n%v\nExpected:\n%v", result, testFullResult)
	}
}

const testDirResult = "├───\033[1;34mproject\033[0;m\n" +
	"├───\033[1;34mstatic\033[0;m\n" +
	"│	├───\033[1;34ma_lorem\033[0;m\n" +
	"│	│	└───\033[1;34mipsum\033[0;m\n" +
	"│	├───\033[1;34mcss\033[0;m\n" +
	"│	├───\033[1;34mhtml\033[0;m\n" +
	"│	├───\033[1;34mjs\033[0;m\n" +
	"│	└───\033[1;34mz_lorem\033[0;m\n" +
	"│		└───\033[1;34mipsum\033[0;m\n" +
	"└───\033[1;34mzline\033[0;m\n" +
	"	└───\033[1;34mlorem\033[0;m\n" +
	"		└───\033[1;34mipsum\033[0;m\n"

func TestRecursiveTreeDir(t *testing.T) {
	out := new(bytes.Buffer)
	err := dirTree(out, "testdata", false)
	if err != nil {
		t.Errorf("test for OK Failed - error")
	}
	result := out.String()
	if result != testDirResult {
		t.Errorf("test for OK Failed - results not match\nGot:\n%v\nExpected:\n%v", result, testDirResult)
	}
}

func TestIterativeTreeDir(t *testing.T) {
	out := new(bytes.Buffer)
	err := dirTreeIterative(out, "testdata", false)
	if err != nil {
		t.Errorf("test for OK Failed - error")
	}
	result := out.String()
	if result != testDirResult {
		t.Errorf("test for OK Failed - results not match\nGot:\n%v\nExpected:\n%v", result, testDirResult)
	}
}

const filesOnlyResult = "├───\033[1;32mfile.txt\033[0;m \033[1;35m(19b)\033[0;m\n" +
	"└───\033[1;32mgopher.png\033[0;m \033[1;35m(70372b)\033[0;m\n"

func TestRecursiveFilesOnly(t *testing.T) {
	out := new(bytes.Buffer)
	err := dirTree(out, "testdata/project", true)
	if err != nil {
		t.Errorf("test for OK Failed - error")
	}
	result := out.String()
	if result != filesOnlyResult {
		t.Errorf("test for OK Failed - results not match\nGot:\n%v\nExpected:\n%v", result, filesOnlyResult)
	}
}

func TestIterativeFilesOnly(t *testing.T) {
	out := new(bytes.Buffer)
	err := dirTreeIterative(out, "testdata/project", true)
	if err != nil {
		t.Errorf("test for OK Failed - error")
	}
	result := out.String()
	if result != filesOnlyResult {
		t.Errorf("test for OK Failed - results not match\nGot:\n%v\nExpected:\n%v", result, filesOnlyResult)
	}
}
