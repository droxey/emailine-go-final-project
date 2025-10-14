Description:
This project is a cli tool written in Go that fetches live weather data for one or more cities concurrently using the wttr.in API. It displays temperature, conditions, humidity, wind speed, and a 5-day forecast. It also saves a local history of the last 20 weather searches for easy reference in a local json file.

How to use:

Install Go 1.20 or higher

`git clone https://github.com/emalinegayhart/go-final-project.git`

`cd go-final-project`

Get weather for one or more cities:
`go run main.go London Paris Tokyo`

View search history:
`go run main.go --history`

<img width="1278" height="918" alt="image" src="https://github.com/user-attachments/assets/1a3ed409-6688-423c-95ef-2274a51d4956" />

