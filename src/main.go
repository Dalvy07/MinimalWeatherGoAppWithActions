package main

import (
	"embed"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

// Embed static files directly into the binary
//go:embed static/*
var staticFiles embed.FS

// Structures for storing weather data
type CurrentWeather struct {
	Condition  struct{ Text string `json:"text"` } `json:"condition"`
	TempC      float64 `json:"temp_c"`
	Humidity   int     `json:"humidity"`
	WindKph    float64 `json:"wind_kph"`
	PressureMb float64 `json:"pressure_mb"`
	FeelslikeC float64 `json:"feelslike_c"`
	LastUpdate string  `json:"last_updated"`
}

type Location struct {
	Name    string `json:"name"`
	Country string `json:"country"`
}

type WeatherResponse struct {
	Location Location      `json:"location"`
	Current  CurrentWeather `json:"current"`
}

type WeatherData struct {
	City        string  `json:"city"`
	Country     string  `json:"country"`
	Condition   string  `json:"condition"`
	Temperature float64 `json:"temperature"`
	Humidity    int     `json:"humidity"`
	WindSpeed   float64 `json:"windSpeed"`
	Pressure    float64 `json:"pressure"`
	Feelslike   float64 `json:"feelslike"`
	LastUpdated string  `json:"last_updated"`
}

// Dictionary with locations
var locations = map[string][]string{
	"Poland":        {"Warsaw", "Krakow", "Gdansk", "Wroclaw", "Poznan", "Lublin"},
	"Germany":       {"Berlin", "Munich", "Hamburg", "Frankfurt", "Cologne"},
	"France":        {"Paris", "Marseille", "Lyon", "Toulouse", "Nice"},
	"Great Britain": {"London", "Manchester", "Liverpool", "Birmingham", "Glasgow"},
	"Italy":         {"Rome", "Milan", "Naples", "Florence", "Venice"},
}

// API key is obtained from environment variable
var apiKey string
// Author information from environment or default value if not set
var authorName = os.Getenv("APP_AUTHOR")

func init() {
    // Check that apiKey was set during compilation
    if apiKey == "" {
        log.Fatal("API key not set. Incorrect application build.")
    }

	// Check that authorName was set during compilation
	if authorName == "" {
		authorName = "Unknown Author"
	}
}

func main() {
	// Get port from environment variable or use 3000 as default
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	// Get current date and time
	startTime := time.Now().Format("2006-01-02 15:04:05")

	// Configure routes
	http.HandleFunc("/", serveHome)
	http.HandleFunc("/api/countries", getCountries)
	http.HandleFunc("/api/cities/", getCities)
	http.HandleFunc("/api/weather", getWeather)

	// Start the server with enhanced logging
	fmt.Printf("Application started at: %s\n", startTime)
	fmt.Printf("Author: %s\n", authorName)
	fmt.Printf("Server listening on TCP port: %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

// Handler for root route - serves HTML page
func serveHome(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	
	content, err := staticFiles.ReadFile("static/index.html")
	if err != nil {
		http.Error(w, "Failed to read index.html file", http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write(content)
}

// Handler for static CSS and JS files
func serveStatic(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path[1:] // Remove leading slash
	content, err := staticFiles.ReadFile(path)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	
	// Determine MIME type based on file extension
	switch {
	case strings.HasSuffix(path, ".css"):
		w.Header().Set("Content-Type", "text/css")
	case strings.HasSuffix(path, ".js"):
		w.Header().Set("Content-Type", "application/javascript")
	}
	
	w.Write(content)
}

// Handler for API request to list of countries
func getCountries(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	
	countries := make([]string, 0, len(locations))
	for country := range locations {
		countries = append(countries, country)
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(countries)
}

// Handler for API request to list of cities for a specific country
func getCities(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	
	country := r.URL.Path[len("/api/cities/"):]
	cities, exists := locations[country]
	
	w.Header().Set("Content-Type", "application/json")
	if !exists {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "Country not found"})
		return
	}
	
	json.NewEncoder(w).Encode(cities)
}

// Handler for API request for weather
func getWeather(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	
	city := r.URL.Query().Get("city")
	country := r.URL.Query().Get("country")
	
	if city == "" || country == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "City and country parameters are required"})
		return
	}
	
	weatherData, err := getWeatherFromAPI(city)
	if err != nil {
		log.Printf("Error while fetching weather data: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Error while fetching weather data"})
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(weatherData)
}

// Function for getting weather data from API
func getWeatherFromAPI(city string) (WeatherData, error) {
	var weatherData WeatherData
	
	apiURL := fmt.Sprintf("http://api.weatherapi.com/v1/current.json?key=%s&q=%s&aqi=no", 
		apiKey, url.QueryEscape(city))
	
	resp, err := http.Get(apiURL)
	if err != nil {
		return weatherData, err
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return weatherData, fmt.Errorf("HTTP error! Status: %d", resp.StatusCode)
	}
	
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return weatherData, err
	}
	
	var weatherResp WeatherResponse
	if err := json.Unmarshal(body, &weatherResp); err != nil {
		return weatherData, err
	}
	
	weatherData = WeatherData{
		City:        weatherResp.Location.Name,
		Country:     weatherResp.Location.Country,
		Condition:   weatherResp.Current.Condition.Text,
		Temperature: weatherResp.Current.TempC,
		Humidity:    weatherResp.Current.Humidity,
		WindSpeed:   weatherResp.Current.WindKph,
		Pressure:    weatherResp.Current.PressureMb,
		Feelslike:   weatherResp.Current.FeelslikeC,
		LastUpdated: weatherResp.Current.LastUpdate,
	}
	
	return weatherData, nil
}

// Function to enable CORS
func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	(*w).Header().Set("Access-Control-Allow-Headers", "Content-Type")
}