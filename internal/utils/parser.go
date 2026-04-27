package utils

import "time"

func ParseToDate(s string) *time.Time {
	var (
		t time.Time
		layout = "01-2006"
	)

	if s != "" {
		t, _ = time.Parse(layout, s)
	}

	return &t
}