package main

import (
	"flag"
	"fmt"
	"os"
	"regexp"
	"strings"
)

func main() {

	var filePath string

	flag.Parse()
	if flag.NArg() == 1 {
		filePath = flag.Arg(0)
	} else {
		panic("invalid args")
	}

	fileInfo, err := os.Stat(filePath)

	if os.IsNotExist(err) {
		panic("input file or directory path does not exists.")
	}

	switch mode := fileInfo.Mode(); {
	case mode.IsDir():
		// do directory stuff
		fmt.Println("directory")
	case mode.IsRegular():
		// do file stuff
		fmt.Println("file")
		srcFile := fileInfo.Name()
		srcDir := strings.TrimRight(filePath, srcFile)
		doGetTitleAndPics(srcDir, srcFile)
	}
}

func doGetTitleAndPics(srcDir string, srcFileName string) {
	fmt.Println(srcDir)
	fmt.Println(srcFileName)

	if checkIgnoreFile(srcFileName) {
		return
	}

	newFileName := srcFileName

	itemCodeMatcher := regexp.MustCompile(`^[0-9a-zA-Z][^ _\.☆\(\)]+`)
	itemCode := itemCodeMatcher.FindString(newFileName)

	if len(itemCode) == 0 {
		return
	}
	fmt.Println(itemCode)

}

func checkIgnoreFile(srcFileName string) bool {
	if strings.HasSuffix(srcFileName, ".crdownload") {
		return true
	}
	if strings.Contains(srcFileName, ".rar") {
		return true
	}
	if strings.HasPrefix(srcFileName, "[") {
		return true
	}
	if strings.HasPrefix(srcFileName, "+[") {
		return true
	}
	if strings.HasPrefix(srcFileName, "IV)") {
		return true
	}
	return false
}
