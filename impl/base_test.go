package impl

import (
	"github.com/gyokuro/findcab"
	"testing"
)

var (
	cabs = []findcab.Cab{
		findcab.Cab{
			Id:        1,
			Longitude: -77.037852,
			Latitude:  38.898556,
		},
		findcab.Cab{
			Id:        2,
			Longitude: -77.037852,
			Latitude:  39.898557,
		},
	}

	locations = []findcab.Location{
		findcab.Location{
			Longitude: -77.043934,
			Latitude:  38.897147,
		}}
)

func locationOf(cab findcab.Cab) findcab.Location {
	return findcab.Location{
		Latitude:  cab.Latitude,
		Longitude: cab.Longitude}
}

func testUpsert(service findcab.CabService, test *testing.T) {
	for _, c := range cabs {
		err := service.Upsert(c)
		if err != nil {
			test.Error("Got error", err)
		}
	}
}

func testGet(service findcab.CabService, test *testing.T, id findcab.Id, expected *findcab.Cab) {
	found, err := service.Read(id)
	if expected != nil && err != nil {
		test.Error("Expecting", *expected, "but got error", err)
	}
	if expected != nil && *expected != found {
		test.Error("Expecting", *expected, "but found", found)
	}
	if expected == nil && err != findcab.ErrorNotFound {
		test.Error("Nothing found should always return ErrorNotFound")
	}
}

func testQuery(service findcab.CabService, test *testing.T,
	loc findcab.Location, radius float64, expected []findcab.Cab) {
	q := findcab.GeoWithin{
		Center: loc,
		Radius: radius,
		Unit:   findcab.Meters,
	}

	cabs, err := service.Query(q)
	if err != nil {
		test.Error("Got error", err)
	}

	if cabs == nil || len(cabs) != len(expected) {
		test.Error("Expect vs actual", expected, cabs)
	}

	for i, c := range cabs {
		if expected[i] != c {
			test.Error("Expecting", expected[i], "got", c)
		}
	}
}

func testDelete(service findcab.CabService, test *testing.T, id findcab.Id) {
	err := service.Delete(id)
	if err != nil {
		test.Error("Got error", err)
	}
}

func testDeleteAll(service findcab.CabService, test *testing.T) {
	err := service.DeleteAll()
	if err != nil {
		test.Error("Got error", err)
	}
}
