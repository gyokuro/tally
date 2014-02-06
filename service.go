package tally

import (
	"errors"
)

var (
	ErrorNotFound = errors.New("Not found")
	ErrorBadParam = errors.New("Bad parameter")
)

// Location by lat, lng
type Location struct {
	Latitude  float64
	Longitude float64
}

// Typedef of Id, using 64 bit unsigned int.
type Id uint64

// Enumeration of units for distance
type DistanceUnit int

// Distance units
const (
	Meters DistanceUnit = iota
	Kilometers
	Feet
	Miles
)

// Cab strcture for minimum of id and location
type Cab struct {
	Id        Id      `json:"id"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

// Structure for capturing the query parameters for within or proximity computation
type GeoWithin struct {
	Center Location
	Radius float64
	Unit   DistanceUnit
	Limit  int
}

// Service interface implemented by various backend datastores.
// The http server requires an implementation of this interface.
type CabService interface {

	// Load a Cab by id.  If not found, a ErrorNotFound must be returned.
	Read(id Id) (Cab, error)

	// Insert or update a cab
	Upsert(cab Cab) error

	// Deletes the cab by Id
	Delete(id Id) error

	// Queries for list of cabs by location and radius.  If none, return empty list.
	Query(query GeoWithin) ([]Cab, error)

	// Delete all cabs
	DeleteAll() error

	// Performs any necessary clean up
	Close()
}

// Sanitizes the input based on the spec. For example,
// the default value of limit is 8, and default unit of
// distance is meters.
func Sanitize(q *GeoWithin) *GeoWithin {
	if q.Limit == 0 {
		q.Limit = 8
	}
	if q.Unit == 0 {
		q.Unit = Meters
	}
	return q
}
