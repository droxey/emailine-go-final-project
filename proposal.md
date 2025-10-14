# Proposal: Weather CLI Tool

I want to create a simple command-line tool that fetches weather data from a free public API.

It will:
- Get current weather for any city
- Display temperature, humidity, and weather conditions
- Show a 5-day forecast
- Save recent searches to a local file
- Use Go's unique features like goroutines for concurrent API calls

## Tooling
- Go
- `net/http` for making API calls  
- `encoding/json` for handling responses
- `golang.org/x/oauth2` for authentication (if needed later)
- Public weather API (no signup required)

## Go Features to Showcase
- Goroutines for fetching multiple cities concurrently
- Channels for coordinating API responsesini
- JSON unmarshaling for API data
- File I/O for saving search history
- Error handling and graceful fallbacks