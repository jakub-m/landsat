package assets_test

import (
	"landsat/nasaapi/assets"
	"testing"
	"time"
)

func TestUnmarshalAssetsResponse(t *testing.T) {
	raw := `
    {
        "count": 3,
        "results": [
        {
            "date": "2014-02-04T03:30:01",
            "id": "LC8_L1T_TOA/LC81270592014035LGN00"
        },
        {
            "date": "2014-02-20T03:29:47",
            "id": "LC8_L1T_TOA/LC81270592014051LGN00"
        },
        {
            "date": "2014-03-08T03:29:33",
            "id": "LC8_L1T_TOA/LC81270592014067LGN00"
        }]
    }
    `

	r, err := assets.UnmarshallResponse([]byte(raw))
	if err != nil {
		t.Fatal(err)
	}
	if r == nil {
		t.Fail()
	}

	if r.Count != 3 {
		t.Error("Count", r.Count)
	}

	if len(r.Results) != 3 {
		t.Error("len results", len(r.Results))
	}

	result := r.Results[0]
	if result.ID != "LC8_L1T_TOA/LC81270592014035LGN00" {
		t.Error("ID", result.ID)
	}

	expectedDate, _ := time.Parse(time.RFC3339, "2014-02-04T03:30:01Z")
	if result.Date != expectedDate {
		t.Error("Date", result.Date)
	}
}
