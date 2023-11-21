package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-redis/redis"
	"github.com/izaakdale/vancouver-conditions-backend/pkg/api"
)

func main() {
	chronOpts, err := redis.ParseURL(os.Getenv("REDIS_URL"))
	if err != nil {
		log.Printf("error trying to connect to redis\n")
		return
	}
	chronCli := redis.NewClient(chronOpts)

	weatherApiEndpoint := os.Getenv("WEATHER_API_ENDPOINT")
	apiKey := os.Getenv("WEATHER_API_KEY")

	rec := api.Record{
		Data: []api.FullBody{},
	}

	for searchQuery, adds := range searchParams {
		s := searchQuery
		a := adds

		u := fmt.Sprintf("%s/%s?unitGroup=metric&key=%s&contentType=json", weatherApiEndpoint, s, apiKey)

		req, err := http.NewRequest(http.MethodGet, u, nil)
		if err != nil {
			log.Printf("error creating request for %s: %+v\n", s, err)
			return
		}
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Printf("error when fetching data from %s: %+v\n", s, err)
			return
		}
		var fb api.FullBody
		err = json.NewDecoder(resp.Body).Decode(&fb)
		if err != nil {
			log.Printf("error decoding response from %s: %+v\n", s, err)
			return
		}

		fb.ForecastUrl = a.forecastUrl
		fb.WebCamUrl = a.webCamUrl
		fb.Title = a.title
		fb.GoogleMapsUrl = a.googleMapsUrl

		rec.Data = append(rec.Data, fb)
	}

	bytes, err := json.Marshal(rec)
	if err != nil {
		log.Printf("error marshalling responses to bytes: %+v\n", err)
		return
	}

	err = chronCli.Set("latest-conditions", bytes, 0).Err()
	if err != nil {
		log.Printf("error setting data in redis: %+v\n", err)
		return
	}
}

// search params maps the api search query to any additional data
// I want included in the record which is not supplied by the api
var searchParams = map[string]additionalData{
	"whistler-blackcomb-mountain": {
		title:         "Whistler Blackcomb",
		webCamUrl:     "https://www.whistlerblackcomb.com/the-mountain/mountain-conditions/mountain-cams.aspx",
		forecastUrl:   "https://www.snow-forecast.com/resorts/Whistler-Blackcomb/6day/mid",
		googleMapsUrl: "https://maps.app.goo.gl/7YTvXnCQPS32mxE9A",
	},
	"mt-baker-washington": {
		title:         "Mount Baker",
		webCamUrl:     "https://www.snowstash.com/usa/washington/mt-baker/snow-cams",
		forecastUrl:   "https://www.snow-forecast.com/resorts/Mount-Baker/6day/mid",
		googleMapsUrl: "https://maps.app.goo.gl/gaqSji8YiTb8RacY6",
	},
	"20955-hemlock-valley-rd": {
		title:         "Sasquatch Mountain Resort",
		webCamUrl:     "https://sasquatchmountain.ca/weather-and-conditions/webcams/",
		forecastUrl:   "https://www.snow-forecast.com/resorts/HemlockResort/6day/mid",
		googleMapsUrl: "https://maps.app.goo.gl/o5CWVongU85nwqhT7",
	},
	"cypress-mountain-vancouver": {
		title:         "Cypress Mountain",
		webCamUrl:     "https://cypressmountain.com/downhill-conditions-and-cams",
		forecastUrl:   "https://www.snow-forecast.com/resorts/Cypress-Mountain/6day/mid",
		googleMapsUrl: "https://maps.app.goo.gl/pJkSrmDLMb4RikAd8",
	},
	// omitted due to weather reports being almost identical to cypress.
	// "seymour-mountain-vancouver": {
	// 	webCamUrl:   "https://www.youtube.com/watch?v=vLawo-FrBKk",
	// 	forecastUrl: "https://www.snow-forecast.com/resorts/Mount-Seymour/6day/mid",
	// },
	// "grouse-mountain-vancouver": {
	// 	webCamUrl:   "https://www.grousemountain.com/web-cams",
	// 	forecastUrl: "https://www.snow-forecast.com/resorts/Grouse-Mountain/6day/mid",
	// },
}

type additionalData struct {
	title         string
	webCamUrl     string
	forecastUrl   string
	googleMapsUrl string
}
