package impl

import (
	"github.com/gyokuro/tally"
	"math"
	"testing"
)

// Compares two values up to p places after the decimal point.
// This is effectively comparing two values rounded to p places after the decimal.
func check(a, b float64, p int) bool {
	scale := math.Pow10(p + 1)
	return math.Abs(a*scale-b*scale) < 5.
}

func TestHaversine(test *testing.T) {

	p1 := tally.Location{
		Longitude: -77.037852,
		Latitude:  38.898556,
	}

	p2 := tally.Location{
		Longitude: -77.043934,
		Latitude:  38.897147,
	}

	expected := 0.341
	d := Haversine(p1, p2, tally.Miles)
	if !check(expected, d, 3) {
		test.Error("Expect in miles", expected, d)
	}

	d = Haversine(p2, p1, tally.Miles)
	if !check(expected, d, 3) {
		test.Error("Expect in miles", expected, d)
	}

	expected = d * 5280.
	d = Haversine(p1, p2, tally.Feet)
	if !check(expected, d, 0) {
		test.Error("Expect in feet", expected, d)
	}

	expected = 0.549
	d = Haversine(p1, p2, tally.Kilometers)
	if !check(expected, d, 3) {
		test.Error("Expect in km", expected, d)
	}

	expected = 549.
	d = Haversine(p1, p2, tally.Meters)
	if !check(expected, d, 0) {
		test.Error("Expect in m", expected, d)
	}
}
