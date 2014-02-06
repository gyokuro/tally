package impl

import (
	"github.com/gyokuro/tally"
	"testing"
)

var (
	simple = NewSimpleCabService()
)

func TestSimpleUpsert(test *testing.T) {
	testUpsert(simple, test)
}

func TestSimpleGet(test *testing.T) {
	testGet(simple, test, cabs[0].Id, &cabs[0])
	testGet(simple, test, cabs[1].Id, &cabs[1])
}

func TestSimpleQuery(test *testing.T) {
	testQuery(simple, test, locations[0], 1000., []tally.Cab{cabs[0]})
	testQuery(simple, test, locations[0], 500., []tally.Cab{})
}

func TestSimpleDelete(test *testing.T) {
	testUpsert(simple, test) // make sure data is there
	testQuery(simple, test, locationOf(cabs[1]), 1000., []tally.Cab{cabs[1]})

	// Now delete cab
	testDelete(simple, test, cabs[1].Id)
	testGet(simple, test, cabs[1].Id, nil)

	// Shouldn't be able to get cabs[1] after it's deleted.
	testQuery(simple, test, locationOf(cabs[1]), 1000., []tally.Cab{})
}

func TestSimpleDeleteAll(test *testing.T) {
	testDeleteAll(simple, test)
	testQuery(simple, test, locationOf(cabs[0]), 1000., []tally.Cab{})
	testQuery(simple, test, locationOf(cabs[1]), 1000., []tally.Cab{})
	testGet(simple, test, cabs[0].Id, nil)
	testGet(simple, test, cabs[1].Id, nil)
}
