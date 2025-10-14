package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

const (
	historyFile = "weather_history.json"
	baseURL     = "https://wttr.in/"
)

type WeatherData struct {
	Location    string    `json:"location"`
	Temperature string    `json:"temperature"`
	Condition   string    `json:"condition"`
	Humidity    string    `json:"humidity"`
	Wind        string    `json:"wind"`
	Forecast    []string  `json:"forecast"`
	Timestamp   time.Time `json:"timestamp"`
}

type SearchHistory struct {
	Searches []WeatherData `json:"searches"`
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <city1> [city2] [city3]...")
		fmt.Println("Example: go run main.go London Paris Tokyo")
		return
	}

	cities := os.Args[1:]
	
	if len(cities) == 1 && cities[0] == "--history" {
		showHistory()
		return
	}

	var wg sync.WaitGroup
	var mu sync.Mutex
	var allWeather []WeatherData

	for _, city := range cities {
		wg.Add(1)
		go func(cityName string) {
			defer wg.Done()
			weather := fetchWeather(cityName)
			if weather != nil {
				mu.Lock()
				allWeather = append(allWeather, *weather)
				mu.Unlock()
			}
		}(city)
	}

	wg.Wait()

	for _, weather := range allWeather {
		displayWeather(weather)
		fmt.Println()
	}

	if len(allWeather) > 0 {
		saveToHistory(allWeather)
	}
}

func fetchWeather(city string) *WeatherData {
	url := baseURL + city + "?format=j1"
	
	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("Error fetching weather for %s: %v\n", city, err)
		return nil
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading response for %s: %v\n", city, err)
		return nil
	}

	var apiResponse map[string]interface{}
	if err := json.Unmarshal(body, &apiResponse); err != nil {
		fmt.Printf("Error parsing JSON for %s: %v\n", city, err)
		return nil
	}

	current := apiResponse["current_condition"].([]interface{})[0].(map[string]interface{})
	weather := apiResponse["weather"].([]interface{})[0].(map[string]interface{})

	temp := current["temp_C"].(string)
	condition := current["weatherDesc"].([]interface{})[0].(map[string]interface{})["value"].(string)
	humidity := current["humidity"].(string)
	windSpeed := current["windspeedKmph"].(string)
	windDir := current["winddir16Point"].(string)

	var forecast []string
	for i := 1; i <= 5 && i < len(weather["hourly"].([]interface{})); i++ {
		hourly := weather["hourly"].([]interface{})[i].(map[string]interface{})
		time := hourly["time"].(string)
		temp := hourly["tempC"].(string)
		desc := hourly["weatherDesc"].([]interface{})[0].(map[string]interface{})["value"].(string)
		forecast = append(forecast, fmt.Sprintf("%s: %s°C, %s", time, temp, desc))
	}

	return &WeatherData{
		Location:    city,
		Temperature: temp + "°C",
		Condition:   condition,
		Humidity:    humidity + "%",
		Wind:        windSpeed + " km/h " + windDir,
		Forecast:    forecast,
		Timestamp:   time.Now(),
	}
}

func displayWeather(weather WeatherData) {
	fmt.Printf("Weather for %s\n", strings.Title(weather.Location))
	fmt.Println(strings.Repeat("=", 30))
	fmt.Printf("Current: %s, %s\n", weather.Temperature, weather.Condition)
	fmt.Printf("Humidity: %s\n", weather.Humidity)
	fmt.Printf("Wind: %s\n", weather.Wind)
	
	if len(weather.Forecast) > 0 {
		fmt.Println("\n5-Day Forecast:")
		for i, forecast := range weather.Forecast {
			if i >= 5 {
				break
			}
			fmt.Printf("  %s\n", forecast)
		}
	}
}

func saveToHistory(weatherData []WeatherData) {
	var history SearchHistory
	
	if data, err := os.ReadFile(historyFile); err == nil {
		json.Unmarshal(data, &history)
	}

	history.Searches = append(history.Searches, weatherData...)
	
	if len(history.Searches) > 20 {
		history.Searches = history.Searches[len(history.Searches)-20:]
	}

	data, err := json.MarshalIndent(history, "", "  ")
	if err != nil {
		fmt.Printf("Error saving history: %v\n", err)
		return
	}

	os.WriteFile(historyFile, data, 0644)
}

func showHistory() {
	data, err := os.ReadFile(historyFile)
	if err != nil {
		fmt.Println("No search history found.")
		return
	}

	var history SearchHistory
	if err := json.Unmarshal(data, &history); err != nil {
		fmt.Printf("Error reading history: %v\n", err)
		return
	}

	fmt.Println("Recent Weather Searches:")
	fmt.Println(strings.Repeat("=", 30))
	
	for i, search := range history.Searches {
		if i >= 10 {
			break
		}
		fmt.Printf("%d. %s - %s (%s)\n", i+1, search.Location, search.Temperature, search.Timestamp.Format("2025-10-13 15:04"))
	}
}
