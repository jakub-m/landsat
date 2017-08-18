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
	list   bool
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
	flag.BoolVar(&args.list, "list", false, "List mode if set, otherwise Get mode")
	flag.Var(&args.begin, "begin", "Begin date")
	flag.Var(&args.end, "end", "End date")
	flag.Parse()

	if args.list {
		log.SetPrefix("list ")
	} else {
		log.SetPrefix("get ")
	}
}

func main() {
	if args.list {
		doList()
	} else {
		doGet()
	}
}

// List mode. List all the assets.
func doList() {
	req := &assets.Request{
		Lat:    float32(args.lat),
		Lon:    float32(args.lon),
		Begin:  args.begin.value,
		End:    args.end.value,
		APIKey: types.APIKey(args.apiKey),
	}

	log.Println("List mode")
	log.Println("API key:", req.APIKey)
	log.Println("Latitude:", req.Lat)
	log.Println("Longitude:", req.Lon)
	log.Println("Begin date:", req.Begin.String())
	log.Println("End date:", req.End.String())

	res, err := assets.Get(req)
	if err != nil {
		log.Fatalln(err)
	}

	log.Printf("Got %d assets", res.Count)
	for _, r := range res.Results {
		fmt.Printf("%s\t%g\t%g\t%s\n", r.ID, req.Lat, req.Lon, r.Date.Format(timeFormat))
	}
}

// Get assets as specified in stdin.
func doGet() {
	log.Println("Get mode")
}
