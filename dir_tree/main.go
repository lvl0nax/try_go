package main

import (
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
)

func prefix(last bool) string {
	sign := "├───"
	if last {
		sign = "└───"
	}

	return sign
}

func isRoot(path string) bool {
	return strings.Count(path, "/") == 0
}

func newOffset(path string, last bool, offset string) string {
	if isRoot(path) {
		return offset
	}

	rep := "│\t"
	if last {
		rep = "\t"
	}
	return offset + rep
}

func getSubnames(dir *os.File, withFiles bool) []string {
	if withFiles {
		names, _ := dir.Readdirnames(0)
		return names
	}

	fileInfos, _ := dir.Readdir(0)
	var names []string
	counter := len(fileInfos)
	for i := 0; i < counter; i++ {
		if fileInfos[i].IsDir() {
			names = append(names, fileInfos[i].Name())
		}
	}

	return names
}

func folderContents(path string, withFiles bool, last bool) ([]string, string) {
	dir, err := os.Open(path)
	defer dir.Close()

	if err != nil {
		return nil, ""
	}

	info, _ := dir.Stat()

	if withFiles && !info.IsDir() {
		size := fmt.Sprint(info.Size()) + "b"
		if size == "0b" {
			size = "empty"
		}
		return nil, fmt.Sprintf("%s%s (%s)", prefix(last), info.Name(), size)
	}

	var str string
	if info.IsDir() && !isRoot(path) {
		str = fmt.Sprintf("%s%s", prefix(last), info.Name())
	}

	return getSubnames(dir, withFiles), str
}

func printTree(path string, offset string, withFiles bool, last bool, output io.Writer) {
	names, str := folderContents(path, withFiles, last)

	if str != "" {
		fmt.Fprintln(output, offset+str)
	}

	if names == nil {
		return
	}

	sort.Strings(names)

	counter := len(names)
	for i := 0; i < counter; i++ {
		isLast := counter == i+1
		printTree(path+"/"+names[i], newOffset(path, last, offset), withFiles, isLast, output)
	}
}

func dirTree(output io.Writer, path string, withFiles bool) error {
	printTree(path, "", withFiles, true, output)

	return nil
}

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
