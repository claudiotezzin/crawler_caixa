package main

import (
	"errors"
	"strings"
)

type RealStateData struct {
	List []RealStateEntry
}

var realStateEntrySize = 11

type RealStateEntry struct {
	Id          string
	State       string
	City        string
	Neiborhood  string
	Address     string
	Price       string
	Evaluation  string
	Discount    string
	Description string
	SalesMode   string
	Url         string
}

func NewRealStateEntry(data []string) (*RealStateEntry, error) {
	rs := &RealStateEntry{}

	if len(data) != realStateEntrySize {
		return nil, errors.New("data is not valid")
	}

	rs.Id = strings.TrimSpace(data[0])
	rs.State = strings.TrimSpace(data[1])
	rs.City = strings.TrimSpace(data[2])
	rs.Neiborhood = strings.TrimSpace(data[3])
	rs.Address = strings.TrimSpace(data[4])
	rs.Price = strings.TrimSpace(data[5])
	rs.Evaluation = strings.TrimSpace(data[6])
	rs.Discount = strings.TrimSpace(data[7])
	rs.Description = strings.TrimSpace(data[8])
	rs.SalesMode = strings.TrimSpace(data[9])
	rs.Url = strings.TrimSpace(data[10])

	return rs, nil
}

func NewRealStateData() *RealStateData {
	return &RealStateData{
		List: make([]RealStateEntry, 0),
	}
}

func (rsd *RealStateData) filterRealStatesByCity() {
	var rse []RealStateEntry

	cities := getDesiredCities()

	for _, v := range rsd.List {
		if contains(cities, v.City) {
			rse = append(rse, v)
		}
	}

	rsd.List = rse
}

func contains(sl []string, name string) bool {
	for _, value := range sl {
		if value == name {
			return true
		}
	}
	return false
}

func (rse RealStateEntry) toStringArray() []string {
	result := []string{
		rse.Id,
		rse.State,
		rse.City,
		rse.Neiborhood,
		rse.Address,
		rse.Price,
		rse.Evaluation,
		rse.Discount,
		rse.Description,
		rse.SalesMode,
		rse.Url,
	}

	return result
}

func getRealStateHeaderAsStringArray() []string {
	return []string{
		"Id",
		"State",
		"City",
		"Neiborhood",
		"Address",
		"Price",
		"Evaluation",
		"Discount",
		"Description",
		"SalesMode",
		"Url",
	}
}
