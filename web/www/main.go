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

type Route struct {
	RouteID        string `json:"route_id"`
	RouteShortName string `json:"route_short_name"`
	RouteLongName  string `json:"route_long_name"`
}

type Stop struct {
	Name string `json:"name"`
	Area string `json:"area"`
	ID   string `json:"id"`
}

type StopWithDistance struct {
	Stop
	Lat    float64 `json:"lat"`
	Lon    float64 `json:"lon"`
	TripID string  `json:"trip_id"`
}

type StopTime struct {
	ArrivalTime string `json:"arrival_time"`
	TripID      string `json:"trip_id"`
	StopID      string `json:"stop_id"`
}

func haversine(lat1, lon1, lat2, lon2 float64) float64 {
	const earthRadius = 6371 // in kilometers

	// Convert latitude and longitude from degrees to radians
	dLat := (lat2 - lat1) * (math.Pi / 180.0)
	dLon := (lon2 - lon1) * (math.Pi / 180.0)

	// Convert latitudes to radians
	lat1 = lat1 * (math.Pi / 180.0)
	lat2 = lat2 * (math.Pi / 180.0)

	// Haversine formula
	a := math.Pow(math.Sin(dLat/2), 2) + math.Pow(math.Sin(dLon/2), 2)*math.Cos(lat1)*math.Cos(lat2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	distance := earthRadius * c

	return distance
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
		rows, err := db.Query("SELECT DISTINCT ON (stop_name) stop_name, stop_area, stop_id FROM stops WHERE stop_area = $1", region)
		if err != nil {
			fmt.Println("Error querying stops:", err)
			return
		}
		defer rows.Close()
		var stops []Stop
		for rows.Next() {
			var stop Stop
			err := rows.Scan(&stop.Name, &stop.Area, &stop.ID)
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
		stopName := r.URL.Query().Get("stop_name")
		stopArea := r.URL.Query().Get("stop_area")

		query := `
			SELECT r.route_id, r.route_short_name, r.route_long_name
			FROM routes r
			JOIN (
				SELECT route_id
				FROM (
					SELECT t.route_id, array_length(array_agg(t.trip_id), 1) AS trip_count
					FROM stops s
					JOIN stop_times st ON s.stop_id = st.stop_id
					JOIN trips t ON st.trip_id = t.trip_id
					WHERE s.stop_name = $1 AND s.stop_area = $2
					GROUP BY t.route_id, t.trip_id
				) AS trip_counts
				GROUP BY route_id
			) AS max_trips ON r.route_id = max_trips.route_id
		`
		rows, err := db.Query(query, stopName, stopArea)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error querying routes: %v", err), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var routes []Route

		for rows.Next() {
			var route Route
			err := rows.Scan(&route.RouteID, &route.RouteShortName, &route.RouteLongName)
			if err != nil {
				http.Error(w, fmt.Sprintf("Error scanning row: %v", err), http.StatusInternalServerError)
				return
			}
			routes = append(routes, route)
		}

		jsonResponse, err := json.Marshal(routes)
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

	http.HandleFunc("/api/nearest_stop", func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		routeID := query.Get("route_id")
		userLatStr := query.Get("user_lat")
		userLonStr := query.Get("user_lon")

		// Convert userLatStr and userLonStr to float64
		userLat, _ := strconv.ParseFloat(userLatStr, 64)
		userLon, _ := strconv.ParseFloat(userLonStr, 64)

		// Query the database to get stops for the given route
		rows, err := db.Query("SELECT trip_id, stop_id, stop_name, lat, lon FROM stops WHERE route_id = $1", routeID)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error querying stops: %v", err), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var closestStop StopWithDistance
		closestDistance := math.MaxFloat64

		for rows.Next() {
			var stop StopWithDistance
			err := rows.Scan(&stop.TripID, &stop.ID, &stop.Name, &stop.Lat, &stop.Lon)
			if err != nil {
				http.Error(w, fmt.Sprintf("Error scanning row: %v", err), http.StatusInternalServerError)
				return
			}

			// Calculate distance using Haversine formula
			distance := haversine(userLat, userLon, stop.Lat, stop.Lon)

			if distance < closestDistance {
				closestStop = stop
				closestDistance = distance
			}
		}

		// Return the closest stop as JSON
		jsonStop, err := json.Marshal(closestStop)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error marshaling JSON: %v", err), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(jsonStop)
	})

	http.HandleFunc("/api/schedule", func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		tripID := query.Get("trip_id")
		stopID := query.Get("stop_id")

		stopTimes, err := db.Query(`
            SELECT arrival_time, trip_id, stop_id
            FROM stop_times
            WHERE trip_id = $1 AND stop_id = $2`, tripID, stopID)

		if err != nil {
			http.Error(w, fmt.Sprintf("Database query error: %s", err), http.StatusInternalServerError)
			return
		}
		defer stopTimes.Close()

		var stopTimesList []StopTime
		for stopTimes.Next() {
			var stopTime StopTime
			if err := stopTimes.Scan(&stopTime.ArrivalTime, &stopTime.TripID, &stopTime.StopID); err != nil {
				http.Error(w, "Error reading stop times data", http.StatusInternalServerError)
				return
			}
			stopTimesList = append(stopTimesList, stopTime)
		}
		json.NewEncoder(w).Encode(stopTimesList)
	})

	fmt.Println("Server listening on port 8080")
	http.ListenAndServe(":8080", nil)
}
