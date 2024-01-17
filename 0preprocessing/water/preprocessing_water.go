package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var relativePath = filepath.Join("..", "..", "Datasets", "water", "swm_trialA_1K")
var rFolderFormat = "swm_trialA_%dK.csv"
var wFolderFormat = "households_%d"
var wFileFormat = filepath.Join("households_%d", "%s.csv")
var wd string

const FILE_ROWS = 40960 // element number within a section

func main() {
	var wd, _ = os.Getwd()
	fmt.Println(wd)
	genCSV(1)
}

//file reading
func ReadCSV(path string) []string {
	// fmt.Println("reading without buffer:")
	data, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}
	raw := strings.Replace(string(data), " \r", "", -1)
	rawArr := strings.Split(raw, "\n")
	arr := rawArr[1 : len(rawArr)-1]
	return arr
}

//trim csv
func genCSV(id int) {
	rFolderPath := fmt.Sprintf(rFolderFormat, id)
	rFolderFullPath := filepath.Join(wd, relativePath, rFolderPath)
	lines := ReadCSV(rFolderFullPath)
	myMap := make(map[string][][]string)

	var file string
	for _, line := range lines {
		line = strings.Replace(line, " ", ";", -1)
		line = strings.Replace(line, "Null\r", "0.001", -1) //fix some broken data

		slices := strings.Split(line, ";")
		file = slices[0]
		if _, ok := myMap[file]; !ok {
			myMap[file] = make([][]string, 0)
		}
		myMap[file] = append(myMap[file], slices[1:])
	}

	//create folder
	wFolderPath := fmt.Sprintf(wFolderFormat, FILE_ROWS)
	wFolderFullPath := filepath.Join(wd, relativePath, wFolderPath)
	_, err := os.Stat(wFolderFullPath)
	if err != nil {
		if os.IsNotExist(err) {
			//create folder if not existed
			err = os.Mkdir(wFolderFullPath, 0755)
			if err != nil {
				log.Fatal("Error creating folder:", err)
			}
		}
	}

	//write files
	for key, value := range myMap {
		if len(value) < FILE_ROWS {
			fmt.Printf("[%s] has %d rows.\n", key, len(value))
			continue
		}
		//create csv file
		wFilePath := fmt.Sprintf(wFileFormat, FILE_ROWS, key)
		wFileFullPath := filepath.Join(wd, relativePath, wFilePath)
		newFile, err := os.Create(wFileFullPath)
		if err != nil {
			log.Fatal("Error creating file:")
		}
		defer newFile.Close()

		//create csv writer
		newWriter := csv.NewWriter(newFile)
		defer newWriter.Flush()

		//write data into csv file
		for _, row := range value {
			err := newWriter.Write(row)
			if err != nil {
				log.Fatal("Error writing to CSV:", err)
			}
		}
	}
	fmt.Printf("len: %d\n", len(myMap))
}
