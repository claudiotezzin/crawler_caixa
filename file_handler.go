package main

import (
	"encoding/csv"
	"errors"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

const FILE_FOLDER = "data"
const FILE_EXTENSION = "csv"
const FILE_NAME = "caixa_sp"
const DATA_FORMAT = "2006_01_02"
const DELIMITER = ';'

func getTodayFileIfExist(t time.Time, fnSufix string) ([]string, error) {
	fName := FILE_NAME
	if fnSufix != "" {
		fName += "_"
		fName += fnSufix
	}

	fn := getFilePath(t, fName)

	b, err := os.ReadFile(fn)
	if err != nil {
		return nil, errors.New("Could not read file. Error: " + err.Error())
	}

	return strings.Split(string(b), "\n"), nil
}

func downloadTodayFile() ([]string, error) {
	l := "https://venda-imoveis.caixa.gov.br/listaweb/Lista_imoveis_SP.csv"

	resp, err := http.Get(l)
	if err != nil {
		return nil, errors.New("Could not reach " + l)
	}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.New("Could not read file. Error: " + err.Error())
	}

	return strings.Split(string(b), "\n"), nil
}

func saveStringsToCsvFile(text string, t time.Time, fnSufix string) (string, error) {
	fName := FILE_NAME
	if fnSufix != "" {
		fName += "_"
		fName += fnSufix
	}

	fn := getFilePath(t, fName)
	f, err := os.Create(fn)

	if err != nil {
		return "", errors.New("Could not create csv file. Error: " + err.Error())
	}

	defer f.Close()

	_, err2 := f.WriteString(text)

	if err2 != nil {
		return "", errors.New("Could not write to csv file. Error: " + err2.Error())
	}

	return fn, nil
}

func csvFileToRealStateData(fn string) (*RealStateData, error) {
	// open CSV file
	f, err := os.Open(fn)
	if err != nil {
		return nil, errors.New("Could not open csv file. Error: " + err.Error())
	}

	defer f.Close()

	fr := getCsvReader(f)

	rsd := NewRealStateData()

	for {
		data, err := fr.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			continue
		}

		rse, err := NewRealStateEntry(data)
		if err != nil {
			return nil, errors.New("Could not create real state data. Error: " + err.Error())
		}

		rsd.List = append(rsd.List, *rse)
	}

	return rsd, nil
}

func saveRecordsToCsvFile(records [][]string, t time.Time, fnSufix string) (string, error) {
	fName := FILE_NAME
	if fnSufix != "" {
		fName += "_"
		fName += fnSufix
	}

	fn := getFilePath(t, fName)
	f, err := os.Create(fn)

	if err != nil {
		return "", errors.New("Could not create csv file. Error: " + err.Error())
	}

	defer f.Close()

	w := getCsvWriter(f)
	w.WriteAll(records)

	if err := w.Error(); err != nil {
		return "", errors.New("Could write to csv file. Error: " + err.Error())
	}

	return fn, nil
}

func (rs RealStateData) toCsvStringWithHeader() [][]string {
	var result = make([][]string, len(rs.List)+1)

	for i, v := range rs.List {
		rsArray := v.toStringArray()
		lineIdx := i + 1
		for j := 0; j < realStateEntrySize; j++ {
			result[lineIdx] = rsArray
		}
	}

	result[0] = getRealStateHeaderAsStringArray()

	return result
}

func getFilePath(t time.Time, fn string) string {
	return FILE_FOLDER + "/" + t.Format(DATA_FORMAT) + "_" + fn + "." + FILE_EXTENSION
}

func getCsvReader(f *os.File) *csv.Reader {
	fr := csv.NewReader(f)
	fr.Comma = DELIMITER
	return fr
}

func getCsvWriter(f *os.File) *csv.Writer {
	fw := csv.NewWriter(f)
	fw.Comma = DELIMITER
	return fw
}
