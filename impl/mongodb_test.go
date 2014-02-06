package impl

import (
	"github.com/gyokuro/tally"
	"testing"
)

var (
	mongodb, dbErr = NewMongoDbCabService("localhost", "test", "cabs")
)

func TestMongoDbUpsert(test *testing.T) {
	testUpsert(mongodb, test)
}

func TestMongoDbGet(test *testing.T) {
	testGet(mongodb, test, cabs[0].Id, &cabs[0])
	testGet(mongodb, test, cabs[1].Id, &cabs[1])
}

func TestMongoDbQuery(test *testing.T) {
	testQuery(mongodb, test, locations[0], 1000., []tally.Cab{cabs[0]})
	testQuery(mongodb, test, locations[0], 500., []tally.Cab{})
}

func TestMongoDbDelete(test *testing.T) {
	testUpsert(mongodb, test) // make sure data is there
	testQuery(mongodb, test, locationOf(cabs[1]), 1000., []tally.Cab{cabs[1]})

	// Now delete cab
	testDelete(mongodb, test, cabs[1].Id)
	testGet(mongodb, test, cabs[1].Id, nil)

	// Shouldn't be able to get cabs[1] after it's deleted.
	testQuery(mongodb, test, locationOf(cabs[1]), 1000., []tally.Cab{})
}

func TestMongoDbDeleteAll(test *testing.T) {
	testDeleteAll(mongodb, test)
	testQuery(mongodb, test, locationOf(cabs[0]), 1000., []tally.Cab{})
	testQuery(mongodb, test, locationOf(cabs[1]), 1000., []tally.Cab{})
	testGet(mongodb, test, cabs[0].Id, nil)
	testGet(mongodb, test, cabs[1].Id, nil)

	testUpsert(mongodb, test)
}
