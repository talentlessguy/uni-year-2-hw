package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"

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

	http.Handle("/", http.FileServer(http.Dir("public")))

	http.HandleFunc("/search", func(w http.ResponseWriter, r *http.Request) {
		queryParams := r.URL.Query()
		region := queryParams.Get("region")
		city := queryParams.Get("city")

		fmt.Fprintf(w, "City: %s, region: %s", city, region)
	})

	http.ListenAndServe(":8080", nil)
}
