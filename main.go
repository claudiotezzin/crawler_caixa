package main

import (
	"fmt"
	"os"
	"strings"
	"time"
)

func main() {
	lines, err1 := getTodayFileIfExist(time.Now(), "")

	if err1 != nil {
		lines, err1 = downloadTodayFile()
		if err1 != nil {
			fmt.Println("Error getting new file. Error: " + err1.Error())
			os.Exit(1)
		}
	}

	if len(lines) < 4 {
		fmt.Println("Error getting new file. File is empty.")
		os.Exit(1)
	}

	//Removing strange lines
	lines = lines[4:]

	newFn, err2 := saveStringsToCsvFile(strings.Join(lines, "\n"), time.Now(), "")
	if err2 != nil {
		fmt.Println("Error creating new csv file. Error: " + err2.Error())
		os.Exit(1)
	}

	//Read csv
	newRsd, err3 := csvFileToRealStateData(newFn)
	if err3 != nil {
		fmt.Println("Error reading new csv file. Error: " + err3.Error())
		os.Exit(1)
	}

	newRsd.filterRealStatesByCity()
	_, err4 := saveRecordsToCsvFile(newRsd.toCsvStringWithHeader(), time.Now(), "filtered")
	if err4 != nil {
		fmt.Println("Error creating new filtered csv file. Error: " + err4.Error())
		os.Exit(1)
	}

	yesterday := time.Now().AddDate(0, 0, -1)

	var oldFileName = getFilePath(yesterday, FILE_NAME)
	oldRsd, err5 := csvFileToRealStateData(oldFileName)
	if err5 != nil {
		fmt.Println("Error reading old csv file. Error: " + err5.Error())
		os.Exit(1)
	}

	dif := difference(newRsd.List, oldRsd.List)
	difRsd := NewRealStateData()
	difRsd.List = dif

	//save dif
	_, err6 := saveRecordsToCsvFile(difRsd.toCsvStringWithHeader(), time.Now(), "dif")
	if err6 != nil {
		fmt.Println("Error creating new filtered csv file. Error: " + err6.Error())
		os.Exit(1)
	}
}

func difference(slice1 []RealStateEntry, slice2 []RealStateEntry) []RealStateEntry {
	diffStr := []RealStateEntry{}

	for _, val1 := range slice1 { //nested for loop to check if two values are equal
		found := false
		// Iterate over slice2
		for _, val2 := range slice2 {
			if val1.Id == val2.Id {
				found = true
				break
			}
		}
		if !found {
			diffStr = append(diffStr, val1)
		}
	}

	return diffStr
}
