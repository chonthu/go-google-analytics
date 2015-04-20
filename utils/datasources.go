package utils

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
)

func ProcessCSV(filename string) (output []string, ok bool) {
	file, err := os.Open(filename)
	if err != nil {
		return
	}
	defer file.Close()
	reader := csv.NewReader(file)
	reader.Comma = ','
	reader.FieldsPerRecord = 5
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		fmt.Println("Record", record, "and has", len(record), "fields")
		// and we can iterate on top of that
		for a, b := range record {
			if a == 0 || a%5 == 0 {
				output = append(output, b)

			}

		}
	}

	ok = true
	return
}

func ToCSV(file *os.File, data [][]string) bool {
	writer := csv.NewWriter(file)
	if err := writer.WriteAll(data); err != nil {
		fmt.Println("Error:", err)
		return false
	}
	writer.Flush()
	return true
}
