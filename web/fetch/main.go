package main

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	_ "github.com/lib/pq"
)

func main() {
	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		fmt.Println("Error connecting to the database:", err)
		return
	}
	defer db.Close()

	files := []string{}

	err = filepath.Walk("data/gtfs", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && filepath.Ext(path) == ".txt" {
			files = append(files, path)
		}
		return nil
	})

	if err != nil {
		fmt.Println("Error walking through directory:", err)
		return
	}

	for _, file := range files {
		f, err := os.Open(file)
		if err != nil {
			fmt.Println("Error opening file:", err)
			return
		}
		defer f.Close()

		reader := csv.NewReader(f)
		records, err := reader.ReadAll()
		if err != nil {
			fmt.Println("Error reading CSV:", err)
			return
		}

		tableName := strings.TrimSuffix(filepath.Base(file), filepath.Ext(file))

		// Creating a table based on CSV headers
		createTableStatement := "CREATE TABLE IF NOT EXISTS " + tableName + "("
		for i, header := range records[0] {
			createTableStatement += header + " TEXT"
			if i < len(records[0])-1 {
				createTableStatement += ", "
			}
		}
		createTableStatement += ");"

		_, err = db.Exec(createTableStatement)
		if err != nil {
			fmt.Println("Error creating table:", err)
			return
		}

		// Inserting data into the table
		for _, record := range records[1:] {
			insertStatement := "INSERT INTO " + tableName + " VALUES ('" + strings.Join(record, "','") + "');"
			_, err = db.Exec(insertStatement)
			if err != nil {
				fmt.Println("Error inserting record into table:", err)
				return
			}
		}
		fmt.Println("Data inserted into table:", tableName)
	}
}
