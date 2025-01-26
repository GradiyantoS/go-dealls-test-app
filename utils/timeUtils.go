package utils

import "time"

// TimePtr converts a time.Time value to a pointer
func TimePtr(t time.Time) *time.Time {
	return &t
}
