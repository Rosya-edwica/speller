package db

import (
	"encoding/csv"
	"os"
)

type Element struct {
	Id      string
	Name    string
	NewName string
}

func GetItems() (items []Element) {
	file := GetFile("data/input.csv")
	defer file.Close()
	reader := csv.NewReader(file)
	reader.FieldsPerRecord = -1
	data, err := reader.ReadAll()
	checkErr(err)
	for j, row := range data {
		if j == 0 {
			continue
		}
		items = append(items, Element{Id: row[0], Name: row[1]})
	}
	return
}

func GetFile(path string) *os.File {
	file, err := os.Open(path)
	checkErr(err)
	return file
}

func SaveResult(data []Element) {
	file, err := os.Create("data/output.csv")
	checkErr(err)
	defer file.Close()
	writer := csv.NewWriter(file)
	writer.Comma = '|'
	defer writer.Flush()

	headers := []string{"Id", "До исправления", "После исправления"}

	writer.Write(headers)

	for _, row := range data {
		writer.Write([]string{row.Id, row.Name, row.NewName})
	}
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
