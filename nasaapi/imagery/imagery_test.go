package imagery_test

import (
	"jakub-m/landsat/nasaapi/imagery"
	"net/url"
	"testing"
	"time"
)

func TestUnmarshallResponse(t *testing.T) {
	raw := `
    {
        "id": "LC8_L1T_TOA/LC81870232017095LGN00",
        "cloud_score": 0.0007715016371718557,
        "date": "2017-04-05T09:25:01",
        "url": "http://localhost"
    }
    `
	r, err := imagery.UnmarshallResponse([]byte(raw))
	if err != nil {
		t.Fatal(err)
	}
	if r == nil {
		t.Fail()
	}

	if r.CloudScore != 0.0007715016371718557 {
		t.Error("CloudScore", r.CloudScore)
	}

	if expectedDate, err := time.Parse(time.RFC3339, "2017-04-05T09:25:01Z"); err == nil {
		if r.Date != expectedDate {
			t.Error("Date", r.Date)
		}
	} else {
		t.Fatal(err)
	}
	if expectedUrl, err := url.Parse("http://localhost"); err == nil {
		if r.URL.String() != expectedUrl.String() {
			t.Error("Url", r.URL)
		}
	} else {
		t.Fatal(err)
	}
	if r.ID != "LC8_L1T_TOA/LC81870232017095LGN00" {
		t.Error("Id", r.ID)
	}
}
