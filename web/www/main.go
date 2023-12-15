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
	URL    string
	Title  string
	Action string
}

type Trip struct {
	TripID        string `json:"trip_id"`
	TripHeadsign  string `json:"trip_headsign"`
	TripLongName  string `json:"trip_long_name"`
	RouteID       string `json:"route_id"`
	ArrivalTime   string `json:"arrival_time"`
	DepartureTime string `json:"departure_time"`
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
				Action: "search",
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

	http.HandleFunc("/api/stops", func(w http.ResponseWriter, r *http.Request) {
		region := r.URL.Query().Get("region")
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

	http.HandleFunc("/api/buses", func(w http.ResponseWriter, r *http.Request) {
		stop := r.URL.Query().Get("stop")

		rows, err := db.Query("SELECT t.trip_id, t.trip_headsign, t.trip_long_name, t.route_id, st.arrival_time, st.departure_time FROM stops s JOIN stop_times st ON s.stop_id = st.stop_id JOIN trips t ON st.trip_id = t.trip_id WHERE s.stop_name = $1", stop)

		if err != nil {
			http.Error(w, fmt.Sprintf("Error querying buses: %v", err), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var trips []Trip

		for rows.Next() {
			var trip Trip
			err := rows.Scan(&trip.TripID, &trip.TripHeadsign, &trip.TripLongName, &trip.RouteID, &trip.ArrivalTime, &trip.DepartureTime)
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

	http.HandleFunc("/search", func(w http.ResponseWriter, r *http.Request) {
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

	http.HandleFunc("/buses", func(w http.ResponseWriter, r *http.Request) {
		stop := r.URL.Query().Get("value")

		tmpl, err := template.ParseFiles("templates/search.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		data := Template{
			URL:    fmt.Sprintf("/api/buses?stop=%s", stop),
			Title:  "Select a bus",
			Action: "idk",
		}
		err = tmpl.Execute(w, data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	http.ListenAndServe(":8080", nil)
}
