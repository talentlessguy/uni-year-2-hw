package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"os"

	_ "github.com/lib/pq"
)

type Template struct {
	URL   string
	Title string
}

func main() {

	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		fmt.Println("Error connecting to the database:", err)
		return
	}
	defer db.Close()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			tmpl, err := template.ParseFiles("templates/search.html")
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			data := Template{
				URL:   "/regions.json",
				Title: "Select a Region and City in Estonia",
			}
			err = tmpl.Execute(w, data)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		} else {
			fs := http.FileServer(http.Dir("public"))
			fs.ServeHTTP(w, r)
		}
	})

	http.HandleFunc("/stops", func(w http.ResponseWriter, r *http.Request) {
		queryParams := r.URL.Query()
		region := queryParams.Get("region")
		rows, err := db.Query("SELECT stop_name FROM stops WHERE stop_area = $1", region)
		if err != nil {
			fmt.Println("Error querying stops:", err)
			return
		}
		defer rows.Close()
		var stops []string
		for rows.Next() {
			var stop string
			err := rows.Scan(&stop)
			if err != nil {
				fmt.Println("Error scanning stops:", err)
				return
			}
			if err != nil {
				fmt.Println("Error scanning stops:", err)
				return
			}
			// Append only unique stop names
			unique := true
			for _, s := range stops {
				if s == stop {
					unique = false
					break
				}
			}
			if unique {
				stops = append(stops, stop)
			}
		}
		jsonData, err := json.Marshal(stops)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(jsonData)
	})

	http.HandleFunc("/search", func(w http.ResponseWriter, r *http.Request) {
		queryParams := r.URL.Query()
		region := queryParams.Get("region")

		tmpl, err := template.ParseFiles("templates/search.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		data := Template{
			URL:   fmt.Sprintf("/stops?region=%s", region),
			Title: "Select a stop",
		}
		err = tmpl.Execute(w, data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	http.ListenAndServe(":8080", nil)
}
