package main

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"path/filepath"
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

		files, err := ioutil.ReadDir(filePath)
		if err != nil {
			panic(err)
		}

		for index, file := range files {
			fmt.Println("---> [" + strconv.Itoa(index+1) + "/" + strconv.Itoa(len(files)) + "]")
			if file.IsDir() == false {
				doGetTitleAndPics(filePath, file.Name())
			}
		}
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
		fmt.Println("skip")
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
	newFileName := buildFileName(srcFileName, itemCode, itemInfo)

	fmt.Println(newFileName)

	downloadPackageImage(srcDir, newFileName, itemInfo)

	err = os.Rename(
		srcDir+string(os.PathSeparator)+srcFileName,
		srcDir+string(os.PathSeparator)+newFileName)

	if err != nil {
		panic("rename fail")
	}
}

func checkIgnoreFile(srcFileName string) bool {
	if strings.HasPrefix(srcFileName, ".") {
		return true
	}
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
		fmt.Println(value.ItemDetailURL)
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

func buildFileName(srcFileName string, itemCode string, itemInfo *goScrapeDmmCoJp.ItemOfDmmCoJp) string {
	var newFileName string

	newFileName = "[" + strings.ToUpper(itemCode) + "]" + itemInfo.Title
	if len(itemInfo.ActorList) > 0 {
		newFileName += " {"

		actorNameList := []string{}
		for _, value := range itemInfo.ActorList {
			actorNameList = append(actorNameList, value.Name)
		}
		newFileName += strings.Join(actorNameList, ", ")

		newFileName += "}"
		newFileName += strings.TrimLeft(srcFileName, itemCode)
	}

	return newFileName
}

func downloadPackageImage(srcDir string, newFileName string, itemInfo *goScrapeDmmCoJp.ItemOfDmmCoJp) {
	response, err := http.Get(itemInfo.PackageImageURL)
	if err != nil {
		fmt.Println(err)
		panic("package image download err in get")
	}

	fmt.Println("status:", response.Status)

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println(err)
		panic("package image download err in read")
	}

	downloadFileName := srcDir
	downloadFileName += string(os.PathSeparator)
	downloadFileName += strings.TrimRight(newFileName, filepath.Ext(newFileName))
	downloadFileName += ".jpg"

	file, err := os.OpenFile(downloadFileName, os.O_CREATE|os.O_WRONLY, 0666)

	if err != nil {
		fmt.Println(err)
		panic("output file ")
	}

	defer func() {
		file.Close()
	}()

	file.Write(body)
}
