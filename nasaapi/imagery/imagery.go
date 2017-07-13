package imagery

import (
	"encoding/json"
	"landsat/nasaapi/types"
	"net/url"
	"time"
)

// Response is a response from NASA Imagery API.
type Response struct {
	CloudScore float32
	Date       time.Time
	URL        *url.URL
	ID         types.ID
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
		ID         ID `json:"id"`
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
