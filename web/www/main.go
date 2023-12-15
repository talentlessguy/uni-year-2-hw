package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"html/template"
	"math"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	_ "github.com/lib/pq"
)

type Template struct {
	URL    string
	Title  string
	Action string
}

type Trip struct {
	RouteID        string `json:"route_id"`
	RouteShortName string `json:"route_short_name"`
	TripID         string `json:"trip_id"`
	TripLongName   string `json:"trip_long_name"`
	ArrivalTime    string `json:"arrival_time"`
	DepartureTime  string `json:"departure_time"`
	Lat            string `json:"stop_lat"`
	Lon            string `json:"stop_lon"`
}

type Stop struct {
	Name string  `json:"name"`
	ID   string  `json:"id"`
	Lat  float64 `json:"lat"`
	Lon  float64 `json:"lon"`
}

func haversine(lat1, lon1, lat2, lon2 float64) float64 {
	const R = 6371 // Radius of the Earth in km
	dLat := (lat2 - lat1) * math.Pi / 180.0
	dLon := (lon2 - lon1) * math.Pi / 180.0

	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(lat1*math.Pi/180.0)*math.Cos(lat2*math.Pi/180.0)*
			math.Sin(dLon/2)*math.Sin(dLon/2)

	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	return R * c // Distance in km
}

func findNearestStop(stops []Stop, userLat, userLon float64) Stop {
	var nearest Stop
	minDistance := math.MaxFloat64

	for _, stop := range stops {
		distance := haversine(userLat, userLon, stop.Lat, stop.Lon)
		if distance < minDistance {
			minDistance = distance
			nearest = stop
		}
	}

	return nearest
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

		rows, err := db.Query("SELECT DISTINCT ON (r.route_short_name) r.route_id, r.route_short_name, t.trip_id, t.trip_long_name, st.arrival_time, st.departure_time, s.stop_lat, s.stop_lon FROM stops s JOIN stop_times st ON s.stop_id = st.stop_id JOIN trips t ON st.trip_id = t.trip_id JOIN routes r ON t.route_id = r.route_id WHERE s.stop_name = $1 ORDER BY r.route_short_name, st.arrival_time", stop)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error querying buses: %v", err), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var trips []Trip

		for rows.Next() {
			var trip Trip
			err := rows.Scan(&trip.RouteID, &trip.RouteShortName, &trip.TripID, &trip.TripLongName, &trip.ArrivalTime, &trip.DepartureTime, &trip.Lat, &trip.Lon)
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

	http.HandleFunc("/api/arrivals", func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		routeID := query.Get("route_id")
		userLatStr := query.Get("user_lat")
		userLonStr := query.Get("user_lon")

		userLat, err := strconv.ParseFloat(userLatStr, 64)
		if err != nil {
			http.Error(w, "Invalid user latitude", http.StatusBadRequest)
			return
		}

		userLon, err := strconv.ParseFloat(userLonStr, 64)
		if err != nil {
			http.Error(w, "Invalid user longitude", http.StatusBadRequest)
			return
		}

		rows, err := db.Query(`
            SELECT DISTINCT s.stop_name, s.stop_id, s.stop_lat, s.stop_lon
            FROM stops s
            JOIN stop_times st ON s.stop_id = st.stop_id
            JOIN trips t ON st.trip_id = t.trip_id
            JOIN routes r ON t.route_id = r.route_id
            WHERE r.route_id = $1`, routeID)

		if err != nil {
			http.Error(w, fmt.Sprintf("Database query error: %s", err), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var stops []Stop
		for rows.Next() {
			var stop Stop
			if err := rows.Scan(&stop.Name, &stop.ID, &stop.Lat, &stop.Lon); err != nil {
				http.Error(w, "Error reading stops data", http.StatusInternalServerError)
				return
			}
			stops = append(stops, stop)
		}
		nearestStop := findNearestStop(stops, userLat, userLon)

		json.NewEncoder(w).Encode(nearestStop)
	})

	fmt.Println("Server listening on port 8080")
	http.ListenAndServe(":8080", nil)
}
