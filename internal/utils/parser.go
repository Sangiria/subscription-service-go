package utils

import "time"

const layout = "01-2006"

func ParseToDate(s string) *time.Time {
	if s == "" {
		return nil
	}

	t, err := time.Parse(layout, s)
	if err != nil {
		return nil
	}

	return &t
}