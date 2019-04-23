package main

import (
	"fmt"
	"sync"
	"strings"
	"strconv"
	// "io"
//	"os"
//	"crypto/sha256"
//	"encoding/hex"

	"github.com/dutchcoders/goftp"
)

var wg = sync.WaitGroup{}

func main() {
	var ftp *goftp.FTP
	var err error

	//Connect to ftp server
	ftp, err = goftp.Connect("192.168.8.101:2121")
	if err != nil {
		panic(err)
	}

	//Authenticate
	err = ftp.Login("prince", "123") 
	if err != nil {
		panic(err)
	}

	wg.Add(1)
	var filesList []string
	go func() {
		filesList, err = ftp.List("/WhatsApp/media/WhatsApp Video")
		if err != nil {
			panic(err)
		}
		ftp.Close()
		wg.Done()
	}()
	
	wg.Wait()

	orderedList, totalFileSize := CreateDirectory(filesList)
	if err != nil {
		panic(err)
	}

// Create a handler that would handle the download.

//	PATH := "/WhatsApp/media/WhatsApp Video/VID-20180819-WA0003.mp4"
//	download := func(r io.Reader,info os.FileMode, err error) error {
//		var hasher = sha256.New()
//		if _,err = io.Copy(hasher, r); err != nil {
//			return(err)
//		}
//
//		hash := fmt.Sprintf("%s %x", PATH, hex.EncodeToString(hasher.Sum(nil)))
//		return nil
//	}
//
//	_,err = ftp.Walk("/", dowload)
//	if err != nil {
//		panic(err)
//	}

	// download := func(r io.Reader) error {
	
	// }



	// datafile, err := ftp.Retr(fmt.Sprintf("%s%s", "/WhatsApp/media/WhatsApp Video/", "VID-20180819-WA0003.mp4"), )
	// if err != nil {
	// 	panic(err)
	// }

	fmt.Println(totalFileSize)
	for file := range orderedList {
		fmt.Println(file)
	}
}

func CreateDirectory(filesList []string) (map[string]int, int) {
	var cache []string
	var directoryMap map[string]int
	var err error

	if filesList == nil {
		panic("no files")
	}	

	directoryMap = make(map[string]int)
	for _, eachList := range filesList {
		sepSlice := strings.Split(eachList, ";")
		for _, eachSliceElem := range sepSlice {
			cache = append(cache, eachList)
			if ok := strings.Contains(eachSliceElem, "VID"); ok {
				fileName, size := FilenameAndSize(sepSlice)
				directoryMap[fileName] = size
				if err != nil {
					panic(err)
				}
			}
		}
	}
	return directoryMap, SumOfFileSizes(directoryMap)
}

func FilenameAndSize(data []string) (string, int) {
	var fileSize string
	var fileName string
	var convFileSize int
	var err error

	for _, eachData := range data {
		if ok := strings.Contains(eachData, "size"); ok {
			fileSize = strings.Trim(eachData, "size=")
			convFileSize, err = strconv.Atoi(fileSize)
			if err != nil {
				panic(err)
			}
		}
		if ok := strings.Contains(eachData, "VID"); ok {
			fileName = eachData
		}
	}
	return fileName, convFileSize
}

func SumOfFileSizes(data map[string]int) int {
	var totalSize int

	for _, size := range data {
		totalSize = totalSize + size
	}
	return totalSize
}