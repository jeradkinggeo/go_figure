package main

import (
	"encoding/json"
	"log"
	"math"
	"net/http"
	"strconv"
)


func haversine(lat1, lon1, lat2, lon2 float64) float64 {
	const R = 6371 

	lat1Rad := lat1 * math.Pi / 180
	lat2Rad := lat2 * math.Pi / 180
	dLat := (lat2 - lat1) * math.Pi / 180
	dLon := (lon2 - lon1) * math.Pi / 180

	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(lat1Rad)*math.Cos(lat2Rad)*
			math.Sin(dLon/2)*math.Sin(dLon/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return R * c
}

func haversineHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	lat1Str := query.Get("lat1")
	lon1Str := query.Get("lon1")
	lat2Str := query.Get("lat2")
	lon2Str := query.Get("lon2")

	if lat1Str == "" || lon1Str == "" || lat2Str == "" || lon2Str == "" {
		http.Error(w, "Missing query parameters. Expected lat1, lon1, lat2, lon2.", http.StatusBadRequest)
		return
	}

	lat1, err := strconv.ParseFloat(lat1Str, 64)
	if err != nil {
		http.Error(w, "Invalid value for lat1", http.StatusBadRequest)
		return
	}
	lon1, err := strconv.ParseFloat(lon1Str, 64)
	if err != nil {
		http.Error(w, "Invalid value for lon1", http.StatusBadRequest)
		return
	}
	lat2, err := strconv.ParseFloat(lat2Str, 64)
	if err != nil {
		http.Error(w, "Invalid value for lat2", http.StatusBadRequest)
		return
	}
	lon2, err := strconv.ParseFloat(lon2Str, 64)
	if err != nil {
		http.Error(w, "Invalid value for lon2", http.StatusBadRequest)
		return
	}

	distance := haversine(lat1, lon1, lat2, lon2)
	response := map[string]float64{"distance": distance}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func inputHandler(w http.ResponseWriter, r *http.Request) {
	html := `<!DOCTYPE html>
<html>
<head>
  <title>Haversine Calculator</title>
  <style>
    body { font-family: Arial, sans-serif; margin: 40px; }
    label { display: block; margin-top: 10px; }
    input[type="text"] { padding: 5px; width: 200px; }
    input[type="submit"] { margin-top: 20px; padding: 10px 15px; }
  </style>
</head>
<body>
  <h1>Haversine Calculator</h1>
  <p>Enter the coordinates for two points to calculate the distance between them.</p>
  <form action="/haversine" method="get">
    <label for="lat1">Latitude 1:</label>
    <input type="text" id="lat1" name="lat1" required>
    
    <label for="lon1">Longitude 1:</label>
    <input type="text" id="lon1" name="lon1" required>
    
    <label for="lat2">Latitude 2:</label>
    <input type="text" id="lat2" name="lat2" required>
    
    <label for="lon2">Longitude 2:</label>
    <input type="text" id="lon2" name="lon2" required>
    
    <br>
    <input type="submit" value="Calculate Distance">
  </form>
</body>
</html>`
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(html))
}

func main() {
	http.HandleFunc("/haversine", haversineHandler)
	http.HandleFunc("/input", inputHandler)

	log.Println("Server running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
