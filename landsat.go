package landsat

import (
	"fmt"
	"landsat/http"
	"landsat/nasaapi/imagery"
	"log"
	"time"
)

const (
	apiKey   = "2YzaD3F2NoQusPVLco3ETnhbAR9ehVngscQfgUsZ"
	dateOnly = "2006-01-02"
)

type scanRequest struct {
	lat float32
	lon float32
}

func GenerateScanChan() chan scanRequest {
	scanChan := make(chan scanRequest)

	go func() {
		scanChan <- scanRequest{
			lat: 54.401325,
			lon: 18.572122,
		}
		scanChan <- scanRequest{
			lat: 54.2278043,
			lon: 19.9084272,
		}
		close(scanChan)
	}()

	return scanChan
}

func ProcessScan(in <-chan scanRequest) {
	i := 0
	for req := range in {
		i++
		date, _ := time.Parse(dateOnly, "2016-04-15")
		imagery, err := imagery.Get(&imagery.Request{Lat: req.lat, Lon: req.lon, APIKey: apiKey, Date: date})
		if err != nil {
			log.Println(err)
			continue
		}
		log.Println(imagery)
		http.FetchImage(imagery.URL, fmt.Sprintf("image_%d.png", i))
	}
}
