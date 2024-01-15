package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
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

type Bus struct {
	ServiceID      string `json:"service_id"`
	TripLongName   string `json:"trip_long_name"`
	RouteShortName string `json:"route_short_name"`
	Departure      string `json:"departure"`
	Arrival        string `json:"arrival"`
}

type Stop struct {
	Name string `json:"name"`
	Area string `json:"area"`
	ID   string `json:"id"`
}

type StopWithDistance struct {
	Stop
	Lat      float64 `json:"lat"`
	Lon      float64 `json:"lon"`
	TripID   string  `json:"trip_id"`
	Distance float64 `json:"distance"`
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
		userLatStr := r.URL.Query().Get("user_lat")
		userLonStr := r.URL.Query().Get("user_lon")

		userTime := r.URL.Query().Get("user_time")

		userLat, _ := strconv.ParseFloat(userLatStr, 64)
		userLon, _ := strconv.ParseFloat(userLonStr, 64)

		query := `
		SELECT DISTINCT
			s.stop_name,
			s.stop_area,
			s.stop_lat,
			s.stop_lon,
			t.trip_id
		FROM
			stops s
		JOIN
			stop_times st ON s.stop_id = st.stop_id
		JOIN
			trips t ON st.trip_id = t.trip_id
		JOIN
			stop_times stTempo ON t.trip_id = stTempo.trip_id
		JOIN
			stops sTempo ON stTempo.stop_id = sTempo.stop_id
		WHERE
			t.trip_id IN (
				SELECT
					stInner.trip_id
				FROM
					stop_times stInner
				JOIN
					stops sInner ON stInner.stop_id = sInner.stop_id
				WHERE
					sInner.stop_name = $1
			)
			AND s.stop_area = $2
			AND st.stop_sequence < stTempo.stop_sequence
	`

		rows, err := db.Query(query, stopName, stopArea)
		if err != nil {
			log.Println("Error fetching closest stops:", err)
			http.Error(w, "Error fetching closest stops", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var stops []StopWithDistance

		for rows.Next() {
			var stop StopWithDistance
			err := rows.Scan(&stop.Name, &stop.Area, &stop.Lat, &stop.Lon, &stop.TripID)
			if err != nil {
				log.Println("Error scanning rows:", err)
				http.Error(w, "Error processing rows", http.StatusInternalServerError)
				return
			}

			stop.Distance = haversine(stop.Lat, stop.Lon, userLat, userLon)
			stops = append(stops, stop)
		}

		if len(stops) == 0 {
			http.Error(w, "No stops found for the given stop_name", http.StatusNotFound)
			return
		}

		closestStop := stops[0]
		for _, s := range stops {
			if s.Distance < closestStop.Distance {
				closestStop = s
			}
		}

		query = `
		SELECT DISTINCT t.service_id, t.trip_long_name, r.route_short_name, st1.arrival_time, st2.arrival_time
		FROM stop_times st1
		JOIN stops s1 ON st1.stop_id = s1.stop_id
		JOIN stop_times st2 ON st1.trip_id = st2.trip_id
		JOIN stops s2 ON st2.stop_id = s2.stop_id
		JOIN trips t ON st1.trip_id = t.trip_id
		JOIN routes r ON t.route_id = r.route_id
		WHERE s1.stop_name = $1
			AND s2.stop_name = $2
			AND s1.stop_area = $3
			AND s2.stop_area = $4
			AND st1.stop_sequence > st2.stop_sequence
			AND st1.departure_time >= $5
		ORDER BY st1.arrival_time
		LIMIT 5;
	`
		rows, err = db.Query(query, closestStop.Name, stopName, stopArea, closestStop.Area, userTime)
		fmt.Println(closestStop.Name, closestStop.Area, stopName, stopArea)
		if err != nil {
			log.Println("Error fetching bus data from the database:", err)
			http.Error(w, "Error fetching bus data", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var buses []Bus
		for rows.Next() {
			var bus Bus
			err := rows.Scan(&bus.ServiceID, &bus.TripLongName, &bus.RouteShortName, &bus.Arrival, &bus.Departure)
			if err != nil {
				log.Println("Error scanning rows:", err)
				http.Error(w, "Error processing rows", http.StatusInternalServerError)
				return
			}
			buses = append(buses, bus)
		}

		response := map[string]interface{}{
			"closest_stop": closestStop,
			"buses":        buses,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
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

	fmt.Println("Server listening on port 8080")
	http.ListenAndServe(":8080", nil)
}
