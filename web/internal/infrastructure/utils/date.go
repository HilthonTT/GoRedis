package utils

import "time"

func FormatDue(d time.Time) string {
	// ISO date (UTC) + short time when present
	if d.IsZero() {
		return "No due date"
	}
	return d.UTC().Format("2006-01-02 15:04")
}
