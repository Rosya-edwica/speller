package db

import (
	"encoding/csv"
	"os"
)

type Skill struct {
	CurrentName string
	FixedName string
}

func GetSkills() (skills []string) {
	file := GetFile("data/input.csv")
	defer file.Close()
	reader := csv.NewReader(file)
	reader.FieldsPerRecord = -1
	data, err := reader.ReadAll()
	checkErr(err)
	for _, row := range data {
		for _, col := range row {
			skills = append(skills, col)
		}
	}
	return
}


func GetFile(path string) *os.File {
	file, err := os.Open(path)
	checkErr(err)
	return file
}


func SaveResult(data []Skill) {
	file, err := os.Create("data/output.csv")
	checkErr(err)
	defer file.Close()
	writer := csv.NewWriter(file)
	writer.Comma = '|'
	defer writer.Flush()

	headers := []string{"До исправления", "После исправления"}

	writer.Write(headers)

	for _, row := range data {
		writer.Write([]string{row.CurrentName, row.FixedName})
	}
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}