package main

import (
	"fmt"
	"jakub-m/landsat/http"
	"jakub-m/landsat/nasaapi/assets"
	"jakub-m/landsat/nasaapi/imagery"
	"log"
	"time"
)

const (
	apiKey = "2YzaD3F2NoQusPVLco3ETnhbAR9ehVngscQfgUsZ"
)

func main() {
	req := &assets.Request{
		Lat:    54.401325,
		Lon:    18.572122,
		Begin:  time.Date(2014, time.August, 1, 0, 0, 0, 0, time.UTC),
		End:    time.Date(2016, time.August, 1, 0, 0, 0, 0, time.UTC),
		APIKey: apiKey,
	}
	res, err := assets.Get(req)
	if err != nil {
		panic(err)
	}
	for i, r := range res.Results {
		log.Println(r.Date)
		err := fetchForDate(req.Lat, req.Lon, r.Date, i)
		if err != nil {
			panic(err)
		}
	}
	// scanChan := landsat.GenerateScanChan()
	// landsat.ProcessScan(scanChan)
}

func fetchForDate(lat, lon float32, date time.Time, i int) error {
	req := &imagery.Request{
		Lat:    lat,
		Lon:    lon,
		Date:   date,
		APIKey: apiKey,
	}
	res, err := imagery.Get(req)
	if err != nil {
		return err
	}
	fileName := fmt.Sprintf("image_%d.png", i+1)
	log.Println(date, fileName)
	http.FetchImage(res.URL, fileName)
	return nil
}
