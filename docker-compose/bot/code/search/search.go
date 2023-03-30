package search

import (
	"fmt"
	"log"
	"github.com/gocolly/colly"
)

type Weather struct { 
	Temperature, Wind, Humidity, ChanceOfPrecipitation string 
} 
 
func Search(request string) (answer string, err error) { 
	
	c := colly.NewCollector()

	c.UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109.0.0.0 Safari/537.36" 

	weather := Weather{}
	c.OnHTML("div.UQt4rd", func(e *colly.HTMLElement) { 
		
		weather.Temperature = e.ChildText(".q8U8x")
		weather.ChanceOfPrecipitation = e.ChildText("#wob_pp")
		weather.Humidity = e.ChildText("#wob_hm")
		weather.Wind = e.ChildText("#wob_ws")
		
		answer = fmt.Sprintf("Температура воздуха %v °C,\n Вероятность выпадения осадков %v,\n Влажность %v,\n Скорость ветра %v  ", weather.Temperature, weather.ChanceOfPrecipitation, weather.Humidity, weather.Wind)

	})

	err = c.Visit(request) 
	if err != nil {
		log.Printf("failed to visit url: %v\n", err)
		return 
	}

	return
}