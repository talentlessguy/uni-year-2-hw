package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	_ "github.com/lib/pq"
)

type Template struct {
	URL    string
	Title  string
	Action string
}

type Trip struct {
	RouteShortName string `json:"route_short_name"`
	TripID         string `json:"trip_id"`
	TripLongName   string `json:"trip_long_name"`
	ArrivalTime    string `json:"arrival_time"`
	DepartureTime  string `json:"departure_time"`
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
				URL:    "/regions.json",
				Title:  "Select a Region and City in Estonia",
				Action: "stops",
			}
			err = tmpl.Execute(w, data)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		} else {
			baseDir := "public"
			path := strings.TrimPrefix(r.URL.Path, "/")
			fullPath := filepath.Join(baseDir, path)

			if !strings.Contains(path, ".") {
				htmlPath := fullPath + ".html"
				_, err := os.Stat(htmlPath)
				if err == nil {
					http.ServeFile(w, r, htmlPath)
					return
				}
			}

			http.ServeFile(w, r, fullPath)
		}
	})

	http.HandleFunc("/api/stops", func(w http.ResponseWriter, r *http.Request) {
		region := r.URL.Query().Get("region")
		rows, err := db.Query("SELECT DISTINCT ON (stop_name) stop_name FROM stops WHERE stop_area = $1", region)
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

			stops = append(stops, stop)
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

	http.HandleFunc("/api/buses", func(w http.ResponseWriter, r *http.Request) {
		stop := r.URL.Query().Get("stop")

		rows, err := db.Query("SELECT DISTINCT ON (r.route_short_name) r.route_short_name, t.trip_id, t.trip_long_name, st.arrival_time, st.departure_time FROM stops s JOIN stop_times st ON s.stop_id = st.stop_id JOIN trips t ON st.trip_id = t.trip_id JOIN routes r ON t.route_id = r.route_id WHERE s.stop_name = $1 ORDER BY r.route_short_name, st.arrival_time", stop)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error querying buses: %v", err), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var trips []Trip

		for rows.Next() {
			var trip Trip
			err := rows.Scan(&trip.RouteShortName, &trip.TripID, &trip.TripLongName, &trip.ArrivalTime, &trip.DepartureTime)
			if err != nil {
				http.Error(w, fmt.Sprintf("Error scanning row: %v", err), http.StatusInternalServerError)
				return
			}
			trips = append(trips, trip)
		}

		if err = rows.Err(); err != nil {
			http.Error(w, fmt.Sprintf("Error iterating rows: %v", err), http.StatusInternalServerError)
			return
		}

		jsonResponse, err := json.Marshal(trips)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error marshaling JSON: %v", err), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(jsonResponse)

	})

	http.HandleFunc("/stops", func(w http.ResponseWriter, r *http.Request) {
		region := r.URL.Query().Get("value")

		tmpl, err := template.ParseFiles("templates/search.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		data := Template{
			URL:    fmt.Sprintf("/api/stops?region=%s", region),
			Title:  "Select a stop",
			Action: "buses",
		}
		err = tmpl.Execute(w, data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	http.ListenAndServe(":8080", nil)
}
