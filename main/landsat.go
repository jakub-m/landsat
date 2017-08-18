package main

import (
	"flag"
	"fmt"
	"io"
	"jakub-m/landsat/http"
	"jakub-m/landsat/nasaapi/assets"
	"jakub-m/landsat/nasaapi/imagery"
	"jakub-m/landsat/nasaapi/types"
	"log"
	"time"
)

const (
	dateFormat      = "2006-01-02"
	timeFormat      = "2006-01-02T15:04:05"
	assetLineFormat = "%s\t%g\t%g\t%s\n"
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

type coord struct {
	lat float64
	lon float64
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
	log.Println("List mode")
	req := &assets.Request{
		Lat:    float32(args.lat),
		Lon:    float32(args.lon),
		Begin:  args.begin.value,
		End:    args.end.value,
		APIKey: types.APIKey(args.apiKey),
	}

	log.Println("API key:", req.APIKey)
	log.Println("Lat:", req.Lat)
	log.Println("Lon:", req.Lon)
	log.Println("Begin date:", req.Begin.String())
	log.Println("End date:", req.End.String())

	res, err := assets.Get(req)
	if err != nil {
		log.Fatalln(err)
	}

	log.Printf("Got %d assets", res.Count)
	for _, r := range res.Results {
		printAsset(req, r)
	}
}

// Get assets as specified in stdin.
func doGet() {
	log.Println("Get mode")
	for {
		coord, time, err := scanAsset()
		if err == io.EOF {
			return
		}
		if err != nil {
			log.Fatalln(err)
		}
		err = getSingleAsset(*coord, *time)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func printAsset(req *assets.Request, res assets.Result) {
	fmt.Printf(assetLineFormat, res.ID, req.Lat, req.Lon, res.Date.Format(timeFormat))
}

func scanAsset() (*coord, *time.Time, error) {
	var idString, dateString string
	var lat, lon float64
	_, err := fmt.Scanf(assetLineFormat, &idString, &lat, &lon, &dateString)
	if err != nil {
		return nil, nil, err
	}
	t, err := time.Parse(timeFormat, dateString)
	if err != nil {
		return nil, nil, err
	}

	return &coord{lat: lat, lon: lon}, &t, nil
}

func getSingleAsset(coord coord, time time.Time) error {
	req := &imagery.Request{
		APIKey: types.APIKey(args.apiKey),
		Lat:    float32(coord.lat),
		Lon:    float32(coord.lon),
		Date:   time,
	}
	log.Println("Lat:", req.Lat)
	log.Println("Lon:", req.Lon)
	log.Println("API key:", req.APIKey)
	log.Println("Date:", req.Date)
	res, err := imagery.Get(req)
	if err != nil {
		return err
	}

	targetFileName := fmt.Sprintf("lat_%.4f_lon_%.4f_%s_cs_%.2f.png", req.Lat, req.Lon, req.Date.Format(timeFormat), res.CloudScore)
	log.Printf("%s %s cloud score %g file %s", res.ID, res.Date.Format(timeFormat), res.CloudScore, targetFileName)
	err = http.FetchImage(res.URL, targetFileName)
	return err
}
