package time

import "time"

const (
	dateOnly = "2006-01-02"
)

func ParseDate(raw string) (time.Time, error) {
	t, err := time.Parse(dateOnly, raw)
	if err != nil {
		return time.Time{}, err
	}
	return t, nil
}

func FormatDate(t time.Time) string {
	return t.Format(dateOnly)
}
