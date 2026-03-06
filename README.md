# torii-weather

A [Torii](https://github.com/haratosan/torii) extension that returns current weather and a 3-day forecast using the free [Open-Meteo](https://open-meteo.com/) API.

## Requirements

- Go 1.24+

## Installation

Clone this repo into your Torii extensions directory and build:

```sh
cd torii/extensions
git clone https://github.com/haratosan/torii-weather.git
cd torii-weather && go build .
```

Torii will automatically detect the extension on the next start.

## Usage

The extension accepts latitude and longitude coordinates and returns current conditions (temperature, wind speed, humidity) plus a 3-day forecast.

No API key required -- Open-Meteo is free and open.

## License

MIT
