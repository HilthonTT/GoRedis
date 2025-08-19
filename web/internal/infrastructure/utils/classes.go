package utils

import (
	"goredis-web/internal/domain"
	"time"
)

func DueClasses(d time.Time, completed bool) string {
	if completed {
		return "text-gray-400 line-through"
	}
	now := time.Now().UTC()
	if d.Before(now) {
		return "text-red-600"
	}
	// within 48h → warn
	if d.Before(now.Add(48 * time.Hour)) {
		return "text-amber-600"
	}
	return "text-gray-600"
}

func PriorityBadgeClasses(p domain.Priority) string {
	cl := "bg-gray-100 text-gray-700 border border-gray-200"

	switch p {
	case domain.Low:
		cl = "bg-emerald-50 text-emerald-700 border border-emerald-200"
	case domain.Medium:
		cl = "bg-amber-50 text-amber-700 border border-amber-200"
	case domain.High:
		cl = "bg-orange-50 text-orange-700 border border-orange-200"
	case domain.Top:
		cl = "bg-red-50 text-red-700 border border-red-200"
	}
	return "inline-flex items-center rounded-full px-2.5 py-0.5 text-xs font-medium " + cl
}

func PriorityLabel(p domain.Priority) string {
	switch p {
	case domain.Low:
		return "Low"
	case domain.Medium:
		return "Medium"
	case domain.High:
		return "High"
	case domain.Top:
		return "Top"
	default:
		return "Normal"
	}
}
