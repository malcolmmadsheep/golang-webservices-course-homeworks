package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
)

func main() {
	out := os.Stdout
	if !(len(os.Args) == 2 || len(os.Args) == 3) {
		panic("usage go run main.go . [-f]")
	}
	path := os.Args[1]
	printFiles := len(os.Args) == 3 && os.Args[2] == "-f"
	err := dirTree(out, path, printFiles)
	if err != nil {
		panic(err.Error())
	}
}

func dirTree(out io.Writer, path string, printFiles bool) error {
	return recursiveDirTree(out, path, printFiles, "")
}

func filterOutFiles(files []os.FileInfo) []os.FileInfo {
	folders := make([]os.FileInfo, 0)

	for _, file := range files {
		if !file.IsDir() {
			continue
		}

		folders = append(folders, file)
	}

	return folders
}

func recursiveDirTree(out io.Writer, path string, keepFiles bool, prefix string) error {
	folder, err := os.Open(path)

	if err != nil {
		return err
	}

	defer folder.Close()

	files, err := folder.Readdir(-1)

	if err != nil {
		return err
	}

	if !keepFiles {
		files = filterOutFiles(files)
	}

	sort.SliceStable(files, func(i, j int) bool {
		return files[i].Name() < files[j].Name()
	})

	for index, file := range files {
		name := file.Name()
		isDir := file.IsDir()
		isLast := index == len(files)-1

		fmt.Fprintf(out, prefix)

		if isLast {
			fmt.Fprintf(out, "└───")
		} else {
			fmt.Fprintf(out, "├───")
		}

		if isDir {
			fmt.Fprintf(out, "\033[1;34m")
		} else {
			fmt.Fprintf(out, "\033[1;32m")
		}

		fmt.Fprintf(out, name)

		if !isDir {
			fileSize := file.Size()

			fmt.Fprintf(out, "\033[1;35m")

			if fileSize == 0 {
				fmt.Fprintf(out, " (empty)")
			} else {
				fmt.Fprintf(out, " (%db)", fileSize)
			}
		}

		fmt.Fprintf(out, "\033[0m\n")

		if isDir {
			if isLast {
				recursiveDirTree(out, filepath.Join(path, name), keepFiles, prefix+"\t")

				continue
			}

			recursiveDirTree(out, filepath.Join(path, name), keepFiles, prefix+"│\t")
		}
	}

	return nil
}
