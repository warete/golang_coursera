package main

import (
	"io"
	"io/ioutil"
	"os"
	"sort"
	"strconv"
)

type MyFile struct {
	IsDir bool
	Size  int
}

func getDirData(path string, needFiles bool) (error, map[string]MyFile) {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return err, nil
	}

	dirContent := make(map[string]MyFile)

	i := 0
	for _, file := range files {
		i++
		if !needFiles && !file.IsDir() {
			continue
		}
		var fileSize int = 0
		if !file.IsDir() {
			fileSize = int(file.Size())
		}
		dirContent[file.Name()] = MyFile{
			file.IsDir(),
			fileSize,
		}
	}

	return nil, dirContent
}

func readDir(out io.Writer, path string, printFiles bool, level int, dirPrefix string) error {
	err, dirContent := getDirData(path, printFiles)
	if err != nil {
		return err
	}

	keys := make([]string, 0, len(dirContent))
	for k := range dirContent {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	i := 0
	for _, k := range keys {
		var filePrefix string
		var dirPrefixLocal string
		if i == len(keys)-1 {
			filePrefix = "└───"
			dirPrefixLocal = "\t"
		} else {
			filePrefix = "├───"
			dirPrefixLocal = "│\t"
		}

		if dirContent[k].IsDir {
			//fmt.Println(dirPrefix + filePrefix + k)
			out.Write([]byte(dirPrefix + filePrefix + k + "\n"))
			err := readDir(out, path+string(os.PathSeparator)+k, printFiles, level+1, dirPrefix+dirPrefixLocal)
			if err != nil {
				return err
			}
		} else {
			var fileSizeData string
			if dirContent[k].Size > 0 {
				fileSizeData = strconv.Itoa(dirContent[k].Size) + "b"
			} else {
				fileSizeData = "empty"
			}
			//fmt.Println()
			out.Write([]byte(dirPrefix + filePrefix + k + " (" + fileSizeData + ")" + "\n"))
		}
		i++
	}

	return nil
}

func dirTree(out io.Writer, path string, printFiles bool) error {
	err := readDir(out, path, printFiles, 0, "")
	if err != nil {
		return err
	}
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
