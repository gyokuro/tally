
# Implementations

## Backends

Currently there are two implementations of the `CabService` interface as defined in `service.go`:

*  Simple/ naive implementation: all cabs are stored in memory in a hashmap. Computation of the
nearest cab within location of radius R requires O(N) computations of the haversine distance, where
N is the number of cabs for a shard (we can shard by city and that will sufficiently limit the size of
N to orders of at most 10's of thousands for even a city like NYC).

*  MongoDb implementation: this uses mongodb as the database backend and spatial index. The
`2dsphere` index is used on a GeoJSON representation of the Cab struct in the database and
proximity queries are used for the 'within' computations.

MongoDb is used for the following reasons:
*  This application is actually write heavy because each cab is expected to send an update of
its locations at frequent intervals.  Because it's write heavy, backend datastores that also support
spatial indexing (e.g. CouchDb, Lucene, ElasticSearch) are not ideal candidates:
    *  CouchDb with lots of updates, would require frequent database truncation of the append-only write log.
    *  Lucene, ElasticSearch are optimized for reads of infrequently updated documents.
* The servers are likely to be sharded by cities serviced, so each database (corresponding to the backing b-tree
files) there aren't really going to be a massive number of objects in the collections.
* There's enough flexibility in creating different kinds of indexes that this design can evolve as requirements
change in the future.
* There are likely other kinds of objects that need to be stored in the database, like user accounts, ride information,
cab profiles, etc. that makes using a general purpose database with spatial indexing attractive.


## Alternative implementations

An alternative implementation involves computation of geohashes from the `(lng, lat)` coordinates of the
cabs positions.  The geohash can then be used as key in a sorted list or prefix tree (trie) that to index
the cabs by location.  Then a search of nearby cabs invole computation of a set of possible geohash prefixes
and use that for querying for a set of candidates:

* One can use MySQL to store the cabs indexed by the geohash. The query for proximity becomes
`select ... where left(key, N) in ('prefix1', 'prefix2'...).  This implementation using a SQL database is
not used because
    * Tedius object/relational mapping - extra code
    * Search by matching string prefixes is potentially not much faster than a native spatial index like MongoDb.
    * Computation of adjacent quads need more time to implement and verify and account for corner cases.
    * After candiates are selected, we potentially still have to filter by computing haversine distance.

* One can compute geohash and instead of encoding in alphanumeric (base32) representation, simply use a long (up
to 64 bit) to represent the hash.  This is then used as a key in a database where the index keys are sorted:
    * In-memory trie or a ordered tree that allows range scan to find all candidates inside a list of quads.
    * Redis sorted set where the score is the hash value
    * LevelDb which uses SS table (sorted string) as the backing store.  So the hashes are the keys in the db.

The problem with in-memory implementation is that they can be complex, since the api for GET and QUERY effectively
require *two* indexes: one by id and one by location.  This makes the implementation a bit more complex for the
purpose of the coding homework.  Also, as future requirements change, additional indexes may be required (e.g. by
different kinds cabs - town cars or cheap ones) which would make a hand-written implementation harder to maintain.
