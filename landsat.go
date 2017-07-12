package landsat

import (
	"fmt"
	"io/ioutil"
	"landsat/nasaapi/imagery"
	"log"
	"net/http"
	"net/url"
)

const (
	apiKey = "2YzaD3F2NoQusPVLco3ETnhbAR9ehVngscQfgUsZ"
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
		close(scanChan)
	}()

	return scanChan
}

func ProcessScan(in <-chan scanRequest) {
	i := 0
	for req := range in {
		i++
		imagery, err := getImagery(req)
		if err != nil {
			log.Println(err)
			continue
		}
		log.Println(imagery)
		fetchImage(imagery.URL, fmt.Sprintf("image_%d.png", i))
	}
}

func getImagery(req scanRequest) (*imagery.Response, error) {
	urlString := fmt.Sprintf("https://api.nasa.gov/planetary/earth/imagery?lat=%f&lon=%f&date=2017-04-05&cloud_score=True&api_key=%s", req.lat, req.lon, apiKey)
	log.Println(urlString)
	url, err := url.Parse(urlString)
	if err != nil {
		return nil, err
	}
	body, err := httpGet(url)
	if err != nil {
		return nil, err
	}
	response, err := imagery.UnmarshallResponse(body)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func fetchImage(url *url.URL, filename string) error {
	body, err := httpGet(url)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(filename, body, 0644)
	if err != nil {
		return err
	}
	return nil
}

func httpGet(url *url.URL) ([]byte, error) {
	res, err := http.Get(url.String())
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		return nil, err
	}
	return body, nil
}
