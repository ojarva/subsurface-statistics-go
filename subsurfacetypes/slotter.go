package subsurfacetypes

import "time"

func DurationToSlot(duration time.Duration) string {
	switch {
	case duration == time.Duration(0):
		return "unknown"
	case duration < time.Duration(10*time.Minute):
		return "<10min"
	case duration < time.Duration(20*time.Minute):
		return "<20min"
	case duration < time.Duration(30*time.Minute):
		return "<30min"
	case duration < time.Duration(40*time.Minute):
		return "<40min"
	case duration < time.Duration(50*time.Minute):
		return "<50min"
	case duration < time.Duration(60*time.Minute):
		return "<1h"
	case duration < time.Duration(70*time.Minute):
		return "<1h10min"
	case duration < time.Duration(80*time.Minute):
		return "<1h20min"
	case duration < time.Duration(90*time.Minute):
		return "<1h30min"
	default:
		return ">1h30min"
	}
}

func MaxDepthToSlot(depth float64) string {
	switch {
	case depth == 0:
		return "unknown"
	case depth < 19:
		return "P1"
	case depth < 33:
		return "P2"
	case depth < 48:
		return "rec tmx"
	case depth < 56:
		return "nmx tmx"
	default:
		return "hypo tmx"
	}
}
func MeanDepthToSlot(depth float64) string {
	switch {
	case depth == 0:
		return "unknown"
	case depth < 10:
		return "<10m"
	case depth < 20:
		return "<20m"
	case depth < 30:
		return "<30m"
	case depth < 40:
		return "<40m"
	case depth < 50:
		return "<50m"
	case depth < 56:
		return "<56m"
	default:
		return ">56m"
	}
}

func TemperatureToSlot(temperature float64) string {
	switch {
	case temperature < 0:
		return "<0c"
	case temperature < 5:
		return "<5c"
	case temperature < 10:
		return "<10c"
	case temperature < 15:
		return "<15c"
	case temperature < 20:
		return "<20c"
	default:
		return ">20c"
	}
}
