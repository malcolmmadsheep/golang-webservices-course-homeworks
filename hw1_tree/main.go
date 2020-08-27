package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

func main() {
	out := os.Stdout
	if !(len(os.Args) == 2 || len(os.Args) == 3) {
		panic("usage go run main.go . [-f]")
	}
	path := os.Args[1]
	keepFiles := len(os.Args) == 3 && os.Args[2] == "-f"
	// err := dirTree(out, path, keepFiles)
	err := dirTreeIterative(out, path, keepFiles)
	if err != nil {
		panic(err.Error())
	}
}

func dirTree(out io.Writer, path string, keepFiles bool) error {
	return recursiveDirTree(out, path, keepFiles, []string{})
}

func readDir(path string, keepFiles bool) ([]os.FileInfo, error) {
	folder, err := os.Open(path)

	if err != nil {
		return nil, err
	}

	files, err := folder.Readdir(-1)

	if err != nil {
		return nil, err
	}

	folder.Close()

	if !keepFiles {
		files = filterOutFiles(files)
	}

	sortFilesByName(files)

	return files, nil
}

func sortFilesByName(files []os.FileInfo) {
	sort.SliceStable(files, func(i, j int) bool {
		return files[i].Name() < files[j].Name()
	})
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

func printOutLine(out io.Writer, file os.FileInfo, prefix []string, isLast bool) {
	fmt.Fprintf(out, strings.Join(prefix, ""))

	if isLast {
		fmt.Fprintf(out, "└───")
	} else {
		fmt.Fprintf(out, "├───")
	}

	if file.IsDir() {
		fmt.Fprintf(out, "\033[1;34m")
	} else {
		fmt.Fprintf(out, "\033[1;32m")
	}

	fmt.Fprintf(out, "%s%s", file.Name(), "\033[0;m")
}

func printOutFileSize(out io.Writer, file os.FileInfo) {
	fileSize := file.Size()

	fmt.Fprintf(out, " \033[1;35m")

	if fileSize == 0 {
		fmt.Fprintf(out, "(empty)")
	} else {
		fmt.Fprintf(out, "(%db)", fileSize)
	}

	fmt.Fprintf(out, "\033[0;m")
}

func dirTreeIterative(out io.Writer, root string, keepFiles bool) error {
	path := []string{root}
	prefix := []string{""}

	rootFolderFiles, err := readDir(root, keepFiles)

	if err != nil {
		return err
	}

	files := [][]os.FileInfo{rootFolderFiles}

	for len(files) != 0 {
		curDirFiles := files[len(files)-1]

		if len(curDirFiles) == 0 {
			files = files[:len(files)-1]
			prefix = prefix[:len(prefix)-1]
			path = path[:len(path)-1]

			continue
		}

		for len(curDirFiles) != 0 {
			file := curDirFiles[0]
			curDirFiles = curDirFiles[1:]
			files[len(files)-1] = curDirFiles

			isLast := len(curDirFiles) == 0

			printOutLine(out, file, prefix, isLast)

			if !file.IsDir() {
				printOutFileSize(out, file)
			}

			fmt.Fprintf(out, "\n")

			if file.IsDir() {
				path = append(path, file.Name())

				dirFiles, err := readDir(filepath.Join(path...), keepFiles)

				if err != nil {
					return err
				}

				files = append(files, dirFiles)

				if isLast {
					prefix = append(prefix, "\t")
				} else {
					prefix = append(prefix, "│\t")
				}

				break
			}
		}
	}

	return nil
}

func recursiveDirTree(out io.Writer, path string, keepFiles bool, prefix []string) error {
	files, err := readDir(path, keepFiles)

	if err != nil {
		return err
	}

	for index, file := range files {
		name := file.Name()
		isDir := file.IsDir()
		isLast := index == len(files)-1

		printOutLine(out, file, prefix, isLast)

		if !isDir {
			printOutFileSize(out, file)
		}

		fmt.Fprintf(out, "\n")

		if isDir {
			parentPath := filepath.Join(path, name)

			if isLast {
				recursiveDirTree(out, parentPath, keepFiles, append(prefix, "\t"))

				continue
			}

			recursiveDirTree(out, parentPath, keepFiles, append(prefix, "│\t"))
		}
	}

	return nil
}
