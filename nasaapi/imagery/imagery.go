package imagery

import (
	"encoding/json"
	"fmt"
	"jakub-m/landsat/http"
	"jakub-m/landsat/nasaapi/types"
	"log"
	"net/url"
	"time"
)

// Request is NASA imagery API request.
type Request struct {
	Lat    float32
	Lon    float32
	APIKey types.APIKey
	Date   time.Time
}

// Response is a response from NASA Imagery API.
type Response struct {
	CloudScore float32
	Date       time.Time
	URL        *url.URL
	ID         types.ID
}

func (r *Response) String() string {
	return fmt.Sprintf("CloudScore: %f, Date: %s, URL: %s, ID: %s", r.CloudScore, r.Date, r.URL, r.ID)
}

func Get(req *Request) (*Response, error) {
	urlValues := url.Values{}
	urlValues.Set("lat", fmt.Sprintf("%f", req.Lat))
	urlValues.Set("lon", fmt.Sprintf("%f", req.Lon))
	urlValues.Set("api_key", fmt.Sprint(req.APIKey))
	urlValues.Set("date", req.Date.Format("2006-01-02"))
	urlValues.Set("cloud_score", "True")
	url, err := url.Parse("https://api.nasa.gov/planetary/earth/imagery?" + urlValues.Encode())
	log.Println(url)
	if err != nil {
		return nil, err
	}
	body, err := http.GetBody(url)
	if err != nil {
		return nil, err
	}
	response, err := UnmarshallResponse(body)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func UnmarshallResponse(data []byte) (*Response, error) {
	var response Response
	err := json.Unmarshal(data, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

func (i *Response) UnmarshalJSON(data []byte) error {
	var tmp struct {
		CloudScore float32 `json:"cloud_score"`
		Date       string
		URL        string
		ID         types.ID `json:"id"`
	}
	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}

	i.CloudScore = tmp.CloudScore
	if url, err := url.Parse(tmp.URL); err == nil {
		i.URL = url
	} else {
		return err
	}
	if d, err := time.Parse("2006-01-02T15:04:05", tmp.Date); err == nil {
		i.Date = d
	} else {
		return err
	}
	i.ID = tmp.ID
	return nil
}
