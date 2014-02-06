package impl

import (
	"github.com/gyokuro/tally"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
)

type MongoDbCabService struct {
	Url, Db, Collection string
	db                  *mgo.Database
	session             *mgo.Session
	collection          *mgo.Collection
}

// Special struct that accounts for how mongo indexes spatial data.
// Conversions to and from this struct are required to satisfy the service api.
type mgo_record struct {
	Id  tally.Id `bson:"_id"`
	Loc []float64  `bson:"loc"`
}

// Converts input Cab struct into a mongodb record
func to_mgo(cab *tally.Cab) *mgo_record {
	return &mgo_record{
		Id:  cab.Id,
		Loc: []float64{cab.Longitude, cab.Latitude},
	}
}

// Converts output Cab from input mongodb record
func from_mgo(r *mgo_record) tally.Cab {
	return tally.Cab{
		Id:        r.Id,
		Longitude: r.Loc[0],
		Latitude:  r.Loc[1],
	}
}

// Constructor, returns an instance of the service backed by mongodb using a 2dsphere spatial index.
// This also connects to the database and ensures that the spatial indexed is used.
func NewMongoDbCabService(url, db, collection string) (service *MongoDbCabService, err error) {
	service = &MongoDbCabService{
		Url:        url,
		Db:         db,
		Collection: collection,
	}

	// Connect to db
	service.session, err = mgo.Dial(url)
	if err != nil {
		return
	}
	service.db = service.session.DB(service.Db)
	service.collection = service.db.C(service.Collection)
	service.ensure2dIndex()
	return
}

// Makes sure the index is maintained on the loc property
func (s *MongoDbCabService) ensure2dIndex() {
	// 2d spatial index on 'loc'
	s.collection.EnsureIndex(mgo.Index{
		Key:      []string{"loc"},
		Unique:   false,
		DropDups: false,
		Name:     "2dsphere",
	})
}

// Implements CabService
func (s *MongoDbCabService) Read(id tally.Id) (found tally.Cab, err error) {
	result := mgo_record{}
	err = s.collection.FindId(id).One(&result)
	switch err {
	case mgo.ErrNotFound:
		err = tally.ErrorNotFound
		return
	case nil:
		found = from_mgo(&result)
		return
	default:
		return
	}
}

// Implements CabService
func (s *MongoDbCabService) Upsert(cab tally.Cab) (err error) {
	r := to_mgo(&cab)
	_, err = s.collection.UpsertId(cab.Id, *r)
	return
}

// Implements CabService
func (s *MongoDbCabService) Query(q tally.GeoWithin) (cabs []tally.Cab, err error) {
	return s.QueryIndexed(q)
}

// Implements CabService
func (s *MongoDbCabService) Delete(id tally.Id) (err error) {
	return s.collection.RemoveId(id)
}

// Implements CabService
func (s *MongoDbCabService) DeleteAll() (err error) {
	_, err = s.collection.RemoveAll(nil)
	s.ensure2dIndex() // Make sure we have the index!
	return
}

// Implements CabService
func (s *MongoDbCabService) Close() {
	s.session.Close()
}

// Uses the spatial index '2dsphere' for fast lookup of cabs by proximity.
func (s *MongoDbCabService) QueryIndexed(q tally.GeoWithin) (cabs []tally.Cab, err error) {
	tally.Sanitize(&q)
	cabs = make([]tally.Cab, 0)

	distance := 0. // in meters, per mongo api
	switch q.Unit {
	case tally.Kilometers:
		distance = q.Radius * 1000.
	case tally.Meters:
		distance = q.Radius
	case tally.Miles:
		distance = q.Radius * 1609.34
	case tally.Feet:
		distance = q.Radius * 0.3048
	}
	query := bson.M{
		"loc": bson.M{
			"$near": bson.M{
				"$geometry": bson.M{
					"type":        "Point",
					"coordinates": []float64{q.Center.Longitude, q.Center.Latitude},
				},
				"$maxDistance": distance,
			},
		},
	}

	itr := s.collection.Find(query).Limit(q.Limit).Iter()
	if itr.Err() != nil {
		err = itr.Err()
		return
	}
	defer itr.Close()

	cab := mgo_record{}
	for itr.Next(&cab) {
		cabs = append(cabs, from_mgo(&cab))
		if len(cabs) > q.Limit {
			return
		}
	}
	return
}

// Slow version where the spatial index isn't used.  This involves the scan of the entire collection
// with calculation of haversine distance for each point.
// Included here for testing.
func (s *MongoDbCabService) QueryUnindexed(q tally.GeoWithin) (cabs []tally.Cab, err error) {
	tally.Sanitize(&q)
	cabs = make([]tally.Cab, 0)

	itr := s.collection.Find(nil).Iter()
	if itr.Err() != nil {
		err = itr.Err()
		return
	}
	defer itr.Close()

	cab := mgo_record{}
	for itr.Next(&cab) {
		distance := Haversine(q.Center, tally.Location{
			Latitude:  cab.Loc[1],
			Longitude: cab.Loc[0],
		}, q.Unit)
		if distance <= q.Radius {
			cabs = append(cabs, from_mgo(&cab))
		}
		if len(cabs) == q.Limit {
			return
		}
	}
	return
}
