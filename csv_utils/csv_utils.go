package csvutils

import (
	"encoding/csv"
	"fmt"
	"os"
)

const (
	iFPath = "xinhuar.csv"
	oFPath = "xinhuar_no_duplicates.csv"
)

func RemoveDuplicates() {
	file, err := os.Open(iFPath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		panic(err)
	}
	if len(records) == 0 {
		panic("No records found in CSV file")
	}

	uniqueRecords := make(map[string][]string)
	header := records[0]

	for _, record := range records[1:] {
		title := record[0]
		if _, exists := uniqueRecords[title]; !exists {
			uniqueRecords[title] = record
		}
	}

	outFile, err := os.Create(oFPath)
	if err != nil {
		panic(err)
	}
	defer outFile.Close()

	writer := csv.NewWriter(outFile)
	defer writer.Flush()

	writer.Write(header)

	for _, record := range uniqueRecords {
		writer.Write(record)
	}

	fmt.Printf("Removed duplicates in %s and wrote to %s\n", iFPath, oFPath)
}
