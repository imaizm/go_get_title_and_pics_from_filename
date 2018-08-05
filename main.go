package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/imaizm/go_scrape_dmm.co.jp"
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

	//newFileName := srcFileName

	itemCode, err := getItemCodeFromFilename(srcFileName)
	if err != nil {
		panic(err)
	}

	searchResultList := goScrapeDmmCoJp.Search(itemCode)
	if len(searchResultList) == 0 {
		fmt.Println("item not found at dmm.co.jp")
		return
	}

	var searchResult *goScrapeDmmCoJp.SearchListItem
	if len(searchResultList) == 1 {
		searchResult = searchResultList[0]
	} else {
		for index, value := range searchResultList {
			fmt.Println(strconv.Itoa(index) + " : " + value.Title)
		}

		indexFromScan := 0
		scanComplete := false
		for scanComplete == false {
			var stdin string
			fmt.Scan(&stdin)
			input, err := strconv.Atoi(stdin)
			if err != nil {
				fmt.Println("is not number : " + stdin)
			} else if input < 0 || input >= len(searchResultList) {
				fmt.Println("is not between 0 and " + strconv.Itoa(len(searchResultList)-1) + " : " + stdin)
			} else {
				indexFromScan = input
				scanComplete = true
			}
		}
		searchResult = searchResultList[indexFromScan]
	}

	fmt.Println(searchResult.Title)

	result := goScrapeDmmCoJp.New(searchResult.ItemDetailURL)

	fmt.Println("ItemCode : " + result.ItemCode)
	fmt.Println("Title : " + result.Title)
	fmt.Println("PackageImageThumbURL : " + result.PackageImageThumbURL)
	fmt.Println("PackageImageURL : " + result.PackageImageURL)
	fmt.Println("ActorList :")
	for index, value := range result.ActorList {
		fmt.Println("\t" + strconv.Itoa(index) + " : " + value.Name + " : " + value.ListPageURL)
	}
	fmt.Println("SampleImageList :")
	for index, value := range result.SampleImageList {
		fmt.Println("\t" + strconv.Itoa(index) + " : " + value.ImageThumbURL + " : " + value.ImageURL)
	}

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

func getItemCodeFromFilename(fileName string) (string, error) {
	itemCodeMatcher := regexp.MustCompile(`^[0-9a-zA-Z][^ _\.☆\(\)]+`)
	itemCode := itemCodeMatcher.FindString(fileName)

	if len(itemCode) == 0 {
		err := errors.New("itemCode not found in filename")
		return "", err
	}

	return itemCode, nil
}
