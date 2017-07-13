package assets

import (
	"encoding/json"
	"landsat/nasaapi/types"
	"time"
)

type Response struct {
	Count   int
	Results []Result
}

type Result struct {
	Date time.Time
	ID   types.ID
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
