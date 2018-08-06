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

	searchResult := selectFromSearchResultList(searchResultList)

	fmt.Println(searchResult.Title)

	itemInfo := goScrapeDmmCoJp.GetItemInfoFromURL(searchResult.ItemDetailURL)
	newFileName := buildFilenameFromItemCodeAndItemInfo(itemCode, itemInfo)

	fmt.Println(newFileName)

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

func selectFromSearchResultList(searchResultList []*goScrapeDmmCoJp.SearchListItem) *goScrapeDmmCoJp.SearchListItem {

	if len(searchResultList) == 1 {
		return searchResultList[0]
	}

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
	return searchResultList[indexFromScan]
}

func buildFilenameFromItemCodeAndItemInfo(itemCode string, itemInfo *goScrapeDmmCoJp.ItemOfDmmCoJp) string {
	var filename string

	filename = "[" + itemCode + "]" + itemInfo.Title
	if len(itemInfo.ActorList) > 0 {
		filename += " {"

		actorNameList := []string{}
		for _, value := range itemInfo.ActorList {
			actorNameList = append(actorNameList, value.Name)
		}
		filename += strings.Join(actorNameList, ", ")

		filename += "}"
	}

	return filename
}
