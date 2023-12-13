package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
)

func main() {
	// Open the CSV file
	file, err := os.Open("../fetch/data/gtfs/stops.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// Create a CSV reader
	reader := csv.NewReader(file)

	// Read the first line to find the index of "stop_area"
	header, err := reader.Read()
	if err != nil {
		panic(err)
	}

	stopAreaIndex := -1
	for i, column := range header {
		if column == "stop_area" {
			stopAreaIndex = i
			break
		}
	}

	if stopAreaIndex == -1 {
		fmt.Println("Column 'stop_area' not found.")
		return
	}

	// Map to store unique values
	uniqueValues := make(map[string]bool)

	// Read the rest of the rows
	for {
		record, err := reader.Read()
		if err != nil {
			break
		}
		uniqueValues[record[stopAreaIndex]] = true
	}

	// Convert the map's keys to a slice
	var result []string
	for key := range uniqueValues {
		result = append(result, key)
	}

	// Convert to JSON
	jsonData, err := json.Marshal(result)
	if err != nil {
		panic(err)
	}

	// Save output to a file
	err = os.WriteFile("output.json", jsonData, 0644)
	if err != nil {
		panic(err)
	}

	fmt.Println("Unique values saved to output.json")
}
