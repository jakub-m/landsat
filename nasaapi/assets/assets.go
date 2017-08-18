package assets

import (
	"encoding/json"
	"fmt"
	"jakub-m/landsat/http"
	"jakub-m/landsat/nasaapi/types"
	l_time "jakub-m/landsat/time"
	"net/url"
	"time"
)

type Request struct {
	Lat    float32
	Lon    float32
	Begin  time.Time
	End    time.Time
	APIKey types.APIKey
}

type Response struct {
	Count   int
	Results []Result
}

type Result struct {
	Date time.Time
	ID   types.ID
}

func Get(req *Request) (*Response, error) {
	urlValues := url.Values{}
	urlValues.Add("lat", fmt.Sprintf("%f", req.Lat))
	urlValues.Add("lon", fmt.Sprintf("%f", req.Lon))
	urlValues.Add("begin", l_time.FormatDate(req.Begin))
	urlValues.Add("end", l_time.FormatDate(req.End))
	urlValues.Add("api_key", fmt.Sprint(req.APIKey))

	url, err := url.Parse("https://api.nasa.gov/planetary/earth/assets?" + urlValues.Encode())
	if err != nil {
		return nil, err
	}
	body, err := http.GetBody(url)
	if err != nil {
		return nil, err
	}
	resp, err := UnmarshallResponse(body)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func UnmarshallResponse(data []byte) (*Response, error) {
	var response Response
	err := json.Unmarshal(data, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

func (r *Response) UnmarshalJSON(data []byte) error {
	type tmpResult struct {
		Date string   `json:"date"`
		ID   types.ID `json:"id"`
	}
	var tmp struct {
		Count   int         `json:"count"`
		Results []tmpResult `json:"results"`
	}

	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}

	r.Count = tmp.Count
	r.Results = make([]Result, len(tmp.Results))
	for i, s := range tmp.Results {
		d, err := time.Parse("2006-01-02T15:04:05", s.Date)
		if err != nil {
			return err
		}

		r.Results[i] = Result{
			Date: d,
			ID:   s.ID,
		}
	}
	return nil
}
