package utils

import "time"

func IsValidDate(layout string, date ...string) bool {
	for _, d := range date {
		if !ContainsEmptyString(d) {
			_, err := time.Parse(layout, d)
			if err != nil {
				return false
			}
		}
	}
	return true
}
