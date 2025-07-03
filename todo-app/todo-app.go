package main

import (
	"encoding/csv"
	"fmt"
	"os"
)

func main() {
	file, err := os.Open("todo.csv")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, record := range records {
		for _, item := range record {
			fmt.Println(item)
		}
	}

	file, err = os.Create("todo.csv")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	data := [][]string{
		{"id", "description", "status"},
		{"1", "Learn go", "0"},
		{"2", "Working on projects", "1"},
	}

	if err := writer.WriteAll(data); err != nil {
		fmt.Println(err)
		return
	}
}
