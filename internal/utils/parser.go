package utils

import "time"

const layout = "01-2006"

func ParseToDate(s string) *time.Time {
	var t time.Time

	if s != "" {
		t, _ = time.Parse(layout, s)
	}

	return &t
}