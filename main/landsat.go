package main

import (
	"flag"
	"fmt"
	"jakub-m/landsat/nasaapi/assets"
	"jakub-m/landsat/nasaapi/types"
	"log"
	"time"
)

const (
	dateFormat = "2006-01-02"
	timeFormat = "2006-01-02T15:04:05"
)

var args struct {
	apiKey string
	begin  date
	end    date
	lat    float64
	lon    float64
}

type date struct {
	value time.Time
}

func (d *date) Set(value string) error {
	time, err := time.Parse(dateFormat, value)
	if err != nil {
		return err
	}
	d.value = time
	return nil
}

func (d *date) String() string {
	return d.value.Format(dateFormat)
}

func init() {
	flag.StringVar(&args.apiKey, "api-key", "", "API key")
	flag.Float64Var(&args.lat, "lat", 0, "Latitude")
	flag.Float64Var(&args.lon, "lon", 0, "Longitude")
	flag.Var(&args.begin, "begin", "Begin date")
	flag.Var(&args.end, "end", "End date")
	flag.Parse()
}

func main() {
	req := &assets.Request{
		Lat:    float32(args.lat),
		Lon:    float32(args.lon),
		Begin:  args.begin.value,
		End:    args.end.value,
		APIKey: types.APIKey(args.apiKey),
	}

	res, err := assets.Get(req)
	if err != nil {
		log.Fatalln(err)
	}

	log.Println("API key:", args.apiKey)
	log.Println("Latitude:", args.lat)
	log.Println("Longitude:", args.lon)
	log.Println("Begin date:", args.begin.String())
	log.Println("End date:", args.end.String())
	log.Printf("Got %d assets: ", res.Count)
	for _, r := range res.Results {
		fmt.Printf("%s\t%g\t%g\t%s\n", r.ID, req.Lat, req.Lon, r.Date.Format(timeFormat))
	}
}
