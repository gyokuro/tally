package impl

import (
	"github.com/gyokuro/tally"
)

// Simple implementation of the CabService interface
// This implementation uses a hashmap and does a O(N) scan of all entries
// when computing the nearest neighbor.
type simpleCabService struct {
	cabs map[tally.Id]tally.Cab
}

// Constructor method.  Returns an instance of the simple service
func NewSimpleCabService() *simpleCabService {
	return &simpleCabService{
		cabs: make(map[tally.Id]tally.Cab),
	}
}

// Implements CabService
func (s *simpleCabService) Read(id tally.Id) (result tally.Cab, err error) {
	var exists bool
	if result, exists = s.cabs[id]; exists {
		return
	}
	err = tally.ErrorNotFound
	return
}

// Implements CabService
func (s *simpleCabService) Upsert(cab tally.Cab) (err error) {
	s.cabs[cab.Id] = cab
	return nil
}

// Implements CabService
func (s *simpleCabService) Delete(id tally.Id) (err error) {
	delete(s.cabs, id)
	return nil
}

// Implements CabService
func (s *simpleCabService) Query(q tally.GeoWithin) (cabs []tally.Cab, err error) {
	tally.Sanitize(&q)
	cabs = make([]tally.Cab, 0)
	for _, cab := range s.cabs {
		distance := Haversine(q.Center, tally.Location{
			Latitude:  cab.Latitude,
			Longitude: cab.Longitude,
		}, q.Unit)
		if distance <= q.Radius {
			cabs = append(cabs, cab)
		}
		if len(cabs) == q.Limit {
			return
		}
	}
	return
}

// Implements CabService
func (s *simpleCabService) DeleteAll() (err error) {
	s.cabs = make(map[tally.Id]tally.Cab)
	return
}

// Implements CabService
func (s *simpleCabService) Close() {
	// no op
}
